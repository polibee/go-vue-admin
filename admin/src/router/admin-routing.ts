import { isAdminRoute } from '../core/routing/url-namespaces.ts'

export const adminHomePath = '/admin'
export const adminLoginPath = '/admin/login'

export function safeAdminRedirect(value?: string) {
  return value && isAdminRoute(value) ? value : adminHomePath
}
