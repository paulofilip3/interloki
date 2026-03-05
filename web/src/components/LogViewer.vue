<script setup lang="ts">
import { ref, watch, nextTick, onMounted } from 'vue'
import { useLogsStore } from '../stores/logs'
import { useSettingsStore } from '../stores/settings'
import LogRow from './LogRow.vue'
import type { LogMessage } from '../types'

defineEmits<{
  'row-click': [message: LogMessage]
}>()

const logsStore = useLogsStore()
const settings = useSettingsStore()

const container = ref<HTMLElement | null>(null)
const isUserScrolledUp = ref(false)

function onScroll() {
  if (!container.value) return
  const { scrollTop, scrollHeight, clientHeight } = container.value
  // Consider "at bottom" if within 30px of the end
  const atBottom = scrollHeight - scrollTop - clientHeight < 30
  if (atBottom) {
    isUserScrolledUp.value = false
    if (!settings.autoFollow) {
      settings.setAutoFollow(true)
    }
  } else {
    isUserScrolledUp.value = true
    if (settings.autoFollow) {
      settings.setAutoFollow(false)
    }
  }
}

function scrollToBottom() {
  if (!container.value) return
  container.value.scrollTop = container.value.scrollHeight
}

watch(
  () => logsStore.filteredMessages.length,
  async () => {
    if (settings.autoFollow && !isUserScrolledUp.value) {
      await nextTick()
      scrollToBottom()
    }
  },
)

onMounted(() => {
  if (settings.autoFollow) {
    scrollToBottom()
  }
})
</script>

<template>
  <div class="log-viewer" ref="container" @scroll="onScroll">
    <div v-if="logsStore.filteredMessages.length === 0" class="log-viewer__empty">
      Waiting for log messages...
    </div>
    <LogRow
      v-for="msg in logsStore.filteredMessages"
      :key="msg.id"
      :message="msg"
      @click="$emit('row-click', msg)"
    />
  </div>
</template>

<style scoped>
.log-viewer {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  background-color: var(--interloki-bg);
}

.log-viewer__empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--interloki-fg-secondary);
  font-family: var(--interloki-font-family);
  font-size: var(--interloki-font-size);
}
</style>
