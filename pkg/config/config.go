package config

import (
	"Q-Solver/pkg/shortcut"
	"encoding/json"
)

type Config struct {
	APIKey             *string                        `json:"apiKey,omitempty"`
	Provider           *string                        `json:"provider,omitempty"`
	Model              *string                        `json:"model,omitempty"`
	BaseURL            *string                        `json:"baseURL,omitempty"`
	Prompt             *string                        `json:"prompt,omitempty"`
	Opacity            *float64                       `json:"opacity,omitempty"`
	NoCompression      *bool                          `json:"noCompression,omitempty"`
	CompressionQuality *int                           `json:"compressionQuality,omitempty"`
	Sharpening         *float64                       `json:"sharpening,omitempty"`
	Grayscale          *bool                          `json:"grayscale,omitempty"`
	KeepContext        *bool                          `json:"keepContext,omitempty"`
	InterruptThinking  *bool                          `json:"interruptThinking,omitempty"`
	ScreenshotMode     *string                        `json:"screenshotMode,omitempty"`
	ResumePath         *string                        `json:"resumePath,omitempty"`
	ResumeBase64       *string                        `json:"-"`
	ResumeContent      *string                        `json:"resumeContent,omitempty"`
	UseMarkdownResume  *bool                          `json:"useMarkdownResume,omitempty"`
	Shortcuts          map[string]shortcut.KeyBinding `json:"shortcuts,omitempty"`
}

const DefaultModel = "gemini-2.5-flash"

func NewDefaultConfig() Config {
	return Config{
		APIKey:             ptr(""),
		Model:              ptr(DefaultModel),
		BaseURL:            ptr(""),
		ResumePath:         ptr(""),
		Prompt:             ptr(""),
		Opacity:            ptr(1.0),
		KeepContext:        ptr(false),
		InterruptThinking:  ptr(false),
		ScreenshotMode:     ptr("window"),
		NoCompression:      ptr(false),
		CompressionQuality: ptr(80),
		Sharpening:         ptr(0.0),
		Grayscale:          ptr(false),
		UseMarkdownResume:  ptr(false),
		ResumeBase64:       ptr(""),
		ResumeContent:      ptr(""),

		Shortcuts: map[string]shortcut.KeyBinding{
			"solve":        {ComboID: "119", KeyName: "F8"},
			"toggle":       {ComboID: "120", KeyName: "F9"},
			"clickthrough": {ComboID: "121", KeyName: "F10"},
			"move_up":      {ComboID: "38+164", KeyName: "Alt+↑"},
			"move_down":    {ComboID: "40+164", KeyName: "Alt+↓"},
			"move_left":    {ComboID: "37+164", KeyName: "Alt+←"},
			"move_right":   {ComboID: "39+164", KeyName: "Alt+→"},
			"scroll_up":    {ComboID: "33+164", KeyName: "Alt+PgUp"},
			"scroll_down":  {ComboID: "34+164", KeyName: "Alt+PgDn"},
		},
		Provider: ptr("google"), // Default provider
	}
}

func (c *Config) ToJSON() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}

func (c *Config) Validate() error {
	if c.ScreenshotMode != nil {
		if *c.ScreenshotMode != "fullscreen" && *c.ScreenshotMode != "window" {
			return &ValidationError{Field: "screenshotMode", Message: "截图模式必须是 'fullscreen' 或 'window'"}
		}
	}
	if c.Opacity != nil {
		if *c.Opacity < 0 || *c.Opacity > 1 {
			return &ValidationError{Field: "opacity", Message: "透明度必须在 0-1 之间"}
		}
	}
	if c.CompressionQuality != nil {
		if *c.CompressionQuality < 1 || *c.CompressionQuality > 100 {
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

// Helpers
func ptr[T any](v T) *T {
	return &v
}

// Getters
func (c Config) GetAPIKey() string {
	if c.APIKey == nil {
		return ""
	}
	return *c.APIKey
}
func (c Config) GetProvider() string {
	if c.Provider == nil {
		return ""
	}
	return *c.Provider
}
func (c Config) GetModel() string {
	if c.Model == nil {
		return ""
	}
	return *c.Model
}
func (c Config) GetBaseURL() string {
	if c.BaseURL == nil {
		return ""
	}
	return *c.BaseURL
}
func (c Config) GetPrompt() string {
	if c.Prompt == nil {
		return ""
	}
	return *c.Prompt
}
func (c Config) GetOpacity() float64 {
	if c.Opacity == nil {
		return 1.0
	}
	return *c.Opacity
}
func (c Config) GetNoCompression() bool {
	if c.NoCompression == nil {
		return false
	}
	return *c.NoCompression
}
func (c Config) GetCompressionQuality() int {
	if c.CompressionQuality == nil {
		return 80
	}
	return *c.CompressionQuality
}
func (c Config) GetSharpening() float64 {
	if c.Sharpening == nil {
		return 0.0
	}
	return *c.Sharpening
}
func (c Config) GetGrayscale() bool {
	if c.Grayscale == nil {
		return false
	}
	return *c.Grayscale
}
func (c Config) GetKeepContext() bool {
	if c.KeepContext == nil {
		return false
	}
	return *c.KeepContext
}
func (c Config) GetInterruptThinking() bool {
	if c.InterruptThinking == nil {
		return false
	}
	return *c.InterruptThinking
}
func (c Config) GetScreenshotMode() string {
	if c.ScreenshotMode == nil {
		return "window"
	}
	return *c.ScreenshotMode
}
func (c Config) GetResumePath() string {
	if c.ResumePath == nil {
		return ""
	}
	return *c.ResumePath
}
func (c Config) GetResumeBase64() string {
	if c.ResumeBase64 == nil {
		return ""
	}
	return *c.ResumeBase64
}
func (c Config) GetResumeContent() string {
	if c.ResumeContent == nil {
		return ""
	}
	return *c.ResumeContent
}
func (c Config) GetUseMarkdownResume() bool {
	if c.UseMarkdownResume == nil {
		return false
	}
	return *c.UseMarkdownResume
}
