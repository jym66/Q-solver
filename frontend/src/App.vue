<template>
  <Transition name="fade">
    <!-- <InitLoading v-if="initStatus !== 'ready'" :status="initStatus" /> -->
  </Transition>
  <TopBar :shortcuts="shortcuts" :activeButtons="activeButtons" :isClickThrough="isClickThrough"
    :statusIcon="statusIcon" :statusText="statusText" :balance="balance" :isRefreshingBalance="isRefreshingBalance"
    :settings="settings" :isStealthMode="isStealthMode" @openSettings="openSettings" @refreshBalance="refreshBalance"
    @quit="quit" />

  <WelcomeView v-if="!hasStarted" :solveShortcut="solveShortcut"
    :toggleShortcut="shortcuts.toggle?.keyName || 'Alt+H'" :initStatus="initStatus" />

  <div v-else id="main-interface" class="main-interface" :class="{ visible: mainVisible }">
    <div class="left-panel" id="history-list">
      <div v-if="history.length === 0" class="history-item placeholder">
        <div class="history-tag">历史记录</div>
        <div class="history-preview">暂无记录</div>
      </div>
      <div v-for="(h, idx) in history" :key="idx" :class="['history-item', { active: idx === activeHistoryIndex }]"
        @click="selectHistory(idx)">
        <div class="history-tag">{{ idx === 0 ? '当前问题' : '历史问题' }}</div>
        <div class="history-preview" v-html="renderMarkdown(h.summary)"></div>
        <div class="history-time">{{ h.time }}</div>
      </div>
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
            账户</div>
        </div>
        <span class="close-btn" @click="closeSettings">&times;</span>
      </div>
      <div class="modal-body">
        <div v-show="uiState.activeTab === 'account'">
          <div class="account-card"
            style="background: rgba(30,32,36,0.92); border-radius: 16px; box-shadow: 0 4px 24px rgba(0,0,0,0.12); padding: 32px 28px; border: 1px solid rgba(255,255,255,0.04);">
            <div class="account-header" style="display: flex; align-items: center; gap: 16px; margin-bottom: 28px;">
              <span class="account-icon"
                style="font-size: 32px; background: rgba(255,255,255,0.08); border-radius: 50%; padding: 10px; color: #fff; box-shadow: 0 2px 8px rgba(0,0,0,0.18);">🔑</span>
              <div>
                <div class="account-title"
                  style="font-size: 22px; font-weight: 700; color: rgba(255,255,255,0.92); letter-spacing: 1px;">账户设置
                </div>
                <div class="account-desc" style="font-size: 14px; color: rgba(255,255,255,0.48); margin-top: 4px;">配置
                  API 相关信息与代理地址</div>
              </div>
            </div>


            <div class="form-group" style="margin-bottom: 22px;">
              <label style="font-weight: 600; color: rgba(255,255,255,0.72); font-size: 15px; margin-bottom: 8px; display: block;">Base URL</label>
              <div class="input-group" style="margin-top: 0;">
                <input type="text" v-model="tempSettings.baseURL" placeholder="https://api.openai.com/v1" style="border-radius: 10px; border: 1.5px solid rgba(255,255,255,0.12); padding: 12px; background: rgba(60,62,68,0.92); color: #fff; font-size: 15px; width: 100%; outline: none; transition: box-shadow 0.2s, border-color 0.2s; box-shadow: none;" @focus="(e)=>{e.target.style.boxShadow='0 0 0 2px #4CAF50';e.target.style.borderColor='#4CAF50'}" @blur="(e)=>{e.target.style.boxShadow='none';e.target.style.borderColor='rgba(255,255,255,0.12)' }" />
              </div>
              <p class="hint-text" style="color: rgba(255,255,255,0.38); margin-left: 0; margin-top: 8px; font-size: 13px;">如用自建代理或替换 API 域名，请填写完整地址。</p>
            </div>

            <div class="form-group" style="margin-bottom: 22px;">
              <label
                style="font-weight: 600; color: rgba(255,255,255,0.72); font-size: 15px; margin-bottom: 8px; display: block;">API
                Key</label>
              <div v-if="!uiState.isEditingKey" class="input-group" style="margin-top: 0;">
                <input type="text" :value="maskedKey" disabled
                  style="border-radius: 10px; border: 1.5px solid rgba(255,255,255,0.12); padding: 12px; background: rgba(60,62,68,0.92); color: #fff; font-size: 15px; width: 100%; outline: none; transition: box-shadow 0.2s, border-color 0.2s; box-shadow: none;" />
                <button class="btn-secondary" @click="uiState.isEditingKey = true"
                  style="margin-left: 10px; background: linear-gradient(90deg,#4CAF50,#43e97b); color: #fff; border: none; border-radius: 10px; padding: 10px 22px; font-weight: 600; cursor: pointer; box-shadow: 0 2px 8px rgba(76,175,80,0.10);">更换</button>
              </div>
              <div v-else class="input-group" style="margin-top: 0;">
                <input type="password" v-model="tempSettings.apiKey" placeholder="sk-..."
                  style="border-radius: 10px; border: 1.5px solid rgba(255,255,255,0.12); padding: 12px; background: rgba(60,62,68,0.92); color: #fff; font-size: 15px; width: 100%; outline: none; transition: box-shadow 0.2s, border-color 0.2s; box-shadow: none;"
                  @focus="(e) => { e.target.style.boxShadow = '0 0 0 2px #4CAF50'; e.target.style.borderColor = '#4CAF50' }"
                  @blur="(e) => { e.target.style.boxShadow = 'none'; e.target.style.borderColor = 'rgba(255,255,255,0.12)' }" />
                <button class="btn-secondary" @click="verifyKey"
                  :disabled="uiState.isValidatingKey || !tempSettings.apiKey"
                  style="margin-left: 10px; background: linear-gradient(90deg,#4CAF50,#43e97b); color: #fff; border: none; border-radius: 10px; padding: 10px 22px; font-weight: 600; cursor: pointer; box-shadow: 0 2px 8px rgba(76,175,80,0.10);">
                  {{ uiState.isValidatingKey ? '验证中...' : '验证' }}
                </button>
              </div>
              <p v-if="uiState.keyValidationError && uiState.isEditingKey" class="error-text"
                style="color: #ff5252; margin-top: 8px; font-size: 13px;">{{ uiState.keyValidationError }}</p>
            </div>

            <div class="account-info-row" v-if="tempSettings.apiKey && uiState.isKeyValid"
              style="margin-top: 14px; display: flex; align-items: center; gap: 16px; background: rgba(60,62,68,0.92); border-radius: 10px; padding: 12px 20px; border: 1px solid rgba(255,255,255,0.08);">
              <span class="account-info-label" style="font-weight: 600; color: #4CAF50; font-size: 15px;">API状态</span>
              <div class="balance-display" style="font-size: 16px; font-weight: bold; color: rgba(255,255,255,0.92);">
                <span v-if="tempBalance === 0">✅ 验证通过</span>
                 <span v-else-if="tempBalance === -1">🚫 无效Key</span>
                <span v-else-if="uiState.isRefreshingBalance">验证中...</span>
                <span v-else>--</span>
              </div>
            </div>
          </div>
        </div>

        <div v-show="uiState.activeTab === 'model'">
          <div class="form-group">
            <label>
              模型选择
              <span v-if="uiState.isLoadingModels" class="loading-text">加载中...</span>
            </label>

            <div class="custom-select" :class="{ open: uiState.isModelDropdownOpen, disabled: uiState.isLoadingModels }"
              @click="toggleModelDropdown" ref="modelSelectRef">
              <div class="selected-value">
                {{ tempSettings.model || '请选择模型' }}
                <span class="arrow">▼</span>
              </div>
              <div class="options-list" v-show="uiState.isModelDropdownOpen">
                <div v-for="m in uiState.availableModels" :key="m" class="option-item"
                  :class="{ selected: tempSettings.model === m }" @click.stop="selectModel(m)">
                  {{ m }}
                </div>
              </div>
            </div>
          </div>

          <div class="form-group">
            <div class="prompt-header">
              <label for="prompt-text" style="margin-bottom: 0">系统提示词 (Prompt)</label>
              <div class="prompt-tabs">
                <div class="prompt-tab" :class="{ active: uiState.promptTab === 'edit' }"
                  @click="uiState.promptTab = 'edit'">编辑</div>
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

              <!-- 先把这个配置取消掉 -->
              <!-- <div class="setting-row" v-if="tempSettings.voiceListening">
                <div class="setting-info">
                  <span class="setting-title">语音打断</span>
                  <span class="setting-desc">思考时检测到新语音，立即打断并重新发送</span>
                </div>
                <label class="switch">
                  <input type="checkbox" v-model="tempSettings.interruptThinking">
                  <span class="slider round"></span>
                </label>
              </div> -->
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
            :isParsing="resumeState.isParsing" v-model:useMarkdownResume="tempSettings.useMarkdownResume"
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
    <div v-for="(t, i) in toasts" :key="t.id || i" class="toast" :class="[t.type, { show: t.show }]">{{ t.text }}</div>
  </div>


