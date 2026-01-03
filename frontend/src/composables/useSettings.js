import { reactive, computed, watch } from 'vue'
import { marked } from 'marked'
import { GetSettings, SyncSettingsToDefaultSettings, GetModels, TestConnection } from '../../wailsjs/go/main/App'

/**
 * 配置管理 composable
 * 配置完全由后端管理，前端只负责展示和临时编辑
 */
export function useSettings(shortcuts, tempShortcuts, uiState, callbacks) {
  // 当前生效的配置（从后端同步）
  const settings = reactive({
    apiKey: '',
    provider: 'google',
    baseURL: '',
    model: '',
    prompt: '',
    transparency: 1.0,
    mode: 'interview',
    keepContext: false,
    screenshotMode: 'window',
    resumePath: '',
    resumeContent: '',
    useMarkdownResume: false,
    compressionQuality: 80,
    sharpening: 0,
    grayscale: true,
    noCompression: false
  })

  // 临时编辑的配置（用于设置面板）
  const tempSettings = reactive({ ...settings })

  // 计算属性
  const renderedPrompt = computed(() => marked.parse(tempSettings.prompt || ''))
  const maskedKey = computed(() => {
    if (!settings.apiKey) return ''
    if (settings.apiKey.length < 8) return settings.apiKey
    return settings.apiKey.substring(0, 3) + '****' + settings.apiKey.substring(settings.apiKey.length - 4)
  })

  // 监听透明度变化（仅更新 UI，不通知后端）
  watch(() => tempSettings.transparency, (newVal) => {
    const opacity = 1.0 - newVal
    // 使用与设计系统搭配的深色调
    const app = document.getElementById('app')
    if (app) {
      // 基础色使用设计系统的 bg-base 颜色 (rgb(17, 24, 39))
      app.style.backgroundColor = `rgba(17, 24, 39, ${0.85 + opacity * 0.15})`
    }
  })

  // 监听 API Key 变化（只有真正变化时才重置状态）
  let lastApiKey = ''
  watch(() => tempSettings.apiKey, (newVal) => {
    // 只有当 API Key 真正变化时才重置状态
    if (newVal !== lastApiKey) {
      lastApiKey = newVal
      // 清空连通性测试结果
      uiState.connectionStatus = null
    }
  })

  /**
   * 从后端加载配置
   */
  async function loadSettings() {
    try {
      // 从后端获取配置
      const backendConfig = await GetSettings()

      // 应用配置到本地状态
      applyConfig(backendConfig)

      // 同步快捷键
      if (backendConfig.shortcuts) {
        Object.assign(shortcuts, backendConfig.shortcuts)
      }

      // 如果有 API Key，标记为已验证
      if (settings.apiKey) {
        uiState.isKeyValid = true

      } else {

      }
    } catch (e) {
      console.error('loadSettings error', e)
    }
  }

  /**
   * 应用配置到本地状态
   */
  function applyConfig(config) {
    settings.apiKey = config.apiKey || ''
    settings.provider = config.provider || 'google'
    settings.baseURL = config.baseURL || ''
    settings.model = config.model || 'gemini-2.5-flash'
    settings.prompt = config.prompt || ''
    settings.compressionQuality = config.compressionQuality || 80
    settings.sharpening = config.sharpening || 0
    settings.grayscale = config.grayscale !== undefined ? config.grayscale : true
    settings.noCompression = config.noCompression || false
    settings.keepContext = config.keepContext || false
    settings.resumePath = config.resumePath || ''
    settings.resumeContent = config.resumeContent || ''
    settings.useMarkdownResume = config.useMarkdownResume || false
    settings.screenshotMode = config.screenshotMode || 'window'

    // 透明度转换
    if (config.opacity !== undefined) {
      settings.transparency = 1.0 - config.opacity
      const app = document.getElementById('app')
      if (app) {
        // 基础色使用设计系统的 bg-base 颜色 (rgb(17, 24, 39))
        app.style.backgroundColor = `rgba(17, 24, 39, ${0.85 + config.opacity * 0.15})`
      }
    }

    // 同步到 tempSettings，确保设置面板显示正确的值
    Object.assign(tempSettings, JSON.parse(JSON.stringify(settings)))
  }


  /**
   * 刷新模型列表
   */
  async function refreshModels() {
    if (!tempSettings.apiKey) {
      if (callbacks.showToast) callbacks.showToast('请先填写 API Key', 'warning')
      return
    }
    await fetchModels(tempSettings.apiKey, tempSettings.baseURL)
    if (uiState.availableModels.length > 0) {
      if (callbacks.showToast) callbacks.showToast(`已加载 ${uiState.availableModels.length} 个模型`, 'success')
    }
  }

  /**
   * 测试模型连通性
   */
  async function testConnection() {
    if (!tempSettings.model) {
      if (callbacks.showToast) callbacks.showToast('请先选择模型', 'warning')
      return
    }

    uiState.isTestingConnection = true
    uiState.connectionStatus = null

    try {
      const result = await TestConnection(tempSettings.apiKey, tempSettings.baseURL, tempSettings.model)
      if (result === '') {
        uiState.connectionStatus = {
          type: 'success',
          icon: '✅',
          message: `模型 ${tempSettings.model} 连接成功`
        }
        if (callbacks.showToast) callbacks.showToast('连接测试成功', 'success')
      } else {
        uiState.connectionStatus = {
          type: 'error',
          icon: '❌',
          message: result
        }
        if (callbacks.showToast) callbacks.showToast('连接测试失败', 'error')
      }
    } catch (e) {
      console.error('连接测试异常:', e)
      uiState.connectionStatus = {
        type: 'error',
        icon: '❌',
        message: e.message || '连接测试失败'
      }
    } finally {
      uiState.isTestingConnection = false
    }
  }

  /**
   * 获取模型列表
   */
  async function fetchModels(apiKey, baseURL) {
    if (!apiKey) return
    uiState.isLoadingModels = true
    try {
      const models = await GetModels(apiKey, baseURL || '')
      if (models && models.length > 0) {
        uiState.availableModels = models
        if (!tempSettings.model || tempSettings.model === 'auto') {
          tempSettings.model = models[0]
        }
      }
    } catch (e) {
      console.error("获取模型列表失败", e)
    } finally {
      uiState.isLoadingModels = false
    }
  }

  /**
   * 保存设置到后端（不再使用 localStorage）
   */
  async function saveSettings() {
    try {
      // 同步快捷键
      Object.assign(shortcuts, JSON.parse(JSON.stringify(tempShortcuts)))

      // 构建要保存的配置
      const configToSave = {
        apiKey: tempSettings.apiKey,
        baseURL: tempSettings.baseURL,
        model: tempSettings.model,
        prompt: tempSettings.prompt,
        opacity: 1.0 - tempSettings.transparency,
        keepContext: tempSettings.keepContext,
        screenshotMode: tempSettings.screenshotMode,
        compressionQuality: tempSettings.compressionQuality,
        sharpening: tempSettings.sharpening,
        grayscale: tempSettings.grayscale,
        noCompression: tempSettings.noCompression,
        resumePath: tempSettings.resumePath,
        resumeContent: tempSettings.resumeContent,
        useMarkdownResume: tempSettings.useMarkdownResume,
        provider: tempSettings.provider,
        shortcuts: tempShortcuts
      }

      // 发送到后端保存（后端会持久化到文件）
      const err = await SyncSettingsToDefaultSettings(JSON.stringify(configToSave))

      if (err) {
        if (callbacks.showToast) callbacks.showToast(err)
      } else {
        if (callbacks.showToast) callbacks.showToast('设置已保存', 'success')
        // 更新本地状态
        Object.assign(settings, tempSettings)
        if (callbacks.resetStatus) callbacks.resetStatus()

        if (callbacks.closeSettings) callbacks.closeSettings()
      }
    } catch (e) {
      console.error('保存设置失败', e)
      if (callbacks.showToast) callbacks.showToast('保存失败', 'error')
    }
  }

  /**
   * 重置临时设置为当前生效的设置
   * 用于取消编辑时恢复原值
   */
  function resetTempSettings() {
    Object.assign(tempSettings, settings)
    // 恢复 UI 透明度
    const opacity = 1.0 - settings.transparency
    const app = document.getElementById('app')
    if (app) app.style.backgroundColor = `rgba(0, 0, 0, ${opacity * 0.8})`
  }

  /**
   * 打开设置面板时调用
   * 复制当前配置到临时变量，初始化状态
   */
  function openSettings() {
    // 复制配置到临时变量
    Object.assign(tempSettings, JSON.parse(JSON.stringify(settings)))
    Object.assign(tempShortcuts, JSON.parse(JSON.stringify(shortcuts)))

    // 更新 lastApiKey 避免触发 watch
    lastApiKey = settings.apiKey

    // 清空连通性状态
    uiState.connectionStatus = null

    // 如果有 API Key，自动加载模型列表
    if (settings.apiKey) {
      fetchModels(settings.apiKey)
    }
  }

  return {
    settings,
    tempSettings,
    renderedPrompt,
    maskedKey,
    loadSettings,
    refreshModels,
    testConnection,
    fetchModels,
    saveSettings,
    resetTempSettings,
    openSettings
  }
}
