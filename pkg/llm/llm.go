package llm

import (
	"context"

	openai "github.com/openai/openai-go"
)

type Provider interface {
	GenerateContentStream(ctx context.Context, history []openai.ChatCompletionMessageParamUnion, onToken func(string)) (string, error)
	GetBalance(ctx context.Context) (float64, error)
	Validate(ctx context.Context) error
	GetModels(ctx context.Context) ([]string, error)
	ParseResume(ctx context.Context, resumeBase64 string) (string, error)
}
