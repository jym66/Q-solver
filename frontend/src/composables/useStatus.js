import { ref } from 'vue'

export function useStatus(settings) {
  const statusText = ref('就绪')
  const statusIcon = ref('📝')

  function resetStatus() {
    if (!settings.apiKey) {
      statusText.value = '未配置Key'
      statusIcon.value = '⚠️'
      return
    }

    statusText.value = '就绪'
    statusIcon.value = '📝'
  }

  return {
    statusText,
    statusIcon,
    resetStatus
  }
}
