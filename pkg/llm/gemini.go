package llm

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"Q-Solver/pkg/config"
	"Q-Solver/pkg/logger"

	"google.golang.org/genai"
)

// GeminiAdapter Gemini 适配器
type GeminiAdapter struct {
	client *genai.Client
	config *config.Config
}

// NewGeminiAdapter 创建 Gemini 适配器
func NewGeminiAdapter(cfg *config.Config) (*GeminiAdapter, error) {

	clientConfig := &genai.ClientConfig{
		APIKey:  cfg.GetAPIKey(),
		Backend: genai.BackendGeminiAPI,
	}
	//自定义的话就用自定义的URL
	if *cfg.Provider == "custom" {
		baseUrl := strings.TrimSuffix(cfg.GetBaseURL(),"/v1")
		logger.Println("配置Gemini自定义URL",baseUrl)
		clientConfig.HTTPOptions = genai.HTTPOptions{
			BaseURL: baseUrl,
		}
	}
	client, err := genai.NewClient(context.Background(), clientConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiAdapter{
		client: client,
		config: cfg,
	}, nil
}

// ==================== 类型转换方法 ====================

// toGeminiContents 将统一格式转换为 Gemini SDK 格式
func (a *GeminiAdapter) toGeminiContents(messages []Message) ([]*genai.Content, string) {
	var contents []*genai.Content
	var systemInstruction string

	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			systemInstruction = msg.Content

		case RoleUser:
			parts := a.toGeminiParts(msg)
			contents = append(contents, &genai.Content{
				Role:  "user",
				Parts: parts,
			})

		case RoleAssistant:
			contents = append(contents, &genai.Content{
				Role:  "model",
				Parts: []*genai.Part{{Text: msg.Content}},
			})
		}
	}

	return contents, systemInstruction
}

// toGeminiParts 将 Message 转换为 Gemini Parts
func (a *GeminiAdapter) toGeminiParts(msg Message) []*genai.Part {
	if len(msg.Parts) == 0 {
		return []*genai.Part{{Text: msg.Content}}
	}

	var parts []*genai.Part
	for _, p := range msg.Parts {
		switch p.Type {
		case ContentText:
			parts = append(parts, &genai.Part{Text: p.Text})

		case ContentImage:
			// 解析 base64 图片
			mimeType, data := parseBase64DataURL(p.Base64)
			if data != nil {
				parts = append(parts, &genai.Part{
					InlineData: &genai.Blob{
						MIMEType: mimeType,
						Data:     data,
					},
				})
			}

		case ContentPDF:
			// 解析 base64 PDF
			_, data := parseBase64DataURL(p.Base64)
			if data != nil {
				parts = append(parts, &genai.Part{
					InlineData: &genai.Blob{
						MIMEType: "application/pdf",
						Data:     data,
					},
				})
			}
		}
	}

	return parts
}

// parseBase64DataURL 解析 data:xxx;base64,... 格式
func parseBase64DataURL(dataURL string) (mimeType string, data []byte) {
	// 格式: data:image/png;base64,xxxxxx
	if !strings.HasPrefix(dataURL, "data:") {
		return "", nil
	}

	// 找到 base64, 的位置
	commaIdx := strings.Index(dataURL, ",")
	if commaIdx == -1 {
		return "", nil
	}

	// 解析 MIME 类型
	header := dataURL[5:commaIdx] // 去掉 "data:"
	semicolonIdx := strings.Index(header, ";")
	if semicolonIdx != -1 {
		mimeType = header[:semicolonIdx]
	} else {
		mimeType = header
	}

	// 解码 base64
	base64Data := dataURL[commaIdx+1:]
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", nil
	}

	return mimeType, decoded
}

// ==================== Provider 接口实现 ====================

// GenerateContentStream 流式生成内容
func (a *GeminiAdapter) GenerateContentStream(ctx context.Context, messages []Message, onChunk StreamCallback) (Message, error) {
	contents, systemInstruction := a.toGeminiContents(messages)

	model := a.config.GetModel()
	if model == "" {
		model = "gemini-2.0-flash"
	}

	maxTokens := int32(a.config.GetMaxTokens())
	temp := float32(a.config.GetTemperature())
	topP := float32(a.config.GetTopP())
	topK := float32(a.config.GetTopK())
	thinkingBudget := int32(a.config.GetThinkingBudget())

	genConfig := &genai.GenerateContentConfig{
		MaxOutputTokens: maxTokens,
		Temperature:     &temp,
		TopP:            &topP,
		TopK:            &topK,
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
			ThinkingBudget:  &thinkingBudget,
		},
	}
	if systemInstruction != "" {
		genConfig.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: systemInstruction}},
		}
	}

	var fullContent strings.Builder
	var fullThinking strings.Builder
	// s,_:=a.client.Live.Connect(ctx,"",nil)
	// s.Receive()
	for resp := range a.client.Models.GenerateContentStream(ctx, model, contents, genConfig) {
		if resp == nil {
			continue
		}

		// 遍历所有候选响应
		for _, candidate := range resp.Candidates {
			if candidate.Content == nil {
				continue
			}

			for _, part := range candidate.Content.Parts {
				if part == nil {
					continue
				}

				if part.Thought {
					fullThinking.WriteString(part.Text)
					if onChunk != nil {
						onChunk(StreamChunk{
							Type:    ChunkThinking,
							Content: part.Text,
						})
					}
				} else if part.Text != "" {
					fullContent.WriteString(part.Text)
					if onChunk != nil {
						onChunk(StreamChunk{
							Type:    ChunkContent,
							Content: part.Text,
						})
					}
				}
			}
		}
	}

	// 返回最终结果
	return Message{
		Role:     RoleAssistant,
		Content:  fullContent.String(),
		Thinking: fullThinking.String(),
	}, nil
}

// TestChat 测试连通性
func (a *GeminiAdapter) TestChat(ctx context.Context) error {
	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: []*genai.Part{{Text: "hi"}},
		},
	}

	config := &genai.GenerateContentConfig{
		MaxOutputTokens: 1,
	}

	model := a.config.GetModel()
	if model == "" {
		model = "gemini-2.0-flash"
	}

	_, err := a.client.Models.GenerateContent(ctx, model, contents, config)
	return err
}

// GetModels 获取模型列表
func (a *GeminiAdapter) GetModels(ctx context.Context) ([]string, error) {
	page, err := a.client.Models.List(ctx, &genai.ListModelsConfig{})
	if err != nil {
		return nil, err
	}

	var models []string
	for _, model := range page.Items {
		if model != nil {
			models = append(models, model.Name)
		}
	}

	return models, nil
}
