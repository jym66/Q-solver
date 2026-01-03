<template>
  <Transition name="fade">
    <!-- <InitLoading v-if="initStatus !== 'ready'" :status="initStatus" /> -->
  </Transition>
  <TopBar :shortcuts="shortcuts" :activeButtons="activeButtons" :isClickThrough="isClickThrough"
    :statusIcon="statusIcon" :statusText="statusText" :balance="balance" :isRefreshingBalance="isRefreshingBalance"
    :settings="settings" :isStealthMode="isStealthMode" @openSettings="openSettings" @refreshBalance="refreshBalance"
    @quit="quit" />

  <WelcomeView v-if="!hasStarted" :solveShortcut="solveShortcut" :toggleShortcut="shortcuts.toggle?.keyName || 'Alt+H'"
    :initStatus="initStatus" />

  <div v-else id="main-interface" class="main-interface" :class="{ visible: mainVisible }">
    <div class="left-panel" id="history-list">
      <div v-if="history.length === 0" class="history-empty">
        <span class="empty-icon">📝</span>
        <span class="empty-text">暂无记录</span>
      </div>
      <HistoryItem v-for="(h, idx) in history" :key="idx" :summary="h.summary" :time="h.time"
        :isActive="idx === activeHistoryIndex" :isFirst="idx === 0" :previewHtml="renderMarkdown(h.summary)"
        @select="selectHistory(idx)" @delete="deleteHistory(idx)" @export-markdown="exportMarkdown(idx)"
        @export-image="exportImage(idx)" />
    </div>
    <div class="right-panel">
      <ErrorView v-if="errorState.show" :errorState="errorState" :solveShortcut="solveShortcut" />
      <LoadingView v-else-if="isLoading" />
      <div v-else id="content" class="markdown-body">
        <div v-html="renderedContent"></div>
        <div v-if="isAppending" class="append-loading">
          <div class="ai-icon">
            <div class="ai-icon-inner"></div>
          </div>
          <span class="text">AI 正在思考</span>
          <div class="wave-dots">
            <span></span><span></span><span></span>
          </div>
        </div>
      </div>
    </div>
  </div>


  <!-- Settings Modal -->
  <div v-if="uiState.showSettings" class="modal" id="settings-modal" style="display: flex">

    <div class="modal-content">
      <div class="modal-warning-banner"
        style="background: rgba(255, 169, 64, 0.15); border: 1px solid rgba(255, 169, 64, 0.3); border-radius: 50px; padding: 6px 20px; color: #ffc069; font-size: 12px; display: flex; align-items: center; justify-content: center; margin: 12px auto 4px auto; width: fit-content;">
        ⚠️ 当前窗口已获取焦点，关闭设置后将自动恢复防抢焦模式
      </div>
      <div class="modal-header">
        <div class="tabs">
          <div class="tab" :class="{ active: uiState.activeTab === 'general' }" @click="uiState.activeTab = 'general'">
            常规设置</div>
          <div class="tab" :class="{ active: uiState.activeTab === 'model' }" @click="uiState.activeTab = 'model'">模型设置
          </div>
          <div class="tab" :class="{ active: uiState.activeTab === 'screenshot' }"
            @click="uiState.activeTab = 'screenshot'">截图设置</div>
          <div class="tab" :class="{ active: uiState.activeTab === 'resume' }" @click="uiState.activeTab = 'resume'">
            简历设置</div>
          <div class="tab" :class="{ active: uiState.activeTab === 'account' }" @click="uiState.activeTab = 'account'">
            提供商</div>
        </div>
        <span class="close-btn" @click="closeSettings">&times;</span>
      </div>
      <div class="modal-body">
        <div v-show="uiState.activeTab === 'account'">
          <ProviderSelect v-model:provider="tempSettings.provider" v-model:apiKey="tempSettings.apiKey"
            v-model:baseURL="tempSettings.baseURL" />
        </div>

        <div v-show="uiState.activeTab === 'model'">
          <div class="form-group">
            <div class="model-header">
              <label>模型选择</label>
              <div class="model-actions">
                <button class="btn-icon" @click="refreshModels"
                  :disabled="uiState.isLoadingModels || !tempSettings.apiKey" title="刷新模型列表">
                  <span :class="{ spin: uiState.isLoadingModels }">🔄</span>
                </button>
                <button class="btn-icon" @click="testConnection"
                  :disabled="uiState.isTestingConnection || !tempSettings.model" title="测试模型连通性">
                  <span :class="{ spin: uiState.isTestingConnection }">{{ uiState.isTestingConnection ? '⏳' : '🔗'
                  }}</span>
                </button>
              </div>
            </div>
            <ModelSelect v-model="tempSettings.model" :models="uiState.availableModels"
              :loading="uiState.isLoadingModels" />

            <!-- 连通性测试结果 -->
            <div v-if="uiState.connectionStatus" class="connection-status" :class="uiState.connectionStatus.type">
              <span class="status-icon">{{ uiState.connectionStatus.icon }}</span>
              <span class="status-text">{{ uiState.connectionStatus.message }}</span>
            </div>

            <p v-if="!tempSettings.apiKey" class="hint-text" style="color: #ff9800; margin-top: 8px;">
              ⚠️ 请先填写 API Key
            </p>
          </div>

          <div class="form-group">
            <div class="prompt-header">
              <label for="prompt-text" style="margin-bottom: 0">系统提示词 (Prompt)</label>
              <div class="prompt-tabs">
                <div class="prompt-tab" :class="{ active: uiState.promptTab === 'edit' }"
                  @click="uiState.promptTab = 'edit'">编辑
                </div>
                <div class="prompt-tab" :class="{ active: uiState.promptTab === 'preview' }"
                  @click="uiState.promptTab = 'preview'">预览</div>
              </div>
            </div>

            <textarea v-show="uiState.promptTab === 'edit'" id="prompt-text" class="prompt-textarea" rows="10"
              v-model="tempSettings.prompt" placeholder="请输入提示词 (支持 Markdown)..."></textarea>

            <div v-show="uiState.promptTab === 'preview'" class="prompt-preview markdown-body" v-html="renderedPrompt">
            </div>
          </div>
        </div>

        <div v-show="uiState.activeTab === 'general'">
          <div class="form-group">
            <div class="context-setting">
              <div class="setting-row">
                <div class="setting-info">
                  <span class="setting-title">保存上下文</span>
                  <span class="setting-desc">开启后，每次对话将包含之前的历史记录</span>
                </div>
                <label class="switch">
                  <input type="checkbox" v-model="tempSettings.keepContext">
                  <span class="slider round"></span>
                </label>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label>快捷键配置 (点击录制)</label>
            <div class="shortcut-list">
              <div class="shortcut-item" v-for="key in shortcutActions" :key="key.action">
                <span>{{ key.label }}</span>
                <button class="btn-record" :class="{ recording: recordingAction === key.action }"
                  @click="recordKey(key.action)">
                  {{ recordingAction === key.action ? recordingText : (tempShortcuts[key.action]?.keyName ||
                    key.default) }}
                </button>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label for="opacity-slider">窗口透明度: <span>{{ Math.round(tempSettings.transparency * 100) }}%</span></label>
            <input type="range" id="opacity-slider" min="0.0" max="1.0" step="0.05"
              v-model.number="tempSettings.transparency" />
          </div>
        </div>

        <div v-show="uiState.activeTab === 'screenshot'">
          <ScreenshotSettings :modelValue="tempSettings" @update:modelValue="Object.assign(tempSettings, $event)" />
        </div>

        <div v-show="uiState.activeTab === 'resume'" style="height: 100%">
          <ResumeImport :resumePath="tempSettings.resumePath" :rawContent="resumeState.rawContent"
            :isParsing="resumeState.isParsing" :currentModel="tempSettings.model"
            v-model:useMarkdownResume="tempSettings.useMarkdownResume"
            @update:rawContent="val => resumeState.rawContent = val" @select-resume="selectResume"
            @clear-resume="clearResume" @parse-resume="parseResume" />
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-primary" @click="saveSettings">保存</button>
      </div>
    </div>
  </div>

  <div id="toast-container">
    <div v-for="(t, i) in toasts" :key="t.id || i" class="toast" :class="[t.type, { show: t.show }]">{{ t.text }}
    </div>
  </div>


