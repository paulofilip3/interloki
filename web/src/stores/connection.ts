import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { StatusData } from '../types'

export const useConnectionStore = defineStore('connection', () => {
  const status = ref<'connecting' | 'connected' | 'disconnected'>('disconnected')
  const clientId = ref<string | null>(null)
  const bufferSize = ref(0)
  const serverStats = ref<StatusData | null>(null)

  function setConnected(id: string, buffer: number) {
    status.value = 'connected'
    clientId.value = id
    bufferSize.value = buffer
  }

  function setDisconnected() {
    status.value = 'disconnected'
    clientId.value = null
  }

  function updateStats(stats: StatusData) {
    serverStats.value = stats
  }

  return {
    status,
    clientId,
    bufferSize,
    serverStats,
    setConnected,
    setDisconnected,
    updateStats,
  }
})
