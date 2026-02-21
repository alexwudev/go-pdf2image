<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/appStore'
const { t } = useI18n()
const store = useAppStore()
const label = computed(() => {
  if (store.progress.percent >= 100) return t('progress.done')
  if (store.progress.current > 0) return t('progress.converting', { current: store.progress.current, total: store.progress.total })
  return ''
})
const show = computed(() => store.isConverting || store.progress.percent >= 100)
</script>

<template>
  <div v-if="show" class="prog">
    <div class="row">
      <span class="lbl">{{ label }}</span>
      <div class="row-r">
        <span v-if="store.convertElapsed" class="elapsed">{{ store.convertElapsed }}</span>
        <span class="pct">{{ Math.round(store.progress.percent) }}%</span>
      </div>
    </div>
    <div class="track"><div class="fill" :class="{ done: store.progress.percent >= 100 }" :style="{ width: store.progress.percent + '%' }"/></div>
  </div>
</template>

<style scoped>
.prog{margin-bottom:12px}
.row{display:flex;justify-content:space-between;align-items:center;margin-bottom:4px}
.row-r{display:flex;align-items:center;gap:8px}
.lbl{font-size:12px;color:#a1a1aa}
.elapsed{font-size:12px;color:#71717a}
.pct{font-size:12px;color:#71717a}
.track{height:4px;background:#27272a;border-radius:2px;overflow:hidden}
.fill{height:100%;background:#2563eb;border-radius:2px;transition:width .3s}
.fill.done{background:#22c55e}
</style>
