<script lang="ts" setup>
import { watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from './stores/appStore'
import { WindowSetTitle } from '../wailsjs/runtime/runtime'
import PdfImport from './components/PdfImport.vue'
import SettingsPanel from './components/SettingsPanel.vue'
import ActionBar from './components/ActionBar.vue'
import PreviewPanel from './components/PreviewPanel.vue'
import ConvertProgress from './components/ConvertProgress.vue'

const { t, locale } = useI18n()
const store = useAppStore()

const languages = [
  { code: 'zh-TW', label: '繁體中文' },
  { code: 'en', label: 'English' },
]

function switchLang(e: Event) {
  const val = (e.target as HTMLSelectElement).value
  locale.value = val
  localStorage.setItem('pdf2image-lang', val)
}

// Restore saved language
const saved = localStorage.getItem('pdf2image-lang')
if (saved && languages.some(l => l.code === saved)) {
  locale.value = saved
}

// Update window title with conversion progress
watch(() => store.progress.percent, (pct) => {
  if (store.isConverting && pct > 0 && pct < 100) {
    WindowSetTitle(`${Math.round(pct)}% - PDF2Image`)
  }
})
watch(() => store.isConverting, (converting) => {
  if (!converting) {
    WindowSetTitle('PDF2Image')
  }
})
</script>

<template>
  <div class="app">
    <header class="hdr">
      <div v-if="store.isConverting || store.progress.percent > 0"
           class="hdr-fill" :class="{ done: store.progress.percent >= 100 }"
           :style="{ width: store.progress.percent + '%' }"/>
      <div class="hdr-content">
        <div class="hdr-l">
          <h1 class="brand">{{ t('app.title') }}</h1>
          <span class="sub">{{ store.isConverting ? Math.round(store.progress.percent) + '%' : t('app.subtitle') }}</span>
        </div>
        <select class="lang-sel" :value="locale" @change="switchLang">
          <option v-for="lang in languages" :key="lang.code" :value="lang.code">{{ lang.label }}</option>
        </select>
      </div>
    </header>
    <main class="main">
      <aside class="left">
        <PdfImport/>
        <SettingsPanel/>
        <ActionBar/>
      </aside>
      <section class="right">
        <ConvertProgress/>
        <PreviewPanel/>
      </section>
    </main>
  </div>
</template>

<style>
.app{display:flex;flex-direction:column;height:100vh;background:#18181b;color:#e4e4e7;font-family:'Segoe UI',system-ui,sans-serif}
.hdr{position:relative;overflow:hidden;background:#09090b;border-bottom:1px solid #27272a}
.hdr-fill{position:absolute;top:0;left:0;height:100%;background:rgba(37,99,235,.25);transition:width .3s;pointer-events:none}
.hdr-fill.done{background:rgba(34,197,94,.2)}
.hdr-content{position:relative;display:flex;align-items:center;justify-content:space-between;padding:12px 24px}
.hdr-l{display:flex;align-items:baseline;gap:12px}
.brand{font-size:20px;font-weight:700;color:#f4f4f5;margin:0}
.sub{font-size:13px;color:#71717a}
.lang-sel{padding:4px 8px;border:1px solid #3f3f46;border-radius:6px;background:#27272a;color:#d4d4d8;font-size:13px;cursor:pointer;outline:none}
.lang-sel:hover{border-color:#52525b}
.lang-sel:focus{border-color:#3b82f6}
.lang-sel option{background:#27272a;color:#d4d4d8}
.main{display:flex;flex:1;overflow:hidden}
.left{width:360px;min-width:320px;overflow-y:auto;padding:16px;border-right:1px solid #27272a;display:flex;flex-direction:column;gap:12px}
.right{flex:1;overflow:hidden;display:flex;flex-direction:column;padding:16px}
</style>
