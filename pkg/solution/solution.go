package solution

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/llm"
	"Q-Solver/pkg/logger"
	"context"

	"github.com/openai/openai-go"
)

type Callbacks struct {
	EmitEvent func(event string, data ...interface{})
}

type Request struct {
	Config           config.Config
	ScreenshotBase64 string
	ResumeBase64     string
}

type Solver struct {
	llmProvider llm.Provider
	chatHistory []openai.ChatCompletionMessageParamUnion
}

func NewSolver(provider llm.Provider) *Solver {
	return &Solver{
		llmProvider: provider,
		chatHistory: make([]openai.ChatCompletionMessageParamUnion, 0),
	}
}

func (s *Solver) SetProvider(provider llm.Provider) {
	s.llmProvider = provider
}

func (s *Solver) ClearHistory() {
	s.chatHistory = make([]openai.ChatCompletionMessageParamUnion, 0)
}

func (s *Solver) Solve(ctx context.Context, req Request, cb Callbacks) bool {
	if req.Config.GetAPIKey() == "" {
		if cb.EmitEvent != nil {
			cb.EmitEvent("require-login")
		}
		return false
	}

	logger.Println("开始解题流程...")

	InitParts := []openai.ChatCompletionContentPartUnionParam{}
	if req.Config.GetPrompt() != "" {
		InitParts = append(InitParts, openai.TextContentPart(req.Config.GetPrompt()))
	}

	if req.Config.GetUseMarkdownResume() && req.Config.GetResumeContent() != "" {
		logger.Println("使用 Markdown 简历内容")
		InitParts = append(InitParts, openai.TextContentPart("\n\n# 候选人简历内容：\n"+req.Config.GetResumeContent()))
	} else if req.ResumeBase64 != "" {
		InitParts = append(InitParts, openai.TextContentPart("\n\n# 候选人简历已作为附件发送，请参考简历内容回答。"))
		InitParts = append(InitParts, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL: "data:application/pdf;base64," + req.ResumeBase64,
		}))
		logger.Println("已注入简历附件 (PDF)")
	}

	var messagesToSend []openai.ChatCompletionMessageParamUnion
	historyParts := []openai.ChatCompletionContentPartUnionParam{}
	historyParts = append(historyParts, openai.TextContentPart("[用户上传了一张图片，但是为了防止请求体过大，这里用文本替代，**你可以在之前的对话中找到图片内容的描述**。]"))

	userParts := []openai.ChatCompletionContentPartUnionParam{}
	userParts = append(userParts, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
		URL: req.ScreenshotBase64,
	}))

	currentUserMsg := openai.UserMessage(userParts)

	if len(InitParts) > 0 {
		messagesToSend = append(messagesToSend, openai.UserMessage(InitParts))
	}

	if req.Config.GetKeepContext() {
		if len(s.chatHistory) > 0 {
			messagesToSend = append(messagesToSend, s.chatHistory...)
		}
	} else {
		s.chatHistory = nil
	}
	messagesToSend = append(messagesToSend, currentUserMsg)

	if cb.EmitEvent != nil {
		cb.EmitEvent("solution-stream-start")
	}

	answer, err := s.llmProvider.GenerateContentStream(ctx, messagesToSend, func(token string) {
		if cb.EmitEvent != nil {
			cb.EmitEvent("solution-stream-chunk", token)
		}
	})

	if err != nil {
		if ctx.Err() == context.Canceled {
			logger.Println("当前任务已中断 (用户产生新输入)")
			if cb.EmitEvent != nil {
				cb.EmitEvent("solution-error", "context canceled")
			}
			return false
		}
		logger.Printf("LLM 请求失败: %v\n", err)
		if cb.EmitEvent != nil {
			cb.EmitEvent("solution-error", err.Error())
		}
		return false
	}

	if !req.Config.GetKeepContext() {
		if cb.EmitEvent != nil {
			cb.EmitEvent("solution", answer)
		}
		s.chatHistory = make([]openai.ChatCompletionMessageParamUnion, 0)
		return true
	}

	cleanUserMsg := openai.UserMessage(historyParts)
	s.chatHistory = append(s.chatHistory, cleanUserMsg)
	s.chatHistory = append(s.chatHistory, openai.AssistantMessage(answer))

	if cb.EmitEvent != nil {
		cb.EmitEvent("solution", answer)
	}
	return true
}
