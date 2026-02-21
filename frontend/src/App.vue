<script lang="ts" setup>
import { watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from './stores/appStore'
import { WindowSetTitle, WindowMinimise, WindowToggleMaximise, Quit } from '../wailsjs/runtime/runtime'
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

// Update window title with conversion progress (visible in taskbar)
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
    <!-- Custom title bar (frameless window) -->
    <div class="titlebar" style="--wails-draggable: drag">
      <div v-if="store.isConverting || store.progress.percent > 0"
           class="titlebar-fill" :class="{ done: store.progress.percent >= 100 }"
           :style="{ width: store.progress.percent + '%' }"/>
      <div class="titlebar-content">
        <span class="titlebar-title">
          <template v-if="store.isConverting">{{ Math.round(store.progress.percent) }}% - PDF2Image</template>
          <template v-else>PDF2Image</template>
        </span>
        <div class="titlebar-btns" style="--wails-draggable: none">
          <button class="tb-btn" @click="WindowMinimise" title="Minimize">&#xE921;</button>
          <button class="tb-btn" @click="WindowToggleMaximise" title="Maximize">&#xE922;</button>
          <button class="tb-btn tb-close" @click="Quit" title="Close">&#xE8BB;</button>
        </div>
      </div>
    </div>
    <!-- App header -->
    <header class="hdr">
      <div class="hdr-l">
        <h1 class="brand">{{ t('app.title') }}</h1>
        <span class="sub">{{ t('app.subtitle') }}</span>
      </div>
      <select class="lang-sel" :value="locale" @change="switchLang">
        <option v-for="lang in languages" :key="lang.code" :value="lang.code">{{ lang.label }}</option>
      </select>
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

/* Custom title bar */
.titlebar{position:relative;height:32px;background:#0a0a0a;flex-shrink:0;user-select:none;overflow:hidden}
.titlebar-fill{position:absolute;top:0;left:0;height:100%;background:rgba(37,99,235,.35);transition:width .3s;pointer-events:none}
.titlebar-fill.done{background:rgba(34,197,94,.3)}
.titlebar-content{position:relative;display:flex;align-items:center;justify-content:space-between;height:100%;padding-left:12px}
.titlebar-title{font-size:12px;color:#a1a1aa;font-weight:500}
.titlebar-btns{display:flex;height:100%}
.tb-btn{width:46px;height:100%;border:none;background:transparent;color:#a1a1aa;font-family:'Segoe MDL2 Assets','Segoe Fluent Icons';font-size:10px;cursor:pointer;display:flex;align-items:center;justify-content:center;transition:background .1s}
.tb-btn:hover{background:rgba(255,255,255,.1);color:#e4e4e7}
.tb-close:hover{background:#c42b1c;color:#fff}

/* App header */
.hdr{display:flex;align-items:center;justify-content:space-between;padding:12px 24px;background:#09090b;border-bottom:1px solid #27272a}
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
