<template>
    <div class="model-select-container">
        <div class="model-select" :class="{ open: isOpen, disabled: disabled }" @click="toggle" ref="selectRef">

            <!-- 选中项 -->
            <div class="selected-item">
                <template v-if="modelValue">
                    <div class="provider-logo" v-html="getLogo(modelValue)"></div>
                    <div class="model-info">
                        <span class="model-name">{{ getName(modelValue) }}</span>
                    </div>
                </template>
                <span v-else class="placeholder">请选择提供商</span>
                <span class="arrow" :class="{ rotated: isOpen }">
                    <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                        <path d="M2.5 4.5L6 8L9.5 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"
                            stroke-linejoin="round" />
                    </svg>
                </span>
            </div>

            <!-- 下拉列表 -->
            <Transition name="dropdown">
                <div v-if="isOpen" class="dropdown-list">
                    <div v-for="provider in providers" :key="provider.value" class="dropdown-item"
                        :class="{ selected: modelValue === provider.value }"
                        @click.stop="selectProvider(provider.value)">
                        <div class="provider-logo" v-html="provider.logo"></div>
                        <div class="model-info">
                            <span class="model-name">{{ provider.label }}</span>
                        </div>
                        <span v-if="modelValue === provider.value" class="check-icon">✓</span>
                    </div>
                </div>
            </Transition>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { PROVIDER_LOGOS } from '../utils/modelCapabilities'

const props = defineProps({
    modelValue: {
        type: String,
        default: 'google'
    },
    disabled: {
        type: Boolean,
        default: false
    }
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)
const selectRef = ref(null)

// Define providers using imported logos (single source of truth)
const providers = [
    { value: 'google', label: 'Google Gemini', logo: PROVIDER_LOGOS.google },
    { value: 'openai', label: 'OpenAI', logo: PROVIDER_LOGOS.openai },
    { value: 'anthropic', label: 'Anthropic (Claude)', logo: PROVIDER_LOGOS.anthropic },
    { value: 'deepseek', label: 'DeepSeek', logo: PROVIDER_LOGOS.deepseek },
    { value: 'custom', label: '自定义 (Custom)', logo: PROVIDER_LOGOS.custom }
]

function toggle() {
    if (props.disabled) return
    isOpen.value = !isOpen.value
}

function selectProvider(value) {
    emit('update:modelValue', value)
    isOpen.value = false
}

function getLogo(value) {
    // If value not in our list, try to find in PROVIDER_LOGOS, else default
    return PROVIDER_LOGOS[value] || PROVIDER_LOGOS.default
}

function getName(value) {
    const p = providers.find(p => p.value === value)
    return p ? p.label : value
}

// 点击外部关闭
function handleClickOutside(event) {
    if (selectRef.value && !selectRef.value.contains(event.target)) {
        isOpen.value = false
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
.model-select-container {
    width: 100%;
}

.model-select {
    position: relative;
    background: rgba(30, 30, 36, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    /* Slightly reduced radius for compact feel */
    cursor: pointer;
    transition: all 0.2s ease;
}

.model-select:hover {
    border-color: rgba(255, 255, 255, 0.2);
    background: rgba(40, 40, 48, 0.95);
}

.model-select.open {
    border-color: #4CAF50;
    box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.2);
}

.model-select.disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.selected-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 12px;
    /* Reduced padding */
    height: 40px;
    /* Reduced height from 52px */
}

.provider-logo {
    width: 20px;
    /* Reduced size */
    height: 20px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
}

.provider-logo :deep(svg) {
    width: 100%;
    height: 100%;
}

.model-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

.model-name {
    font-size: 14px;
    /* Slightly adjusted font */
    font-weight: 500;
    color: #fff;
    white-space: nowrap;
}

.placeholder {
    color: rgba(255, 255, 255, 0.4);
    font-size: 13px;
}

.arrow {
    display: flex;
    align-items: center;
    justify-content: center;
    color: rgba(255, 255, 255, 0.5);
    transition: transform 0.2s ease;
    margin-left: auto;
}

.arrow.rotated {
    transform: rotate(180deg);
}

/* 下拉列表 */
.dropdown-list {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    max-height: 280px;
    overflow-y: auto;
    background: rgba(28, 28, 34, 0.98);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
    z-index: 100;
    backdrop-filter: blur(20px);
}

.dropdown-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 12px;
    /* Reduced padding */
    cursor: pointer;
    transition: background 0.15s ease;
    border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.dropdown-item:last-child {
    border-bottom: none;
}

.dropdown-item:hover {
    background: rgba(255, 255, 255, 0.08);
}

.dropdown-item.selected {
    background: rgba(76, 175, 80, 0.15);
}

.dropdown-item .provider-logo {
    width: 20px;
    height: 20px;
}

.dropdown-item .model-name {
    font-size: 13px;
}

.check-icon {
    color: #4CAF50;
    font-weight: bold;
    margin-left: auto;
    font-size: 12px;
}
</style>