</template>

<script setup>
import { reactive, ref, onMounted, watch, nextTick, computed } from 'vue'
import ResumeImport from './components/ResumeImport.vue'
import ScreenshotSettings from './components/ScreenshotSettings.vue'
import ProviderSelect from './components/ProviderSelect.vue'
import WelcomeView from './components/WelcomeView.vue'
import ErrorView from './components/ErrorView.vue'
import LoadingView from './components/LoadingView.vue'
// import InitLoading from './components/InitLoading.vue'
import TopBar from './components/TopBar.vue'
import ModelSelect from './components/ModelSelect.vue'
import HistoryItem from './components/HistoryItem.vue'
import { EventsOn, Quit } from '../wailsjs/runtime/runtime'
import { StopRecordingKey, SelectResume, ClearResume, RestoreFocus, RemoveFocus, ParseResume, GetInitStatus } from '../wailsjs/go/main/App'

import { useUI } from './composables/useUI'
import { useStatus } from './composables/useStatus'
import { useBalance } from './composables/useBalance'
import { useShortcuts } from './composables/useShortcuts'
import { useSettings } from './composables/useSettings'
import { useSolution } from './composables/useSolution'

// 样式导入
import './App.global.css'
import './App.scoped.css'

const uiState = reactive({
  showSettings: false,
  activeTab: 'general',
  availableModels: [],
  isLoadingModels: false,
  isModelDropdownOpen: false,
  promptTab: 'edit',
  isTestingConnection: false,
  connectionStatus: null,
})

