<script setup lang="ts">
import { ref, onMounted } from 'vue'
import LogViewer from './components/LogViewer.vue'
import LogDetail from './components/LogDetail.vue'
import SearchBar from './components/SearchBar.vue'
import ColumnConfig from './components/ColumnConfig.vue'
import StatusBar from './components/StatusBar.vue'
import ThemeToggle from './components/ThemeToggle.vue'
import { useWebSocket } from './composables/useWebSocket'
import { useTheme } from './composables/useTheme'
import type { LogMessage } from './types'

useTheme()

const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
const wsUrl = `${wsProtocol}//${window.location.host}/ws`
const { connect } = useWebSocket(wsUrl)

const selectedMessage = ref<LogMessage | null>(null)
const detailVisible = ref(false)

function onRowClick(msg: LogMessage) {
  selectedMessage.value = msg
  detailVisible.value = true
}

function closeDetail() {
  detailVisible.value = false
}

onMounted(() => {
  connect()
})
</script>

<template>
  <div class="interloki-app">
    <header class="interloki-header">
      <h1 class="interloki-header__title">interloki</h1>
      <div class="interloki-header__actions">
        <ColumnConfig />
        <ThemeToggle />
      </div>
    </header>
    <SearchBar />
    <LogViewer class="interloki-main" :selected-message="selectedMessage" @row-click="onRowClick" />
    <LogDetail :message="selectedMessage" :visible="detailVisible" @close="closeDetail" />
    <StatusBar />
  </div>
</template>

<style scoped>
.interloki-app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--interloki-bg);
  color: var(--interloki-fg);
  font-family: var(--interloki-font-family);
}

.interloki-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  height: 36px;
  border-bottom: 1px solid var(--interloki-border);
  background-color: var(--interloki-bg-secondary);
  box-shadow: 0 1px 3px var(--interloki-shadow);
  flex-shrink: 0;
  z-index: 10;
}

.interloki-header__title {
  font-size: 14px;
  font-weight: 600;
  margin: 0;
  letter-spacing: 0.5px;
  color: var(--interloki-accent);
}

.interloki-header__actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.interloki-main {
  flex: 1;
  min-height: 0;
}
</style>
