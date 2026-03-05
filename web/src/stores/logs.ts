import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LogMessage } from '../types'

export const useLogsStore = defineStore('logs', () => {
  const messages = ref<LogMessage[]>([])
  const maxMessages = ref(10000)
  const filter = ref('')
  const filterMode = ref<'text' | 'regex'>('text')

  const filteredMessages = computed(() => {
    if (!filter.value) return messages.value

    if (filterMode.value === 'regex') {
      try {
        const re = new RegExp(filter.value, 'i')
        return messages.value.filter((msg) => re.test(msg.content))
      } catch {
        return messages.value // invalid regex, show all
      }
    }

    const term = filter.value.toLowerCase()
    return messages.value.filter((msg) =>
      msg.content.toLowerCase().includes(term),
    )
  })

  function addMessages(msgs: LogMessage[]) {
    messages.value.push(...msgs)
    if (messages.value.length > maxMessages.value) {
      messages.value = messages.value.slice(-maxMessages.value)
    }
  }

  function clear() {
    messages.value = []
  }

  function setFilter(text: string) {
    filter.value = text
  }

  function setFilterMode(mode: 'text' | 'regex') {
    filterMode.value = mode
  }

  return {
    messages,
    maxMessages,
    filter,
    filterMode,
    filteredMessages,
    addMessages,
    clear,
    setFilter,
    setFilterMode,
  }
})
