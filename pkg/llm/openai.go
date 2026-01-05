package llm

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"Q-Solver/pkg/config"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// OpenAIAdapter OpenAI 适配器
type OpenAIAdapter struct {
	client     *openai.Client
	httpClient *http.Client
	config     *config.Config
}

// NewOpenAIAdapter 创建 OpenAI 适配器
func NewOpenAIAdapter(cfg *config.Config) *OpenAIAdapter {
	model := cfg.GetModel()
	if model == "" {
		model = openai.ChatModelGPT4o
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	opts := []option.RequestOption{
		option.WithAPIKey(cfg.GetAPIKey()),
		option.WithHTTPClient(httpClient),
	}

	if cfg.GetBaseURL() != "" {
		opts = append(opts, option.WithBaseURL(cfg.GetBaseURL()))
	}

	client := openai.NewClient(opts...)

	return &OpenAIAdapter{
		client:     &client,
		httpClient: httpClient,
		config:     cfg,
	}
}

// ==================== 类型转换方法 ====================

// toOpenAIMessages 将统一格式转换为 OpenAI SDK 格式
func (a *OpenAIAdapter) toOpenAIMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
	result := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			result = append(result, openai.SystemMessage(msg.Content))

		case RoleUser:
			if len(msg.Parts) > 0 {
				result = append(result, openai.UserMessage(a.toOpenAIParts(msg.Parts)))
			} else {
				result = append(result, openai.UserMessage(msg.Content))
			}

		case RoleAssistant:
			result = append(result, openai.AssistantMessage(msg.Content))
		}
	}

	return result
}

// toOpenAIParts 将 ContentPart 转换为 OpenAI 格式
func (a *OpenAIAdapter) toOpenAIParts(parts []ContentPart) []openai.ChatCompletionContentPartUnionParam {
	result := make([]openai.ChatCompletionContentPartUnionParam, 0, len(parts))

	for _, part := range parts {
		switch part.Type {
		case ContentText:
			result = append(result, openai.TextContentPart(part.Text))
		case ContentImage, ContentPDF:
			result = append(result, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
				URL: part.Base64,
			}))
		}
	}

	return result
}

// ==================== Provider 接口实现 ====================

// GenerateContentStream 流式生成内容
func (a *OpenAIAdapter) GenerateContentStream(ctx context.Context, messages []Message, onChunk StreamCallback) (Message, error) {
	openaiMessages := a.toOpenAIMessages(messages)

	stream := a.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:       a.config.GetModel(),
		Messages:    openaiMessages,
		Temperature: openai.Float(a.config.GetTemperature()),
		TopP:        openai.Float(a.config.GetTopP()),
		MaxTokens:   openai.Int(int64(a.config.GetMaxTokens())),
	})

	defer stream.Close()

	var fullContent strings.Builder
	var fullThinking strings.Builder

	for stream.Next() {
		evt := stream.Current()

		if len(evt.Choices) > 0 {
			delta := evt.Choices[0].Delta
			content := delta.Content

			if content != "" {
				fullContent.WriteString(content)

				if onChunk != nil {
					onChunk(StreamChunk{
						Type:    ChunkContent,
						Content: content,
					})
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		return Message{}, a.parseError(err)
	}

	return Message{
		Role:     RoleAssistant,
		Content:  fullContent.String(),
		Thinking: fullThinking.String(),
	}, nil
}

// parseError 解析错误信息
func (a *OpenAIAdapter) parseError(err error) error {
	errStr := err.Error()

	startIndex := strings.Index(errStr, "{")
	if startIndex == -1 {
		return fmt.Errorf("未知错误: %s", errStr)
	}

	jsonPart := errStr[startIndex:]
	var response struct {
		StatusCode int    `json:"statusCode"`
		Code       string `json:"code"`
		Message    string `json:"message"`
		Type       string `json:"type"`
	}

	_ = json.Unmarshal([]byte(jsonPart), &response)
	headerPart := errStr[:startIndex]

	lastColon := strings.LastIndex(headerPart, ":")
	if lastColon != -1 {
		statusPart := headerPart[lastColon+1:]
		if _, scanErr := fmt.Sscanf(statusPart, "%d", &response.StatusCode); scanErr != nil {
			response.StatusCode = 500
		}
	} else {
		response.StatusCode = -1
	}

	finalJsonBytes, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		return fmt.Errorf("解析错误: %s", response.Message)
	}

	return fmt.Errorf("%s", string(finalJsonBytes))
}

// TestChat 测试连通性
func (a *OpenAIAdapter) TestChat(ctx context.Context) error {
	_, err := a.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: a.config.GetModel(),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("hi"),
		},
		MaxTokens: openai.Int(1),
	})
	return err
}

// GetModels 获取模型列表
func (a *OpenAIAdapter) GetModels(ctx context.Context) ([]string, error) {
	resp, err := a.client.Models.List(ctx)
	if err != nil {
		return nil, err
	}
	var models []string
	for _, m := range resp.Data {
		models = append(models, m.ID)
	}
	return models, nil
}
