<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/appStore'
import { SelectOutputDir } from '../../wailsjs/go/main/App'

const { t } = useI18n()
const store = useAppStore()

async function pickDir() {
  try {
    const dir = await SelectOutputDir()
    if (dir) store.outputDir = dir
  } catch (e) { console.error(e) }
}
</script>

<template>
  <div class="card">
    <h3 class="card-title">{{ t('settings.title') }}</h3>

    <label class="lbl">{{ t('settings.format') }}</label>
    <div class="seg">
      <button :class="['seg-btn', { on: store.config.format === 'jpg' }]" @click="store.config.format = 'jpg'">JPG</button>
      <button :class="['seg-btn', { on: store.config.format === 'png' }]" @click="store.config.format = 'png'">PNG</button>
    </div>

    <label class="lbl">{{ t('settings.dpi') }}: {{ store.config.dpi }}</label>
    <input type="range" class="slider" v-model.number="store.config.dpi" min="72" max="600" step="1" />

    <template v-if="store.config.format === 'jpg'">
      <label class="lbl">{{ t('settings.quality') }}: {{ store.config.quality }}%</label>
      <input type="range" class="slider" v-model.number="store.config.quality" min="10" max="100" step="1" />
    </template>

    <label class="lbl">{{ t('settings.pages') }}</label>
    <div class="seg">
      <button :class="['seg-btn', { on: store.config.pages === 'all' }]" @click="store.config.pages = 'all'">{{ t('settings.pagesAll') }}</button>
      <button :class="['seg-btn', { on: store.config.pages !== 'all' }]" @click="store.config.pages = ''">自訂</button>
    </div>
    <input v-if="store.config.pages !== 'all'" type="text" class="input" v-model="store.config.pages" :placeholder="t('settings.pagesHint')" />

    <label class="lbl">{{ t('settings.workers') }}: {{ store.config.workers }}</label>
    <input type="range" class="slider" v-model.number="store.config.workers" min="1" max="20" step="1" />

    <label class="lbl">{{ t('settings.outputDir') }}</label>
    <div class="dir-row">
      <span class="dir-path">{{ store.outputDir || t('settings.sameAsPdf') }}</span>
      <button class="dir-btn" @click="pickDir">{{ t('settings.selectDir') }}</button>
    </div>

    <label class="chk-row">
      <input type="checkbox" v-model="store.config.zip_output" class="chk" />
      <span class="chk-label">{{ t('settings.zipOutput') }}</span>
    </label>
  </div>
</template>

<style scoped>
.card{background:#1f1f23;border:1px solid #27272a;border-radius:8px;padding:14px}
.card-title{font-size:14px;font-weight:600;margin:0 0 10px;color:#d4d4d8}
.lbl{display:block;font-size:12px;color:#a1a1aa;margin:8px 0 4px}
.lbl:first-of-type{margin-top:0}
.seg{display:flex;gap:0;border:1px solid #3f3f46;border-radius:6px;overflow:hidden}
.seg-btn{flex:1;padding:6px 12px;border:none;background:#27272a;color:#a1a1aa;font-size:13px;cursor:pointer;transition:.15s}
.seg-btn.on{background:#3b82f6;color:#fff}
.seg-btn:hover:not(.on){background:#3f3f46;color:#e4e4e7}
.slider{width:100%;accent-color:#3b82f6;margin:2px 0}
.input{width:100%;padding:6px 10px;border:1px solid #3f3f46;border-radius:6px;background:#27272a;color:#e4e4e7;font-size:13px;outline:none}
.input:focus{border-color:#3b82f6}
.dir-row{display:flex;align-items:center;gap:8px}
.dir-path{flex:1;font-size:12px;color:#71717a;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.dir-btn{padding:5px 12px;border:1px solid #3f3f46;border-radius:6px;background:transparent;color:#d4d4d8;font-size:12px;cursor:pointer;white-space:nowrap}
.dir-btn:hover{background:#27272a}
.chk-row{display:flex;align-items:center;gap:8px;margin-top:10px;cursor:pointer}
.chk{accent-color:#3b82f6;width:16px;height:16px;cursor:pointer}
.chk-label{font-size:13px;color:#d4d4d8}
</style>
