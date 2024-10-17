import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

import type { User } from '@/types/ProfileTypes'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token'))
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const user = ref<User>(localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user') as string) : {} as User)

  const setUser = (newUser: User) => {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  const isAuthenticated = computed(() => !!token.value)

  const logout = () => {
    token.value = null
    user.value = {} as User
    localStorage.removeItem('token')
    localStorage.removeItem('user')

    window.location.reload()
  }

  return {
    token,
    setToken,
    user,
    setUser,
    isAuthenticated,
    logout,
  }
})
