<div align="center">
  <!-- 建议替换为实际的 Logo 或 Banner -->
  <img src="assets/banner.jpg" alt="Q-Solver Banner" width="100%" style="border-radius: 12px; margin-bottom: 20px; box-shadow: 0 8px 24px rgba(0,0,0,0.2);">

  <h1 style="font-size: 3rem; margin-bottom: 0;">Q-Solver</h1>
  <p style="font-size: 1.2rem; color: #666;">
    👻 <b>你的“隐形”桌面 AI 智囊团</b>
  </p>

  <p>
    <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white" alt="Vue">
    <img src="https://img.shields.io/badge/Wails-v2-E30613?logo=wails&logoColor=white" alt="Wails">
    <img src="https://img.shields.io/badge/AI-Gemini%202.5-8E75B2?logo=google-gemini&logoColor=white" alt="Gemini">
    <img src="https://img.shields.io/badge/Platform-Windows-0078D6?logo=windows&logoColor=white" alt="Windows">
  </p>

  <p>
    <a href="#-演示">演示</a> •
    <a href="#-核心功能">功能</a> •
    <a href="#-应用场景">场景</a> •
    <a href="#-下载与安装">下载</a> •
    <a href="#-配置指南">配置</a>
  </p>
</div>

---

> **💡 为什么开发这个工具？**
>
> 在面试、在线会议、网课甚至某些“特殊测试”场景中，我们常常需要快速获取信息，但频繁切换窗口（Alt+Tab）既危险又低效。
> **Q-Solver** 就是为了解决这个问题而生——它像一层透明的魔法皮肤覆盖在你的屏幕上，**看你所看，听你所听**，并在你需要的时候悄无声息地给出答案。

---

## 📺 演示 (Demo)

<!-- 强烈建议在此处放置一个 GIF 动图，展示：透明模式 -> 截图 -> AI 回答 -> 隐藏 的全过程 -->
<div align="center">
  <p align="center">
    <img src="assets/img1.png" width="49%" style="border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1);" />
    <img src="assets/img2.png" width="49%" style="border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1);" />
  </p>
  <p align="center" style="margin-top: 10px;">
    <img src="assets/img3.png" width="49%" style="border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1);" />
    <img src="assets/img4.png" width="49%" style="border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.1);" />
  </p>

</div>

## ✨ 核心功能：不仅仅是截图工具

### 1. 👁️ 上帝视角 (God Mode Vision)
*   **全屏/区域智能识别**：不再是简单的 OCR。基于 **Gemini 2.5** 多模态大模型，它能理解复杂的图表、代码逻辑、数学公式甚至 UI 布局。
*   **所见即所得**：遇到不懂的代码报错？复杂的物理题？直接截图，AI 瞬间解析。

### 2. 👻 幽灵模式 (Stealth & Camouflage)
*   **透明伪装**：界面可调节至**全透明**，文字像水印一样浮在屏幕上，旁人难以察觉。
*   **鼠标穿透**：开启后，鼠标点击直接穿透浮窗作用于底层应用，**完全不影响你正常操作**。
*   **进程伪装**：在任务管理器中伪装成系统组件或媒体渲染器（如 `GS Video Renderer`），深藏功与名。

### 3. ⚡ 极速交互 (Speed & Flow)
*   **全局快捷键**：一键唤醒、一键截图、一键隐藏（Boss Key）。
*   **极简 UI**：没有复杂的菜单，只有你需要的信息。

## 🎯 应用场景

| 场景 | Q-Solver 能做什么？ |
| :--- | :--- |
| **💻 笔试考试** | 截图题目，AI 实时分析并给出答案，助你从容应对。 |
| **🎓 网课/讲座** | 自动记录老师讲的重点，遇到不懂的 PPT 直接截图询问，生成复习笔记。 |
| **🐞 编程开发** | 遇到报错弹窗无法复制？直接截图，AI 帮你分析堆栈信息并给出修复建议。 |
| **📝 复杂办公** | 财务报表、数据分析图表看不懂？让 AI 帮你解读数据趋势。 |

## 🚀 下载与安装

### 方式一：直接下载 (推荐)
前往 [Releases](https://github.com/jym66/Q-Solver/releases) 页面下载最新版本的 `Q-Solver-Setup.exe` 或绿色版压缩包。

### 方式二：源码编译 (开发者)

如果你想自己定制功能，可以按以下步骤编译：

1.  **环境准备**
    *   Go 1.25+
    *   Node.js 22+
    *   Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

2.  **克隆与构建**
    ```bash
    git clone https://github.com/jym66/Q-Solver.git
    cd Q-Solver
    
    # 开发模式 (热重载)
    wails dev
    
    # 编译生产版本
    wails build
    ```


## ⚙️ 配置指南

为了让 Q-Solver 发挥最大威力，你需要配置 AI 模型：

1.  启动软件，点击右上角设置图标（或使用快捷键）。
2.  **API Key**：输入你的 Google Gemini API Key（[点击获取](https://aistudio.google.com/app/apikey)）。
3.  **模型选择**：推荐使用 Gemini系列。
4.  **快捷键**：根据习惯设置截图和隐藏的快捷键。

## ⌨️ 默认快捷键

| 快捷键 | 功能 |
|--------|------|
| `F8` | 截取当前屏幕区域并发送分析请求 |
| `F9` | 快速隐藏或呼出主窗口 |
| `F10` | 开启/关闭鼠标穿透 |
| `Alt + ↑/↓/←/→` | 像素级移动窗口位置 |
| `Alt + PgUp/PgDn` | 翻页查看长文本回复 |

> 💡 以上快捷键可在设置中自定义


## ⚠️ 免责声明

本项目仅供技术研究与学习交流。

🚫 严禁用于考试作弊、学术不端或违反法律法规的场景。

开发者不对使用本软件产生的任何后果负责。

## 📄 License

Based on CC BY-NC 4.0.

✅ 允许个人学习、修改、二次开发。

❌ 禁止商业用途（包括但不限于售卖、内置广告、付费订阅）。

---

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/jym66">jym66</a></p>
</div>
