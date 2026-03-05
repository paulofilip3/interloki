<script setup lang="ts">
import { onMounted } from 'vue'
import LogViewer from './components/LogViewer.vue'
import StatusBar from './components/StatusBar.vue'
import ThemeToggle from './components/ThemeToggle.vue'
import { useWebSocket } from './composables/useWebSocket'
import { useTheme } from './composables/useTheme'

useTheme()

const wsUrl = `ws://${window.location.host}/ws`
const { connect } = useWebSocket(wsUrl)

onMounted(() => {
  connect()
})
</script>

<template>
  <div class="interloki-app">
    <header class="interloki-header">
      <h1 class="interloki-header__title">interloki</h1>
      <ThemeToggle />
    </header>
    <LogViewer class="interloki-main" />
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
  flex-shrink: 0;
}

.interloki-header__title {
  font-size: 14px;
  font-weight: 600;
  margin: 0;
  letter-spacing: 0.5px;
}

.interloki-main {
  flex: 1;
  min-height: 0;
}
</style>
