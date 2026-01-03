package config

import (
	"Q-Solver/pkg/shortcut"
	"encoding/json"
)

type Config struct {
	APIKey             string                         `json:"apiKey"`
	Provider           string                         `json:"provider"`
	Model              string                         `json:"model"`
	BaseURL            string                         `json:"baseURL"`
	Prompt             string                         `json:"prompt"`
	Opacity            float64                        `json:"opacity"`
	NoCompression      bool                           `json:"noCompression"`
	CompressionQuality int                            `json:"compressionQuality"`
	Sharpening         float64                        `json:"sharpening"`
	Grayscale          bool                           `json:"grayscale"`
	KeepContext        bool                           `json:"keepContext"`
	InterruptThinking  bool                           `json:"interruptThinking"`
	ScreenshotMode     string                         `json:"screenshotMode"`
	ResumePath         string                         `json:"resumePath"`
	ResumeBase64       string                         `json:"-"`
	ResumeContent      string                         `json:"resumeContent"`
	UseMarkdownResume  bool                           `json:"useMarkdownResume"`
	Shortcuts          map[string]shortcut.KeyBinding `json:"shortcuts"`
}

const DefaultModel = "gemini-2.5-flash"

func NewDefaultConfig() Config {
	return Config{
		APIKey:            "",
		Model:             DefaultModel,
		BaseURL:           "",
		ResumePath:        "",
		Prompt:            "",
		Opacity:           1.0,
		KeepContext:       false,
		InterruptThinking: false,
		ScreenshotMode:    "window",

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
		ResumeBase64: "",
		Provider:     "google", // Default provider
	}
}

func (c *Config) ToJSON() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}
