package config

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sync"

	"Q-Solver/pkg/logger"
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
		cm.config.APIKey = apiKey
	}
	if baseURL := os.Getenv("GHOST_BASE_URL"); baseURL != "" {
		cm.config.BaseURL = baseURL
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
	if src.APIKey != "" && os.Getenv("GHOST_API_KEY") == "" {
		cm.config.APIKey = src.APIKey
	}
	if src.BaseURL != "" && os.Getenv("GHOST_BASE_URL") == "" {
		cm.config.BaseURL = src.BaseURL
	}
	if src.Model != "" {
		cm.config.Model = src.Model
	}
	if src.Prompt != "" {
		cm.config.Prompt = src.Prompt
	}
	if src.Opacity != 0 {
		cm.config.Opacity = src.Opacity
	}
	if src.ScreenshotMode != "" {
		cm.config.ScreenshotMode = src.ScreenshotMode
	}
	if src.CompressionQuality != 0 {
		cm.config.CompressionQuality = src.CompressionQuality
	}
	if src.Sharpening != 0 {
		cm.config.Sharpening = src.Sharpening
	}
	if src.ResumePath != "" {
		cm.config.ResumePath = src.ResumePath
	}
	if src.ResumeContent != "" {
		cm.config.ResumeContent = src.ResumeContent
	}

	cm.config.Grayscale = src.Grayscale
	cm.config.NoCompression = src.NoCompression
	cm.config.KeepContext = src.KeepContext
	cm.config.InterruptThinking = src.InterruptThinking
	cm.config.UseMarkdownResume = src.UseMarkdownResume

	if len(src.Shortcuts) > 0 {
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

func (cm *ConfigManager) Update(patch ConfigPatch) error {
	cm.mu.Lock()

	cm.applyPatch(patch)

	configCopy := cm.config
	subscribers := cm.subscribers

	cm.mu.Unlock()

	for _, sub := range subscribers {
		sub(configCopy)
	}

	return cm.Save()
}

func (cm *ConfigManager) UpdateFromJSON(jsonStr string) error {
	var patch ConfigPatch
	if err := json.Unmarshal([]byte(jsonStr), &patch); err != nil {
		return fmt.Errorf("解析配置 JSON 失败: %w", err)
	}

	return cm.Update(patch)
}

func (cm *ConfigManager) applyPatch(patch ConfigPatch) {
	if patch.APIKey != nil {
		cm.config.APIKey = *patch.APIKey
	}
	if patch.BaseURL != nil {
		cm.config.BaseURL = *patch.BaseURL
	}
	if patch.Model != nil {
		cm.config.Model = *patch.Model
	}
	if patch.Prompt != nil {
		cm.config.Prompt = *patch.Prompt
	}
	if patch.Opacity != nil {
		cm.config.Opacity = *patch.Opacity
	}
	if patch.ScreenshotMode != nil {
		cm.config.ScreenshotMode = *patch.ScreenshotMode
	}
	if patch.CompressionQuality != nil {
		cm.config.CompressionQuality = *patch.CompressionQuality
	}
	if patch.Sharpening != nil {
		cm.config.Sharpening = *patch.Sharpening
	}
	if patch.Grayscale != nil {
		cm.config.Grayscale = *patch.Grayscale
	}
	if patch.NoCompression != nil {
		cm.config.NoCompression = *patch.NoCompression
	}
	if patch.KeepContext != nil {
		cm.config.KeepContext = *patch.KeepContext
	}
	if patch.InterruptThinking != nil {
		cm.config.InterruptThinking = *patch.InterruptThinking
	}
	if patch.ResumePath != nil {
		cm.config.ResumePath = *patch.ResumePath
	}
	if patch.ResumeContent != nil {
		cm.config.ResumeContent = *patch.ResumeContent
	}
	if patch.UseMarkdownResume != nil {
		cm.config.UseMarkdownResume = *patch.UseMarkdownResume
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
	cm.config.ResumePath = ""
	cm.config.ResumeBase64 = ""
	cm.config.ResumeContent = ""
}

func (cm *ConfigManager) SetResumePath(path string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.ResumePath = path
}

func (cm *ConfigManager) SetResumeBase64(base64 string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.ResumeBase64 = base64
}

func (cm *ConfigManager) SetResumeContent(content string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.ResumeContent = content
}
