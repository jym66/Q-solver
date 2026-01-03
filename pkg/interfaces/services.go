package interfaces

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/screen"
	"Q-Solver/pkg/shortcut"
	"context"

	openai "github.com/openai/openai-go"
)

type LLMProvider interface {
	GenerateContentStream(ctx context.Context, history []openai.ChatCompletionMessageParamUnion, onToken func(string)) (string, error)
	GetBalance(ctx context.Context) (float64, error)
	Validate(ctx context.Context) error
	GetModels(ctx context.Context) ([]string, error)
	ParseResume(ctx context.Context, resumeBase64 string) (string, error)
}

type LLMService interface {
	GetProvider() LLMProvider
	UpdateProvider()
	GetBalance(ctx context.Context, apiKey string) (float64, error)
	ValidateAPIKey(ctx context.Context, apiKey string) string
	GetModels(ctx context.Context, apiKey string) ([]string, error)
}

type ScreenService interface {
	Startup(ctx context.Context)
	CapturePreview(quality int, sharpen float64, grayscale bool, noCompression bool, mode string) (screen.PreviewResult, error)
}

type ResumeService interface {
	SelectResume(ctx context.Context) string
	ClearResume()
	GetResumeBase64() (string, error)
	ParseResume(ctx context.Context, provider LLMProvider) (string, error)
}

type ShortcutDelegate interface {
	TriggerSolve()
	ToggleVisibility()
	ToggleClickThrough()
	MoveWindow(dx, dy int)
	ScrollContent(direction string)
	EmitEvent(eventName string, data ...interface{})
}

type ShortcutService interface {
	Start()
	Stop()
	GetShortcuts() map[string]shortcut.KeyBinding
	SetShortcuts(shortcuts map[string]shortcut.KeyBinding)
	StartRecording(action string)
	StopRecording()
}

type ConfigService interface {
	Load() error
	Save() error
	Get() config.Config
	GetPtr() *config.Config
	UpdateFromJSON(jsonStr string) error
	Subscribe(callback func(config.Config))
}

type StateService interface {
	GetInitStatusString() string
	UpdateInitStatus(status interface{})
	IsReady() bool
	ToggleVisibility() bool
	ToggleClickThrough() bool
	MoveWindow(dx, dy int)
	RestoreFocus()
	RemoveFocus()
}

type TaskService interface {
	StartTask(name string) (context.Context, int64)
	CompleteTask(taskID int64)
	CancelCurrentTask() bool
	IsTaskRunning(taskID int64) bool
	HasRunningTask() bool
}

type SolveRequest struct {
	Config           config.Config
	ScreenshotBase64 string
	ResumeBase64     string
}

type SolveCallbacks struct {
	EmitEvent func(event string, data ...interface{})
}

type Solver interface {
	Solve(ctx context.Context, req SolveRequest, cb SolveCallbacks) bool
	SetProvider(provider LLMProvider)
	ClearHistory()
}
