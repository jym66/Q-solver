//go:build windows

package platform

import (
	"Q-Solver/pkg/logger"
	"os"
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	ntdll    = syscall.NewLazyDLL("ntdll.dll")

	procFindWindowW                    = user32.NewProc("FindWindowW")
	procSetWindowDisplayAffinity       = user32.NewProc("SetWindowDisplayAffinity")
	procGetWindowLongW                 = user32.NewProc("GetWindowLongW")
	procSetWindowLongW                 = user32.NewProc("SetWindowLongW")
	procRegisterHotKey                 = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey               = user32.NewProc("UnregisterHotKey")
	procGetMessageW                    = user32.NewProc("GetMessageW")
	procSetWindowsHookExW              = user32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx            = user32.NewProc("UnhookWindowsHookEx")
	procCallNextHookEx                 = user32.NewProc("CallNextHookEx")
	procGetModuleHandleW               = kernel32.NewProc("GetModuleHandleW")
	procSetWindowPos                   = user32.NewProc("SetWindowPos")
	procSetLayeredWindowAttributes     = user32.NewProc("SetLayeredWindowAttributes")
	procGetAsyncKeyState               = user32.NewProc("GetAsyncKeyState")
	procEnumWindows                    = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId       = user32.NewProc("GetWindowThreadProcessId")
	procKeybdEvent                     = user32.NewProc("keybd_event")
	procSetWindowCompositionAttribute  = user32.NewProc("SetWindowCompositionAttribute")
	procRtlGetVersion                  = ntdll.NewProc("RtlGetVersion")
)

// WindowHandle 窗口句柄类型（Windows 为 HWND）
type WindowHandle uintptr

// 窗口级别常量
const (
	WindowLevelNormal   = 0 // 正常窗口级别（Windows 不使用）
	WindowLevelFloating = 3 // 置顶窗口级别（Windows 不使用）
)

// Windows 窗口和样式常量
const (
	WDA_NONE               = 0x00000000
	WDA_MONITOR            = 0x00000001
	WDA_EXCLUDEFROMCAPTURE = 0x00000011

	GWL_STYLE   = -16
	GWL_EXSTYLE = -20

	WS_CAPTION     = 0x00C00000
	WS_THICKFRAME  = 0x00040000
	WS_MINIMIZEBOX = 0x00020000
	WS_MAXIMIZEBOX = 0x00010000
	WS_SYSMENU     = 0x00080000

	WS_EX_TOOLWINDOW  = 0x00000080
	WS_EX_APPWINDOW   = 0x00040000
	WS_EX_TRANSPARENT = 0x00000020
	WS_EX_LAYERED     = 0x00080000

	SWP_NOSIZE       = 0x0001
	SWP_NOMOVE       = 0x0002
	SWP_NOZORDER     = 0x0004
	SWP_FRAMECHANGED = 0x0020

	MOD_ALT     = 0x0001
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	MOD_WIN     = 0x0008
	VK_OEM_3    = 0xC0 // `~` 键
	VK_H        = 0x48
	VK_C        = 0x43
	VK_T        = 0x54
	VK_SHIFT    = 0x10
	VK_CONTROL  = 0x11
	VK_MENU     = 0x12 // Alt
	WM_HOTKEY   = 0x0312

	WH_KEYBOARD_LL = 13
	WH_MOUSE_LL    = 14 // 添加鼠标钩子ID

	WM_KEYDOWN    = 0x0100
	WM_SYSKEYDOWN = 0x0104
	WM_KEYUP      = 0x0101
	WM_SYSKEYUP   = 0x0105

	// 添加鼠标消息常量
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
	WM_XBUTTONDOWN = 0x020B
	WM_XBUTTONUP   = 0x020C

	LWA_COLORKEY = 0x00000001
	LWA_ALPHA    = 0x00000002

	// 添加鼠标虚拟键码
	VK_LBUTTON  = 0x01
	VK_RBUTTON  = 0x02
	VK_MBUTTON  = 0x04
	VK_XBUTTON1 = 0x05
	VK_XBUTTON2 = 0x06

	VK_LEFT  = 0x25
	VK_UP    = 0x26
	VK_RIGHT = 0x27
	VK_DOWN  = 0x28
	VK_PRIOR = 0x21 // PageUp
	VK_NEXT  = 0x22 // PageDown

	// 点击窗口不抢焦点
	WS_EX_NOACTIVATE = 0x08000000
	// 保证悬浮在浏览器上面
	WS_EX_TOPMOST = 0x00000008
	//不激活窗口(刷新样式时用)
	SWP_NOACTIVATE = 0x0010
)

