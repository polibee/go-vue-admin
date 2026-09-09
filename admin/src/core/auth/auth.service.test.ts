import { describe, expect, it } from 'vitest'

import { AuthService } from './auth.service'

describe('AuthService', () => {
  it('loads CSRF before login and sends browser credentials', async () => {
    const calls: Array<{ method: string; path: string; body?: unknown; headers?: HeadersInit }> = []
    const service = new AuthService({
      request: async <T>(path: string, options?: RequestInit) => {
        calls.push({ method: options?.method ?? 'GET', path, body: options?.body, headers: options?.headers })
        if (path === '/csrf') return { data: { token: 'csrf-token' } } as T
        return { data: { id: 'user-1', email: 'admin@example.com', name: 'Platform Admin' } } as T
      },
    })

    const user = await service.login({ email: 'admin@example.com', password: 'secret' })

    expect(user.email).toBe('admin@example.com')
    expect(calls).toEqual([
      { method: 'GET', path: '/csrf', body: undefined, headers: undefined },
      {
        method: 'POST',
        path: '/login',
        body: JSON.stringify({ email: 'admin@example.com', password: 'secret' }),
        headers: { 'Content-Type': 'application/json', 'X-CSRF-TOKEN': 'csrf-token' },
      },
    ])
  })

  it('uses the session-backed endpoints for current user and logout', async () => {
    const paths: string[] = []
    const service = new AuthService({
      request: async <T>(path: string) => {
        paths.push(path)
        if (path === '/csrf') return { data: { token: 'csrf-token' } } as T
        return { data: { id: 'user-1', email: 'admin@example.com', name: 'Platform Admin' } } as T
      },
    })

    await service.me()
    await service.logout()

    expect(paths).toEqual(['/me', '/csrf', '/logout'])
  })
})
