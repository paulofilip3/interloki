<script setup lang="ts">
import { useSettingsStore } from '../stores/settings'

const settings = useSettingsStore()
</script>

<template>
  <button class="theme-toggle" @click="settings.toggleTheme()" :title="settings.theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'">
    {{ settings.theme === 'dark' ? '\u2600\uFE0F' : '\uD83C\uDF19' }}
  </button>
</template>

<style scoped>
.theme-toggle {
  background: none;
  border: 1px solid var(--interloki-border);
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  padding: 2px 8px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.15s;
}

.theme-toggle:hover {
  border-color: var(--interloki-accent);
}
</style>