// KBDLLHOOKSTRUCT 键盘钩子结构体
type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

// MSLLHOOKSTRUCT 鼠标钩子结构体
type MSLLHOOKSTRUCT struct {
	Pt          POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

// RTL_OSVERSIONINFOW 系统版本信息结构体
type RTL_OSVERSIONINFOW struct {
	dwOSVersionInfoSize uint32
	dwMajorVersion      uint32
	dwMinorVersion      uint32
	dwBuildNumber       uint32
	dwPlatformId        uint32
	szCSDVersion        [128]uint16
}

// ACCENT_POLICY 窗口合成属性策略
type ACCENT_POLICY struct {
	AccentState   uint32
	AccentFlags   uint32
	GradientColor uint32
	AnimationId   uint32
}

// WINDOWCOMPOSITIONATTRIBDATA 窗口合成属性数据
type WINDOWCOMPOSITIONATTRIBDATA struct {
	Attrib     uint32
	Data       unsafe.Pointer
	SizeOfData uintptr
}

// 窗口合成属性常量
const (
	WCA_ACCENT_POLICY                 = 19
	ACCENT_ENABLE_TRANSPARENTGRADIENT = 2
)

// MSG Windows 消息结构体
type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// POINT 点坐标结构体
type POINT struct {
	X, Y int32
}

// ==================== 平台接口函数 ====================

// GetWindowHandle 获取当前进程主窗口句柄
func GetWindowHandle() (WindowHandle, error) {
	hwnd, err := getHwndByPid(uint32(os.Getpid()))
	return WindowHandle(hwnd), err
}

// ApplyGhostMode 应用幽灵模式（无边框、置顶、防录屏、不抢焦点）
func ApplyGhostMode(hwnd WindowHandle) error {
	applyGhostMode(uintptr(hwnd))
	return nil
}

// SetClickThrough 设置鼠标穿透
func SetClickThrough(hwnd WindowHandle, enabled bool) error {
	return setWindowClickThrough(uintptr(hwnd), enabled)
}

// SetDisplayAffinity 设置防录屏状态
func SetDisplayAffinity(hwnd WindowHandle, hidden bool) error {
	affinity := WDA_NONE
	if hidden {
		affinity = WDA_EXCLUDEFROMCAPTURE
	}
	return setWindowDisplayAffinity(uintptr(hwnd), uint32(affinity))
}

// RestoreFocus 恢复焦点
func RestoreFocus(hwnd WindowHandle) error {
	restoreFocus(uintptr(hwnd))
	return nil
}

// RemoveFocus 移除焦点
func RemoveFocus(hwnd WindowHandle) error {
	removeFocus(uintptr(hwnd))
	return nil
}

// CheckScreenCaptureAccess 检查截图权限 (Windows 直接返回 true)
func CheckScreenCaptureAccess() bool {
	return true // Windows 不需要截图权限
}

// RequestScreenCaptureAccess 请求截图权限 (Windows 直接返回 true)
func RequestScreenCaptureAccess() bool {
	return true // Windows 不需要截图权限
}

// OpenScreenCaptureSettings 打开系统设置的屏幕录制权限页面 (Windows 无操作)
func OpenScreenCaptureSettings() {
	// Windows 不需要截图权限，无操作
}

// SetWindowLevel 设置窗口层级 (Windows 在 ApplyGhostMode 中处理)
func SetWindowLevel(hwnd WindowHandle, level int) error {
	// Windows 使用不同的置顶机制，在 ApplyGhostMode 中处理
	return nil
}

// CheckMicrophoneAccess 检查麦克风权限状态 (Windows 直接返回已授权)
func CheckMicrophoneAccess() int {
	return 1 // Windows 不需要预先请求麦克风权限
}

// RequestMicrophoneAccess 请求麦克风权限 (Windows 无操作)
func RequestMicrophoneAccess() {
	// Windows 不需要预先请求麦克风权限
}

// OpenMicrophoneSettings 打开系统设置的麦克风权限页面 (Windows 无操作)
func OpenMicrophoneSettings() {
	// Windows 不需要预先请求麦克风权限
}

// ==================== Windows API 封装函数 ====================

// RegisterHotKey 注册热键
func RegisterHotKey(hwnd uintptr, id int, fsModifiers uint32, vk uint32) bool {
	ret, _, _ := procRegisterHotKey.Call(
		hwnd,
		uintptr(id),
		uintptr(fsModifiers),
		uintptr(vk),
	)
	return ret != 0
}

// UnregisterHotKey 注销热键
func UnregisterHotKey(hwnd uintptr, id int) bool {
	ret, _, _ := procUnregisterHotKey.Call(
		hwnd,
		uintptr(id),
	)
	return ret != 0
}

// GetMessage 获取消息
func GetMessage(msg *MSG, hwnd uintptr, msgFilterMin, msgFilterMax uint32) int32 {
	ret, _, _ := procGetMessageW.Call(
		uintptr(unsafe.Pointer(msg)),
		hwnd,
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	return int32(ret)
}

// SetWindowsHookEx 安装钩子
func SetWindowsHookEx(idHook int, lpfn uintptr, hMod uintptr, dwThreadId uint32) uintptr {
	ret, _, _ := procSetWindowsHookExW.Call(
		uintptr(idHook),
		lpfn,
		hMod,
		uintptr(dwThreadId),
	)
	return ret
}

// UnhookWindowsHookEx 卸载钩子
func UnhookWindowsHookEx(hhk uintptr) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(hhk)
	return ret != 0
}

// CallNextHookEx 调用下一个钩子
func CallNextHookEx(hhk uintptr, nCode int, wParam uintptr, lParam uintptr) uintptr {
	ret, _, _ := procCallNextHookEx.Call(
		hhk,
		uintptr(nCode),
		wParam,
		lParam,
	)
	return ret
}

// GetModuleHandle 获取模块句柄
func GetModuleHandle(moduleName string) uintptr {
	var m *uint16
	if moduleName != "" {
		m, _ = syscall.UTF16PtrFromString(moduleName)
	}
	ret, _, _ := procGetModuleHandleW.Call(uintptr(unsafe.Pointer(m)))
	return ret
}

// FindWindow 查找窗口
func FindWindow(className, windowName string) uintptr {
	var c *uint16
	var w *uint16
	if className != "" {
		c, _ = syscall.UTF16PtrFromString(className)
	}
	if windowName != "" {
		w, _ = syscall.UTF16PtrFromString(windowName)
	}
	ret, _, _ := procFindWindowW.Call(
		uintptr(unsafe.Pointer(c)),
		uintptr(unsafe.Pointer(w)),
	)
	return ret
}

// setWindowDisplayAffinity 设置窗口防录屏
func setWindowDisplayAffinity(hwnd uintptr, affinity uint32) error {
	ret, _, err := procSetWindowDisplayAffinity.Call(
		hwnd,
		uintptr(affinity),
	)
	if ret == 0 {
		return err
	}
	return nil
}

// SetWindowToolWindow 设置为工具窗口（隐藏任务栏）
func SetWindowToolWindow(hwnd uintptr) error {
	nIndex := int32(GWL_EXSTYLE)
	// 获取当前扩展样式
	style, _, _ := procGetWindowLongW.Call(hwnd, uintptr(nIndex))

	// 添加 WS_EX_TOOLWINDOW 并移除 WS_EX_APPWINDOW
	newStyle := (style | WS_EX_TOOLWINDOW) &^ WS_EX_APPWINDOW

	ret, _, err := procSetWindowLongW.Call(hwnd, uintptr(nIndex), newStyle)
	if ret == 0 {
		return err
	}
	return nil
}

// setWindowClickThrough 设置鼠标穿透
func setWindowClickThrough(hwnd uintptr, enabled bool) error {
	nIndex := int32(GWL_EXSTYLE)
	style, _, _ := procGetWindowLongW.Call(hwnd, uintptr(nIndex))
	var newStyle uintptr
	if enabled {
		newStyle = style | WS_EX_TRANSPARENT | WS_EX_LAYERED
	} else {
		newStyle = style &^ WS_EX_TRANSPARENT
	}
	ret, _, err := procSetWindowLongW.Call(hwnd, uintptr(nIndex), newStyle)
	if ret == 0 {
		return err
	}
	return nil
}

// SetWindowPos 设置窗口位置
func SetWindowPos(hwnd uintptr, hWndInsertAfter uintptr, x, y, cx, cy int, uFlags uint32) bool {
	ret, _, _ := procSetWindowPos.Call(
		hwnd,
		hWndInsertAfter,
		uintptr(x),
		uintptr(y),
		uintptr(cx),
		uintptr(cy),
		uintptr(uFlags),
	)
	return ret != 0
}

// GetWindowLong 获取窗口属性
func GetWindowLong(hwnd uintptr, index int) int32 {
	ret, _, _ := procGetWindowLongW.Call(
		hwnd,
		uintptr(index),
	)
	return int32(ret)
}

// SetWindowLong 设置窗口属性
func SetWindowLong(hwnd uintptr, index int, value int32) int32 {
	ret, _, _ := procSetWindowLongW.Call(
		hwnd,
		uintptr(index),
		uintptr(value),
	)
	return int32(ret)
}

// SetLayeredWindowAttributes 设置分层窗口属性
func SetLayeredWindowAttributes(hwnd uintptr, crKey uint32, bAlpha byte, dwFlags uint32) error {
	ret, _, err := procSetLayeredWindowAttributes.Call(
		hwnd,
		uintptr(crKey),
		uintptr(bAlpha),
		uintptr(dwFlags),
	)
	if ret == 0 {
		return err
	}
	return nil
}

// GetAsyncKeyState 获取异步按键状态
func GetAsyncKeyState(vKey int) uint16 {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vKey))
	return uint16(ret)
}

