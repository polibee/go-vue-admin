export { AuthService } from './auth.service'
export type { AuthUser, LoginCredentials } from './auth.service'
export { AuthGuard, GuestGuard } from './guards'
import { useAuthStore } from './auth.store'

export { useAuthStore }
export type { AuthStatus, AuthStore } from './auth.store'

export function useAuth() {
  return useAuthStore()
}
