<template>
    <div class="account-panel modern-panel">
        <!-- Header Area -->
        <div class="panel-header">
            <div class="header-icon-box">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="url(#header-gradient)" stroke-width="2"
                        stroke-linecap="round" stroke-linejoin="round" />
                    <path d="M2 17L12 22L22 17" stroke="url(#header-gradient)" stroke-width="2" stroke-linecap="round"
                        stroke-linejoin="round" />
                    <path d="M2 12L12 17L22 12" stroke="url(#header-gradient)" stroke-width="2" stroke-linecap="round"
                        stroke-linejoin="round" />
                    <defs>
                        <linearGradient id="header-gradient" x1="2" y1="2" x2="22" y2="22"
                            gradientUnits="userSpaceOnUse">
                            <stop stop-color="#10b981" />
                            <stop offset="1" stop-color="#0ea5e9" />
                        </linearGradient>
                    </defs>
                </svg>
            </div>
            <div class="header-content">
                <h3>模型服务商</h3>
                <p>配置 AI 大脑连接 (Provider Configuration)</p>
            </div>
        </div>

        <div class="config-form">
            <!-- Provider Selection -->
            <div class="form-item">
                <label class="item-label">选择服务商 <span class="sub-label">Model Provider</span></label>
                <div class="control-wrapper provider-wrapper">
                    <ProviderDropdown :modelValue="provider" @update:modelValue="$emit('update:provider', $event)" />
                </div>
            </div>

            <!-- API Key -->
            <div class="form-item">
                <label class="item-label">API 密钥 <span class="sub-label">Secret Key</span></label>
                <div class="input-wrapper">
                    <span class="input-icon">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path
                                d="M21 2L12 11M21 2C21.5523 2 22 2.44772 22 3V6.5C22 6.63261 21.9473 6.75979 21.8536 6.85355L19.4393 9.26777C18.596 10.1111 17.2292 10.1111 16.3858 9.26777L14.7322 7.61421C13.8889 6.77088 13.8889 5.40404 14.7322 4.56071L17.1464 2.14645C17.2402 2.05268 17.3674 2 17.5 2H21ZM10 14C11.6569 14 13 12.6569 13 11C13 9.34315 11.6569 8 10 8C8.34315 8 7 9.34315 7 11C7 12.6569 8.34315 14 10 14ZM10 14C6.68629 14 4 16.6863 4 20C4 21.1046 4.89543 22 6 22H10"
                                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                        </svg>
                    </span>
                    <input type="password" :value="apiKey" @input="$emit('update:apiKey', $event.target.value)"
                        placeholder="sk-..." class="modern-input" />
                </div>
            </div>

            <!-- Base URL -->
            <div class="form-item" v-if="provider === 'custom'">
                <label class="item-label">代理地址 <span class="sub-label">Base URL</span></label>
                <div class="input-wrapper">
                    <span class="input-icon">
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" />
                            <line x1="2" y1="12" x2="22" y2="12" stroke="currentColor" stroke-width="2" />
                            <path
                                d="M12 2C14.5013 4.73835 15.9228 8.29203 16 12C15.9228 15.708 14.5013 19.2616 12 22C9.49872 19.2616 8.07725 15.708 8 12C8.07725 8.29203 9.49872 4.73835 12 2Z"
                                stroke="currentColor" stroke-width="2" />
                        </svg>
                    </span>
                    <input type="text" :value="baseURL" @input="$emit('update:baseURL', $event.target.value)"
                        placeholder="https://api.openai.com/v1" class="modern-input" />
                </div>
            </div>
        </div>

        <div class="panel-footer">
            <span class="status-dot" :class="{ active: apiKey }"></span>
            <span class="footer-text">
                {{ apiKey ? 'API Key 已配置' : '请填写 API Key 以启用服务' }}
            </span>
        </div>
    </div>
</template>

<script setup>
import ProviderDropdown from './ProviderDropdown.vue'

import { watch, onMounted } from 'vue'

import { PROVIDER_BASE_URLS } from '../utils/modelCapabilities'

const props = defineProps({
    provider: String,
    apiKey: String,
    baseURL: String
})

const emit = defineEmits(['update:provider', 'update:apiKey', 'update:baseURL'])

// 当 provider 变化时，自动设置对应的 baseURL
watch(() => props.provider, (newProvider) => {
    if (newProvider && newProvider !== 'custom' && PROVIDER_BASE_URLS[newProvider] !== undefined) {
        emit('update:baseURL', PROVIDER_BASE_URLS[newProvider])
    }
}, { immediate: true })  // immediate: true 确保首次渲染时也会执行

</script>

<style scoped>
.modern-panel {
    background: linear-gradient(165deg, var(--bg-elevated) 0%, var(--bg-base) 100%);
    border-radius: var(--radius-xl);
    padding: var(--space-6);
    border: 1px solid var(--border-subtle);
    box-shadow: var(--shadow-lg);
    backdrop-filter: blur(12px);
    transition: all var(--transition-base);
}

.modern-panel:hover {
    box-shadow: 0 15px 50px rgba(0, 0, 0, 0.35);
    border-color: var(--border-default);
}

.panel-header {
    display: flex;
    align-items: center;
    gap: var(--space-4);
    margin-bottom: var(--space-6);
    padding-bottom: var(--space-5);
    border-bottom: 1px solid var(--border-subtle);
}

.header-icon-box {
    width: 48px;
    height: 48px;
    background: var(--bg-card);
    border-radius: var(--radius-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: inset 0 0 12px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-subtle);
}

.header-content h3 {
    font-size: var(--text-lg);
    font-weight: 700;
    color: var(--text-primary);
    margin: 0 0 6px 0;
    letter-spacing: 0.5px;
}

.header-content p {
    font-size: var(--text-sm);
    color: var(--text-tertiary);
    margin: 0;
}

.config-form {
    display: flex;
    flex-direction: column;
    gap: var(--space-5);
}

.form-item {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
}

.item-label {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--text-secondary);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.sub-label {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
    font-weight: 400;
    font-family: var(--font-mono);
}

.control-wrapper {
    position: relative;
    z-index: 10;
}

.input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
    background: var(--bg-inset);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: all var(--transition-base);
}

.input-wrapper:focus-within {
    background: var(--bg-card);
    border-color: var(--color-primary);
    box-shadow: var(--shadow-focus);
}

.input-icon {
    padding: 0 var(--space-4);
    color: var(--text-tertiary);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color var(--transition-fast);
}

.input-wrapper:focus-within .input-icon {
    color: var(--color-primary);
}

.modern-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: var(--text-primary);
    font-size: var(--text-sm);
    padding: 14px var(--space-4) 14px 0;
    font-family: var(--font-family);
    letter-spacing: 0.5px;
    width: 100%;
}

.modern-input::placeholder {
    color: var(--text-disabled);
}

/* Panel Footer */
.panel-footer {
    margin-top: var(--space-8);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-2);
    opacity: 0.6;
    transition: opacity var(--transition-fast);
}

.panel-footer:hover {
    opacity: 1;
}

.status-dot {
    width: 8px;
    height: 8px;
    background: var(--color-error);
    border-radius: var(--radius-full);
    box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);
    transition: all var(--transition-base);
}

.status-dot.active {
    background: var(--color-primary);
    box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
}

.footer-text {
    font-size: var(--text-xs);
    color: var(--text-tertiary);
}
</style>
