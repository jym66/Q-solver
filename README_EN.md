<div align="center">
  <img src="assets/banner.jpg" alt="Q-Solver Banner" width="100%" style="border-radius: 12px; box-shadow: 0 8px 30px rgba(0,0,0,0.12);">

  <h1 style="font-size: 3rem; margin: 20px 0;">Q-Solver</h1>
  <p style="font-size: 1.2rem; color: #666;">📝 All-in-One Desktop AI Assistant: One-Click Captions, Deep Thinking, Instant Answers</p>

  <p>
    <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white" alt="Go"></a>
    <a href="https://vuejs.org"><img src="https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white" alt="Vue"></a>
    <a href="https://wails.io"><img src="https://img.shields.io/badge/Wails-v2-E30613?logo=wails&logoColor=white" alt="Wails"></a>
    <img src="https://img.shields.io/badge/Platform-Windows-0078D6?logo=windows&logoColor=white" alt="Windows">
    <img src="https://img.shields.io/badge/License-CC_BY--NC_4.0-lightgrey" alt="License">
  </p>

  <p>
    <a href="#✨-Features">Features</a> •
    <a href="#🚀-Quick-Start">Install</a> •
    <a href="#⚙️-Configuration">Config</a> •
    <a href="#⌨️-Shortcuts">Shortcuts</a> •
    <a href="README.md">中文文档</a>
  </p>
</div>

> [!WARNING]
> **Development Phase Note**: This project is currently in early active development (Alpha/Beta) and may contain minor bugs or unstable features. If you encounter any issues, please submit an Issue, and we will fix it ASAP!

<br>

## 📖 Introduction

**Q-Solver** is a desktop AI assistant tailored for written tests and efficient multitasking. It deeply integrates top-tier LLMs like **OpenAI**, **Google Gemini**, and **Anthropic Claude**. With minimal screenshot interactions, it provides real-time code analysis, Q&A, and content creation assistance.

Unlike traditional chatbots, Q-Solver features unique capabilities such as **Reasoning Chain Visualization**, **Stealth/Anti-Recording Mode**, and **Resume Context Awareness**, seamlessly blending into your workflow as a truly "understanding" AI assistant.

---

## ✨ Features

### 🛡️ Stealth Mode
Designed for high-focus or written test environments. When activated:
- **Borderless / Semi-Transparent**: Floats perfectly above other windows.
- **Anti-Focus Stealing**: Does not auto-hide when operating other software; answers remain visible.
- **Anti-Screen Recording**: Special window attributes to evade some screen capture and recording detection.
- **Click-Through**: Click through the window to interact with the application behind it.

### 🔌 Comprehensive Model Support
- **Native SDK Integration**: Built-in native SDKs for Google Gemini, Anthropic Claude, and OpenAI (not just HTTP forwarding) to ensure optimal streaming response and stability.
- **Custom Access**: Supports aggregation services like OneAPI for flexible network adaptation.

### 📄 Context & Resume Awareness
- **Resume Assistant**: Import PDF/Markdown resumes, and the AI generates personalized answers based on your background (ideal for mock interviews and resume polishing).
- **Smart Memory**: Supports multi-turn conversation context, automatically clearing when a new topic starts.

### 🎙️ Gemini Live API (Experimental)
- **Real-time Voice Conversation**: Supports bidirectional real-time voice calls with Gemini models, providing an ultra-low latency interaction experience.
- **Interviewer Mode**: Simulates real technical interview scenarios, transcribing interviewer questions in real-time, with AI providing instant answer suggestions.
- **Note**: This feature is currently in an **experimental stage**. Connection instability or audio interruptions may occur due to network conditions.
- **Supported Models**: Please select `gemini-2.0-flash-exp`.

### ⚡ Extreme Performance
- **Go + Wails**: Modern tech stack with native-level performance and extremely low memory footprint.
- **Vue 3 Frontend**: Responsive, modern UI that is smooth and fluid.
- **Global Shortcuts**: Wake up with `F8` in a split second without switching windows.

---

## 🖼️ Preview

<div align="center">
  <img src="assets/demo.gif" style="border-radius: 8px; width: 100%; box-shadow: 0 4px 20px rgba(0,0,0,0.1);" />
  <p><i>👆 Live Demo: Screenshot -> Thinking -> Answer</i></p>
</div>

| | |
|:---:|:---:|
| <img src="assets/img1.png" style="border-radius: 8px; width: 100%;" /> | <img src="assets/img2.png" style="border-radius: 8px; width: 100%;" /> |
| <img src="assets/img3.png" style="border-radius: 8px; width: 100%;" /> | <img src="assets/img4.png" style="border-radius: 8px; width: 100%;" /> |

---

## 🚀 Quick Start

### Method 1: Download Installer (Recommended)
Go to the [Releases](https://github.com/jym66/Q-Solver/releases) page to download the latest Windows installer (`.exe`).

### Method 2: Build from Source
If you are a developer, you can clone the source code for secondary development:

```bash
# Prerequisites: Go 1.25+, Node.js 22+, Wails CLI

# 1. Clone repository
git clone https://github.com/jym66/Q-Solver.git
cd Q-Solver

# 2. Run in Dev Mode (Hot Reload)
wails dev

# 3. Build Production Version
wails build
```

---

## ⚙️ Configuration

Launch the software and click the **Settings** icon in the top right corner:

1.  **Select Provider**: Supports OpenAI, Gemini, Claude, or Custom.
2.  **API Key**: Enter the API Key for the corresponding platform.
3.  **Model Parameters**:
    *   **Temperature**: Control answer randomness (0.0 - 2.0).
    *   **Thinking Budget**: Set Token budget for o1/Claude-3.5 thinking process.
4.  **Custom Mode**:
    *   Select `Custom` provider.
    *   Enter API URL (e.g., OneAPI).
    *   Supports automatic protocol switching based on model prefix (Gemini/Claude).

---

## ⌨️ Shortcuts

| Key | Function |
| :--- | :--- |
| **F8** | **capture Screenshot & Ask** (Core Function) |
| **F9** | Show / Hide Main Window |
| **F10** | Toggle Mouse Click-Through |
| **Alt + ⬆️/⬇️/⬅️/➡️** | Fine-tune Window Position |
| **Alt + PgUp/PgDn** | Quick History Navigation |

---

## Support

If you like this project, please give it a ⭐️ **Star**, it helps us a lot!

**Disclaimer**: This project is for technical research and personal learning only. Please do not use it for any illegal purposes.

<br>

<div align="center">
  <p>Made with ❤️ by <a href="https://github.com/jym66">jym66</a></p>
</div>
