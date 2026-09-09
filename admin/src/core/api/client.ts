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
const defaultApiPorts = ['3003', '3000', '3001', '3002']

function runtimeBackendHost(): string {
  if (typeof window === 'undefined') return '127.0.0.1'
  return window.location.hostname || '127.0.0.1'
}

export async function resolveApiBaseUrl(
  fetchImpl: typeof fetch = fetch,
  host = runtimeBackendHost(),
  ports = defaultApiPorts,
): Promise<string> {
  const protocol = typeof window !== 'undefined' && window.location.protocol === 'https:' ? 'https' : 'http'
  for (const port of ports) {
    const candidate = `${protocol}://${host}:${port}`
    try {
      const response = await fetchImpl(`${candidate}/api/health`, {
        credentials: 'include',
        headers: { Accept: 'application/json' },
      })
      if (response.ok) return candidate
    } catch {
      // Try the next conventional local backend port.
    }
  }
  return ''
}

export function createApiClient(fetchImpl: typeof fetch = fetch, baseUrl = apiBaseUrl): ApiClient {
  let runtimeBaseUrl: Promise<string> | null = null

  return {
    async request<T>(path: string, options: RequestInit = {}) {
      let response: Response
      try {
        if (!baseUrl && !runtimeBaseUrl) {
          runtimeBaseUrl = resolveApiBaseUrl(fetchImpl)
        }
        const requestBaseUrl = baseUrl || await runtimeBaseUrl
        response = await fetchImpl(`${requestBaseUrl}${path}`, {
          ...options,
          credentials: 'include',
          headers: {
            Accept: 'application/json',
            ...options.headers,
          },
        })
      } catch {
        throw new ApiError('无法连接后端服务，请确认 API 端口已启动。', 0, 'NETWORK_ERROR')
      }

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
}

export const apiClient = createApiClient()
