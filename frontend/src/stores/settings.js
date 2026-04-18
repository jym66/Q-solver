import { defineStore } from 'pinia'
import { reactive, ref, computed, watch } from 'vue'
import { api } from '../services/api'
import { useUIStore } from './ui'
import { currentTheme, setTheme } from '../services/theme'

export const useSettingsStore = defineStore('settings', () => {
  const ui = useUIStore()

  const settings = reactive({
    apiKey: '',
    baseURL: 'https://api.openai.com/v1',
    model: '',
    assistantModel: '',
    prompt: '',
    domainId: 'general-assistant',
    transparency: 0,
    keepContext: false,
    screenshotMode: 'window',
    resumePath: '',
    resumeContent: '',
    compressionQuality: 80,
    sharpening: 0,
    grayscale: true,
    noCompression: false,
  })

  const tempSettings = reactive({ ...settings })
  const shortcuts = reactive({})
  const tempShortcuts = reactive({})
  const domainCategories = ref([])
  const resumeRawContent = ref('')
  const isResumeParsing = ref(false)

  const isMacOS = ref(
    typeof navigator !== 'undefined' &&
    (navigator.platform.toLowerCase().includes('mac') ||
      navigator.userAgent.toLowerCase().includes('mac'))
  )

  const recordingAction = ref(null)
  const recordingText = ref('')
  const shortcutActions = [
    { action: 'screenshot', label: '截图', default: 'F8', macDefault: 'Cmd+1' },
    { action: 'send', label: '发送解题', default: 'Ctrl+J', macDefault: 'Cmd+J' },
    { action: 'delete', label: '删除截图', default: 'Ctrl+D', macDefault: 'Cmd+D' },
    { action: 'toggle', label: '隐藏/显示', default: 'F9', macDefault: 'Cmd+2' },
    { action: 'clickthrough', label: '鼠标穿透', default: 'F10', macDefault: 'Cmd+3' },
    { action: 'move_up', label: '向上移动', default: 'Alt+Up', macDefault: 'Cmd+Option+Up' },
    { action: 'move_down', label: '向下移动', default: 'Alt+Down', macDefault: 'Cmd+Option+Down' },
    { action: 'move_left', label: '向左移动', default: 'Alt+Left', macDefault: 'Cmd+Option+Left' },
    { action: 'move_right', label: '向右移动', default: 'Alt+Right', macDefault: 'Cmd+Option+Right' },
    { action: 'scroll_up', label: '向上滚动', default: 'Alt+PgUp', macDefault: 'Cmd+Option+Shift+Up' },
    { action: 'scroll_down', label: '向下滚动', default: 'Alt+PgDn', macDefault: 'Cmd+Option+Shift+Down' },
  ]

  const maskedKey = computed(() => {
    if (!settings.apiKey) return ''
    if (settings.apiKey.length < 8) return settings.apiKey
    return settings.apiKey.substring(0, 3) + '****' + settings.apiKey.substring(settings.apiKey.length - 4)
  })

  const solveShortcut = computed(() => shortcuts.screenshot?.keyName || shortcuts.solve?.keyName || 'F8')
  const sendShortcut = computed(() => shortcuts.send?.keyName || 'Ctrl+J')
  const deleteShortcut = computed(() => shortcuts.delete?.keyName || 'Ctrl+D')
  const toggleShortcut = computed(() => shortcuts.toggle?.keyName || 'F9')

  const statusText = ref('就绪')
  const statusIcon = ref('●')

  function resetStatus() {
    if (!settings.apiKey) {
      statusText.value = '未配置'
      statusIcon.value = '!'
    } else {
      statusText.value = '已连接'
      statusIcon.value = '✓'
    }
  }

  watch(() => settings.apiKey, () => resetStatus(), { immediate: true })

  watch(() => tempSettings.transparency, (newVal) => {
    applyTransparency(1.0 - newVal)
  })

  window.addEventListener('theme-changed', () => {
    requestAnimationFrame(() => {
      applyTransparency(1.0 - (tempSettings.transparency ?? 0))
    })
  })

  watch(() => resumeRawContent.value, (newVal) => {
    tempSettings.resumeContent = newVal || ''
  })

  let themeInitialized = false
  watch(currentTheme, () => {
    if (!themeInitialized) {
      themeInitialized = true
      return
    }
    saveSettingsSilent()
  })

  function applyTransparency(opacity) {
    const root = document.documentElement
    if (!root) return
    const style = getComputedStyle(root)
    const base = style.getPropertyValue('--surface-base').trim()
    let r, g, b
    if (base.startsWith('rgba')) {
      const parts = base.match(/rgba\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*[\d.]+\s*\)/)
      if (parts) {
        r = parts[1]
        g = parts[2]
        b = parts[3]
      }
    }
    if (!r) {
      const theme = root.getAttribute('data-theme')
      if (theme === 'light') {
        r = '252'
        g = '252'
        b = '255'
      } else {
        r = '10'
        g = '10'
        b = '14'
      }
    }
    root.style.setProperty('--app-bg-r', r)
    root.style.setProperty('--app-bg-g', g)
    root.style.setProperty('--app-bg-b', b)
    root.style.setProperty('--app-bg-a', String(opacity))
    root.style.setProperty('--app-panel-a', String(Math.max(0, opacity * 0.5)))
    root.style.setProperty('--text-shadow-a', String(Math.max(0, (1 - opacity) * 0.8)))

    const isLight = root.getAttribute('data-theme') === 'light'
    root.style.setProperty('--surface-elevated-a', String(opacity * (isLight ? 0.96 : 0.92)))
    root.style.setProperty('--surface-card-a', String(opacity * (isLight ? 0.78 : 0.80)))
    root.style.setProperty('--surface-card-hover-a', String(opacity * (isLight ? 0.86 : 0.88)))
  }

  async function loadSettings() {
    try {
      const [backendConfig, categories] = await Promise.all([
        api.getSettings(),
        api.getDomainCategories(),
      ])
      if (categories) domainCategories.value = categories
      applyConfig(backendConfig)
      if (backendConfig.shortcuts) Object.assign(shortcuts, backendConfig.shortcuts)
      if (backendConfig.theme) setTheme(backendConfig.theme)
      if (settings.apiKey && (!settings.model || settings.model === 'auto')) {
        await fetchModels(settings.apiKey)
        if (ui.availableModels.length > 0 && !settings.model) {
          settings.model = ui.availableModels[0]
          tempSettings.model = ui.availableModels[0]
        }
      }
    } catch (e) {
      console.error('loadSettings error', e)
    }
  }

  function applyConfig(config) {
    settings.apiKey = config.apiKey || ''
    settings.baseURL = config.baseURL || 'https://api.openai.com/v1'
    settings.model = config.model || ''
    settings.assistantModel = config.assistantModel || ''
    settings.prompt = config.prompt || ''
    settings.domainId = config.domainId || 'general-assistant'
    settings.compressionQuality = config.compressionQuality || 80
    settings.sharpening = config.sharpening || 0
    settings.grayscale = config.grayscale !== undefined ? config.grayscale : true
    settings.noCompression = config.noCompression || false
    settings.keepContext = config.keepContext || false
    settings.resumePath = config.resumePath || ''
    settings.resumeContent = config.resumeContent || ''
    settings.screenshotMode = config.screenshotMode || 'window'

    const opacity = config.opacity !== undefined ? config.opacity : 1.0
    settings.transparency = 1.0 - opacity
    applyTransparency(opacity)
    Object.assign(tempSettings, JSON.parse(JSON.stringify(settings)))
  }

  async function fetchModels(apiKey) {
    if (!apiKey) return false
    ui.isLoadingModels = true
    try {
      const models = await api.getModels(apiKey)
      if (models && models.length > 0) {
        ui.availableModels = models
        if (!tempSettings.model || !models.includes(tempSettings.model)) {
          tempSettings.model = models[0]
        }
        return true
      }
      ui.availableModels = []
      return false
    } catch (e) {
      console.error('获取模型列表失败', e)
      ui.availableModels = []
      let errorMsg = '获取模型列表失败'
      try {
        const errObj = JSON.parse(e.message || e)
        if (errObj.message) errorMsg = errObj.message
      } catch (_) {
        errorMsg = e.message || '获取模型列表失败'
      }
      ui.showToast(errorMsg, 'error')
      return false
    } finally {
      ui.isLoadingModels = false
    }
  }

  async function refreshModels() {
    if (!tempSettings.apiKey) {
      ui.showToast('请先填写 API Key', 'warning')
      return
    }
    const success = await fetchModels(tempSettings.apiKey)
    if (success && ui.availableModels.length > 0) {
      ui.showToast(`已加载 ${ui.availableModels.length} 个模型`, 'success')
    }
  }

  async function testConnection() {
    if (!tempSettings.model) {
      ui.showToast('请先选择模型', 'warning')
      return
    }
    ui.isTestingConnection = true
    ui.connectionStatus = null
    try {
      const currentConfig = JSON.parse(JSON.stringify(settings))
      settings.baseURL = tempSettings.baseURL
      const result = await api.testConnection(tempSettings.apiKey, tempSettings.model)
      Object.assign(settings, currentConfig)
      if (result === '') {
        ui.connectionStatus = { type: 'success', icon: '✓', message: `模型 ${tempSettings.model} 连接成功` }
        ui.showToast('连接测试成功', 'success')
      } else {
        ui.connectionStatus = { type: 'error', icon: '!', message: result }
        ui.showToast('连接测试失败', 'error')
      }
    } catch (e) {
      ui.connectionStatus = { type: 'error', icon: '!', message: e.message || '连接测试失败' }
    } finally {
      ui.isTestingConnection = false
    }
  }

  function buildConfigToSave(sourceSettings, sourceShortcuts) {
    return {
      apiKey: sourceSettings.apiKey,
      baseURL: sourceSettings.baseURL,
      model: sourceSettings.model,
      assistantModel: sourceSettings.assistantModel,
      prompt: '',
      domainId: sourceSettings.domainId,
      opacity: 1.0 - sourceSettings.transparency,
      keepContext: sourceSettings.keepContext,
      screenshotMode: sourceSettings.screenshotMode,
      compressionQuality: sourceSettings.compressionQuality,
      sharpening: sourceSettings.sharpening,
      grayscale: sourceSettings.grayscale,
      noCompression: sourceSettings.noCompression,
      resumePath: sourceSettings.resumePath,
      resumeContent: sourceSettings.resumeContent,
      shortcuts: sourceShortcuts,
      theme: currentTheme.value,
    }
  }

  async function saveSettings() {
    try {
      if (!tempSettings.model && tempSettings.apiKey) {
        ui.showToast('正在自动获取模型...', 'info')
        await fetchModels(tempSettings.apiKey)
        if (!tempSettings.model && ui.availableModels.length > 0) {
          tempSettings.model = ui.availableModels[0]
        }
      }

      Object.assign(shortcuts, JSON.parse(JSON.stringify(tempShortcuts)))

      const err = await api.syncSettings(JSON.stringify(buildConfigToSave(tempSettings, tempShortcuts)))
      if (err) {
        ui.showToast(err, 'error')
      } else {
        ui.showToast('设置已保存', 'success')
        Object.assign(settings, tempSettings)
        resetStatus()
        closeSettings()
      }
    } catch (e) {
      console.error('保存设置失败', e)
      ui.showToast('保存失败', 'error')
    }
  }

  function openSettings() {
    api.restoreFocus()
    Object.assign(tempSettings, JSON.parse(JSON.stringify(settings)))
    Object.assign(tempShortcuts, JSON.parse(JSON.stringify(shortcuts)))
    ui.connectionStatus = null
    if (settings.resumeContent) resumeRawContent.value = settings.resumeContent
    if (settings.apiKey && ui.availableModels.length === 0) {
      fetchModels(settings.apiKey)
    }
    ui.showSettings = true
  }

  function closeSettings() {
    if (isResumeParsing.value) {
      ui.showToast('简历正在解析中，请稍候...', 'warning')
      return
    }
    api.removeFocus()
    ui.showSettings = false
    if (recordingAction.value) api.stopRecordingKey()
    recordingAction.value = null
    recordingText.value = ''
    resetTempSettings()
  }

  function resetTempSettings() {
    Object.assign(tempSettings, settings)
    applyTransparency(1.0 - settings.transparency)
  }

  function recordKey(action) {
    if (isMacOS.value) return
    recordingAction.value = action
    recordingText.value = '请按键...'
    api.startRecordingKey(action)
  }

  async function selectResume() {
    const path = await api.selectResume()
    if (path) {
      tempSettings.resumePath = path
      resumeRawContent.value = ''
      ui.showToast('已选择文件，正在自动解析...', 'info')
      await parseResume()
    }
  }

  async function clearResume() {
    await api.clearResume()
    tempSettings.resumePath = ''
    resumeRawContent.value = ''
  }

  async function parseResume() {
    if (!tempSettings.resumePath) return
    isResumeParsing.value = true
    try {
      const result = await api.parseResume()
      resumeRawContent.value = result
      tempSettings.resumeContent = result
      ui.showToast('简历解析成功', 'success')
    } catch (e) {
      console.error(e)
      ui.showToast(e?.message || '简历解析失败', 'error')
    } finally {
      isResumeParsing.value = false
    }
  }

  async function saveSettingsSilent() {
    try {
      await api.syncSettings(JSON.stringify(buildConfigToSave(settings, shortcuts)))
    } catch (e) {
      console.error('静默保存失败', e)
    }
  }

  return {
    settings, tempSettings,
    shortcuts, tempShortcuts,
    domainCategories,
    resumeRawContent, isResumeParsing,
    isMacOS,
    recordingAction, recordingText, shortcutActions,
    maskedKey, solveShortcut, sendShortcut, deleteShortcut, toggleShortcut,
    statusText, statusIcon, resetStatus,
    loadSettings, fetchModels, refreshModels, testConnection,
    saveSettings, saveSettingsSilent, openSettings, closeSettings, resetTempSettings,
    recordKey, selectResume, clearResume, parseResume,
    applyTransparency,
  }
})
