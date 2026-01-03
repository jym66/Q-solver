package resume

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/llm"
	"Q-Solver/pkg/logger"
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	config *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

// SelectResume 打开文件对话框选择简历，并返回路径
func (s *Service) SelectResume(ctx context.Context) string {
	selection, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "选择简历 (PDF)",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF Files",
				Pattern:     "*.pdf",
			},
		},
	})

	if err != nil {
		logger.Printf("选择文件失败: %v\n", err)
		return ""
	}

	if selection == "" {
		return "" // 用户取消
	}

	s.config.ResumePath = &selection
	return selection
}

// ClearResume 清除简历
func (s *Service) ClearResume() {
	empty := ""
	s.config.ResumePath = &empty
	s.config.ResumeBase64 = &empty
	logger.Println("简历已清除")
}

// GetResumeBase64 读取简历并转换为 Base64
func (s *Service) GetResumeBase64() (string, error) {
	//读一下缓存
	if len(s.config.GetResumeBase64()) > 0 {
		logger.Println("使用缓存的简历 Base64")
		return s.config.GetResumeBase64(), nil
	}
	if s.config.GetResumePath() == "" {
		return "", nil
	}

	// 检查文件大小
	fileInfo, err := os.Stat(s.config.GetResumePath())
	if err != nil {
		return "", err
	}

	// 限制 5MB
	if fileInfo.Size() > 5*1024*1024 {
		return "", fmt.Errorf("简历文件大小超过 5MB 限制")
	}

	// 读取文件内容
	content, err := os.ReadFile(s.config.GetResumePath())
	if err != nil {
		return "", err
	}

	// 转换为 Base64
	encoded := base64.StdEncoding.EncodeToString(content)
	s.config.ResumeBase64 = &encoded
	return encoded, nil
}

// ParseResume 解析简历为 Markdown
func (s *Service) ParseResume(ctx context.Context, provider llm.Provider) (string, error) {
	// 1. Read Resume
	resumeBase64, err := s.GetResumeBase64()
	if err != nil {
		return "", fmt.Errorf("读取简历失败: %v", err)
	}
	if resumeBase64 == "" {
		return "", fmt.Errorf("未选择简历文件")
	}
	logger.Println("开始解析简历为 Markdown...")
	// 2. Call LLM Provider
	return provider.ParseResume(ctx, resumeBase64)
}
