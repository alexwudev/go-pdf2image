<script lang="ts" setup>
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/appStore'
import { ConvertPDF, CancelConvert } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

const { t } = useI18n()
const store = useAppStore()

onMounted(() => {
  EventsOn('convert:progress', (p: any) => store.updateProgress(p))
})

async function convert() {
  if (!store.pdfPath) return
  store.reset()
  store.isConverting = true
  store.startTimer()
  try {
    const cfg = {
      ...store.config,
      output_dir: store.outputDir,
    }
    const r = await ConvertPDF(store.pdfPath, cfg as any)
    if (r.error) {
      if (r.error === 'cancelled') {
        // User-initiated cancellation — no error display
      } else {
        store.setError(r.error)
      }
    } else {
      store.outputFiles = r.output_files || []
    }
  } catch (e: any) {
    store.setError(e.message || String(e))
  } finally {
    store.stopTimer()
    store.isConverting = false
  }
}

async function cancel() {
  try {
    await CancelConvert()
  } catch (e) {
    console.error('Cancel failed:', e)
  }
}

function copyError() {
  navigator.clipboard.writeText(store.lastError)
}
</script>

<template>
  <div class="bar">
    <button class="btn pri" :disabled="!store.pdfPath || store.isConverting" @click="convert">
      {{ store.isConverting ? t('action.converting') : t('action.convert') }}
    </button>
    <button v-if="store.isConverting" class="btn stop" @click="cancel">
      {{ t('action.stop') }}
    </button>
  </div>
  <div v-if="store.outputFiles.length > 0 && !store.isConverting && !store.lastError" class="suc-box">
    <span class="suc-msg">{{ t('message.convertSuccess', { count: store.outputFiles.length, dir: store.outputDir || '...' }) }}</span>
  </div>
  <div v-if="store.lastError" class="err-box">
    <div class="err-hdr">
      <span class="err-title">{{ t('message.convertError', { error: '' }) }}</span>
      <div class="err-actions">
        <button class="err-btn" @click="copyError">📋</button>
        <button class="err-btn" @click="store.clearError">✕</button>
      </div>
    </div>
    <pre class="err-msg">{{ store.lastError }}</pre>
  </div>
</template>

<style scoped>
.bar{display:flex;gap:8px}
.btn{padding:10px 16px;border:none;border-radius:6px;font-size:14px;font-weight:600;cursor:pointer;transition:.15s}
.btn:disabled{opacity:.5;cursor:not-allowed}
.pri{flex:1;background:#2563eb;color:#fff}
.pri:hover:not(:disabled){background:#1d4ed8}
.stop{flex:none;width:100px;background:#dc2626;color:#fff}
.stop:hover{background:#b91c1c}
.suc-box{margin-top:8px;padding:10px 14px;background:#052e16;border:1px solid #166534;border-radius:8px}
.suc-msg{font-size:13px;color:#86efac}
.err-box{margin-top:8px;background:#1c1012;border:1px solid #7f1d1d;border-radius:8px;overflow:hidden}
.err-hdr{display:flex;justify-content:space-between;align-items:center;padding:8px 12px;background:#291415;border-bottom:1px solid #7f1d1d}
.err-title{font-size:13px;font-weight:600;color:#fca5a5}
.err-actions{display:flex;gap:4px}
.err-btn{background:none;border:none;color:#fca5a5;cursor:pointer;font-size:14px;padding:2px 6px;border-radius:4px}
.err-btn:hover{background:rgba(252,165,165,.15)}
.err-msg{margin:0;padding:10px 12px;font-size:12px;color:#fecaca;font-family:'Cascadia Code','Consolas',monospace;white-space:pre-wrap;word-break:break-all;user-select:text;max-height:200px;overflow-y:auto;line-height:1.5}
</style>
