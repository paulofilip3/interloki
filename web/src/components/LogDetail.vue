<script setup lang="ts">
import { computed } from 'vue'
import type { LogMessage } from '../types'

const props = defineProps<{
  message: LogMessage | null
  visible: boolean
}>()

defineEmits<{
  close: []
}>()

const highlightedJson = computed(() => {
  if (!props.message?.is_json || !props.message.json_content) return ''
  try {
    const raw = JSON.stringify(props.message.json_content, null, 2)
    return raw.replace(
      /("(?:[^"\\]|\\.)*")(\s*:)?|(\b(?:true|false|null)\b)|(-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)/g,
      (_match: string, str?: string, colon?: string, bool?: string, num?: string) => {
        if (str) {
          if (colon) {
            return `<span class="json-key">${str}</span>${colon}`
          }
          return `<span class="json-string">${str}</span>`
        }
        if (bool) {
          return `<span class="json-bool">${bool}</span>`
        }
        if (num) {
          return `<span class="json-number">${num}</span>`
        }
        return _match
      },
    )
  } catch {
    return ''
  }
})
</script>

<template>
  <div v-if="visible && message" class="log-detail">
    <div class="log-detail__header">
      <span class="log-detail__title">Log Detail</span>
      <button class="log-detail__close" @click="$emit('close')">&#x2715;</button>
    </div>
    <div class="log-detail__body">
      <div class="log-detail__section">
        <h3>Metadata</h3>
        <table class="log-detail__meta">
          <tr><td>Timestamp</td><td>{{ message.ts }}</td></tr>
          <tr><td>Source</td><td>{{ message.source }}</td></tr>
          <tr><td>Origin</td><td>{{ message.origin.name }}</td></tr>
          <tr v-if="message.level"><td>Level</td><td>{{ message.level }}</td></tr>
          <tr v-for="(v, k) in message.labels" :key="k"><td>{{ k }}</td><td>{{ v }}</td></tr>
        </table>
      </div>
      <div v-if="message.is_json" class="log-detail__section">
        <h3>JSON Content</h3>
        <!-- eslint-disable-next-line vue/no-v-html -->
        <pre class="log-detail__json" v-html="highlightedJson"></pre>
      </div>
      <div class="log-detail__section">
        <h3>Raw</h3>
        <pre class="log-detail__raw">{{ message.content }}</pre>
      </div>
    </div>
  </div>
</template>

<style scoped>
.log-detail {
  flex-shrink: 0;
  height: 40%;
  min-height: 120px;
  max-height: 50vh;
  display: flex;
  flex-direction: column;
  border-top: 2px solid var(--interloki-accent);
  background-color: var(--interloki-bg);
  font-family: var(--interloki-font-family);
  font-size: var(--interloki-font-size);
}

.log-detail__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  height: 32px;
  background-color: var(--interloki-bg-secondary);
  border-bottom: 1px solid var(--interloki-border);
  flex-shrink: 0;
}

.log-detail__title {
  font-size: 12px;
  font-weight: 600;
  color: var(--interloki-fg);
  letter-spacing: 0.3px;
}

.log-detail__close {
  background: none;
  border: 1px solid var(--interloki-border);
  color: var(--interloki-fg-secondary);
  font-size: 12px;
  padding: 1px 6px;
  border-radius: 3px;
  cursor: pointer;
  line-height: 1;
}

.log-detail__close:hover {
  background-color: var(--interloki-border);
  color: var(--interloki-fg);
}

.log-detail__body {
  flex: 1;
  overflow-y: auto;
  padding: 8px 12px;
}

.log-detail__section {
  margin-bottom: 12px;
}

.log-detail__section h3 {
  margin: 0 0 6px 0;
  font-size: 11px;
  font-weight: 600;
  color: var(--interloki-fg-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.log-detail__meta {
  border-collapse: collapse;
  font-size: 12px;
}

.log-detail__meta td {
  padding: 2px 12px 2px 0;
  vertical-align: top;
}

.log-detail__meta td:first-child {
  color: var(--interloki-fg-secondary);
  font-weight: 600;
  white-space: nowrap;
}

.log-detail__meta td:last-child {
  color: var(--interloki-fg);
  word-break: break-all;
}

.log-detail__json,
.log-detail__raw {
  margin: 0;
  padding: 8px;
  background-color: var(--interloki-bg-secondary);
  border: 1px solid var(--interloki-border);
  border-radius: 4px;
  font-family: var(--interloki-font-family);
  font-size: 12px;
  color: var(--interloki-fg);
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.5;
}

.log-detail__json :deep(.json-key) {
  color: var(--interloki-accent);
}

.log-detail__json :deep(.json-string) {
  color: var(--interloki-level-warn);
}

.log-detail__json :deep(.json-number) {
  color: var(--interloki-level-info);
}

.log-detail__json :deep(.json-bool) {
  color: var(--interloki-level-error);
}
</style>
