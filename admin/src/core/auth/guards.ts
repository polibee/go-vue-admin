import type { RouteLocationNormalized } from 'vue-router'

import { useAuthStore } from './auth.store'

export async function AuthGuard(to: RouteLocationNormalized) {
  const auth = useAuthStore()
  const authenticated = await auth.loadCurrentUser()
  return authenticated || { path: '/login', query: { redirect: to.fullPath } }
}

export async function GuestGuard() {
  const auth = useAuthStore()
  const authenticated = auth.status === 'authenticated' || await auth.loadCurrentUser()
  return authenticated ? { path: '/admin/dashboard' } : true
}
