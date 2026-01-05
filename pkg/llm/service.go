package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"Q-Solver/pkg/config"
	"Q-Solver/pkg/logger"
)

// ProviderType 提供商类型
type ProviderType string

const (
	ProviderOpenAI ProviderType = "openai"
	ProviderGemini ProviderType = "gemini"
	ProviderClaude ProviderType = "claude"
)

// Service LLM 服务
type Service struct {
	config   *config.Config
	provider Provider
}

// NewService 创建 LLM 服务
func NewService(cfg *config.Config) *Service {
	s := &Service{
		config: cfg,
	}
	s.UpdateProvider()
	return s
}

// UpdateProvider 更新 Provider（配置变更时调用）
func (s *Service) UpdateProvider() {
	providerType := DetectProviderType(s.config.GetBaseURL(), s.config.GetModel())
	s.provider = CreateProvider(providerType, s.config)
}

// GetProvider 获取当前 Provider
func (s *Service) GetProvider() Provider {
	return s.provider
}

// DetectProviderType 根据 baseURL 或 model 名称自动识别提供商
func DetectProviderType(baseURL, model string) ProviderType {
	// 1. 优先根据 baseURL 判断
	switch {
	case strings.Contains(baseURL, "generativelanguage.googleapis.com"):
		return ProviderGemini
	case strings.Contains(baseURL, "anthropic.com"):
		return ProviderClaude
	}

	// 2. 根据模型名称判断
	switch {
	case strings.HasPrefix(model, "gemini"):
		return ProviderGemini
	case strings.HasPrefix(model, "claude"):
		return ProviderClaude
	}

	// 3. 默认使用 OpenAI 兼容模式
	return ProviderOpenAI
}

// CreateProvider 工厂函数：根据类型创建对应 Provider
func CreateProvider(providerType ProviderType, cfg *config.Config) Provider {
	switch providerType {
	case ProviderGemini:
		adapter, _ := NewGeminiAdapter(cfg)
		logger.Println("创建GeminiAdapter")
		return adapter
	case ProviderClaude:
		logger.Println("创建ClaudeAdapter")
		return NewClaudeAdapter(cfg)
	default:
		logger.Println("创建OpenAIAdapter")
		return NewOpenAIAdapter(cfg)
	}
}

// TestConnection 测试模型连通性
func (s *Service) TestConnection(ctx context.Context, apiKey, baseURL, model string) string {
	if apiKey == "" {
		return "API Key 不能为空"
	}
	if model == "" {
		return "请选择模型"
	}

	if baseURL == "" {
		baseURL = s.config.GetBaseURL()
	}

	// 创建临时 config 用于测试
	tempConfig := *s.config
	tempConfig.APIKey = &apiKey
	tempConfig.BaseURL = &baseURL
	tempConfig.Model = &model

	providerType := DetectProviderType(baseURL, model)
	tempProvider := CreateProvider(providerType, &tempConfig)

	timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	err := tempProvider.TestChat(timeoutCtx)
	if err != nil {
		return err.Error()
	}

	return ""
}

// GetModels 获取模型列表
func (s *Service) GetModels(ctx context.Context, apiKey string, baseURL string) ([]string, error) {
	if baseURL == "" {
		baseURL = s.config.GetBaseURL()
	}
	if apiKey == "" {
		apiKey = s.config.GetAPIKey()
	}

	// 如果提供了临时参数，使用临时 provider
	if apiKey != s.config.GetAPIKey() || baseURL != s.config.GetBaseURL() {
		tempConfig := *s.config
		tempConfig.APIKey = &apiKey
		tempConfig.BaseURL = &baseURL

		providerType := DetectProviderType(baseURL, s.config.GetModel())
		tempProvider := CreateProvider(providerType, &tempConfig)
		return tempProvider.GetModels(ctx)
	}

	// 使用当前 provider
	if s.provider == nil {
		return nil, fmt.Errorf("provider not initialized")
	}
	return s.provider.GetModels(ctx)
}
