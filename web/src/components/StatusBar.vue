<script setup lang="ts">
import { computed } from 'vue'
import { useConnectionStore } from '../stores/connection'
import { useLogsStore } from '../stores/logs'
import { useSettingsStore } from '../stores/settings'

const connectionStore = useConnectionStore()
const logsStore = useLogsStore()
const settings = useSettingsStore()

const statusLabel = computed(() => {
  switch (connectionStore.status) {
    case 'connected':
      return 'Connected'
    case 'connecting':
      return 'Connecting...'
    default:
      return 'Disconnected'
  }
})

const statusDotClass = computed(() => {
  return `status-bar__dot--${connectionStore.status}`
})

const messageCount = computed(() => logsStore.messages.length)

const bufferUsage = computed(() => {
  const stats = connectionStore.serverStats
  if (!stats) return null
  return `${stats.buffer_used}/${stats.buffer_capacity}`
})
</script>

<template>
  <div class="status-bar">
    <div class="status-bar__item">
      <span class="status-bar__dot" :class="statusDotClass"></span>
      <span>{{ statusLabel }}</span>
    </div>
    <div class="status-bar__divider"></div>
    <div class="status-bar__item">
      <span>Messages: {{ messageCount }}</span>
    </div>
    <div v-if="bufferUsage" class="status-bar__divider"></div>
    <div v-if="bufferUsage" class="status-bar__item">
      <span>Buffer: {{ bufferUsage }}</span>
    </div>
    <div class="status-bar__divider"></div>
    <div class="status-bar__item">
      <button @click="connectionStore.toggleFollowing()" class="status-bar__control">
        {{ connectionStore.following ? '\u23F8' : '\u25B6' }}
      </button>
      <button v-if="!connectionStore.following" @click="connectionStore.loadHistory()" class="status-bar__control">
        Load History
      </button>
    </div>
    <div class="status-bar__spacer"></div>
    <div v-if="!connectionStore.following" class="status-bar__item status-bar__item--subtle status-bar__item--paused">
      PAUSED
    </div>
    <div v-if="settings.autoFollow" class="status-bar__item status-bar__item--subtle">
      AUTO-FOLLOW
    </div>
  </div>
</template>

<style scoped>
.status-bar {
  display: flex;
  align-items: center;
  height: 28px;
  padding: 0 8px;
  background-color: var(--interloki-bg-secondary);
  border-top: 1px solid var(--interloki-border);
  font-family: var(--interloki-font-family);
  font-size: 11px;
  color: var(--interloki-fg-secondary);
  flex-shrink: 0;
  gap: 0;
}

.status-bar__item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 8px;
}

.status-bar__item--subtle {
  opacity: 0.6;
  font-size: 10px;
  letter-spacing: 0.5px;
}

.status-bar__divider {
  width: 1px;
  height: 14px;
  background-color: var(--interloki-border);
}

.status-bar__spacer {
  flex: 1;
}

.status-bar__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-bar__dot--connected {
  background-color: #22c55e;
}

.status-bar__dot--connecting {
  background-color: #eab308;
}

.status-bar__dot--disconnected {
  background-color: #ef4444;
}

.status-bar__control {
  background: none;
  border: 1px solid var(--interloki-border);
  color: var(--interloki-fg-secondary);
  font-family: var(--interloki-font-family);
  font-size: 11px;
  padding: 1px 8px;
  border-radius: 3px;
  cursor: pointer;
  line-height: 1;
}

.status-bar__control:hover {
  background-color: var(--interloki-border);
  color: var(--interloki-fg);
}

.status-bar__item--paused {
  color: #eab308;
}
</style>
