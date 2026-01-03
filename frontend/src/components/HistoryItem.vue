<template>
    <div class="history-item" :class="{ active: isActive }" @click="$emit('select')">
        <div class="history-header">
            <span class="history-tag">{{ isFirst ? '当前问题' : '历史问题' }}</span>
            <div class="menu-trigger" @click.stop="toggleMenu" ref="menuTriggerRef">
                <span class="dots">⋮</span>

                <!-- 下拉菜单 -->
                <Transition name="menu-fade">
                    <div v-if="menuOpen" class="history-menu" @click.stop>
                        <div class="menu-item" @click="handleExportMarkdown">
                            <span class="menu-icon">📋</span>
                            <span>导出 Markdown</span>
                        </div>
                        <div class="menu-item" @click="handleExportImage">
                            <span class="menu-icon">🖼️</span>
                            <span>导出为图片</span>
                        </div>
                        <div class="menu-divider"></div>
                        <div class="menu-item danger" @click="handleDelete">
                            <span class="menu-icon">🗑️</span>
                            <span>删除此会话</span>
                        </div>
                    </div>
                </Transition>
            </div>
        </div>

        <div class="history-preview" v-html="previewHtml"></div>
        <div class="history-time">{{ time }}</div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
    summary: { type: String, default: '' },
    time: { type: String, default: '' },
    isActive: { type: Boolean, default: false },
    isFirst: { type: Boolean, default: false },
    previewHtml: { type: String, default: '' }
})

const emit = defineEmits(['select', 'delete', 'export-markdown', 'export-image'])

const menuOpen = ref(false)
const menuTriggerRef = ref(null)

function toggleMenu() {
    menuOpen.value = !menuOpen.value
}

function closeMenu() {
    menuOpen.value = false
}

function handleDelete() {
    emit('delete')
    closeMenu()
}

function handleExportMarkdown() {
    emit('export-markdown')
    closeMenu()
}

function handleExportImage() {
    emit('export-image')
    closeMenu()
}

// 点击外部关闭菜单
function handleClickOutside(event) {
    if (menuTriggerRef.value && !menuTriggerRef.value.contains(event.target)) {
        closeMenu()
    }
}

onMounted(() => {
    document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.history-item {
    background: rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    padding: 10px;
    cursor: pointer;
    transition: transform 0.2s, background 0.2s;
    border: 1px solid rgba(255, 255, 255, 0.08);
}

.history-item:hover {
    transform: translateY(-1px);
    background: rgba(255, 255, 255, 0.12);
}

.history-item.active {
    border-color: rgba(76, 175, 80, 0.6);
    background: rgba(76, 175, 80, 0.1);
}

.history-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 6px;
}

.history-tag {
    font-size: 9px;
    padding: 1px 5px;
    border-radius: 3px;
    background: rgba(255, 255, 255, 0.15);
    color: rgba(255, 255, 255, 0.7);
    font-weight: 500;
}

.history-item.active .history-tag {
    background: rgba(76, 175, 80, 0.3);
    color: #4CAF50;
}

.history-preview {
    font-size: 11px;
    color: #fff;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
}

.history-time {
    font-size: 9px;
    color: rgba(255, 255, 255, 0.4);
    margin-top: 6px;
    text-align: right;
}

/* 三点菜单 */
.menu-trigger {
    position: relative;
}

.dots {
    font-size: 14px;
    color: rgba(255, 255, 255, 0.4);
    cursor: pointer;
    padding: 2px 4px;
    border-radius: 4px;
    transition: all 0.2s;
    line-height: 1;
}

.dots:hover {
    background: rgba(255, 255, 255, 0.15);
    color: #fff;
}

/* 下拉菜单 */
.history-menu {
    position: absolute;
    top: calc(100% + 4px);
    right: 0;
    background: rgba(35, 35, 40, 0.98);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
    z-index: 200;
    min-width: 150px;
    padding: 4px 0;
    backdrop-filter: blur(10px);
}

.menu-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.85);
    cursor: pointer;
    transition: background 0.15s;
}

.menu-item:hover {
    background: rgba(255, 255, 255, 0.1);
}

.menu-item.danger {
    color: #ff6b6b;
}

.menu-item.danger:hover {
    background: rgba(255, 100, 100, 0.15);
}

.menu-icon {
    font-size: 14px;
}

.menu-divider {
    height: 1px;
    background: rgba(255, 255, 255, 0.1);
    margin: 4px 0;
}

/* 菜单动画 */
.menu-fade-enter-active,
.menu-fade-leave-active {
    transition: all 0.15s ease;
}

.menu-fade-enter-from,
.menu-fade-leave-to {
    opacity: 0;
    transform: translateY(-4px);
}
</style>
