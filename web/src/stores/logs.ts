import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LogMessage } from '../types'

export const useLogsStore = defineStore('logs', () => {
  const messages = ref<LogMessage[]>([])
  const maxMessages = ref(10000)
  const filter = ref('')

  const filteredMessages = computed(() => {
    if (!filter.value) return messages.value
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

  return {
    messages,
    maxMessages,
    filter,
    filteredMessages,
    addMessages,
    clear,
    setFilter,
  }
})
