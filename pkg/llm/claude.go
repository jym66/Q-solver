package llm

import (
	"context"
	"strings"

	"Q-Solver/pkg/config"
	"Q-Solver/pkg/logger"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ClaudeAdapter Claude 适配器
type ClaudeAdapter struct {
	client *anthropic.Client
	config *config.Config
}

// NewClaudeAdapter 创建 Claude 适配器
func NewClaudeAdapter(cfg *config.Config) *ClaudeAdapter {
	client := anthropic.NewClient(
		option.WithAPIKey(cfg.GetAPIKey()),
	)

	return &ClaudeAdapter{
		client: &client,
		config: cfg,
	}
}

// ==================== 类型转换方法 ====================

// toClaudeMessages 将统一格式转换为 Claude SDK 格式
func (a *ClaudeAdapter) toClaudeMessages(messages []Message) ([]anthropic.MessageParam, string) {
	var claudeMessages []anthropic.MessageParam
	var systemPrompt string

	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			systemPrompt = msg.Content

		case RoleUser:
			blocks := a.toClaudeBlocks(msg)
			claudeMessages = append(claudeMessages, anthropic.NewUserMessage(blocks...))

		case RoleAssistant:
			claudeMessages = append(claudeMessages, anthropic.NewAssistantMessage(
				anthropic.NewTextBlock(msg.Content),
			))
		}
	}

	return claudeMessages, systemPrompt
}

// toClaudeBlocks 将 Message 转换为 Claude ContentBlockParamUnion
func (a *ClaudeAdapter) toClaudeBlocks(msg Message) []anthropic.ContentBlockParamUnion {
	if len(msg.Parts) == 0 {
		return []anthropic.ContentBlockParamUnion{
			anthropic.NewTextBlock(msg.Content),
		}
	}

	var blocks []anthropic.ContentBlockParamUnion
	for _, p := range msg.Parts {
		switch p.Type {
		case ContentText:
			blocks = append(blocks, anthropic.NewTextBlock(p.Text))

		case ContentImage:
			// 解析 base64 图片
			mimeType, data := ParseBase64DataURL(p.Base64)
			if data != "" {
				blocks = append(blocks, anthropic.NewImageBlockBase64(mimeType, data))
			}

		case ContentPDF:
			// Claude 支持 PDF 作为文档
			_, data := ParseBase64DataURL(p.Base64)
			if data != "" {
				blocks = append(blocks, anthropic.ContentBlockParamUnion{
					OfDocument: &anthropic.DocumentBlockParam{
						Type: "document",
						Source: anthropic.DocumentBlockParamSourceUnion{
							OfBase64: &anthropic.Base64PDFSourceParam{
								Type:      "base64",
								MediaType: "application/pdf",
								Data:      data,
							},
						},
					},
				})
			}
		}
	}

	return blocks
}

// GenerateContentStream 流式生成内容
func (a *ClaudeAdapter) GenerateContentStream(ctx context.Context, messages []Message, onChunk StreamCallback) (Message, error) {
	claudeMessages, systemPrompt := a.toClaudeMessages(messages)

	model := a.config.GetModel()
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}

	params := anthropic.MessageNewParams{
		Model:       anthropic.Model(model),
		MaxTokens:   int64(a.config.GetMaxTokens()),
		Messages:    claudeMessages,
		Temperature: anthropic.Float(a.config.GetTemperature()),
		TopP:        anthropic.Float(a.config.GetTopP()),
		TopK:        anthropic.Int(int64(a.config.GetTopK())),
		Thinking:    anthropic.ThinkingConfigParamOfEnabled(int64(a.config.GetThinkingBudget())),
	}

	if systemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Type: "text", Text: systemPrompt},
		}
	}

	stream := a.client.Messages.NewStreaming(ctx, params)

	var fullContent strings.Builder
	var fullThinking strings.Builder

	for stream.Next() {
		evt := stream.Current()

		delta := evt.Delta
		if delta.Text != "" {
			fullContent.WriteString(delta.Text)
			if onChunk != nil {
				onChunk(StreamChunk{
					Type:    ChunkContent,
					Content: delta.Text,
				})
			}
		}
		if delta.Thinking != "" {
			fullThinking.WriteString(delta.Thinking)
			if onChunk != nil {
				onChunk(StreamChunk{
					Type:    ChunkThinking,
					Content: delta.Thinking,
				})
			}
		}
	}

	if err := stream.Err(); err != nil {
		return Message{}, err
	}

	return Message{
		Role:     RoleAssistant,
		Content:  fullContent.String(),
		Thinking: fullThinking.String(),
	}, nil
}

// TestChat 测试连通性
func (a *ClaudeAdapter) TestChat(ctx context.Context) error {
	model := a.config.GetModel()
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}

	_, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(model),
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	})
	return err
}

// GetModels 获取模型列表
func (a *ClaudeAdapter) GetModels(ctx context.Context) ([]string, error) {
	page, err := a.client.Models.List(ctx, anthropic.ModelListParams{})
	if err != nil {
		logger.Println("Claude获取模型错误", err)
	}
	var models []string
	for _, v := range page.Data {
		models = append(models, v.ID)
	}
	return models, nil
}
