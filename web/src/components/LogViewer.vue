<script setup lang="ts">
import { ref, watch, nextTick, onMounted } from 'vue'
import { useLogsStore } from '../stores/logs'
import { useSettingsStore } from '../stores/settings'
import LogRow from './LogRow.vue'
import type { LogMessage } from '../types'

const props = defineProps<{
  selectedMessage?: LogMessage | null
}>()

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
      <div class="log-viewer__empty-dot"></div>
      <div class="log-viewer__empty-title">Waiting for log messages...</div>
      <div class="log-viewer__empty-subtitle">Connect a log source to get started</div>
    </div>
    <LogRow
      v-for="msg in logsStore.filteredMessages"
      :key="msg.id"
      :message="msg"
      :selected="props.selectedMessage?.id === msg.id"
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
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--interloki-fg-secondary);
  font-family: var(--interloki-font-family);
  gap: 8px;
}

.log-viewer__empty-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: var(--interloki-fg-secondary);
  opacity: 0.4;
  animation: empty-pulse 2s ease-in-out infinite;
  margin-bottom: 4px;
}

@keyframes empty-pulse {
  0%, 100% { opacity: 0.15; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(1.2); }
}

.log-viewer__empty-title {
  font-size: var(--interloki-font-size);
}

.log-viewer__empty-subtitle {
  font-size: 11px;
  opacity: 0.5;
}
</style>
