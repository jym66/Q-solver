package config

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sync"

	"Q-Solver/pkg/logger"
	"Q-Solver/pkg/shortcut"
)

type ConfigManager struct {
	config      Config
	mu          sync.RWMutex
	configPath  string
	subscribers []func(Config)
}

func NewConfigManager() *ConfigManager {
	cm := &ConfigManager{
		config:      NewDefaultConfig(),
		subscribers: make([]func(Config), 0),
	}
	cm.configPath = cm.getConfigPath()
	return cm
}

func (cm *ConfigManager) getConfigPath() string {
	configDir := "."
	appDir := filepath.Join(configDir, "config")
	_ = os.MkdirAll(appDir, 0755)
	return filepath.Join(appDir, "config.json")
}

func (cm *ConfigManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.config = NewDefaultConfig()
	cm.loadFromEnv()

	if err := cm.loadFromFile(); err != nil {
		logger.Printf("加载配置文件失败 (使用默认配置): %v", err)
	}

	logger.Println("配置已加载")
	return nil
}

func (cm *ConfigManager) loadFromEnv() {
	if apiKey := os.Getenv("GHOST_API_KEY"); apiKey != "" {
		cm.config.APIKey = &apiKey
	}
	if baseURL := os.Getenv("GHOST_BASE_URL"); baseURL != "" {
		cm.config.BaseURL = &baseURL
	}
}

func (cm *ConfigManager) loadFromFile() error {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var fileConfig Config
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	cm.mergeConfig(&fileConfig)
	return nil
}

func (cm *ConfigManager) mergeConfig(src *Config) {
	// src is partial or full config from file (pointers)
	// We want to merge non-nil fields info cm.config (initialized with default)

	// Environment variable overrides file if env is set, but here logic was:
	// if src.APIKey != "" && env == "" -> set
	// This logic seems to say: File has priority UNLESS Env is set?
	// The original code:
	// if src.APIKey != "" && os.Getenv("GHOST_API_KEY") == "" { cm.config.APIKey = src.APIKey }
	// Meaning: Use file value only if Env is NOT set. (Env > File)
	// But `loadFromEnv` already set Env values into `cm.config`.
	// So if Env was set, `cm.config.APIKey` is already the Env value.
	// So we should only overwrite if `cm.config.APIKey` was NOT set by Env?
	// Actually `NewDefaultConfig` sets it to empty.
	// `loadFromEnv` sets it if Env exists.
	// So if Env exists, `cm.config.APIKey` is set.
	// If we blindly apply `src.APIKey` (from file), we overwrite Env.
	// So we check if Env is empty.

	if src.APIKey != nil && os.Getenv("GHOST_API_KEY") == "" {
		cm.config.APIKey = src.APIKey
	}
	if src.BaseURL != nil && os.Getenv("GHOST_BASE_URL") == "" {
		cm.config.BaseURL = src.BaseURL
	}

	// For other fields, just overwrite if non-nil
	if src.Model != nil {
		cm.config.Model = src.Model
	}
	if src.Provider != nil {
		cm.config.Provider = src.Provider
	}
	if src.Prompt != nil {
		cm.config.Prompt = src.Prompt
	}
	if src.Opacity != nil {
		cm.config.Opacity = src.Opacity
	}
	if src.ScreenshotMode != nil {
		cm.config.ScreenshotMode = src.ScreenshotMode
	}
	if src.CompressionQuality != nil {
		cm.config.CompressionQuality = src.CompressionQuality
	}
	if src.Sharpening != nil {
		cm.config.Sharpening = src.Sharpening
	}
	if src.ResumePath != nil {
		cm.config.ResumePath = src.ResumePath
	}
	if src.ResumeContent != nil {
		cm.config.ResumeContent = src.ResumeContent
	}

	if src.Grayscale != nil {
		cm.config.Grayscale = src.Grayscale
	}
	if src.NoCompression != nil {
		cm.config.NoCompression = src.NoCompression
	}
	if src.KeepContext != nil {
		cm.config.KeepContext = src.KeepContext
	}
	if src.InterruptThinking != nil {
		cm.config.InterruptThinking = src.InterruptThinking
	}
	if src.UseMarkdownResume != nil {
		cm.config.UseMarkdownResume = src.UseMarkdownResume
	}

	if len(src.Shortcuts) > 0 {
		if cm.config.Shortcuts == nil {
			cm.config.Shortcuts = make(map[string]shortcut.KeyBinding) // wait, importing shortcut?
			// The original code used "maps.Copy(cm.config.Shortcuts, src.Shortcuts)"
			// NewDefaultConfig initializes Shortcuts map.
		}
		maps.Copy(cm.config.Shortcuts, src.Shortcuts)
	}
}

