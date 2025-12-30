<template>
  <div class="top-bar-wrapper" style="--wails-draggable:drag">
    <div class="top-bar">
      <div class="control-group" :class="{ active: activeButtons.toggle }">
        <span class="key-hint">{{ shortcuts.toggle?.keyName || 'Alt+H' }}</span>
        <span class="label">隐藏/展示</span>
      </div>
      <div class="control-group" :class="{ active: activeButtons.solve }">
        <span class="key-hint">{{ shortcuts.solve?.keyName || 'Alt+~' }}</span>
        <span class="label">一键解题</span>
      </div>
      <div class="control-group" :class="{ active: activeButtons.clickthrough || isClickThrough }">
        <span class="key-hint">{{ shortcuts.clickthrough?.keyName || 'Alt+T' }}</span>
        <span class="label">鼠标穿透</span>
      </div>
      <div class="control-group" style="cursor: default;">
        <span class="key-hint">Alt+Move</span>
        <span class="label">移动/滚动</span>
      </div>
      <div class="divider"></div>
      <div class="control-group" @click="$emit('openSettings')" style="cursor: pointer;"
           @mouseenter="showSettingsTooltip" @mouseleave="hideSettingsTooltip" ref="settingsBtnRef">
        <span class="label">⚙️ 设置</span>
      </div>
      <div class="divider"></div>
      <div class="status-group" ref="statusGroupRef" @mouseenter="showTooltip" @mouseleave="hideTooltip" @click="$emit('refreshBalance')" style="cursor: pointer;">
        <span>{{ statusIcon }}</span>
        <span>{{ statusText }}</span>
      </div>
      <div class="divider"></div>
      <div class="control-group" style="cursor: pointer;" @click="$emit('quit')">
        <span class="label">❌ 退出</span>
      </div>
    </div>
  </div>

  <Teleport to="body">
    <div class="status-tooltip" v-if="showStatusTooltip" :style="tooltipStyle">
        <div class="tooltip-row">
          <span class="tooltip-label">状态:</span>
          <span class="tooltip-value">{{ statusText }}</span>
        </div>
        <div class="tooltip-row">
          <span class="tooltip-label">API状态:</span>
          <span class="tooltip-value">
{{ statusText === '已连接' ? '✅ 接口通畅' : (statusText === 'Key无效' ? '🚫 Key无效' : (statusText === '连接失败' ? '❌ 连接失败' : (isRefreshingBalance ? '验证中...' : '未配置'))) }}          </span>
        </div>
        <div class="tooltip-row">
          <span class="tooltip-label">模型:</span>
          <span class="tooltip-value">{{ settings.model }}</span>
        </div>
        <div class="tooltip-row">
          <span class="tooltip-label">隐身:</span>
          <span class="tooltip-value" :style="{ color: isStealthMode ? '#52c41a' : '#ff4d4f' }">
            {{ isStealthMode ? '已开启' : '已关闭' }}
          </span>
        </div>
    </div>
    <div class="settings-tooltip" v-if="showSettingsTip" :style="settingsTooltipStyle">
        <div class="tooltip-warning">
            ⚠️ 注意：打开设置将获取焦点<br>录屏期间请勿操作
        </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive } from 'vue'

const props = defineProps({
  shortcuts: Object,
  activeButtons: Object,
  isClickThrough: Boolean,
  statusIcon: String,
  statusText: String,
  balance: Number,
  isRefreshingBalance: Boolean,
  settings: Object,
  isStealthMode: Boolean
})

defineEmits(['openSettings', 'refreshBalance', 'quit'])

const showStatusTooltip = ref(false)
const statusGroupRef = ref(null)
const tooltipStyle = reactive({ top: '0px', left: '0px' })

const showSettingsTip = ref(false)
const settingsBtnRef = ref(null)
const settingsTooltipStyle = reactive({ top: '0px', left: '0px' })

function showTooltip() {
  if (statusGroupRef.value) {
    const rect = statusGroupRef.value.getBoundingClientRect()
    tooltipStyle.top = `${rect.bottom + 10}px`
    tooltipStyle.left = `${rect.left + rect.width / 2}px`
    showStatusTooltip.value = true
  }
}

function hideTooltip() {
  showStatusTooltip.value = false
}

function showSettingsTooltip() {
  if (settingsBtnRef.value) {
    const rect = settingsBtnRef.value.getBoundingClientRect()
    settingsTooltipStyle.top = `${rect.bottom + 10}px`
    settingsTooltipStyle.left = `${rect.left + rect.width / 2}px`
    showSettingsTip.value = true
  }
}

function hideSettingsTooltip() {
  showSettingsTip.value = false
}
</script>

<style scoped>
/* Styles from App.vue related to TopBar */
.top-bar-wrapper { pointer-events: auto; }

.status-group {
  position: relative;
}

.status-tooltip {
  position: fixed;
  transform: translateX(-50%);
  background-color: rgba(30, 30, 30, 0.95);
  border: 1px solid #444;
  border-radius: 6px;
  padding: 10px 14px;
  min-width: 160px;
  z-index: 99999;
  box-shadow: 0 4px 16px rgba(0,0,0,0.6);
  backdrop-filter: blur(8px);
  pointer-events: none;
  animation: fadeIn 0.2s ease-out;
}

.settings-tooltip {
  position: fixed;
  transform: translateX(-50%);
  background-color: rgba(40, 35, 30, 0.98);
  border: 1px solid #d4b106;
  border-radius: 6px;
  padding: 10px 14px;
  z-index: 99999;
  box-shadow: 0 4px 16px rgba(0,0,0,0.6);
  backdrop-filter: blur(8px);
  pointer-events: none;
  animation: fadeIn 0.2s ease-out;
  text-align: center;
}

.settings-tooltip::before {
  content: '';
  position: absolute;
  top: -6px;
  left: 50%;
  transform: translateX(-50%);
  border-width: 0 6px 6px 6px;
  border-style: solid;
  border-color: transparent transparent #d4b106 transparent;
}

.tooltip-warning {
  color: #fadb14;
  font-size: 12px;
  line-height: 1.5;
  font-weight: 600;
}

.status-tooltip::before {
  content: '';
  position: absolute;
  top: -6px;
  left: 50%;
  transform: translateX(-50%);
  border-width: 0 6px 6px 6px;
  border-style: solid;
  border-color: transparent transparent rgba(30, 30, 30, 0.95) transparent;
}

.tooltip-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
  font-size: 12px;
  white-space: nowrap;
}

.tooltip-row:last-child {
  margin-bottom: 0;
}

.tooltip-label {
  color: #aaa;
  margin-right: 16px;
}

.tooltip-value {
  color: #eee;
  font-weight: 600;
  font-family: monospace;
}


@keyframes fadeIn {
  from { opacity: 0; transform: translate(-50%, -5px); }
  to { opacity: 1; transform: translate(-50%, 0); }
}
</style>