const {
  toasts, activeButtons, isClickThrough, mainVisible, isStealthMode, hasStarted,
  showToast, flash, quit
} = useUI()

const {
  shortcuts, tempShortcuts, recordingAction, recordingText, shortcutActions, recordKey
} = useShortcuts()

// Settings callbacks placeholder
const settingsCallbacks = {}

const {
  settings, tempSettings, renderedPrompt, maskedKey,
  loadSettings, refreshModels, testConnection, fetchModels, saveSettings, resetTempSettings, openSettings: initSettings
} = useSettings(shortcuts, tempShortcuts, uiState, settingsCallbacks)

const resumeState = reactive({
  rawContent: '',
  isParsing: false
})

watch(() => resumeState.rawContent, (newVal) => {
  tempSettings.resumeContent = newVal || ''
})

async function selectResume() {
  const path = await SelectResume()
  if (path) {
    tempSettings.resumePath = path
    resumeState.rawContent = '' // Reset parsed content on new file
    showToast('简历已选择', 'success')
  }
}

async function clearResume() {
  await ClearResume()
  tempSettings.resumePath = ''
  resumeState.rawContent = ''
}
async function parseResume() {
  if (!tempSettings.resumePath) return

  resumeState.isParsing = true
  try {
    const result = await ParseResume()
    resumeState.rawContent = result
    showToast('简历解析成功', 'success')
  } catch (e) {
    console.error(e)
    showToast('解析失败: ' + e, 'error')
  } finally {
    resumeState.isParsing = false
  }
}

const {
  statusText, statusIcon, resetStatus
} = useStatus(settings)

const {
  balance, tempBalance, isRefreshingBalance, fetchBalance, refreshBalance
} = useBalance(settings, statusText, statusIcon, resetStatus)

const {
  renderedContent, history, activeHistoryIndex, isLoading, isAppending, shouldOverwriteHistory,
  errorState, renderMarkdown, selectHistory, handleStreamStart, handleStreamChunk, handleSolution, setStreamBuffer,
  deleteHistory, exportMarkdown, exportImage
} = useSolution(settings)

