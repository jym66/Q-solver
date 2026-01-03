package main

import (
	"Q-Solver/pkg/config"
	"Q-Solver/pkg/llm"
	"Q-Solver/pkg/logger"
	"Q-Solver/pkg/resume"
	"Q-Solver/pkg/screen"
	"Q-Solver/pkg/shortcut"
	"Q-Solver/pkg/solution"
	"Q-Solver/pkg/state"
	"Q-Solver/pkg/task"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 这个只作为连接前端和后端的中间层
type App struct {
	ctx context.Context

	// 管理器
	configManager *config.ConfigManager
	stateManager  *state.StateManager
	taskManager   *task.TaskCoordinator

	// 业务服务
	llmService      *llm.Service
	resumeService   *resume.Service
	shortcutService *shortcut.Service
	screenService   *screen.Service
	solver          *solution.Solver
}

// NewApp 创建 App 实例
func NewApp() *App {
	configManager := config.NewConfigManager()

	app := &App{
		configManager: configManager,
		stateManager:  state.NewStateManager(),
		taskManager:   task.NewTaskCoordinator(),
		screenService: screen.NewService(),
	}

	return app
}

// Startup Wails 启动回调
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// 加载配置
	if err := a.configManager.Load(); err != nil {
		logger.Printf("加载配置失败: %v", err)
	}

	// 初始化状态管理器
	a.stateManager.Startup(ctx, a.EmitEvent)

	// 初始化屏幕服务
	a.screenService.Startup(ctx)

	// 初始化 LLM 服务
	a.llmService = llm.NewService(a.configManager.GetPtr())
	a.solver = solution.NewSolver(a.llmService.GetProvider())

	// 初始化简历服务
	a.resumeService = resume.NewService(a.configManager.GetPtr())

	// 初始化快捷键服务
	a.shortcutService = shortcut.NewService(a)
	a.shortcutService.SetShortcuts(a.configManager.Get().Shortcuts)
	a.shortcutService.Start()
	logger.Println("快捷键服务已初始化")

	// 订阅配置变更
	a.configManager.Subscribe(a.onConfigChanged)
	logger.Println("配置变更订阅已注册")

	// 直接设置为就绪状态
	a.stateManager.UpdateInitStatus(state.StatusReady)
}

// onConfigChanged 配置变更回调
func (a *App) onConfigChanged(cfg config.Config) {
	// 更新 LLM Provider
	a.llmService.UpdateProvider()
	if a.solver != nil {
		a.solver.SetProvider(a.llmService.GetProvider())
	}

	// 更新快捷键
	a.shortcutService.SetShortcuts(cfg.Shortcuts)

	// 如果关闭了上下文，清空历史
	if !cfg.KeepContext && a.solver != nil {
		a.solver.ClearHistory()
	}

	logger.Println("配置已更新并应用")
}

// OnShutdown Wails 关闭回调
func (a *App) OnShutdown(ctx context.Context) {
	if a.shortcutService != nil {
		a.shortcutService.Stop()
	}
	// 保存配置
	if err := a.configManager.Save(); err != nil {
		logger.Printf("保存配置失败: %v", err)
	}
}

// ==================== 事件与状态 ====================

// EmitEvent 发送事件到前端
func (a *App) EmitEvent(eventName string, data ...interface{}) {
	runtime.EventsEmit(a.ctx, eventName, data...)
}

// GetInitStatus 获取初始化状态
func (a *App) GetInitStatus() string {
	return a.stateManager.GetInitStatusString()
}

// ==================== 配置管理 ====================

// GetSettings 返回当前配置
func (a *App) GetSettings() config.Config {
	return a.configManager.Get()
}

// UpdateSettings 更新配置（从前端 JSON）
func (a *App) UpdateSettings(configJson string) string {
	if err := a.configManager.UpdateFromJSON(configJson); err != nil {
		return err.Error()
	}
	return ""
}

// SyncSettingsToDefaultSettings 兼容旧接口
// Deprecated: 使用 UpdateSettings 替代
func (a *App) SyncSettingsToDefaultSettings(configJson string) string {
	return a.UpdateSettings(configJson)
}

// ==================== 窗口控制 ====================

// ToggleVisibility 切换可见性
func (a *App) ToggleVisibility() {
	a.stateManager.ToggleVisibility()
}

// ToggleClickThrough 切换鼠标穿透
func (a *App) ToggleClickThrough() {
	a.stateManager.ToggleClickThrough()
}

