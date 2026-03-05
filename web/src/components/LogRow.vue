<script setup lang="ts">
import { computed } from 'vue'
import type { LogMessage } from '../types'
import { useSettingsStore } from '../stores/settings'

const props = defineProps<{
  message: LogMessage
}>()

const settings = useSettingsStore()

const formattedTime = computed(() => {
  try {
    const d = new Date(props.message.ts)
    const h = String(d.getHours()).padStart(2, '0')
    const m = String(d.getMinutes()).padStart(2, '0')
    const s = String(d.getSeconds()).padStart(2, '0')
    const ms = String(d.getMilliseconds()).padStart(3, '0')
    return `${h}:${m}:${s}.${ms}`
  } catch {
    return props.message.ts
  }
})

const levelClass = computed(() => {
  const level = props.message.level?.toLowerCase() ?? ''
  if (level.includes('fatal')) return 'log-row__level--fatal'
  if (level.includes('error') || level.includes('err')) return 'log-row__level--error'
  if (level.includes('warn')) return 'log-row__level--warn'
  if (level.includes('info')) return 'log-row__level--info'
  if (level.includes('debug') || level.includes('dbg') || level.includes('trace')) return 'log-row__level--debug'
  return ''
})
</script>

<template>
  <div class="log-row">
    <span v-if="settings.showTimestamp" class="log-row__timestamp">{{ formattedTime }}</span>
    <span
      v-if="settings.showLevel && message.level"
      class="log-row__level"
      :class="levelClass"
    >{{ message.level.toUpperCase() }}</span>
    <span v-if="settings.showSource" class="log-row__source">{{ message.source }}</span>
    <span class="log-row__content">{{ message.content }}</span>
  </div>
</template>

<style scoped>
.log-row {
  display: flex;
  align-items: center;
  height: var(--interloki-row-height);
  padding: 0 8px;
  font-family: var(--interloki-font-family);
  font-size: var(--interloki-font-size);
  color: var(--interloki-fg);
  white-space: nowrap;
  cursor: pointer;
  gap: 8px;
  user-select: none;
}

.log-row:nth-child(even) {
  background-color: var(--interloki-bg-secondary);
}

.log-row:hover {
  background-color: var(--interloki-border);
}

.log-row__timestamp {
  flex-shrink: 0;
  color: var(--interloki-fg-secondary);
  min-width: 95px;
}

.log-row__level {
  flex-shrink: 0;
  min-width: 48px;
  text-align: center;
  font-weight: 600;
  font-size: 11px;
  border-radius: 3px;
  padding: 1px 4px;
}

.log-row__level--debug {
  color: var(--interloki-level-debug);
}

.log-row__level--info {
  color: var(--interloki-level-info);
}

.log-row__level--warn {
  color: var(--interloki-level-warn);
}

.log-row__level--error {
  color: var(--interloki-level-error);
}

.log-row__level--fatal {
  color: var(--interloki-level-fatal);
}

.log-row__source {
  flex-shrink: 0;
  color: var(--interloki-fg-secondary);
  opacity: 0.6;
  font-size: 11px;
  min-width: 40px;
}

.log-row__content {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
