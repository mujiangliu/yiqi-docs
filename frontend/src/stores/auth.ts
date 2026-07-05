// frontend/src/stores/auth.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api/client'
import type { User } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)

  async function fetchMe() {
    loading.value = true
    try {
      const resp = await authApi.me()
      user.value = resp.data.data
    } catch {
      user.value = null
    } finally {
      loading.value = false
    }
  }

  async function login(username: string, password: string) {
    const resp = await authApi.login(username, password)
    user.value = resp.data.data
    return user.value
  }

  async function logout() {
    await authApi.logout()
    user.value = null
  }

  function isSuper() {
    return user.value?.role === 'super_admin'
  }

  return { user, loading, fetchMe, login, logout, isSuper }
})
