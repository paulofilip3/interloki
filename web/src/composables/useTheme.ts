import { watch, onMounted } from 'vue'
import { useSettingsStore } from '../stores/settings'

export function useTheme() {
  const settingsStore = useSettingsStore()

  function applyTheme(theme: 'light' | 'dark') {
    document.documentElement.setAttribute('data-theme', theme)
  }

  onMounted(() => {
    applyTheme(settingsStore.theme)
  })

  watch(
    () => settingsStore.theme,
    (newTheme) => {
      applyTheme(newTheme)
    },
  )

  return {
    theme: settingsStore.theme,
    toggleTheme: settingsStore.toggleTheme,
  }
}
