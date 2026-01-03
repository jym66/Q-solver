import modelData from '../config/model-capabilities.json'

export const PROVIDER_LOGOS = {
    google: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M22.5 12L12 22.5L1.5 12L12 1.5L22.5 12Z" fill="#1A73E8"/><path d="M12 2L15 9.5L22.5 12L15 14.5L12 22L9 14.5L1.5 12L9 9.5L12 2Z" fill="url(#gemini-gradient)"/><defs><linearGradient id="gemini-gradient" x1="1.5" y1="1.5" x2="22.5" y2="22.5" gradientUnits="userSpaceOnUse"><stop stop-color="#4E8EF7"/><stop offset="1" stop-color="#E84A5F"/></linearGradient></defs></svg>`,
    openai: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M20.9 9.7C20.6 8.5 19.9 7.4 18.9 6.6C18 5.7 16.8 5.1 15.6 4.9V4.8C15.6 3.6 15.1 2.3 14.3 1.4C13.4 0.5 12.2 0 11 0C9.8 0 8.6 0.5 7.7 1.4C6.9 2.3 6.4 3.6 6.4 4.8C5.2 5 4 5.6 3.1 6.6C2.1 7.4 1.4 8.6 1.1 9.7C0.8 10.9 1 12.2 1.6 13.3C1.6 13.3 1.7 13.4 1.7 13.5C2 14.7 2.7 15.8 3.6 16.6C4.5 17.5 5.7 18.1 6.9 18.3V18.4C6.9 19.6 7.4 20.8 8.2 21.8C9.1 22.7 10.3 23.2 11.5 23.2C12.7 23.2 13.9 22.7 14.8 21.8C15.6 20.9 16.1 19.6 16.1 18.4C17.3 18.2 18.5 17.6 19.4 16.6C20.4 15.8 21.1 14.6 21.4 13.5C21.7 12.3 21.5 11 20.9 9.7ZM11.1 2.1C11.9 2.1 12.6 2.4 13.1 3C13.7 3.5 14 4.3 14 5V6.3L13 5.7L8.7 3.2C9.4 2.5 10.2 2.1 11.1 2.1ZM13.8 11.4L11.5 12.7V15.4L10 14.5L7.7 13.2L10 11.8L12.3 10.5L13.8 11.4ZM8.8 4.4L10.3 5.3L12.5 6.6V9.3L10.2 8L7.9 6.7L8.8 4.4ZM6.6 6.6C6.6 6.6 6.6 6.6 6.6 6.6C7.2 6.1 7.9 5.8 8.7 5.8L7.3 9.4L6.2 10V8.2C6.2 7.6 6.3 7.1 6.6 6.6ZM4.8 7.9C5.3 7.4 5.9 7 6.6 6.8V10.2L4.3 8.9L3.4 8.4C3.7 7.7 4.2 7.2 4.8 7.9ZM3.3 11.9C3.1 11.3 3.1 10.7 3.3 10.1L6.7 12.1L8.2 13L6.7 13.9L4.5 15.2C4.1 14.9 3.8 14.4 3.5 13.9C3.3 13.3 3.2 12.6 3.3 11.9ZM5.6 16.5C5 16 4.6 15.4 4.3 14.7L8.6 12.2L10 13.1L10 13.2L10 15.7V17C8.2 17.3 6.7 17.3 5.6 16.5ZM10.5 19.3L9 18.4L6.8 17.1V14.5L9.1 15.8L11.4 17.1L10.5 19.3ZM12.9 17.9L11.5 17L9.3 15.7V13L11.6 14.3L13.8 15.7L12.9 17.9ZM14.4 17.3C13.7 17.9 13 18.1 12.2 18.2L13.6 14.5L14.7 13.9V15.6C14.7 16.2 14.6 16.8 14.4 17.3ZM16.2 16.1C15.6 16.7 15 17 14.3 17.3V13.8L16.6 15.1L17.5 15.6C17.2 16.3 16.8 16.8 16.2 16.1ZM17.7 12.1C17.8 12.7 17.8 13.4 17.7 14L14.2 12L12.7 11.1L14.2 10.2L16.4 8.9C16.8 9.3 17.1 9.7 17.4 10.2C17.6 10.8 17.7 11.5 17.7 12.1ZM15.4 6.7C15.9 7.2 16.3 7.8 16.6 8.5L12.3 11L10.9 10.1L12.4 9.2V6.2C13.5 6.2 14.5 6.3 15.4 6.7Z" fill="currentColor"/></svg>`,
    anthropic: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M17.5 19.5H21L12 3.5L3 19.5H6.5L12 9.5L17.5 19.5Z" fill="#D97757"/></svg>`,
    deepseek: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="12" cy="12" r="10" fill="#4D6BFE"/><path d="M7 11.5L9.5 9L17 9M7 11.5L9.5 14L17 14" stroke="white" stroke-width="2" stroke-linecap="round"/></svg>`,
    custom: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M12 15C13.6569 15 15 13.6569 15 12C15 10.3431 13.6569 9 12 9C10.3431 9 9 10.3431 9 12C9 13.6569 10.3431 15 12 15Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M19.4 15.5C19.7 15.2 20.3 15.1 20.7 15.3C21.8 15.9 22 17.4 21.1 18.3C20.2 19.2 18.7 19 18.1 17.9C17.9 17.5 18 16.9 18.4 16.6L19.4 15.5Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/><path d="M20 12C20 16.4183 16.4183 20 12 20C7.58172 20 4 16.4183 4 12C4 7.58172 7.58172 4 12 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>`,
    // Fallback/Generic icon
    default: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M21 16V8L12 3L3 8V16L12 21L21 16Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/></svg>`
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
    google: 'https://generativelanguage.googleapis.com/v1beta/openai/',
    openai: 'https://api.openai.com/v1',
    deepseek: 'https://api.deepseek.com',
    anthropic: 'https://api.anthropic.com/v1',
    custom: ''
}

/**
 * Returns the capabilities object for a given model
 * @param {string} model - The model identifier (e.g. "gemini-1.5-pro")
 * @returns {object} Capability object { image: boolean, pdf: boolean, ... }
 */
export function getModelCapabilities(model) {
    if (!model || !modelData.models[model]) {
        return { text: true, image: false, pdf: false, audio: false }
    }
    return modelData.models[model]
}

/**
 * Returns the friendly provider name for a given model
 * @param {string} model - The model identifier
 * @returns {string} Provider display name
 */
export function getProviderName(model) {
    if (!model || !modelData.models[model]) return 'Unknown Provider'

    const providerCode = modelData.models[model].provider
    return PROVIDER_NAMES[providerCode] || providerCode
}

/**
 * Returns the SVG logo string for a given model's provider
 * @param {string} model - The model identifier
 * @returns {string} SVG string
 */
export function getProviderLogo(model) {
    if (!model || !modelData.models[model]) return PROVIDER_LOGOS.default

    const providerCode = modelData.models[model].provider

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
    if (!model || !modelData.models[model]) return 'unknown'
    return modelData.models[model].provider
}

/**
 * Check if the model supports vision capabilities
 * @param {string} model - The model identifier
 * @returns {boolean} True if the model supports images
 */
export function supportsVision(model) {
    if (!model || !modelData.models[model]) return false
    return !!modelData.models[model].image
}
