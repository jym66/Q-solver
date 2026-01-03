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
                            <stop stop-color="#4CAF50" />
                            <stop offset="1" stop-color="#2196F3" />
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

import { watch } from 'vue'

import { PROVIDER_BASE_URLS } from '../utils/modelCapabilities'

const props = defineProps({
    provider: String,
    apiKey: String,
    baseURL: String
})

const emit = defineEmits(['update:provider', 'update:apiKey', 'update:baseURL'])

watch(() => props.provider, (newProvider) => {
    if (newProvider && PROVIDER_BASE_URLS[newProvider] !== undefined) {
        if (newProvider !== 'custom') {
            emit('update:baseURL', PROVIDER_BASE_URLS[newProvider])
        }
    }
})

</script>

<style scoped>
.modern-panel {
    background: linear-gradient(165deg, rgba(32, 32, 36, 0.8) 0%, rgba(20, 20, 25, 0.95) 100%);
    border-radius: 16px;
    padding: 24px 28px;
    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.25);
    backdrop-filter: blur(12px);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.modern-panel:hover {
    box-shadow: 0 15px 50px rgba(0, 0, 0, 0.35);
    border-color: rgba(255, 255, 255, 0.12);
}

.panel-header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 28px;
    padding-bottom: 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.header-icon-box {
    width: 48px;
    height: 48px;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: inset 0 0 12px rgba(0, 0, 0, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.header-content h3 {
    font-size: 18px;
    font-weight: 700;
    color: #ffffff;
    margin: 0 0 6px 0;
    letter-spacing: 0.5px;
}

.header-content p {
    font-size: 13px;
    color: rgba(255, 255, 255, 0.45);
    margin: 0;
}

.config-form {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.form-item {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.item-label {
    font-size: 14px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.85);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.sub-label {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.35);
    font-weight: 400;
    font-family: 'Consolas', monospace;
}

.control-wrapper {
    position: relative;
    z-index: 10;
}

.input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 12px;
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.input-wrapper:focus-within {
    background: rgba(0, 0, 0, 0.3);
    border-color: #4CAF50;
    box-shadow: 0 0 0 3px rgba(76, 175, 80, 0.15);
}

.input-icon {
    padding: 0 16px;
    color: rgba(255, 255, 255, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.2s;
}

.input-wrapper:focus-within .input-icon {
    color: #4CAF50;
}

.modern-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: #fff;
    font-size: 14px;
    padding: 14px 16px 14px 0;
    font-family: 'Inter', system-ui, sans-serif;
    letter-spacing: 0.5px;
    width: 100%;
}

.modern-input::placeholder {
    color: rgba(255, 255, 255, 0.2);
}

/* Panel Footer */
.panel-footer {
    margin-top: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    opacity: 0.6;
    transition: opacity 0.2s;
}

.panel-footer:hover {
    opacity: 1;
}

.status-dot {
    width: 8px;
    height: 8px;
    background: #ff5252;
    border-radius: 50%;
    box-shadow: 0 0 8px rgba(255, 82, 82, 0.5);
    transition: all 0.3s;
}

.status-dot.active {
    background: #4CAF50;
    box-shadow: 0 0 8px rgba(76, 175, 80, 0.5);
}

.footer-text {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.5);
}
</style>
