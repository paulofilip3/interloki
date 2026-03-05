<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useLogsStore } from '../stores/logs'

const logsStore = useLogsStore()
const inputEl = ref<HTMLInputElement | null>(null)
const searchText = ref('')

const regexMode = computed(() => logsStore.filterMode === 'regex')

const filterActive = computed(() => !!logsStore.filter)
const matchCount = computed(() => logsStore.filteredMessages.length)
const totalCount = computed(() => logsStore.messages.length)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

watch(searchText, (val) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    logsStore.setFilter(val)
  }, 200)
})

function toggleRegex() {
  logsStore.setFilterMode(regexMode.value ? 'text' : 'regex')
}

function clearSearch() {
  searchText.value = ''
  logsStore.setFilter('')
}

function onGlobalKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'f') {
    e.preventDefault()
    inputEl.value?.focus()
  }
}

onMounted(() => {
  document.addEventListener('keydown', onGlobalKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', onGlobalKeydown)
  if (debounceTimer) clearTimeout(debounceTimer)
})
</script>

<template>
  <div class="search-bar">
    <div class="search-bar__icon">&#x1F50D;</div>
    <input
      ref="inputEl"
      v-model="searchText"
      class="search-bar__input"
      placeholder="Filter logs... (Ctrl+F)"
      @keydown.escape="clearSearch"
    />
    <span v-if="filterActive" class="search-bar__count">
      {{ matchCount }} of {{ totalCount }}
    </span>
    <button
      class="search-bar__btn"
      :class="{ 'search-bar__btn--active': regexMode }"
      @click="toggleRegex"
      title="Toggle regex mode"
    >.*</button>
    <button v-if="searchText" class="search-bar__btn" @click="clearSearch" title="Clear">&#x2715;</button>
  </div>
</template>

<style scoped>
.search-bar {
  display: flex;
  align-items: center;
  height: 36px;
  padding: 0 8px;
  border-bottom: 1px solid var(--interloki-border);
  background-color: var(--interloki-bg-secondary);
  font-family: var(--interloki-font-family);
  font-size: var(--interloki-font-size);
  flex-shrink: 0;
  gap: 6px;
}

.search-bar__icon {
  flex-shrink: 0;
  font-size: 14px;
  opacity: 0.5;
}

.search-bar__input {
  flex: 1;
  min-width: 0;
  background: transparent;
  border: none;
  outline: none;
  color: var(--interloki-fg);
  font-family: var(--interloki-font-family);
  font-size: var(--interloki-font-size);
  padding: 4px 0;
}

.search-bar__input::placeholder {
  color: var(--interloki-fg-secondary);
  opacity: 0.6;
}

.search-bar__count {
  flex-shrink: 0;
  color: var(--interloki-fg-secondary);
  font-size: 11px;
  padding: 0 4px;
}

.search-bar__btn {
  flex-shrink: 0;
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

.search-bar__btn:hover {
  background-color: var(--interloki-border);
  color: var(--interloki-fg);
}

.search-bar__btn--active {
  border-color: var(--interloki-accent);
  color: var(--interloki-accent);
  background-color: transparent;
}

.search-bar__btn--active:hover {
  background-color: var(--interloki-accent);
  color: var(--interloki-bg);
}
</style>
