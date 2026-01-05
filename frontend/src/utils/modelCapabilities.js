import modelData from '../config/model-capabilities.json'

// 提供商Logo - 简洁清晰的图标设计
export const PROVIDER_LOGOS = {
    // Google Gemini - 四角星
    google: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#4285F4"/>
        <path d="M12 4L14 10L20 12L14 14L12 20L10 14L4 12L10 10L12 4Z" fill="white"/>
    </svg>`,

    // OpenAI - 简化的花形logo
    openai: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#10a37f"/>
        <circle cx="12" cy="8" r="2.5" fill="white"/>
        <circle cx="8" cy="14" r="2.5" fill="white"/>
        <circle cx="16" cy="14" r="2.5" fill="white"/>
        <circle cx="12" cy="12" r="1.5" fill="white"/>
    </svg>`,

    // Anthropic Claude - A字形
    anthropic: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#D4A27F"/>
        <path d="M12 5L6 19H9L10 16H14L15 19H18L12 5ZM11 13L12 9L13 13H11Z" fill="white"/>
    </svg>`,

    // DeepSeek - 眼睛图标
    deepseek: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#4D6BFE"/>
        <ellipse cx="12" cy="12" rx="7" ry="4" stroke="white" stroke-width="2" fill="none"/>
        <circle cx="12" cy="12" r="2" fill="white"/>
    </svg>`,

    // Alibaba Qwen - Q字
    alibaba: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#FF6A00"/>
        <circle cx="12" cy="11" r="5" stroke="white" stroke-width="2" fill="none"/>
        <line x1="15" y1="14" x2="18" y2="18" stroke="white" stroke-width="2" stroke-linecap="round"/>
    </svg>`,

    // Zhipu GLM - Z字
    zhipu: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#2563EB"/>
        <path d="M7 7H17L7 17H17" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
    </svg>`,

    // Moonshot Kimi - 月亮
    moonshot: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#1E3A8A"/>
        <path d="M12 6C8.7 6 6 8.7 6 12S8.7 18 12 18C12 18 10 15 10 12S12 6 12 6Z" fill="white"/>
        <circle cx="15" cy="9" r="1" fill="white"/>
    </svg>`,

    // Mistral AI - M字
    mistral: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#F97316"/>
        <path d="M6 17V7L12 13L18 7V17" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
    </svg>`,

    // xAI Grok - X
    xai: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#000"/>
        <path d="M7 7L17 17M17 7L7 17" stroke="white" stroke-width="2.5" stroke-linecap="round"/>
    </svg>`,

    // Meta Llama - 羊驼头
    meta: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#0668E1"/>
        <ellipse cx="12" cy="14" rx="5" ry="4" fill="white"/>
        <circle cx="10" cy="13" r="1" fill="#0668E1"/>
        <circle cx="14" cy="13" r="1" fill="#0668E1"/>
        <ellipse cx="9" cy="8" rx="2" ry="3" fill="white"/>
        <ellipse cx="15" cy="8" rx="2" ry="3" fill="white"/>
    </svg>`,

    // 01.AI Yi - Yi文字
    '01ai': `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#8B5CF6"/>
        <path d="M8 7L12 12L16 7" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
        <line x1="12" y1="12" x2="12" y2="18" stroke="white" stroke-width="2.5" stroke-linecap="round"/>
    </svg>`,

    // 自定义 - 齿轮
    custom: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#6B7280"/>
        <circle cx="12" cy="12" r="3" stroke="white" stroke-width="2" fill="none"/>
        <path d="M12 5V7M12 17V19M5 12H7M17 12H19M7.05 7.05L8.46 8.46M15.54 15.54L16.95 16.95M7.05 16.95L8.46 15.54M15.54 8.46L16.95 7.05" stroke="white" stroke-width="2" stroke-linecap="round"/>
    </svg>`,

    // 默认 - 方块
    default: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
        <rect width="24" height="24" rx="6" fill="#4B5563"/>
        <rect x="7" y="7" width="10" height="10" rx="2" stroke="white" stroke-width="2" fill="none"/>
        <circle cx="12" cy="12" r="2" fill="white"/>
    </svg>`
}

