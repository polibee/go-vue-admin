export interface ApiClient {
  request<T>(path: string, options?: RequestInit): Promise<T>
}

export class ApiError extends Error {
  readonly status: number
  readonly code?: string

  constructor(message: string, status: number, code?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
  }
}

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL ?? '').replace(/\/$/, '')

export const apiClient: ApiClient = {
  async request<T>(path: string, options: RequestInit = {}) {
    const response = await fetch(`${apiBaseUrl}${path}`, {
      ...options,
      credentials: 'include',
      headers: {
        Accept: 'application/json',
        ...options.headers,
      },
    })

    const payload = await response.json().catch(() => null)
    if (!response.ok) {
      throw new ApiError(
        payload?.error?.message ?? '请求失败，请稍后重试',
        response.status,
        payload?.error?.code,
      )
    }

    return payload as T
  },
}
