package llm

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/logger"
	"context"

	"google.golang.org/genai"
)

// GeminiLiveSession 封装 Gemini SDK 的 Live 会话
type GeminiLiveSession struct {
	session *genai.Session
}

// ConnectLive 实现 LiveProvider 接口
func (a *GeminiAdapter) ConnectLive(ctx context.Context, cfg *LiveConfig) (LiveSession, error) {
	model := cfg.Model
	if model == "" {
		model = a.config.GetModel()
	}

	// 定义截图工具
	screenshotTool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "get_screenshot",
				Description: "获取用户当前屏幕截图，用于查看题目或界面内容",
			},
		},
	}

	// 连接配置
	connectCfg := &genai.LiveConnectConfig{
		Tools: []*genai.Tool{screenshotTool},
	}

	if cfg.SystemInstruction != "" {
		connectCfg.SystemInstruction = &genai.Content{
			Parts: []*genai.Part{{Text: cfg.SystemInstruction}},
		}
	}

	logger.Printf("LiveAPI: 连接到模型 %s", model)

	session, err := a.client.Live.Connect(ctx, model, connectCfg)
	if err != nil {
		return nil, err
	}

	return &GeminiLiveSession{session: session}, nil
}

// SendAudio 发送音频数据
func (s *GeminiLiveSession) SendAudio(data []byte) error {
	return s.session.SendRealtimeInput(genai.LiveRealtimeInput{
		Media: &genai.Blob{
			MIMEType: "audio/pcm",
			Data:     data,
		},
	})
}

// Receive 接收消息 (阻塞)
func (s *GeminiLiveSession) Receive() (*LiveMessage, error) {
	msg, err := s.session.Receive()
	if err != nil {
		return &LiveMessage{Type: LiveMsgError, Text: err.Error()}, err
	}

	return s.convertMessage(msg), nil
}

// convertMessage 转换 SDK 消息为统一格式
func (s *GeminiLiveSession) convertMessage(msg *genai.LiveServerMessage) *LiveMessage {
	if msg == nil {
		return nil
	}

	// 工具调用
	if msg.ToolCall != nil && len(msg.ToolCall.FunctionCalls) > 0 {
		fc := msg.ToolCall.FunctionCalls[0]
		return &LiveMessage{
			Type:     LiveMsgToolCall,
			ToolName: fc.Name,
			ToolID:   fc.ID,
		}
	}

	// 服务端消息
	if msg.ServerContent != nil {
		// 是否完成
		if msg.ServerContent.TurnComplete {
			return &LiveMessage{Type: LiveMsgDone}
		}

		// 检查 ModelTurn 中的文本
		if msg.ServerContent.ModelTurn != nil {
			for _, part := range msg.ServerContent.ModelTurn.Parts {
				if part != nil && part.Text != "" {
					return &LiveMessage{Type: LiveMsgAIText, Text: part.Text}
				}
			}
		}
	}

	return nil
}

// SendToolResponse 发送工具调用结果
func (s *GeminiLiveSession) SendToolResponse(toolID string, result string) error {
	return s.session.SendToolResponse(genai.LiveToolResponseInput{
		FunctionResponses: []*genai.FunctionResponse{
			{
				ID:       toolID,
				Response: map[string]any{"content": result},
			},
		},
	})
}

// Close 关闭会话
func (s *GeminiLiveSession) Close() error {
	return s.session.Close()
}

// 确保 GeminiAdapter 实现 LiveProvider 接口
var _ LiveProvider = (*GeminiAdapter)(nil)

// SupportsLive 检查是否支持 Live API
func SupportsLive(p Provider) bool {
	_, ok := p.(LiveProvider)
	return ok
}

// GetLiveConfig 从配置创建 LiveConfig
func GetLiveConfig(cfg *config.Config) *LiveConfig {
	return &LiveConfig{
		Model:             cfg.GetModel(),
		SystemInstruction: cfg.GetPrompt(),
	}
}
