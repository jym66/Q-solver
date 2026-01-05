<template>
    <div class="params-settings">
        <!-- Temperature -->
        <div class="setting-item">
            <div class="setting-header">
                <div class="setting-title">
                    <span class="setting-name">Temperature</span>
                    <span class="setting-badge">{{ temperature.toFixed(1) }}</span>
                </div>
                <span class="setting-desc">控制输出的随机性，较低值更精确，较高值更有创意</span>
            </div>
            <div class="setting-control">
                <span class="range-label">0</span>
                <input type="range" :value="temperature"
                    @input="$emit('update:temperature', parseFloat($event.target.value))" min="0" max="2" step="0.1" />
                <span class="range-label">2</span>
            </div>
        </div>

        <!-- Top P -->
        <div class="setting-item">
            <div class="setting-header">
                <div class="setting-title">
                    <span class="setting-name">Top P</span>
                    <span class="setting-badge">{{ topP.toFixed(2) }}</span>
                </div>
                <span class="setting-desc">核采样阈值，控制词汇选择的多样性</span>
            </div>
            <div class="setting-control">
                <span class="range-label">0</span>
                <input type="range" :value="topP" @input="$emit('update:topP', parseFloat($event.target.value))" min="0"
                    max="1" step="0.05" />
                <span class="range-label">1</span>
            </div>
        </div>

        <!-- Top K -->
        <div class="setting-item">
            <div class="setting-header">
                <div class="setting-title">
                    <span class="setting-name">Top K</span>
                    <span class="setting-badge">{{ topK }}</span>
                </div>
                <span class="setting-desc">采样时考虑的候选词数量</span>
            </div>
            <div class="setting-control">
                <span class="range-label">1</span>
                <input type="range" :value="topK" @input="$emit('update:topK', parseInt($event.target.value))" min="1"
                    max="100" step="1" />
                <span class="range-label">100</span>
            </div>
        </div>

        <!-- Max Tokens -->
        <div class="setting-item">
            <div class="setting-header">
                <div class="setting-title">
                    <span class="setting-name">最大输出长度</span>
                    <span class="setting-badge large">{{ formatNumber(maxTokens) }}</span>
                </div>
                <span class="setting-desc">模型生成的最大 Token 数量</span>
            </div>
            <div class="setting-control">
                <span class="range-label">1K</span>
                <input type="range" :value="maxTokens" @input="handleMaxTokensChange" min="1024" max="200000"
                    step="1024" />
                <span class="range-label">200K</span>
            </div>
        </div>

        <!-- Thinking Budget -->
        <div class="setting-item">
            <div class="setting-header">
                <div class="setting-title">
                    <span class="setting-name">思考预算</span>
                    <span class="setting-badge large">{{ formatNumber(thinkingBudget) }}</span>
                </div>
                <span class="setting-desc">思维链推理的最大 Token 数（不能超过最大输出长度）</span>
            </div>
            <div class="setting-control">
                <span class="range-label">1K</span>
                <input type="range" :value="thinkingBudget" @input="handleThinkingBudgetChange" min="1024"
                    :max="maxTokens" step="1024" />
                <span class="range-label">{{ formatNumber(maxTokens) }}</span>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
    temperature: { type: Number, default: 1.0 },
    topP: { type: Number, default: 0.95 },
    topK: { type: Number, default: 40 },
    maxTokens: { type: Number, default: 8192 },
    thinkingBudget: { type: Number, default: 16000 }
})

const emit = defineEmits([
    'update:temperature',
    'update:topP',
    'update:topK',
    'update:maxTokens',
    'update:thinkingBudget'
])

function formatNumber(num) {
    if (num >= 1000) {
        return (num / 1000).toFixed(0) + 'K'
    }
    return num.toString()
}

function handleMaxTokensChange(e) {
    const newMax = parseInt(e.target.value)
    emit('update:maxTokens', newMax)
    // 如果思考预算超过新的最大值，自动调整
    if (props.thinkingBudget > newMax) {
        emit('update:thinkingBudget', newMax)
    }
}

function handleThinkingBudgetChange(e) {
    const newBudget = parseInt(e.target.value)
    // 确保不超过 maxTokens
    emit('update:thinkingBudget', Math.min(newBudget, props.maxTokens))
}
</script>

<style scoped>
.params-settings {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.setting-item {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding: 16px;
    background: rgba(255, 255, 255, 0.02);
    border-radius: 8px;
    border: 1px solid rgba(255, 255, 255, 0.06);
    transition: background 0.2s;
}

.setting-item:hover {
    background: rgba(255, 255, 255, 0.04);
}

.setting-header {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.setting-title {
    display: flex;
    align-items: center;
    gap: 10px;
}

.setting-name {
    font-size: 14px;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.9);
}

.setting-badge {
    font-size: 12px;
    font-weight: 600;
    color: #10b981;
    background: rgba(16, 185, 129, 0.12);
    padding: 2px 8px;
    border-radius: 4px;
    font-family: 'SF Mono', 'Menlo', monospace;
}

.setting-badge.large {
    padding: 3px 10px;
}

.setting-desc {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.45);
    line-height: 1.4;
}

.setting-control {
    display: flex;
    align-items: center;
    gap: 12px;
    padding-top: 4px;
}

.range-label {
    font-size: 11px;
    color: rgba(255, 255, 255, 0.35);
    font-family: 'SF Mono', 'Menlo', monospace;
    min-width: 32px;
    text-align: center;
}

/* Slider Styles */
input[type="range"] {
    -webkit-appearance: none;
    appearance: none;
    flex: 1;
    height: 4px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 2px;
    outline: none;
    cursor: pointer;
}

input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 16px;
    height: 16px;
    background: #10b981;
    border-radius: 50%;
    cursor: pointer;
    border: 2px solid rgba(255, 255, 255, 0.15);
    box-shadow: 0 2px 8px rgba(16, 185, 129, 0.4);
    transition: transform 0.15s, box-shadow 0.15s;
}

input[type="range"]::-webkit-slider-thumb:hover {
    transform: scale(1.1);
    box-shadow: 0 3px 12px rgba(16, 185, 129, 0.5);
}

input[type="range"]::-webkit-slider-thumb:active {
    transform: scale(0.95);
}

input[type="range"]::-moz-range-thumb {
    width: 16px;
    height: 16px;
    background: #10b981;
    border-radius: 50%;
    cursor: pointer;
    border: 2px solid rgba(255, 255, 255, 0.15);
    box-shadow: 0 2px 8px rgba(16, 185, 129, 0.4);
}

input[type="range"]::-moz-range-track {
    height: 4px;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 2px;
}

/* Focus state */
input[type="range"]:focus {
    outline: none;
}

input[type="range"]:focus::-webkit-slider-thumb {
    box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2), 0 2px 8px rgba(16, 185, 129, 0.4);
}
</style>
