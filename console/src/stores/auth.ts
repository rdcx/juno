import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

import type { User } from '@/types/ProfileTypes'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token'))
  const setToken = (newToken: string) => {
    token.value
    localStorage.setItem('token', newToken)
  }

  const user = ref(localStorage.getItem('user') || '{}')

  const setUser = (newUser: User) => {
    user.value
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const isAuthenticated = computed(() => !!token.value)

  return {
    token,
    setToken,
    user,
    setUser,
    isAuthenticated,
  }
})
