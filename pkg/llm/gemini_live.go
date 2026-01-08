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
	// Live API 只支持特定模型，用户配置的模型可能不支持 bidiGenerateContent
	model := cfg.Model
	if model == "" {
		model = a.config.Model
	}
	// model="gemini-2.5-flash-native-audio-preview-12-2025"
	// 定义截图工具
	screenshotTool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "get_screenshot",
				Description: "获取用户当前屏幕截图，用于查看题目或界面内容",
			},
		},
	}

	// 连接配置，使用 LiveConfig 中的参数
	connectCfg := &genai.LiveConnectConfig{
		Tools:                    []*genai.Tool{screenshotTool},
		ResponseModalities:       []genai.Modality{genai.ModalityAudio},
		MaxOutputTokens:          int32(cfg.MaxTokens),
		Temperature:              toFloat32Ptr(cfg.Temperature),
		TopP:                     toFloat32Ptr(cfg.TopP),
		TopK:                     intToFloat32Ptr(cfg.TopK),
		InputAudioTranscription:  &genai.AudioTranscriptionConfig{},
		SpeechConfig:             &genai.SpeechConfig{},
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},
		RealtimeInputConfig: &genai.RealtimeInputConfig{
			ActivityHandling: genai.ActivityHandlingNoInterruption,
		},
	}

	instructionText := cfg.SystemInstruction

	connectCfg.SystemInstruction = &genai.Content{
		Parts: []*genai.Part{{Text: instructionText}},
	}

	session, err := a.client.Live.Connect(ctx, model, connectCfg)
	if err != nil {
		logger.Printf("LiveAPI: 连接到模型 %s 发生错误", err)
		return nil, err
	}
	return &GeminiLiveSession{session: session}, nil
}

// SendAudio 发送音频数据 (16kHz, 16-bit, mono PCM)
func (s *GeminiLiveSession) SendAudio(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	return s.session.SendRealtimeInput(genai.LiveRealtimeInput{
		Media: &genai.Blob{
			MIMEType: "audio/pcm;rate=16000",
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
	// 输入音频转录 (面试官说话的文字)
	if msg.ServerContent != nil && msg.ServerContent.InputTranscription != nil {
		text := msg.ServerContent.InputTranscription.Text
		if text != "" {
			return &LiveMessage{Type: LiveMsgTranscript, Text: text}
		}
	}

	if msg.ServerContent != nil && msg.ServerContent.OutputTranscription != nil {
		text := msg.ServerContent.OutputTranscription.Text
		if text != "" {
			return &LiveMessage{Type: LiveMsgAIText, Text: text}
		}
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

		// 检查 ModelTurn 中的文本 (当 ResponseModalities 为 Text 时)
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

// SendToolResponse 发送工具调用结果 (文本)
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

// SendToolResponseWithImage 发送图片作为工具调用结果
func (s *GeminiLiveSession) SendToolResponseWithImage(toolID string, imageData []byte, mimeType string) error {
	logger.Printf("LiveAPI: 发送图片工具响应 ID=%s, size=%d, mime=%s", toolID, len(imageData), mimeType)
	return s.session.SendToolResponse(genai.LiveToolResponseInput{
		FunctionResponses: []*genai.FunctionResponse{
			{
				ID: toolID,
				Response: map[string]any{
					"image": map[string]any{
						"mimeType": mimeType,
						"data":     imageData,
					},
				},
			},
		},
	})
}

// Close 关闭会话
func (s *GeminiLiveSession) Close() error {
	return s.session.Close()
}

// SupportsLive 检查是否支持 Live API
func SupportsLive(p Provider) bool {
	_, ok := p.(LiveProvider)
	return ok
}

// GetLiveConfig 从配置创建 LiveConfig
func GetLiveConfig(cfg config.Config) *LiveConfig {
	return &LiveConfig{
		Model:             cfg.Model,
		SystemInstruction: cfg.Prompt,
		MaxTokens:         cfg.MaxTokens,
		Temperature:       cfg.Temperature,
		TopP:              cfg.TopP,
		TopK:              cfg.TopK,
	}
}

func toFloat32Ptr(v float64) *float32 {
	f := float32(v)
	return &f
}

func intToFloat32Ptr(v int) *float32 {
	f := float32(v)
	return &f
}
