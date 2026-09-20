const API_BASE_URL = import.meta.env?.VITE_API_BASE_URL ?? 'http://127.0.0.1:3000'

const LOCALIZED_ERROR_CODES = new Set([
  'VALIDATION_ERROR',
  'AUTH_INVALID_CREDENTIALS',
  'AUTH_UNAUTHORIZED',
  'RBAC_FORBIDDEN',
  'RBAC_ROLE_NOT_FOUND',
  'RBAC_PERMISSION_NOT_FOUND',
  'RBAC_USER_NOT_FOUND',
  'RBAC_SYSTEM_ROLE',
  'RBAC_LAST_ADMIN',
  'INTERNAL_ERROR',
])

export function errorMessageKey(code?: string) {
  return code && LOCALIZED_ERROR_CODES.has(code) ? `errors.${code}` : 'errors.unknown'
}

export class ApiError extends Error {
  readonly status: number
  readonly code?: string

  constructor(
    message: string,
    status: number,
    code?: string,
  ) {
    super(message)
    this.status = status
    this.code = code
  }
}

export async function apiFetch<T>(path: string, init: RequestInit = {}, token?: string): Promise<T> {
  const headers = new Headers(init.headers)
  headers.set('Content-Type', 'application/json')
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(`${API_BASE_URL}${path}`, { ...init, headers })
  const payload = await response.json().catch(() => null) as { data?: T; code?: string; message?: string } | null
  if (!response.ok) {
    throw new ApiError(payload?.message ?? 'Request failed', response.status, payload?.code)
  }

  return (payload?.data ?? payload) as T
}