// Populate callbacks
settingsCallbacks.fetchBalance = fetchBalance
settingsCallbacks.resetStatus = resetStatus
settingsCallbacks.showToast = showToast
settingsCallbacks.setBalance = (val) => { balance.value = val }
settingsCallbacks.setTempBalance = (val) => { tempBalance.value = val }
settingsCallbacks.updateBalanceFromTemp = () => { balance.value = tempBalance.value }
settingsCallbacks.onKeyChange = () => { tempBalance.value = null }
settingsCallbacks.closeSettings = closeSettings

function openSettings() {
  RestoreFocus()
  // 初始化临时设置
  initSettings()
  tempBalance.value = balance.value

  // 加载模型列表
  if (settings.apiKey) {
    fetchModels(settings.apiKey)
  }

  // 加载简历内容
  if (settings.resumeContent) {
    resumeState.rawContent = settings.resumeContent
  }

  uiState.showSettings = true
}

function closeSettings() {
  RemoveFocus()
  uiState.showSettings = false
  if (recordingAction.value) {
    StopRecordingKey()
  }
  recordingAction.value = null
  recordingText.value = ''
  // 恢复所有临时设置到原值（包括透明度）
  resetTempSettings()
}

const solveShortcut = computed(() => shortcuts.solve?.keyName || 'F8')

