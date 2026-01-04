<template>
  <div class="welcome-screen">
    <!-- 背景装饰 -->
    <div class="bg-glow"></div>
    <div class="bg-grid"></div>

    <!-- 主内容 -->
    <div class="welcome-content">
      <!-- Logo 区域 -->
      <div class="logo-section">
        <div class="logo-ring"></div>
        <div class="logo-icon">Q</div>
      </div>

      <!-- 标题区域 -->
      <div class="title-section">
        <h1 class="main-title">Q-SOLVER</h1>
        <p class="subtitle">智能答题助手 · 即刻开始</p>
      </div>

      <!-- 状态区域 -->
      <Transition name="fade-slide" mode="out-in">
        <!-- 加载状态 -->
        <div v-if="initStatus !== 'ready'" class="status-loading" key="loading">
          <div class="loading-spinner"></div>
          <span class="loading-text">
            {{ initStatus === 'loading-model' ? '正在加载模型...' : '系统初始化中...' }}
          </span>
        </div>

        <!-- 就绪后显示成功过渡 -->
        <div v-else-if="showSuccess" class="status-success" key="success">
          <div class="success-icon">✓</div>
          <span class="success-text">系统就绪</span>
        </div>

        <!-- 快捷键卡片 -->
        <div v-else class="action-cards" key="actions">
          <div class="action-card" @mouseenter="activeCard = 0" @mouseleave="activeCard = -1">
            <div class="card-glow" :class="{ active: activeCard === 0 }"></div>
            <div class="card-content">
              <div class="card-icon">📸</div>
              <div class="card-info">
                <kbd class="card-key">{{ solveShortcut }}</kbd>
                <span class="card-label">一键解题</span>
              </div>
            </div>
          </div>
          <div class="action-card" @mouseenter="activeCard = 1" @mouseleave="activeCard = -1">
            <div class="card-glow" :class="{ active: activeCard === 1 }"></div>
            <div class="card-content">
              <div class="card-icon">👁</div>
              <div class="card-info">
                <kbd class="card-key">{{ toggleShortcut }}</kbd>
                <span class="card-label">隐藏窗口</span>
              </div>
            </div>
          </div>
        </div>
      </Transition>

      <!-- 底部提示 -->
      <div v-if="initStatus === 'ready' && !showSuccess" class="bottom-hint">
        按快捷键开始使用
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  solveShortcut: { type: String, default: 'F8' },
  toggleShortcut: { type: String, default: 'Alt+H' },
  initStatus: { type: String, default: 'ready' }
})

const showSuccess = ref(false)
const activeCard = ref(-1)

watch(() => props.initStatus, (newVal, oldVal) => {
  if (newVal === 'ready' && oldVal !== 'ready') {
    showSuccess.value = true
    setTimeout(() => {
      showSuccess.value = false
    }, 1500)
  }
})
</script>

<style scoped>
/* ========================================
   Welcome Screen Container
   ======================================== */
.welcome-screen {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  z-index: 10;
}

/* ========================================
   Background Effects
   ======================================== */
.bg-glow {
  position: absolute;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.15) 0%, transparent 70%);
  border-radius: 50%;
  filter: blur(60px);
  animation: breathe 4s ease-in-out infinite;
}

.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.02) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.02) 1px, transparent 1px);
  background-size: 40px 40px;
  mask-image: radial-gradient(ellipse at center, black 30%, transparent 70%);
}

/* ========================================
   Main Content
   ======================================== */
.welcome-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 28px;
  z-index: 1;
  pointer-events: auto;
}

/* ========================================
   Logo Section
   ======================================== */
.logo-section {
  position: relative;
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid rgba(16, 185, 129, 0.3);
  animation: pulse-ring 2.5s ease-out infinite;
}

.logo-ring::before {
  content: '';
  position: absolute;
  inset: -8px;
  border-radius: 50%;
  border: 1px solid rgba(16, 185, 129, 0.15);
  animation: pulse-ring 2.5s ease-out infinite 0.5s;
}

.logo-icon {
  font-size: 36px;
  font-weight: 800;
  color: #fff;
  text-shadow: 0 0 20px rgba(16, 185, 129, 0.5);
  z-index: 1;
}

/* ========================================
   Title Section
   ======================================== */
.title-section {
  text-align: center;
}

.main-title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: 6px;
  color: #fff;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin: 0;
  letter-spacing: 2px;
}

/* ========================================
   Action Cards
   ======================================== */
.action-cards {
  display: flex;
  gap: 16px;
}

.action-card {
  position: relative;
  width: 140px;
  padding: 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: default;
  transition: all 0.3s ease;
  overflow: hidden;
}

.action-card:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(16, 185, 129, 0.3);
  transform: translateY(-2px);
}

.card-glow {
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at 50% 0%, rgba(16, 185, 129, 0.2), transparent 60%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.card-glow.active {
  opacity: 1;
}

.card-content {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  z-index: 1;
}

.card-icon {
  font-size: 24px;
}

.card-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.card-key {
  background: rgba(16, 185, 129, 0.15);
  color: var(--color-primary, #10b981);
  padding: 4px 10px;
  border-radius: 6px;
  font-family: var(--font-mono, monospace);
  font-size: 12px;
  font-weight: 600;
  border: 1px solid rgba(16, 185, 129, 0.25);
}

.card-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

/* ========================================
   Loading Status
   ======================================== */
.status-loading {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top-color: var(--color-primary, #10b981);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.loading-text {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

/* ========================================
   Success Status
   ======================================== */
.status-success {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  background: rgba(16, 185, 129, 0.1);
  border-radius: 24px;
  border: 1px solid rgba(16, 185, 129, 0.25);
}

.success-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary, #10b981);
  color: #fff;
  border-radius: 50%;
  font-size: 12px;
  font-weight: bold;
}

.success-text {
  font-size: 14px;
  color: var(--color-primary, #10b981);
  font-weight: 500;
}

/* ========================================
   Bottom Hint
   ======================================== */
.bottom-hint {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 1px;
}

/* ========================================
   Animations
   ======================================== */
@keyframes breathe {

  0%,
  100% {
    transform: scale(1);
    opacity: 0.6;
  }

  50% {
    transform: scale(1.1);
    opacity: 0.8;
  }
}

@keyframes pulse-ring {
  0% {
    transform: scale(0.95);
    opacity: 0.8;
  }

  100% {
    transform: scale(1.3);
    opacity: 0;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* ========================================
   Transitions
   ======================================== */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.4s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