// MoveWindow 移动窗口
func (a *App) MoveWindow(dx, dy int) {
	a.stateManager.MoveWindow(dx, dy)
}

// RestoreFocus 恢复焦点
func (a *App) RestoreFocus() {
	a.stateManager.RestoreFocus()
}

// RemoveFocus 移除焦点
func (a *App) RemoveFocus() {
	a.stateManager.RemoveFocus()
}

// ==================== 解题相关 ====================

// TriggerSolve 触发解题（快捷键调用）
func (a *App) TriggerSolve() {
	cfg := a.configManager.Get()
	if cfg.APIKey == "" {
		a.EmitEvent("require-login")
		return
	}

	a.EmitEvent("start-solving")

	// 使用任务协调器管理任务
	ctx, taskID := a.taskManager.StartTask("solve")

	go func() {
		success := a.solveInternal(ctx)

		if success {
			a.taskManager.CompleteTask(taskID)
		}
	}()
}

// solveInternal 内部解题逻辑
func (a *App) solveInternal(ctx context.Context) bool {
	cfg := a.configManager.Get()

	if cfg.APIKey == "" {
		a.EmitEvent("require-login")
		return false
	}

	// 读取简历
	resumeBase64, err := a.resumeService.GetResumeBase64()
	if err != nil {
		logger.Printf("读取简历失败: %v\n", err)
	}

	// 获取截图
	previewResult, err := a.GetScreenshotPreview(
		cfg.CompressionQuality,
		cfg.Sharpening,
		cfg.Grayscale,
		cfg.NoCompression,
		cfg.ScreenshotMode,
	)
	if err != nil {
		logger.Printf("图片编码失败: %v\n", err)
		return false
	}

	req := solution.Request{
		Config:           cfg,
		ScreenshotBase64: previewResult.Base64,
		ResumeBase64:     resumeBase64,
	}

	cb := solution.Callbacks{
		EmitEvent: a.EmitEvent,
	}

	return a.solver.Solve(ctx, req, cb)
}

// CancelRunningTask 取消当前运行的任务
func (a *App) CancelRunningTask() bool {
	return a.taskManager.CancelCurrentTask()
}

// IsInterruptThinkingEnabled 是否允许打断思考
func (a *App) IsInterruptThinkingEnabled() bool {
	return a.configManager.Get().InterruptThinking
}

// ==================== 快捷键相关 ====================

// StartRecordingKey 开始录制快捷键
func (a *App) StartRecordingKey(action string) {
	a.shortcutService.StartRecording(action)
}

// StopRecordingKey 停止录制快捷键
func (a *App) StopRecordingKey() {
	if a.shortcutService != nil {
		a.shortcutService.StopRecording()
	}
}

// ScrollContent 滚动内容
func (a *App) ScrollContent(direction string) {
	a.EmitEvent("scroll-content", direction)
}

// CopyCode 复制代码
func (a *App) CopyCode() {
	a.EmitEvent("copy-code")
}

// ==================== 简历相关 ====================

// SelectResume 选择简历文件
func (a *App) SelectResume() string {
	path := a.resumeService.SelectResume(a.ctx)
	if path != "" {
		a.configManager.SetResumePath(path)
	}
	return path
}

// ClearResume 清除简历
func (a *App) ClearResume() {
	a.resumeService.ClearResume()
	a.configManager.ClearResume()
}

// GetResumePDF 获取简历 Base64
func (a *App) GetResumePDF() (string, error) {
	return a.resumeService.GetResumeBase64()
}

// ParseResume 解析简历为 Markdown
func (a *App) ParseResume() (string, error) {
	return a.resumeService.ParseResume(a.ctx, a.llmService.GetProvider())
}

// ==================== 截图相关 ====================

// GetScreenshotPreview 获取截图预览
func (a *App) GetScreenshotPreview(quality int, sharpen float64, grayscale bool, noCompression bool, screenshotMode string) (screen.PreviewResult, error) {
	mode := screenshotMode
	if mode == "" {
		mode = a.configManager.Get().ScreenshotMode
	}
	return a.screenService.CapturePreview(quality, sharpen, grayscale, noCompression, mode)
}

// ==================== LLM 相关 ====================

// TestConnection 测试模型连通性
// 通过发送一个简单的消息来测试 API 是否可用
func (a *App) TestConnection(apiKey, baseURL, model string) string {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return a.llmService.TestConnection(ctx, apiKey, baseURL, model)
}

// GetModels 获取模型列表
func (a *App) GetModels(apiKey string) ([]string, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return a.llmService.GetModels(ctx, apiKey)
}
