<template>
  <div class="live-view">
    <!-- 状态栏 -->
    <div class="status-bar">
      <div class="status-indicator" :class="statusClass">
        <span class="status-dot"></span>
        <span class="status-text">{{ statusText }}</span>
      </div>
    </div>

    <!-- 面试官转录区 -->
    <div class="transcript-section interviewer">
      <div class="section-header">
        <span class="icon">👤</span>
        <span class="label">面试官</span>
      </div>
      <div class="content" ref="interviewerContent">
        <p v-if="transcript">{{ transcript }}</p>
        <p v-else class="placeholder">等待语音输入...</p>
      </div>
    </div>

    <!-- AI 回复区 -->
    <div class="transcript-section ai">
      <div class="section-header">
        <span class="icon">🤖</span>
        <span class="label">AI 助手</span>
      </div>
      <div class="content ai-content" ref="aiContent" v-html="renderedAiText"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { marked } from 'marked'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

const status = ref('disconnected') // disconnected, connecting, connected, error
const transcript = ref('')
const aiText = ref('')
const errorMsg = ref('')

const interviewerContent = ref(null)
const aiContent = ref(null)

const statusClass = computed(() => ({
  'status-disconnected': status.value === 'disconnected',
  'status-connecting': status.value === 'connecting',
  'status-connected': status.value === 'connected',
  'status-error': status.value === 'error',
}))

const statusText = computed(() => {
  switch (status.value) {
    case 'disconnected': return '未连接'
    case 'connecting': return '连接中...'
    case 'connected': return '已连接'
    case 'error': return `连接失败: ${errorMsg.value}`
    default: return '未知状态'
  }
})

const renderedAiText = computed(() => {
  if (!aiText.value) return '<p class="placeholder">等待 AI 回复...</p>'
  return marked.parse(aiText.value)
})

// 自动滚动到底部
function scrollToBottom(el) {
  if (el) {
    nextTick(() => {
      el.scrollTop = el.scrollHeight
    })
  }
}

// 事件监听
function onLiveStatus(newStatus) {
  status.value = newStatus
}

function onLiveTranscript(text) {
  transcript.value += text
  scrollToBottom(interviewerContent.value)
}

function onLiveAiText(text) {
  aiText.value += text
  scrollToBottom(aiContent.value)
}

function onLiveError(err) {
  status.value = 'error'
  errorMsg.value = err
}

function onLiveDone() {
  // 一轮对话完成，可以重置 transcript
  transcript.value = ''
}

onMounted(() => {
  EventsOn('live:status', onLiveStatus)
  EventsOn('live:transcript', onLiveTranscript)
  EventsOn('live:ai-text', onLiveAiText)
  EventsOn('live:error', onLiveError)
  EventsOn('live:done', onLiveDone)
})

onUnmounted(() => {
  EventsOff('live:status')
  EventsOff('live:transcript')
  EventsOff('live:ai-text')
  EventsOff('live:error')
  EventsOff('live:done')
})
</script>

<style scoped>
.live-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 12px;
  gap: 12px;
}

/* 状态栏 */
.status-bar {
  display: flex;
  justify-content: flex-start;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-disconnected {
  background: rgba(255, 100, 100, 0.15);
  color: #ff6b6b;
}
.status-disconnected .status-dot {
  background: #ff6b6b;
}

.status-connecting {
  background: rgba(255, 193, 7, 0.15);
  color: #ffc107;
}
.status-connecting .status-dot {
  background: #ffc107;
  animation: pulse 1s infinite;
}

.status-connected {
  background: rgba(76, 175, 80, 0.15);
  color: #4caf50;
}
.status-connected .status-dot {
  background: #4caf50;
}

.status-error {
  background: rgba(255, 100, 100, 0.15);
  color: #ff6b6b;
}
.status-error .status-dot {
  background: #ff6b6b;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

/* 内容区 */
.transcript-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: rgba(30, 30, 35, 0.6);
  border-radius: 10px;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.05);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  font-size: 13px;
  font-weight: 500;
}

.section-header .icon {
  font-size: 14px;
}

.content {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
  font-size: 14px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.9);
}

.content .placeholder {
  color: rgba(255, 255, 255, 0.4);
  font-style: italic;
}

/* AI 内容特殊样式 */
.ai-content :deep(p) {
  margin: 0 0 8px 0;
}

.ai-content :deep(code) {
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Consolas', monospace;
}

.ai-content :deep(pre) {
  background: rgba(0, 0, 0, 0.3);
  padding: 12px;
  border-radius: 6px;
  overflow-x: auto;
}
</style>
