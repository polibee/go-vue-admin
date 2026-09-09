import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { AuthService, type AuthUser, type LoginCredentials } from './auth.service'
import { can as checkPermission } from '@/core/permissions'

export type AuthStatus = 'idle' | 'loading' | 'authenticated' | 'guest'

export const useAuthStore = defineStore('auth', () => {
  const service = new AuthService()
  const user = ref<AuthUser | null>(null)
  const status = ref<AuthStatus>('idle')
  const error = ref<string | null>(null)
  let loadingPromise: Promise<boolean> | null = null

  const isAuthenticated = computed(() => status.value === 'authenticated' && user.value !== null)
  const can = (permission: string) => checkPermission(user.value?.permissions ?? [], permission)

  async function loadCurrentUser(): Promise<boolean> {
    if (loadingPromise) return loadingPromise
    if (status.value === 'authenticated' && user.value) return true

    status.value = 'loading'
    loadingPromise = service.me()
      .then((currentUser) => {
        user.value = currentUser
        status.value = 'authenticated'
        error.value = null
        return true
      })
      .catch(() => {
        user.value = null
        status.value = 'guest'
        return false
      })
      .finally(() => {
        loadingPromise = null
      })

    return loadingPromise
  }

  async function login(credentials: LoginCredentials): Promise<AuthUser> {
    status.value = 'loading'
    error.value = null
    try {
      user.value = await service.login(credentials)
      status.value = 'authenticated'
      return user.value
    } catch (cause) {
      user.value = null
      status.value = 'guest'
      error.value = cause instanceof Error ? cause.message : '登录失败，请稍后重试'
      throw cause
    }
  }

  async function logout(): Promise<void> {
    try {
      await service.logout()
    } finally {
      user.value = null
      status.value = 'guest'
      error.value = null
    }
  }

  return { user, status, error, isAuthenticated, can, loadCurrentUser, login, logout }
})

export type AuthStore = ReturnType<typeof useAuthStore>
