import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { generatedApi, type AuthUser, type LoginResponse } from '@/generated/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | undefined>()
  const user = ref<AuthUser>()
  const restored = ref(false)

  const isAuthenticated = computed(() => Boolean(token.value && user.value))

  async function login(email: string, password: string) {
    const response: LoginResponse = await generatedApi.login({ email, password })
    token.value = response.access_token
    user.value = response.user
  }

  async function fetchCurrentUser() {
    const tokenValue = token.value
    if (!tokenValue) return
    user.value = await generatedApi.currentUser(tokenValue)
  }

  async function restore() {
    if (restored.value) return
    restored.value = true
    try {
      const response = await generatedApi.refresh()
      token.value = response.access_token
      await fetchCurrentUser()
    } catch {
      token.value = undefined
      user.value = undefined
    }
  }

  async function logout() {
    const tokenValue = token.value
    if (tokenValue) {
      await generatedApi.logout(tokenValue)
    }
    token.value = undefined
    user.value = undefined
  }

  async function logoutAll() {
    const tokenValue = token.value
    if (tokenValue) {
      await generatedApi.logoutAll(tokenValue)
    }
    token.value = undefined
    user.value = undefined
  }

  return { token, user, isAuthenticated, login, fetchCurrentUser, restore, logout, logoutAll }
})
