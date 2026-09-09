import { describe, expect, it } from 'vitest'

import { ApiError, createApiClient } from './client'

describe('ApiClient', () => {
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
