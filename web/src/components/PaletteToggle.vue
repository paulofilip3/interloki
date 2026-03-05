<script setup lang="ts">
import { useSettingsStore } from '../stores/settings'

const settings = useSettingsStore()
</script>

<template>
  <button
    class="palette-toggle"
    :class="{ 'palette-toggle--active': settings.palette === 'catppuccin' }"
    @click="settings.togglePalette()"
    :title="settings.palette === 'catppuccin' ? 'Switch to default palette' : 'Switch to Catppuccin palette'"
  >
    {{ settings.palette === 'catppuccin' ? 'Catppuccin' : 'Default' }}
  </button>
</template>

<style scoped>
.palette-toggle {
  background: none;
  border: 1px solid var(--interloki-border);
  border-radius: 4px;
  cursor: pointer;
  padding: 3px 8px;
  font-size: 11px;
  font-family: var(--interloki-font-family);
  color: var(--interloki-fg-secondary);
  transition: border-color 0.15s, color 0.15s, background-color 0.15s;
  white-space: nowrap;
}

.palette-toggle:hover {
  border-color: var(--interloki-accent);
  color: var(--interloki-fg);
}

.palette-toggle--active {
  border-color: var(--interloki-accent);
  color: var(--interloki-accent);
}
</style>
