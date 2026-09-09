import { describe, expect, it } from 'vitest'

import { ApiError, createApiClient, resolveApiBaseUrl } from './client'

describe('ApiClient', () => {
  it('resolves the first reachable local backend port for a static panel', async () => {
    const attempts: string[] = []
    const baseUrl = await resolveApiBaseUrl(async (input) => {
      attempts.push(String(input))
      if (String(input).includes(':3003')) throw new TypeError('connection refused')
      return new Response(JSON.stringify({ data: { status: 'ok' } }), { status: 200 })
    }, '127.0.0.1', ['3003', '3000'])

    expect(baseUrl).toBe('http://127.0.0.1:3000')
    expect(attempts).toEqual([
      'http://127.0.0.1:3003/api/health',
      'http://127.0.0.1:3000/api/health',
    ])
  })

  it('turns a network failure into a recoverable Chinese error', async () => {
    const client = createApiClient(async () => {
      throw new TypeError('Failed to fetch')
    })

    await expect(client.request('/api/health')).rejects.toMatchObject({
      name: 'ApiError',
      status: 0,
      code: 'NETWORK_ERROR',
      message: '无法连接后端服务，请确认 API 端口已启动。',
    })
  })

  it('keeps structured API errors intact', async () => {
    const client = createApiClient(async () => new Response(JSON.stringify({
      error: { code: 'INVALID_CREDENTIALS', message: '邮箱或密码错误' },
    }), { status: 401, headers: { 'Content-Type': 'application/json' } }))

    await expect(client.request('/api/login')).rejects.toEqual(
      new ApiError('邮箱或密码错误', 401, 'INVALID_CREDENTIALS'),
    )
  })
})