const initStatus = ref('initializing')
// Lifecycle
onMounted(() => {
  // localStorage.clear()
  GetInitStatus().then(status => {
    initStatus.value = status
  })

  EventsOn('init-status', (status) => {
    initStatus.value = status
  })

  loadSettings().then(() => {
    resetStatus()
  })

  // Event Listeners
  EventsOn('key-recorded', (data) => {
    if (data && data.action) {
      if (tempShortcuts[data.action]) {
        tempShortcuts[data.action].keyName = data.keyName
        tempShortcuts[data.action].vkCode = data.comboID
      } else {
        tempShortcuts[data.action] = { keyName: data.keyName, vkCode: data.comboID }
      }

      if (recordingAction.value === data.action) {
        recordingText.value = data.keyName
      }
    }
  })

  EventsOn('shortcut-error', async (msg) => {
    showToast(msg, 'error', 2000)
    const targetAction = recordingAction.value
    recordingAction.value = null
    recordingText.value = ''
    StopRecordingKey()
    if (!targetAction) return

    try {
      if (shortcuts[targetAction] && shortcuts[targetAction].keyName) {
        tempShortcuts[targetAction] = JSON.parse(JSON.stringify(shortcuts[targetAction]))
      } else {
        delete tempShortcuts[targetAction]
      }
    } catch (e) {
      console.error("回滚配置失败", e)
    }
  })

  EventsOn('shortcut-saved', (action) => {
    if (recordingAction.value === action) {
      recordingAction.value = null
      showToast('快捷键已保存', 'success')
    }
  })

  EventsOn('start-solving', () => {
    errorState.show = false
    flash('solve')
    statusText.value = '正在思考...'
    statusIcon.value = '🟡'
    mainVisible.value = true
    hasStarted.value = true

    if (settings.keepContext && history.value.length > 0 && activeHistoryIndex.value === 0) {
      isLoading.value = false
      isAppending.value = true
      nextTick(() => {
        const contentDiv = document.getElementById('content')
        if (contentDiv) {
          contentDiv.scrollTop = contentDiv.scrollHeight
        }
      })
    } else {
      isLoading.value = true
      renderedContent.value = ''
      isAppending.value = false
    }
  })

  EventsOn('toggle-visibility', (isVisibleToCapture) => {
    flash('toggle')
    isStealthMode.value = isVisibleToCapture
    if (isVisibleToCapture) {
      showToast('隐身模式已开启 (录屏不可见)', 'info')
    } else {
      showToast('隐身模式已关闭 (录屏可见)', 'success')
    }
  })

  EventsOn('solution', (data) => {
    statusText.value = '解题完成'
    statusIcon.value = '📝'
    handleSolution(data)
    fetchBalance()
  })

  EventsOn('copy-code', () => {
    const old = statusText.value
    statusText.value = '已复制'
    setTimeout(() => (statusText.value = old), 2000)
  })

  EventsOn('click-through-state', (enabled) => {
    isClickThrough.value = enabled
    const el = document.getElementById('main-interface')
    if (el) el.style.pointerEvents = enabled ? "none" : "auto"
  })

  EventsOn("scroll-content", (direction) => {
    const contentDiv = document.getElementById('content')
    if (!contentDiv) return
    const scrollAmount = 50;
    if (direction === "up") {
      contentDiv.scrollBy({ top: -scrollAmount, behavior: 'smooth' });
    } else if (direction === "down") {
      contentDiv.scrollBy({ top: scrollAmount, behavior: 'smooth' });
    }
  });

  EventsOn('solution-stream-start', () => {
    hasStarted.value = true
    handleStreamStart()
  })

  EventsOn('solution-stream-chunk', (token) => {
    handleStreamChunk(token)
  })

  // 错误处理
  EventsOn('solution-error', (rawErrMsg) => {
    // A. 优先处理：用户取消 (这不是错误，是操作)
    if (rawErrMsg && (rawErrMsg.includes('context canceled') || rawErrMsg.includes('canceled'))) {
      handleUserCancellation()
      return
    }

    // 直接显示上游返回的错误信息
    let title = '请求出错'
    let desc = rawErrMsg || '未知错误'
    let icon = '❌'

    // 尝试解析 JSON 格式的错误
    try {
      const errObj = JSON.parse(rawErrMsg)
      if (errObj.message) {
        desc = errObj.message
      }
      if (errObj.statusCode) {
        title = `错误 ${errObj.statusCode}`
      }
    } catch (e) {
      // 如果不是 JSON，直接使用原始字符串
    }

    // 更新 UI 状态
    statusText.value = '出错'
    statusIcon.value = '🔴'
    errorState.show = true
    errorState.title = title
    errorState.desc = desc
    errorState.icon = icon
    errorState.rawError = rawErrMsg
    errorState.showDetails = false
    isLoading.value = false
    isAppending.value = false
    shouldOverwriteHistory.value = true
  })

  // 抽离取消逻辑
  function handleUserCancellation() {
    console.log('请求已由用户主动取消')

    // 恢复状态
    if (isLoading.value) isLoading.value = true
    if (isAppending.value) isAppending.value = true

    // 回滚历史记录逻辑
    if (history.value.length > 0 && activeHistoryIndex.value === 0) {
      const current = history.value[0]

      if (settings.keepContext) {
        const separator = '\n\n---\n\n'
        const lastIndex = current.full.lastIndexOf(separator)

        if (lastIndex !== -1) {
          current.full = current.full.substring(0, lastIndex)
          current.summary = current.full.substring(0, 30).replace(/\n/g, ' ') + '...'
          setStreamBuffer(current.full)
          renderedContent.value = renderMarkdown(current.full)

          isAppending.value = true
          isLoading.value = false
        } else {
          // 没找到分隔符，重置
          resetCurrentHistory(current)
        }
        shouldOverwriteHistory.value = false
      } else {
        // 不保留上下文，直接重置
        resetCurrentHistory(current)
        shouldOverwriteHistory.value = true
      }
    }
  }

  // 辅助函数
  function resetCurrentHistory(current) {
    current.full = ''
    current.summary = '正在思考...'
    renderedContent.value = ''
    setStreamBuffer('')
    isLoading.value = true
    statusText.value = '正在思考...'
    statusIcon.value = '🟡'
  }

  EventsOn('require-login', () => {
    uiState.showSettings = true
    uiState.activeTab = 'account'
    showToast('请先配置 API Key', 'warning')
  })

  const mainInterface = document.getElementById('main-interface')
  if (mainInterface) mainInterface.style.pointerEvents = 'auto'

  // document.addEventListener('contextmenu', event => event.preventDefault());

  document.addEventListener('keydown', event => {
    if (
      event.key === 'F12' ||
      (event.ctrlKey && event.shiftKey && event.key === 'I') ||
      (event.ctrlKey && event.shiftKey && event.key === 'J') ||
      (event.ctrlKey && event.key === 'U')
    ) {
      event.preventDefault();
    }
  });
})
</script>
