<div align="center">
  <img src="assets/banner.jpg" alt="Q-Solver Banner" width="100%" style="border-radius: 12px; box-shadow: 0 8px 30px rgba(0,0,0,0.12);">

  <h1 style="font-size: 3rem; margin: 20px 0;">Q-Solver</h1>
  <p style="font-size: 1.2rem; color: #666;">📝 您的桌面笔试助手：自动识别题目，深度思考，即刻解答</p>

  <p>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white" alt="Go"></a>
    <a href="https://vuejs.org"><img src="https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white" alt="Vue"></a>
    <a href="https://wails.io"><img src="https://img.shields.io/badge/Wails-v2-E30613?logo=wails&logoColor=white" alt="Wails"></a>
    <img src="https://img.shields.io/badge/Platform-Windows-0078D6?logo=windows&logoColor=white" alt="Windows">
    <img src="https://img.shields.io/badge/License-CC_BY--NC_4.0-lightgrey" alt="License">
  </p>

  <p>
    <a href="#✨-核心特性">特性</a> •
    <a href="#🚀-快速开始">安装</a> •
    <a href="#⚙️-配置指南">配置</a> •
    <a href="#⌨️-快捷键">快捷键</a> •
    <a href="#🛠️-技术栈">开发</a>
  </p>
</div>

<br>

## 📖 简介

**Q-Solver** 是一款专为求职者和学生量身打造的桌面笔试助手。它深度集成了 **OpenAI**、**Google Gemini** 和 **Anthropic Claude** 等顶级大模型，通过极简的截图交互，为您提供实时的代码分析、问题解答和内容创作辅助。

不同于传统的聊天机器人，Q-Solver 具备**深度思考链**可视化、**无感隐身/防录屏**和**简历上下文感知**等独创功能，完美融入您的笔试环境，是真正即刻解答的 AI 助手。


## ✨ 核心特性


### 🛡️ 无感交互模式 (Stealth Mode)
专为笔试/多任务环境打造。激活后：
- **无边框/半透明**：完美悬浮于其他窗口之上
- **防抢焦**：操作其他软件时不会自动隐藏，始终显示答案
- **防录屏/隐身**：特殊窗口属性，避开部分即时通讯软件的截屏和录屏检测
- **鼠标穿透**：点击窗口区域直达后方应用，互不干扰

### 🔌 全面模型支持 (Model Support)
- **原生 SDK 集成**：内置 Google Gemini, Anthropic Claude 和 OpenAI 原生 SDK，非简单的 HTTP 转发，确保最佳的流式响应体验和稳定性。
- **自定义接入**：支持 OneAPI 等聚合服务接入，灵活适配各种网络环境。

### 📄 简历与上下文感知
- **简历助手**：导入 PDF/Markdown 简历，AI 将结合您的背景生成个性化回答（适合模拟面试、文书润色）。
- **智能记忆**：支持多轮对话上下文，并在新话题开始时自动清理。



## 🖼️ 功能预览

<div align="center">
  <img src="assets/demo.gif" style="border-radius: 8px; width: 100%; box-shadow: 0 4px 20px rgba(0,0,0,0.1);" />
  <p><i>👆 实时演示：一键截图 -> 思考 -> 解答</i></p>
</div>

| | |
|:---:|:---:|
| <img src="assets/img1.png" style="border-radius: 8px; width: 100%;" /> | <img src="assets/img2.png" style="border-radius: 8px; width: 100%;" /> |
| <img src="assets/img3.png" style="border-radius: 8px; width: 100%;" /> | <img src="assets/img4.png" style="border-radius: 8px; width: 100%;" /> |

---

## 🚀 快速开始

### 方式一：下载安装包 (推荐)
前往 [Releases](https://github.com/jym66/Q-Solver/releases) 页面下载最新的 Windows 安装包 (`.exe`)。

### 方式二：源码编译
如果您是开发者，可以克隆源码进行二次开发：

```bash
# 前置要求：Go 1.25+, Node.js 22+, Wails CLI (go install github.com/wailsapp/wails/v2/cmd/wails@latest)

# 1. 克隆仓库
git clone https://github.com/jym66/Q-Solver.git
cd Q-Solver

# 2. 开发模式运行 (支持热重载)
wails dev

# 3. 编译生产版本
wails build
```

---

## ⚙️ 配置指南

启动软件后点击右上角 **设置 (Settings)** 图标：

1.  **选择提供商**：支持 OpenAI, Gemini, Claude, 或 Custom。
2.  **API Key**：填入对应平台的 API Key。
3.  **模型参数**：
    *   **Temperature**: 控制回答的随机性 (0.0 - 2.0)
    *   **Thinking Budget**: 设置 o1/Claude-3.5 思考过程的 Token 预算
4.  **Custom 自定义模式**：
    *   选择 `Custom` 提供商
    *   填入聚合 API 地址 (如 OneAPI)
    *   支持根据模型前缀自动切换底层协议 (Gemini/Claude)

---

## ⌨️ 快捷键

| 按键 | 作用 |
| :--- | :--- |
| **F8** | **区域截图并提问** (核心功能) |
| **F9** | 显示 / 隐藏主窗口 |
| **F10** | 开启 / 关闭鼠标穿透 |
| **Alt + ⬆️/⬇️/⬅️/➡️** | 微调窗口位置 |
| **Alt + PgUp/PgDn** | 快速翻页查看历史 |

---

## 支持

如果您喜欢这个项目，请给它一个 ⭐️ **Star**，这对我们非常有帮助！

**免责声明**：本项目仅供技术研究与个人学习使用，请勿用于任何非法用途。

<br>

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/jym66">jym66</a></p>
</div>
