package llm

import (
	"Q-Solver/pkg/config"
	"context"
	"fmt"
	"time"
)

type Service struct {
	config   *config.Config
	provider Provider
}

func NewService(cfg *config.Config) *Service {
	s := &Service{
		config: cfg,
	}
	// 初始化时创建 Provider
	s.UpdateProvider()
	return s
}

func (s *Service) UpdateProvider() {
	s.provider = NewOpenAIProvider(s.config.GetAPIKey(), s.config.GetBaseURL(), s.config.GetModel())
}

func (s *Service) GetProvider() Provider {
	return s.provider
}

// TestConnection 测试模型连通性
// 通过发送一个简单的聊天请求来验证 API Key 和模型是否可用
func (s *Service) TestConnection(ctx context.Context, apiKey, baseURL, model string) string {
	if apiKey == "" {
		return "API Key 不能为空"
	}
	if model == "" {
		return "请选择模型"
	}

	// 使用传入的参数创建临时 provider
	if baseURL == "" {
		baseURL = s.config.GetBaseURL()
	}
	tempProvider := NewOpenAIProvider(apiKey, baseURL, model)

	// 设置超时
	timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 发送一个简单的测试消息
	err := tempProvider.TestChat(timeoutCtx)
	if err != nil {
		return err.Error()
	}

	return "" // 成功返回空字符串
}

// GetModels 获取模型列表
func (s *Service) GetModels(ctx context.Context, apiKey string, baseURL string) ([]string, error) {
	var provider *OpenAIProvider

	// 如果提供了临时的 apiKey 或 baseURL，使用临时 provider
	if (apiKey != "" && apiKey != s.config.GetAPIKey()) || (baseURL != "" && baseURL != s.config.GetBaseURL()) {
		if baseURL == "" {
			baseURL = s.config.GetBaseURL()
		}
		if apiKey == "" {
			apiKey = s.config.GetAPIKey()
		}
		provider = NewOpenAIProvider(apiKey, baseURL, s.config.GetModel())
	} else {
		// 尝试将接口转换为具体类型
		var ok bool
		provider, ok = s.provider.(*OpenAIProvider)
		if !ok {
			return nil, fmt.Errorf("current provider does not support listing models")
		}
	}

	return provider.GetModels(ctx)
}
