<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/appStore'
import { OpenPDFDialog } from '../../wailsjs/go/main/App'

const { t } = useI18n()
const store = useAppStore()

async function browse() {
  try {
    const path = await OpenPDFDialog()
    if (path) store.setPdfPath(path)
  } catch (e) { console.error(e) }
}
</script>

<template>
  <div class="card">
    <h3 class="card-title">{{ t('import.title') }}</h3>
    <div class="drop-zone" :class="{ active: store.pdfPath }" @click="browse">
      <template v-if="store.pdfPath">
        <span class="file-name">{{ store.pdfPath.split(/[\\/]/).pop() }}</span>
        <span v-if="store.pageCount" class="page-info">{{ t('import.pageCount', { count: store.pageCount }) }}</span>
      </template>
      <template v-else>
        <p class="hint-text">{{ t('import.dragDrop') }}</p>
        <p class="hint-sub">{{ t('import.or') }}</p>
        <button class="browse-btn" @click.stop="browse">{{ t('import.browse') }}</button>
      </template>
    </div>
    <p class="supported">{{ t('import.supported') }}</p>
  </div>
</template>

<style scoped>
.card{background:#1f1f23;border:1px solid #27272a;border-radius:8px;padding:14px}
.card-title{font-size:14px;font-weight:600;margin:0 0 10px;color:#d4d4d8}
.drop-zone{border:2px dashed #3f3f46;border-radius:8px;padding:24px;text-align:center;cursor:pointer;transition:.15s}
.drop-zone:hover{border-color:#3b82f6;background:rgba(59,130,246,.05)}
.drop-zone.active{border-style:solid;border-color:#3f3f46;padding:12px}
.file-name{font-size:13px;color:#d4d4d8;word-break:break-all;display:block}
.page-info{font-size:12px;color:#3b82f6;margin-top:4px;display:block}
.hint-text{color:#a1a1aa;font-size:13px;margin:0}
.hint-sub{color:#52525b;font-size:12px;margin:4px 0}
.browse-btn{padding:6px 16px;border:1px solid #3b82f6;border-radius:6px;background:transparent;color:#3b82f6;font-size:13px;cursor:pointer}
.browse-btn:hover{background:rgba(59,130,246,.1)}
.supported{font-size:11px;color:#52525b;margin:6px 0 0}
</style>