// getHwndByPid 根据进程ID获取窗口句柄
func getHwndByPid(pid uint32) (uintptr, error) {
	var hwnd uintptr
	cb := syscall.NewCallback(func(h uintptr, p uintptr) uintptr {
		var pId uint32
		procGetWindowThreadProcessId.Call(h, uintptr(unsafe.Pointer(&pId)))
		if pId == pid {
			hwnd = h
			return 0 // Stop enumeration
		}
		return 1 // Continue enumeration
	})
	procEnumWindows.Call(cb, 0)
	if hwnd == 0 {
		return 0, syscall.Errno(0)
	}
	return hwnd, nil
}

// isWindows11OrLater 检测是否为 Windows 11 及以上版本
func isWindows11OrLater() bool {
	var info RTL_OSVERSIONINFOW
	info.dwOSVersionInfoSize = uint32(unsafe.Sizeof(info))
	ret, _, _ := procRtlGetVersion.Call(uintptr(unsafe.Pointer(&info)))
	if ret != 0 {
		return false
	}
	return info.dwBuildNumber >= 22000
}

// enableTransparentGradient 在 Win10 上使用 SetWindowCompositionAttribute
// 覆盖 DWM 的模糊效果，实现无模糊的真透明背景
func enableTransparentGradient(hwnd uintptr) {
	accent := ACCENT_POLICY{
		AccentState:   ACCENT_ENABLE_TRANSPARENTGRADIENT,
		AccentFlags:   2,
		GradientColor: 0x00000000,
		AnimationId:   0,
	}
	data := WINDOWCOMPOSITIONATTRIBDATA{
		Attrib:     WCA_ACCENT_POLICY,
		Data:       unsafe.Pointer(&accent),
		SizeOfData: unsafe.Sizeof(accent),
	}
	procSetWindowCompositionAttribute.Call(hwnd, uintptr(unsafe.Pointer(&data)))
}

