import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { GetPagePreview, GetPDFInfo } from '../../wailsjs/go/main/App'

export interface ConvertConfig {
  dpi: number
  quality: number
  format: string
  pages: string
  output_dir: string
  workers: number
}

export interface ProgressInfo {
  current: number
  total: number
  page: number
  percent: number
}

export const useAppStore = defineStore('app', () => {
  const pdfPath = ref('')
  const pageCount = ref(0)
  const currentPage = ref(0)
  const previewUrl = ref('')
  const isConverting = ref(false)
  const outputDir = ref('')
  const outputFiles = ref<string[]>([])
  const lastError = ref('')
  const progress = reactive<ProgressInfo>({ current: 0, total: 0, page: 0, percent: 0 })

  const config = reactive<ConvertConfig>({
    dpi: 300,
    quality: 90,
    format: 'jpg',
    pages: 'all',
    output_dir: '',
    workers: 4,
  })

  async function setPdfPath(path: string) {
    pdfPath.value = path
    pageCount.value = 0
    currentPage.value = 0
    previewUrl.value = ''
    outputFiles.value = []
    lastError.value = ''

    if (path) {
      try {
        const info = await GetPDFInfo(path)
        if (info.error) {
          lastError.value = info.error
          return
        }
        pageCount.value = info.page_count
        currentPage.value = 0
        await loadPreview(0)
      } catch (e: any) {
        lastError.value = e.message || String(e)
      }
    }
  }

  async function loadPreview(page: number) {
    if (!pdfPath.value || page < 0 || page >= pageCount.value) return
    try {
      previewUrl.value = await GetPagePreview(pdfPath.value, page)
      currentPage.value = page
    } catch (e: any) {
      console.error('Failed to load preview:', e)
    }
  }

  function updateProgress(p: ProgressInfo) {
    progress.current = p.current
    progress.total = p.total
    progress.page = p.page
    progress.percent = p.percent
  }

  function setError(msg: string) { lastError.value = msg }
  function clearError() { lastError.value = '' }

  function reset() {
    isConverting.value = false
    progress.current = 0
    progress.total = 0
    progress.page = 0
    progress.percent = 0
    lastError.value = ''
  }

  return {
    pdfPath, pageCount, currentPage, previewUrl,
    isConverting, outputDir, outputFiles, lastError, progress, config,
    setPdfPath, loadPreview, updateProgress, reset, setError, clearError,
  }
})