</template>

<script setup>
import { reactive, ref, onMounted, watch, nextTick, computed } from 'vue'
import ResumeImport from './components/ResumeImport.vue'
import ScreenshotSettings from './components/ScreenshotSettings.vue'
import WelcomeView from './components/WelcomeView.vue'
import ErrorView from './components/ErrorView.vue'
import LoadingView from './components/LoadingView.vue'
// import InitLoading from './components/InitLoading.vue'
import TopBar from './components/TopBar.vue'
import { EventsOn, Quit } from '../wailsjs/runtime/runtime'
import { StopRecordingKey, SelectResume, ClearResume, RestoreFocus, RemoveFocus, ParseResume, GetInitStatus } from '../wailsjs/go/main/App'

import { useUI } from './composables/useUI'
import { useStatus } from './composables/useStatus'
import { useBalance } from './composables/useBalance'
import { useShortcuts } from './composables/useShortcuts'
import { useSettings } from './composables/useSettings'
import { useSolution } from './composables/useSolution'

const uiState = reactive({
  showSettings: false,
  activeTab: 'general',
  isEditingKey: false,
  availableModels: [],
  isLoadingModels: false,
  isModelDropdownOpen: false,
  promptTab: 'edit',
  keyValidationError: '',
  isKeyValid: false,
  isValidatingKey: false,
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
  loadSettings, verifyKey, fetchModels, saveSettings, resetTempSettings, openSettings: initSettings
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
  errorState, renderMarkdown, selectHistory, handleStreamStart, handleStreamChunk, handleSolution, setStreamBuffer
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

const modelSelectRef = ref(null)

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

function toggleModelDropdown() {
  if (uiState.isLoadingModels) return
  uiState.isModelDropdownOpen = !uiState.isModelDropdownOpen
}

function selectModel(model) {
  tempSettings.model = model
  uiState.isModelDropdownOpen = false
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

  // 1. 定义错误映射表 (配置驱动)
  // key 可以是 HTTP 状态码(number) 或者 错误关键词(string)
  const ERROR_MAP = {
    // 401 我们先定义为鉴权失败，但在逻辑里会进行二次检查
    401: { title: '鉴权失败', desc: 'API Key 无效或过期，请检查设置。', icon: '🔑' },
    403: { title: '权限不足', desc: '您的分组无权使用该模型，请联系管理员。', icon: '🚫' },
    404: { title: '模型不可用', desc: '当前模型不存在或已被移除，请切换模型。', icon: '🤖' },
    429: { title: '请求太频繁', desc: '系统繁忙，请稍作休息。', icon: '⏱️' },
    500: { title: '服务异常', desc: '上游服务暂时不可用。', icon: '🌩️' },
    502: { title: '网关错误', desc: '网关错误，上游服务无响应。', icon: '🔌' },
    504: { title: '请求超时', desc: 'AI 响应时间过长。', icon: '⏳' },

    // 定义一个专门的“余额不足”对象，方便后面复用
    'QUOTA_EXHAUSTED': { title: '额度已用尽', desc: '余额不足以支付当前长文本产生的费用。” “当前输入内容较长，您的余额不足以支付预计费用。', icon: '💸', }
  }
  //解决401问题，因为余额不足也是这个状态码
  function handle401Ambiguity(errObj) {
    const msg = (errObj.message || '').toLowerCase()
    const code = (errObj.code || '').toLowerCase()

    // 检查是否是余额
    if (
      msg.includes('quota') ||
      msg.includes('额度') ||
      msg.includes('余额') ||
      msg.includes('balance') ||
      code.includes('not enough')
    ) {
      return ERROR_MAP['QUOTA_EXHAUSTED']
    }

    // 如果不是余额问题，那就真的是 Key 错了
    return ERROR_MAP[401]
  }

  // 2. 主处理逻辑
  EventsOn('solution-error', (rawErrMsg) => {
    // A. 优先处理：用户取消 (这不是错误，是操作)
    // 这里的 canceled 通常是 Go context 返回的纯字符串，不是 JSON
    if (rawErrMsg && (rawErrMsg.includes('context canceled') || rawErrMsg.includes('canceled'))) {
      handleUserCancellation() // 把那一大坨回滚逻辑抽离出去，下面会写
      return
    }

    // B. 初始化默认错误 UI
    let errorUI = {
      title: 'AI 陷入了沉思',
      desc: '思考过程中遇到了一些未知阻碍，请稍后重试',
      icon: '🤯',
      raw: rawErrMsg
    }

    // C. 尝试解析后端传来的 JSON 错误
    try {
      // 尝试解析我们刚才在 Go 里封装的 JSON: {"statusCode": 403, "message": "..."}
      const errObj = JSON.parse(rawErrMsg)
      if (errObj.statusCode === 403 || errObj.statusCode === 401) {
        const specificError = handle401Ambiguity(errObj)
        Object.assign(errorUI, specificError)
        // if (errObj.message) errorUI.desc = ` (${errObj.message})`
        errorUI.desc += "当前余额 : " + balance.value
      }

      else if (errObj.statusCode && ERROR_MAP[errObj.statusCode]) {
        Object.assign(errorUI, ERROR_MAP[errObj.statusCode])
        // 如果后端有返回具体的 message，且不是默认的，覆盖 desc
        if (errObj.message) errorUI.desc = ` (${errObj.message})`

      }
      //用 API 返回的 code 字符串匹配 (比如 "model_not_found")
      else if (errObj.code && ERROR_MAP[errObj.code]) {
        Object.assign(errorUI, ERROR_MAP[errObj.code])
      }
      // 如果没匹配上，显示后端返回的 message
      else if (errObj.message) {
        errorUI.desc = errObj.message
      }
    } catch (e) {
      // D. 解析失败（说明是纯字符串错误，比如网络断了，或者 Go panic）
      // 降级使用关键词匹配
      const lowerMsg = rawErrMsg ? rawErrMsg.toLowerCase() : ''
      for (const key in ERROR_MAP) {
        // 忽略数字键，只匹配字符串键
        if (isNaN(key) && lowerMsg.includes(key)) {
          Object.assign(errorUI, ERROR_MAP[key])
          break
        }
      }
    }

    // E. 执行特殊副作用 (比如余额不足清空余额)
    // if (errorUI.action === 'clearBalance') {
    //   balance.value = 0.0
    //   statusText.value = '余额不足'
    //   statusIcon.value = '💸'
    // } else {
    //   statusText.value = '出错'
    //   statusIcon.value = '🔴'
    // }
    statusText.value = '出错'
    statusIcon.value = '🔴'

    // F. 更新 UI 状态
    errorState.show = true
    errorState.title = errorUI.title
    errorState.desc = errorUI.desc
    errorState.icon = errorUI.icon
    errorState.rawError = errorUI.raw
    errorState.showDetails = false

    isLoading.value = false
    isAppending.value = false
    shouldOverwriteHistory.value = true
  })

  // 3. 抽离取消逻辑 (让主函数更干净)
  function handleUserCancellation() {
    console.log('请求已由用户主动取消')

    // 如果当前正在录音（说明是因为新语音打断的），则不恢复“思考中”状态
    if (isRecording.value) {
      console.log('检测到新录音，跳过状态恢复')
      return
    }

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

  document.addEventListener('click', (e) => {
    if (modelSelectRef.value && !modelSelectRef.value.contains(e.target)) {
      uiState.isModelDropdownOpen = false
    }
  })
})
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

<style scoped>
.modal-body {
  max-height: 60vh;
  overflow-y: auto;
}

.markdown-body {
  pointer-events: auto;
  overflow-y: auto;
}

.main-interface {
  pointer-events: none;
  overflow: hidden;
}

.main-interface.visible {
  pointer-events: auto;
}

.tabs {
  display: flex;
  gap: 20px;
}

.tab {
  cursor: pointer;
  padding-bottom: 5px;
  color: #888;
  font-weight: 600;
  transition: color 0.2s;
}

.tab.active {
  color: #fff;
  border-bottom: 2px solid #4CAF50;
}

.tab:hover {
  color: #ccc;
}

.input-group {
  display: flex;
  gap: 10px;
}

.input-group input {
  flex: 1;
}

.error-text {
  color: #ff5252;
  font-size: 12px;
  margin-top: 5px;
}

.hint-text {
  font-size: 12px;
  color: #888;
  margin-top: 4px;
  margin-left: 24px;
}

.balance-display {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: bold;
  color: #4CAF50;
}

.radio-group {
  display: flex;
  gap: 15px;
  margin-top: 5px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
  font-size: 14px;
  color: #e0e0e0;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: bold;
  color: #e0e0e0;
}

.context-setting {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.context-setting:first-child {
  margin-top: 0;
  padding-top: 0;
  border-top: none;
}

.setting-row.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.setting-row.disabled .switch {
  pointer-events: none;
}

.mode-toggle-group {
  display: flex;
  gap: 15px;
  margin-top: 10px;
}

.mode-btn {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.mode-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.mode-btn.active {
  background: rgba(76, 175, 80, 0.2);
  border-color: #4CAF50;
}

.mode-icon {
  font-size: 24px;
}

.mode-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.mode-title {
  font-weight: bold;
  color: #fff;
  font-size: 14px;
}

.mode-desc {
  font-size: 12px;
  color: #aaa;
}

.mode-btn.active .mode-title {
  color: #4CAF50;
}

.append-loading {
  padding: 12px 0;
  margin-top: 15px;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #b2ebf2;
  font-size: 13px;
  font-family: 'Segoe UI', sans-serif;
  letter-spacing: 0.5px;
  animation: fadeInSimple 0.5s ease-out;
}

.append-loading .ai-icon {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: rgba(0, 188, 212, 0.08);
  border: 1px solid rgba(0, 188, 212, 0.3);
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}

.append-loading .ai-icon-inner {
  width: 8px;
  height: 8px;
  background: #00bcd4;
  border-radius: 50%;
  box-shadow: 0 0 8px rgba(0, 188, 212, 0.6);
  animation: pulse-glow 3s infinite ease-in-out;
}

.append-loading .text {
  font-weight: 500;
  background: linear-gradient(90deg, #80deea, #ffffff, #80deea);
  background-size: 200% 100%;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  animation: shimmer 4s infinite linear;
}

.wave-dots {
  display: flex;
  gap: 4px;
  align-items: center;
  height: 100%;
  margin-top: 4px;
}

.wave-dots span {
  width: 3px;
  height: 3px;
  background-color: #4dd0e1;
  border-radius: 50%;
  box-shadow: 0 0 4px rgba(77, 208, 225, 0.4);
  animation: wave 1.5s infinite ease-in-out both;
}

.wave-dots span:nth-child(1) {
  animation-delay: -0.32s;
}

.wave-dots span:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes pulse-glow {

  0%,
  100% {
    transform: scale(0.85);
    opacity: 0.6;
  }

  50% {
    transform: scale(1.15);
    opacity: 1;
    box-shadow: 0 0 15px rgba(0, 188, 212, 0.8);
  }
}

@keyframes shimmer {
  0% {
    background-position: 100% 0;
  }

  100% {
    background-position: -100% 0;
  }
}

@keyframes wave {

  0%,
  80%,
  100% {
    transform: translateY(0);
  }

  40% {
    transform: translateY(-4px);
  }
}

@keyframes fadeInSimple {
  from {
    opacity: 0;
  }

  to {
    opacity: 1;
  }
}

.setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.setting-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.setting-title {
  font-size: 14px;
  font-weight: bold;
  color: #e0e0e0;
}

.setting-desc {
  font-size: 12px;
  color: #888;
}

/* Switch Toggle Styles */
.switch {
  position: relative;
  display: inline-block;
  width: 46px;
  height: 24px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .4s;
}

input:checked+.slider {
  background-color: #4CAF50;
}

input:focus+.slider {
  box-shadow: 0 0 1px #4CAF50;
}

input:checked+.slider:before {
  transform: translateX(22px);
}

.slider.round {
  border-radius: 24px;
}

.slider.round:before {
  border-radius: 50%;
}

/* Resume Styles */
.resume-card {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.resume-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.resume-icon {
  font-size: 24px;
}

.resume-title {
  font-size: 16px;
  font-weight: bold;
  color: #e0e0e0;
}

.resume-desc {
  font-size: 13px;
  color: #aaa;
  margin-bottom: 20px;
  line-height: 1.5;
}

.resume-upload-area {
  border: 2px dashed rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 20px;
  text-align: center;
  transition: all 0.3s ease;
}

.resume-upload-area:hover {
  border-color: #4CAF50;
  background: rgba(76, 175, 80, 0.05);
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 20px;
}

.upload-icon {
  font-size: 32px;
  opacity: 0.7;
}

.upload-text {
  font-size: 14px;
  color: #4CAF50;
  font-weight: 500;
}

.file-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(0, 0, 0, 0.2);
  padding: 10px 15px;
  border-radius: 6px;
}

.file-details {
  display: flex;
  align-items: center;
  gap: 10px;
}

.file-icon {
  font-size: 20px;
}

.file-name {
  font-size: 14px;
  color: #e0e0e0;
  word-break: break-all;
}

.file-actions {
  display: flex;
  gap: 10px;
}

.btn-secondary.small,
.btn-danger.small {
  padding: 4px 10px;
  font-size: 12px;
}

.btn-danger {
  background-color: rgba(244, 67, 54, 0.1);
  color: #f44336;
  border: 1px solid rgba(244, 67, 54, 0.3);
}

.btn-danger:hover {
  background-color: rgba(244, 67, 54, 0.2);
}

.interrupted-placeholder {
  color: #888;
  font-style: italic;
  text-align: center;
  padding: 20px;
  font-size: 14px;
  animation: fadeIn 0.3s ease;
}
</style>