// applyGhostMode 应用幽灵模式
func applyGhostMode(hwnd uintptr) {
	style := GetWindowLong(hwnd, int(GWL_STYLE))
	// 移除 标题栏 | 调大小边框 | 最小化 | 最大化 | 系统菜单
	newStyle := style &^ (WS_CAPTION | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX | WS_SYSMENU)
	SetWindowLong(hwnd, int(GWL_STYLE), int32(newStyle))

	exStyle := GetWindowLong(hwnd, int(GWL_EXSTYLE))

	//TOOLWINDOW (隐藏任务栏) | NOACTIVATE (不抢焦点) | TOPMOST (置顶)
	newExStyle := (exStyle &^ WS_EX_APPWINDOW) | WS_EX_TOOLWINDOW | WS_EX_NOACTIVATE | WS_EX_TOPMOST
	SetWindowLong(hwnd, int(GWL_EXSTYLE), int32(newExStyle))

	// 优先尝试"隐身(透明)"，如果失败降级为"黑屏"
	err := setWindowDisplayAffinity(hwnd, WDA_EXCLUDEFROMCAPTURE)
	if err != nil {
		// 如果系统不支持透明(如Win7/旧Win10)，降级为黑屏模式
		setWindowDisplayAffinity(hwnd, WDA_MONITOR)
		logger.Println("[GhostMode] 降级为黑屏防御模式")
	} else {
		logger.Println("[GhostMode] 隐身防御模式已开启")
	}

	// 必须带上 SWP_FRAMECHANGED 让去边框生效
	// 必须带上 SWP_NOACTIVATE 防止刷新时抢走焦点
	SetWindowPos(hwnd, 0, 0, 0, 0, 0,
		uint32(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED|SWP_NOACTIVATE))

	// Win10 上 WindowIsTranslucent 会触发 DWM 模糊效果，
	// 用 SetWindowCompositionAttribute 覆盖为透明渐变，消除模糊
	if !isWindows11OrLater() {
		enableTransparentGradient(hwnd)
		logger.Println("[GhostMode] Win10 透明渐变模式已启用（消除 DWM 模糊）")
	}
}

