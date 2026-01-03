package config

import "Q-Solver/pkg/shortcut"

type ConfigPatch struct {
	APIKey             *string                        `json:"apiKey,omitempty"`
	Provider           *string                        `json:"provider,omitempty"`
	BaseURL            *string                        `json:"baseURL,omitempty"`
	Model              *string                        `json:"model,omitempty"`
	Prompt             *string                        `json:"prompt,omitempty"`
	Opacity            *float64                       `json:"opacity,omitempty"`
	ScreenshotMode     *string                        `json:"screenshotMode,omitempty"`
	CompressionQuality *int                           `json:"compressionQuality,omitempty"`
	Sharpening         *float64                       `json:"sharpening,omitempty"`
	Grayscale          *bool                          `json:"grayscale,omitempty"`
	NoCompression      *bool                          `json:"noCompression,omitempty"`
	KeepContext        *bool                          `json:"keepContext,omitempty"`
	InterruptThinking  *bool                          `json:"interruptThinking,omitempty"`
	ResumePath         *string                        `json:"resumePath,omitempty"`
	ResumeContent      *string                        `json:"resumeContent,omitempty"`
	UseMarkdownResume  *bool                          `json:"useMarkdownResume,omitempty"`
	Shortcuts          map[string]shortcut.KeyBinding `json:"shortcuts,omitempty"`
}

func (p *ConfigPatch) Validate() error {
	if p.ScreenshotMode != nil {
		if *p.ScreenshotMode != "fullscreen" && *p.ScreenshotMode != "window" {
			return &ValidationError{Field: "screenshotMode", Message: "截图模式必须是 'fullscreen' 或 'window'"}
		}
	}
	if p.Opacity != nil {
		if *p.Opacity < 0 || *p.Opacity > 1 {
			return &ValidationError{Field: "opacity", Message: "透明度必须在 0-1 之间"}
		}
	}
	if p.CompressionQuality != nil {
		if *p.CompressionQuality < 1 || *p.CompressionQuality > 100 {
			return &ValidationError{Field: "compressionQuality", Message: "压缩质量必须在 1-100 之间"}
		}
	}
	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
