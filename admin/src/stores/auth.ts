import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { apiFetch } from '@/lib/api'

export interface AuthUser {
  id: number
  name: string
  email: string
  is_active: boolean
  locale: string
}

interface LoginResponse {
  access_token: string
  token_type: string
  user: AuthUser
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>()
  const user = ref<AuthUser>()

  const isAuthenticated = computed(() => Boolean(token.value && user.value))

  async function login(email: string, password: string) {
    const response = await apiFetch<LoginResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    token.value = response.access_token
    user.value = response.user
  }

  async function fetchCurrentUser() {
    if (!token.value) return
    user.value = await apiFetch<AuthUser>('/api/v1/auth/me', {}, token.value)
  }

  async function logout() {
    if (token.value) {
      await apiFetch<void>('/api/v1/auth/logout', { method: 'POST' }, token.value)
    }
    token.value = undefined
    user.value = undefined
  }

  return { token, user, isAuthenticated, login, fetchCurrentUser, logout }
})
