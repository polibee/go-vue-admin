import { apiClient, type ApiClient } from '@/core/api/client'

export interface AuthUser {
  id: string
  email: string
  name: string
  permissions: string[]
}

export interface LoginCredentials {
  email: string
  password: string
}

interface ApiEnvelope<T> {
  data: T
}

interface CsrfPayload {
  token: string
}

export class AuthService {
  constructor(private readonly client: ApiClient = apiClient) {}

  async login(credentials: LoginCredentials): Promise<AuthUser> {
    const csrf = await this.csrf()
    return this.client
      .request<ApiEnvelope<AuthUser>>('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-TOKEN': csrf,
        },
        body: JSON.stringify(credentials),
      })
      .then(({ data }) => data)
  }

  async me(): Promise<AuthUser> {
    const { data } = await this.client.request<ApiEnvelope<AuthUser>>('/api/me')
    return data
  }

  async logout(): Promise<void> {
    const csrf = await this.csrf()
    await this.client.request('/api/logout', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-TOKEN': csrf,
      },
    })
  }

  private async csrf(): Promise<string> {
    const { data } = await this.client.request<ApiEnvelope<CsrfPayload>>('/api/csrf')
    return data.token
  }
}