// restoreFocus 恢复焦点
func restoreFocus(hwnd uintptr) {
	exStyle := GetWindowLong(hwnd, int(GWL_EXSTYLE))
	// 移除 WS_EX_NOACTIVATE 以允许获取焦点
	newExStyle := exStyle &^ WS_EX_NOACTIVATE
	SetWindowLong(hwnd, int(GWL_EXSTYLE), int32(newExStyle))

	// 刷新窗口样式，不带 SWP_NOACTIVATE 以允许激活
	SetWindowPos(hwnd, 0, 0, 0, 0, 0,
		uint32(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED))
}

// removeFocus 移除焦点
func removeFocus(hwnd uintptr) {
	exStyle := GetWindowLong(hwnd, int(GWL_EXSTYLE))
	// 添加 WS_EX_NOACTIVATE 以防止抢焦点
	newExStyle := exStyle | WS_EX_NOACTIVATE
	SetWindowLong(hwnd, int(GWL_EXSTYLE), int32(newExStyle))

	// 刷新窗口样式
	SetWindowPos(hwnd, 0, 0, 0, 0, 0,
		uint32(SWP_NOMOVE|SWP_NOSIZE|SWP_NOZORDER|SWP_FRAMECHANGED|SWP_NOACTIVATE))
}

// KeybdEvent 模拟键盘事件
func KeybdEvent(bVk byte, bScan byte, dwFlags uint32, dwExtraInfo uintptr) {
	procKeybdEvent.Call(
		uintptr(bVk),
		uintptr(bScan),
		uintptr(dwFlags),
		dwExtraInfo,
	)
}