// Maps provider codes to display names
const PROVIDER_NAMES = {
    google: 'Google Gemini',
    openai: 'OpenAI',
    anthropic: 'Anthropic',
    deepseek: 'DeepSeek',
    alibaba: 'Alibaba Cloud (Qwen)',
    zhipu: 'Zhipu AI (GLM)',
    moonshot: 'Moonshot AI (Kimi)',
    mistral: 'Mistral AI',
    xai: 'xAI (Grok)',
    meta: 'Meta (Llama)',
    '01ai': '01.AI (Yi)'
}

export const PROVIDER_BASE_URLS = {
    google: 'https://generativelanguage.googleapis.com',
    openai: 'https://api.openai.com/v1',
    deepseek: 'https://api.deepseek.com',
    anthropic: 'https://api.anthropic.com',
    custom: ''
}

/**
 * Returns the capabilities object for a given model
 * @param {string} model - The model identifier (e.g. "gemini-1.5-pro")
 * @returns {object} Capability object { image: boolean, pdf: boolean, ... }
 */
export function getModelCapabilities(model) {
    if (!model) return { text: true, image: false, pdf: false, audio: false }

    // 先尝试精确匹配
    if (modelData.models[model]) {
        return modelData.models[model]
    }

    // 关键字匹配
    const lowerModel = model.toLowerCase()
    for (const [key, value] of Object.entries(modelData.models)) {
        if (lowerModel.includes(key.toLowerCase()) || key.toLowerCase().includes(lowerModel)) {
            return value
        }
    }

    return { text: true, image: false, pdf: false, audio: false }
}

/**
 * Returns the friendly provider name for a given model
 * @param {string} model - The model identifier
 * @returns {string} Provider display name
 */
export function getProviderName(model) {
    if (!model) return 'Unknown Provider'

    // 获取 provider code
    const providerCode = getProviderCodeFromModel(model)
    return PROVIDER_NAMES[providerCode] || providerCode
}

/**
 * 根据模型名称推断 provider code
 */
function getProviderCodeFromModel(model) {
    if (!model) return 'unknown'

    const lowerModel = model.toLowerCase()

    // 关键字匹配规则
    if (lowerModel.includes('gemini')) return 'google'
    if (lowerModel.includes('gpt') || lowerModel.includes('o1') || lowerModel.includes('o3') || lowerModel.includes('o4')) return 'openai'
    if (lowerModel.includes('claude')) return 'anthropic'
    if (lowerModel.includes('deepseek')) return 'deepseek'
    if (lowerModel.includes('qwen')) return 'alibaba'
    if (lowerModel.includes('glm') || lowerModel.includes('chatglm')) return 'zhipu'
    if (lowerModel.includes('moonshot') || lowerModel.includes('kimi')) return 'moonshot'
    if (lowerModel.includes('mistral') || lowerModel.includes('mixtral')) return 'mistral'
    if (lowerModel.includes('grok')) return 'xai'
    if (lowerModel.includes('llama')) return 'meta'
    if (lowerModel.includes('yi-')) return '01ai'

    // 尝试从配置中查找
    if (modelData.models[model]) {
        return modelData.models[model].provider
    }

    // 模糊匹配配置
    for (const [key, value] of Object.entries(modelData.models)) {
        if (lowerModel.includes(key.toLowerCase()) || key.toLowerCase().includes(lowerModel)) {
            return value.provider
        }
    }

    return 'unknown'
}

/**
 * Returns the SVG logo string for a given model's provider
 * @param {string} model - The model identifier
 * @returns {string} SVG string
 */
export function getProviderLogo(model) {
    const providerCode = getProviderCodeFromModel(model)

    // Check if we have a specific logo for this provider
    if (PROVIDER_LOGOS[providerCode]) {
        return PROVIDER_LOGOS[providerCode]
    }

    // Fallback to default
    return PROVIDER_LOGOS.default
}

/**
 * Get just the provider code for a model
 */
export function getProviderCode(model) {
    return getProviderCodeFromModel(model)
}

/**
 * Check if the model supports vision capabilities
 * @param {string} model - The model identifier
 * @returns {boolean} True if the model supports images
 */
export function supportsVision(model) {
    const capabilities = getModelCapabilities(model)
    return !!capabilities.image
}

/**
 * Check if the model supports PDF capabilities
 * @param {string} model - The model identifier
 * @returns {boolean} True if the model supports PDF
 */
export function supportsPDF(model) {
    const capabilities = getModelCapabilities(model)
    return !!capabilities.pdf
}
