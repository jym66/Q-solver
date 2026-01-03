import { reactive, computed, watch } from 'vue'
import { marked } from 'marked'
import { GetSettings, SyncSettingsToDefaultSettings, ValidateAPIKey, GetModels, GetBalance } from '../../wailsjs/go/main/App'

/**
 * 配置管理 composable
 * 配置完全由后端管理，前端只负责展示和临时编辑
 */
export function useSettings(shortcuts, tempShortcuts, uiState, callbacks) {
  // 当前生效的配置（从后端同步）
  const settings = reactive({
    apiKey: '',
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
    // 只更新 UI 样式，不调用后端
    const app = document.getElementById('app')
    if (app) app.style.backgroundColor = `rgba(0, 0, 0, ${opacity * 0.8})`
  })

  // 监听 API Key 变化（只有真正变化时才重置验证状态）
  let lastApiKey = ''
  watch(() => tempSettings.apiKey, (newVal) => {
    // 只有当 API Key 真正变化时才重置验证状态
    if (newVal !== lastApiKey) {
      lastApiKey = newVal
      // 如果新值等于已保存的 settings.apiKey，保持验证状态
      if (newVal === settings.apiKey && settings.apiKey) {
        // 不重置，保持原有验证状态
      } else {
        uiState.keyValidationError = ''
        uiState.isKeyValid = false
        if (callbacks.onKeyChange) callbacks.onKeyChange()
      }
    }
  })

  /**
   * 从后端加载配置
   * 先尝试迁移 localStorage，然后完全依赖后端
   */
  async function loadSettings() {
    try {
      // 先迁移 localStorage（如果存在）
      await migrateLocalStorage()

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
        if (callbacks.fetchBalance) callbacks.fetchBalance()
      } else {
        if (callbacks.setBalance) callbacks.setBalance(-1)
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
      if (app) app.style.backgroundColor = `rgba(0,0,0,${config.opacity * 0.8})`
    }

    // 同步到 tempSettings，确保设置面板显示正确的值
    Object.assign(tempSettings, JSON.parse(JSON.stringify(settings)))
  }

  /**
   * 迁移旧的 localStorage 配置到后端（只执行一次）
   */
  async function migrateLocalStorage() {
    const localConfigStr = localStorage.getItem('ghost_solver_config')
    if (!localConfigStr) return

    try {
      const localConfig = JSON.parse(localConfigStr)
      // 发送到后端
      await SyncSettingsToDefaultSettings(JSON.stringify(localConfig))
      // 删除 localStorage
      localStorage.removeItem('ghost_solver_config')
      console.log('已迁移 localStorage 配置到后端')
    } catch (e) {
      console.error('迁移 localStorage 配置失败', e)
    }
  }

  /**
   * 验证 API Key
   */
  async function verifyKey() {
    if (!tempSettings.apiKey) return false
    uiState.isValidatingKey = true
    uiState.keyValidationError = ''

    try {
      const err = await ValidateAPIKey(tempSettings.apiKey)
      if (err) {
        const msg = err.toLowerCase()
        if (msg.includes('令牌额度已用尽') || msg.includes('quota_exceeded') || msg.includes('余额不足')) {
          uiState.keyValidationError = '当前API-KEY余额不足'
        } else if (msg.includes('401')) {
          uiState.keyValidationError = 'API-KEY无效'
        } else {
          uiState.keyValidationError = 'API-KEY验证失败，请检查您的输入或网络连接。'
        }
        uiState.isKeyValid = false
        return false
      } else {
        if (callbacks.showToast) callbacks.showToast('验证成功', 'success')
        uiState.isKeyValid = true
        uiState.keyValidationError = ''
        try {
          const b = await GetBalance(tempSettings.apiKey)
          if (callbacks.setTempBalance) callbacks.setTempBalance(b)
          fetchModels(tempSettings.apiKey)
        } catch (e) {
          console.error('GetBalance error', e)
          if (callbacks.setTempBalance) callbacks.setTempBalance(null)
        }
        return true
      }
    } catch (e) {
      console.error("验证异常:", e)
      uiState.keyValidationError = '验证过程发生错误，请稍后重试。'
      uiState.isKeyValid = false
      return false
    } finally {
      uiState.isValidatingKey = false
    }
  }

  /**
   * 获取模型列表
   */
  async function fetchModels(apiKey) {
    if (!apiKey) return
    uiState.isLoadingModels = true
    try {
      const models = await GetModels(apiKey)
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
      // 如果 API Key 改变了且新 Key 未验证，先验证
      const keyChanged = tempSettings.apiKey !== settings.apiKey
      if (keyChanged && tempSettings.apiKey && !uiState.isKeyValid) {
        const isValid = await verifyKey()
        if (!isValid) return
      }

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
        if (callbacks.updateBalanceFromTemp) callbacks.updateBalanceFromTemp()
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

    // 如果有 API Key，标记为已验证
    if (settings.apiKey) {
      uiState.isKeyValid = true
      uiState.isEditingKey = false
    } else {
      uiState.isEditingKey = true
    }
  }

  return {
    settings,
    tempSettings,
    renderedPrompt,
    maskedKey,
    loadSettings,
    verifyKey,
    fetchModels,
    saveSettings,
    resetTempSettings,
    openSettings
  }
}
