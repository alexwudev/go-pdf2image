<script lang="ts" setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../stores/appStore'
const { t } = useI18n()
const store = useAppStore()

const zoom = ref(1)
const pan = ref({ x: 0, y: 0 })
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })
const panStart = ref({ x: 0, y: 0 })

function onWheel(e: WheelEvent) {
  e.preventDefault()
  const factor = e.deltaY > 0 ? 0.85 : 1.18
  zoom.value = Math.max(0.25, Math.min(30, zoom.value * factor))
}

function onMouseDown(e: MouseEvent) {
  isDragging.value = true
  dragStart.value = { x: e.clientX, y: e.clientY }
  panStart.value = { x: pan.value.x, y: pan.value.y }
}

function onMouseMove(e: MouseEvent) {
  if (isDragging.value) {
    pan.value = {
      x: panStart.value.x + (e.clientX - dragStart.value.x),
      y: panStart.value.y + (e.clientY - dragStart.value.y),
    }
  }
}

function onMouseUp() {
  isDragging.value = false
}

function resetView() {
  zoom.value = 1
  pan.value = { x: 0, y: 0 }
}

function imgStyle() {
  return {
    transform: `translate(${pan.value.x}px, ${pan.value.y}px) scale(${zoom.value})`,
    cursor: isDragging.value ? 'grabbing' : 'grab',
  }
}

function prevPage() {
  if (store.currentPage > 0) store.loadPreview(store.currentPage - 1)
}

function nextPage() {
  if (store.currentPage < store.pageCount - 1) store.loadPreview(store.currentPage + 1)
}
</script>

<template>
  <div class="preview" @mousemove="onMouseMove" @mouseup="onMouseUp" @mouseleave="onMouseUp">
    <div class="hdr-row">
      <h3 class="title">{{ t('preview.title') }}</h3>
      <div v-if="store.pageCount > 0" class="page-nav">
        <button class="nav-btn" :disabled="store.currentPage <= 0" @click="prevPage">&lt;</button>
        <span class="page-label">{{ t('preview.page', { page: store.currentPage + 1 }) }} / {{ store.pageCount }}</span>
        <button class="nav-btn" :disabled="store.currentPage >= store.pageCount - 1" @click="nextPage">&gt;</button>
        <span v-if="zoom !== 1" class="badge zoom" @click="resetView">{{ Math.round(zoom * 100) }}%</span>
      </div>
    </div>
    <div class="pane">
      <div class="cnt" @wheel.prevent="onWheel" @mousedown.prevent="onMouseDown" @dblclick="resetView">
        <img v-if="store.previewUrl" :src="store.previewUrl" class="img" :style="imgStyle()" draggable="false" />
        <div v-else class="ph">{{ t('preview.noFile') }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.preview{flex:1;display:flex;flex-direction:column;overflow:hidden}
.hdr-row{display:flex;justify-content:space-between;align-items:center;margin-bottom:12px}
.title{font-size:14px;font-weight:600;margin:0;color:#d4d4d8}
.page-nav{display:flex;align-items:center;gap:8px}
.nav-btn{width:28px;height:28px;border:1px solid #3f3f46;border-radius:6px;background:#27272a;color:#d4d4d8;font-size:14px;cursor:pointer;display:flex;align-items:center;justify-content:center}
.nav-btn:disabled{opacity:.3;cursor:not-allowed}
.nav-btn:hover:not(:disabled){background:#3f3f46}
.page-label{font-size:13px;color:#a1a1aa;min-width:80px;text-align:center}
.badge{background:#27272a;padding:2px 6px;border-radius:4px;font-size:11px}
.zoom{cursor:pointer;color:#3b82f6}
.zoom:hover{background:#3f3f46}
.pane{flex:1;display:flex;flex-direction:column;border:1px solid #27272a;border-radius:8px;overflow:hidden;background:#09090b}
.cnt{flex:1;display:flex;align-items:center;justify-content:center;overflow:hidden;position:relative;user-select:none}
.img{max-width:100%;max-height:100%;object-fit:contain;transform-origin:center center;will-change:transform;pointer-events:none}
.ph{color:#3f3f46;font-size:14px}
</style>
