import { ref, reactive, nextTick } from 'vue'
import { marked } from 'marked'

export function useSolution(settings) {
  const renderedContent = ref('')
  const history = ref([])
  const activeHistoryIndex = ref(0)
  const isLoading = ref(false)
  const isAppending = ref(false)
  const shouldOverwriteHistory = ref(false)
  let streamBuffer = ''

  const errorState = reactive({
    show: false,
    icon: '⚠️',
    title: '出错了',
    desc: '发生了一个未知错误',
    rawError: '',
    showDetails: false
  })

  function renderMarkdown(md) {
    if (!md) return ''
    return marked.parse(md)
  }

  function selectHistory(idx) {
    const item = history.value[idx]
    if (item) {
      renderedContent.value = renderMarkdown(item.full)
      activeHistoryIndex.value = idx
    }
  }

  function handleStreamStart() {
    if (settings.keepContext && history.value.length > 0 && !shouldOverwriteHistory.value) {
      const separator = '\n\n---\n\n'
      streamBuffer = history.value[0].full + separator
      activeHistoryIndex.value = 0
    } else {
      streamBuffer = ''
      renderedContent.value = ''

      if (shouldOverwriteHistory.value && history.value.length > 0) {
        history.value[0] = {
          summary: '正在思考...',
          full: '',
          time: new Date().toLocaleTimeString()
        }
        shouldOverwriteHistory.value = false
      } else {
        history.value.unshift({
          summary: '正在思考...',
          full: '',
          time: new Date().toLocaleTimeString()
        })
      }
      activeHistoryIndex.value = 0
    }
  }

  function handleStreamChunk(token) {
    if (isLoading.value) isLoading.value = false
    if (isAppending.value) isAppending.value = false

    streamBuffer += token
    renderedContent.value = renderMarkdown(streamBuffer)

    if (history.value.length > 0) {
      history.value[0].full = streamBuffer
      history.value[0].summary = streamBuffer.substring(0, 30).replace(/\n/g, ' ') + '...'
    }

    nextTick(() => {
      const contentDiv = document.getElementById('content')
      if (contentDiv) {
        contentDiv.scrollTop = contentDiv.scrollHeight
      }
    })
  }

  function handleSolution(data) {
    isLoading.value = false

    if (settings.keepContext && history.value.length > 0) {
      renderedContent.value = renderMarkdown(history.value[0].full)
    } else {
      renderedContent.value = renderMarkdown(data)
      if (history.value.length > 0) {
        history.value[0].full = data
        history.value[0].summary = data.substring(0, 30).replace(/\n/g, ' ') + '...'
      }
    }
  }

  function setStreamBuffer(val) {
    streamBuffer = val
  }

  /**
   * 删除指定索引的历史记录
   */
  function deleteHistory(index) {
    if (index < 0 || index >= history.value.length) return

    history.value.splice(index, 1)

    // 调整活动索引
    if (history.value.length === 0) {
      renderedContent.value = ''
      activeHistoryIndex.value = 0
    } else if (index <= activeHistoryIndex.value) {
      activeHistoryIndex.value = Math.max(0, activeHistoryIndex.value - 1)
      selectHistory(activeHistoryIndex.value)
    }
  }

  /**
   * 导出为 Markdown 文件
   */
  function exportMarkdown(index) {
    const item = history.value[index]
    if (!item) return

    const content = item.full || item.summary || ''
    const blob = new Blob([content], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `answer-${new Date().toISOString().slice(0, 10)}.md`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  /**
   * 导出为图片
   */
  async function exportImage(index) {
    // 先选中该历史记录以显示其内容
    selectHistory(index)

    // 等待渲染完成
    await nextTick()

    const contentEl = document.getElementById('content')
    if (!contentEl) return

    try {
      // 动态导入 html2canvas
      const { default: html2canvas } = await import('html2canvas')

      const canvas = await html2canvas(contentEl, {
        backgroundColor: '#1e1e1e',
        scale: 2,
        useCORS: true
      })

      const url = canvas.toDataURL('image/png')
      const a = document.createElement('a')
      a.href = url
      a.download = `answer-${new Date().toISOString().slice(0, 10)}.png`
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
    } catch (e) {
      console.error('导出图片失败:', e)
      alert('导出图片需要安装 html2canvas，请运行: npm install html2canvas')
    }
  }

  return {
    renderedContent,
    history,
    activeHistoryIndex,
    isLoading,
    isAppending,
    shouldOverwriteHistory,
    errorState,
    renderMarkdown,
    selectHistory,
    handleStreamStart,
    handleStreamChunk,
    handleSolution,
    setStreamBuffer,
    deleteHistory,
    exportMarkdown,
    exportImage
  }
}
