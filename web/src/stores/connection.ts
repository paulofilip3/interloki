import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { StatusData, LogMessage } from '../types'
import { useLogsStore } from './logs'

export const useConnectionStore = defineStore('connection', () => {
  const status = ref<'connecting' | 'connected' | 'disconnected'>('disconnected')
  const clientId = ref<string | null>(null)
  const bufferSize = ref(0)
  const serverStats = ref<StatusData | null>(null)
  const following = ref(true)

  const pauseFn = ref<(() => void) | null>(null)
  const resumeFn = ref<(() => void) | null>(null)

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

  function setFollowing(val: boolean) {
    following.value = val
  }

  function registerControls(pause: () => void, resume: () => void) {
    pauseFn.value = pause
    resumeFn.value = resume
  }

  function toggleFollowing() {
    if (following.value) {
      pauseFn.value?.()
    } else {
      resumeFn.value?.()
    }
  }

  async function loadHistory() {
    const res = await fetch(`/api/client/load?start=0&count=${bufferSize.value}`)
    const data = await res.json() as { messages: LogMessage[]; total: number }
    const logsStore = useLogsStore()
    logsStore.clear()
    logsStore.addMessages(data.messages)
  }

  return {
    status,
    clientId,
    bufferSize,
    serverStats,
    following,
    setConnected,
    setDisconnected,
    updateStats,
    setFollowing,
    registerControls,
    toggleFollowing,
    loadHistory,
  }
})
