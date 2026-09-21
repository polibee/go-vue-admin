import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { generatedApi, type AuthUser, type LoginResponse } from '@/generated/api'

const STORAGE_KEY = 'go-vue-admin.access-token'
function readStoredToken() { return typeof window === 'undefined' ? undefined : sessionStorage.getItem(STORAGE_KEY) ?? undefined }
function writeStoredToken(value: string | undefined) { if (typeof window === 'undefined') return; if (value) sessionStorage.setItem(STORAGE_KEY, value); else sessionStorage.removeItem(STORAGE_KEY) }

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | undefined>(readStoredToken())
  const user = ref<AuthUser>()
  const restored = ref(false)

  const isAuthenticated = computed(() => Boolean(token.value && user.value))

  async function login(email: string, password: string) {
    const response: LoginResponse = await generatedApi.login({ email, password })
    token.value = response.access_token
    writeStoredToken(token.value)
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
    if (!token.value) return
    try {
      await fetchCurrentUser()
    } catch {
      token.value = undefined
      user.value = undefined
      writeStoredToken(undefined)
    }
  }

  async function logout() {
    const tokenValue = token.value
    if (tokenValue) {
      await generatedApi.logout(tokenValue)
    }
    token.value = undefined
    user.value = undefined
    writeStoredToken(undefined)
  }

  return { token, user, isAuthenticated, login, fetchCurrentUser, restore, logout }
})
