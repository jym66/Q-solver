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

// copyNonNilFields 将 src 中非 nil 的字段复制到 dst
// checkEnv 为 true 时，APIKey 和 BaseURL 仅在环境变量未设置时才复制
func copyNonNilFields(src, dst *Config, checkEnv bool) {
	if src.APIKey != nil {
		if !checkEnv || os.Getenv("GHOST_API_KEY") == "" {
			dst.APIKey = src.APIKey
		}
	}
	if src.BaseURL != nil {
		if !checkEnv || os.Getenv("GHOST_BASE_URL") == "" {
			dst.BaseURL = src.BaseURL
		}
	}
	if src.Provider != nil {
		dst.Provider = src.Provider
	}
	if src.Model != nil {
		dst.Model = src.Model
	}
	if src.Prompt != nil {
		dst.Prompt = src.Prompt
	}
	if src.Opacity != nil {
		dst.Opacity = src.Opacity
	}
	if src.ScreenshotMode != nil {
		dst.ScreenshotMode = src.ScreenshotMode
	}
	if src.CompressionQuality != nil {
		dst.CompressionQuality = src.CompressionQuality
	}
	if src.Sharpening != nil {
		dst.Sharpening = src.Sharpening
	}
	if src.Grayscale != nil {
		dst.Grayscale = src.Grayscale
	}
	if src.NoCompression != nil {
		dst.NoCompression = src.NoCompression
	}
	if src.KeepContext != nil {
		dst.KeepContext = src.KeepContext
	}
	if src.InterruptThinking != nil {
		dst.InterruptThinking = src.InterruptThinking
	}
	if src.ResumePath != nil {
		dst.ResumePath = src.ResumePath
	}
	if src.ResumeContent != nil {
		dst.ResumeContent = src.ResumeContent
	}
	if src.UseMarkdownResume != nil {
		dst.UseMarkdownResume = src.UseMarkdownResume
	}
	if len(src.Shortcuts) > 0 {
		if dst.Shortcuts == nil {
			dst.Shortcuts = make(map[string]shortcut.KeyBinding)
		}
		maps.Copy(dst.Shortcuts, src.Shortcuts)
	}
}

func (cm *ConfigManager) mergeConfig(src *Config) {
	copyNonNilFields(src, &cm.config, true)
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
	copyNonNilFields(&patch, &cm.config, false)
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