func (cm *ConfigManager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	logger.Printf("配置已保存到: %s", cm.configPath)
	return nil
}

func (cm *ConfigManager) Get() Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

func (cm *ConfigManager) GetPtr() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return &cm.config
}

func (cm *ConfigManager) Update(patch Config) error {
	cm.mu.Lock()

	cm.applyConfig(patch)

	configCopy := cm.config
	subscribers := cm.subscribers

	cm.mu.Unlock()

	for _, sub := range subscribers {
		sub(configCopy)
	}

	return cm.Save()
}

func (cm *ConfigManager) UpdateFromJSON(jsonStr string) error {
	var patch Config
	if err := json.Unmarshal([]byte(jsonStr), &patch); err != nil {
		return fmt.Errorf("解析配置 JSON 失败: %w", err)
	}

	return cm.Update(patch)
}

func (cm *ConfigManager) applyConfig(patch Config) {
	if patch.APIKey != nil {
		cm.config.APIKey = patch.APIKey
	}
	if patch.Provider != nil {
		cm.config.Provider = patch.Provider
	}
	if patch.BaseURL != nil {
		cm.config.BaseURL = patch.BaseURL
	}
	if patch.Model != nil {
		cm.config.Model = patch.Model
	}
	if patch.Prompt != nil {
		cm.config.Prompt = patch.Prompt
	}
	if patch.Opacity != nil {
		cm.config.Opacity = patch.Opacity
	}
	if patch.ScreenshotMode != nil {
		cm.config.ScreenshotMode = patch.ScreenshotMode
	}
	if patch.CompressionQuality != nil {
		cm.config.CompressionQuality = patch.CompressionQuality
	}
	if patch.Sharpening != nil {
		cm.config.Sharpening = patch.Sharpening
	}
	if patch.Grayscale != nil {
		cm.config.Grayscale = patch.Grayscale
	}
	if patch.NoCompression != nil {
		cm.config.NoCompression = patch.NoCompression
	}
	if patch.KeepContext != nil {
		cm.config.KeepContext = patch.KeepContext
	}
	if patch.InterruptThinking != nil {
		cm.config.InterruptThinking = patch.InterruptThinking
	}
	if patch.ResumePath != nil {
		cm.config.ResumePath = patch.ResumePath
	}
	if patch.ResumeContent != nil {
		cm.config.ResumeContent = patch.ResumeContent
	}
	if patch.UseMarkdownResume != nil {
		cm.config.UseMarkdownResume = patch.UseMarkdownResume
	}
	if patch.Shortcuts != nil {
		for k, v := range patch.Shortcuts {
			cm.config.Shortcuts[k] = v
		}
	}
}

func (cm *ConfigManager) Subscribe(callback func(Config)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.subscribers = append(cm.subscribers, callback)
}

func (cm *ConfigManager) ClearResume() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	empty := ""
	cm.config.ResumePath = &empty
	cm.config.ResumeBase64 = &empty
	cm.config.ResumeContent = &empty
}

func (cm *ConfigManager) SetResumePath(path string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.ResumePath = &path
}

func (cm *ConfigManager) SetResumeBase64(base64 string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.ResumeBase64 = &base64
}

func (cm *ConfigManager) SetResumeContent(content string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.ResumeContent = &content
}
