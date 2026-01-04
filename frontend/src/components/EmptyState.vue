<template>
    <div class="empty-state">
        <div class="empty-glow"></div>
        <div class="empty-content">
            <div class="empty-icon">💡</div>
            <h3 class="empty-title">准备就绪</h3>
            <p class="empty-desc">
                按 <kbd class="shortcut-key">{{ shortcut }}</kbd> 截图开始提问
            </p>
        </div>
    </div>
</template>

<script setup>
defineProps({
    shortcut: { type: String, default: 'F8' }
})
</script>

<style scoped>
.empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    position: relative;
    overflow: hidden;
}

.empty-glow {
    position: absolute;
    width: 200px;
    height: 200px;
    background: radial-gradient(circle, rgba(16, 185, 129, 0.08) 0%, transparent 70%);
    border-radius: 50%;
    filter: blur(40px);
    animation: breathe 4s ease-in-out infinite;
}

.empty-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    text-align: center;
    z-index: 1;
}

.empty-icon {
    font-size: 40px;
    margin-bottom: 4px;
    animation: float 3s ease-in-out infinite;
}

.empty-title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary, #fff);
}

.empty-desc {
    margin: 0;
    font-size: 13px;
    color: var(--text-muted, rgba(255, 255, 255, 0.5));
    display: flex;
    align-items: center;
    gap: 6px;
}

.shortcut-key {
    background: rgba(16, 185, 129, 0.15);
    color: var(--color-primary, #10b981);
    padding: 3px 8px;
    border-radius: 4px;
    font-family: var(--font-mono, monospace);
    font-size: 12px;
    font-weight: 600;
    border: 1px solid rgba(16, 185, 129, 0.25);
}

@keyframes breathe {

    0%,
    100% {
        transform: scale(1);
        opacity: 0.6;
    }

    50% {
        transform: scale(1.15);
        opacity: 0.8;
    }
}

@keyframes float {

    0%,
    100% {
        transform: translateY(0);
    }

    50% {
        transform: translateY(-6px);
    }
}
</style>
