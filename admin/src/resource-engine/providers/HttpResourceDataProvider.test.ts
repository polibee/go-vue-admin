import { describe, expect, it } from 'vitest'

import { HttpResourceDataProvider } from './HttpResourceDataProvider'

describe('HttpResourceDataProvider', () => {
  it('maps resource queries to the shared API contract', async () => {
    const calls: Array<{ path: string; options?: RequestInit }> = []
    const provider = new HttpResourceDataProvider<{ id: string; name: string }>(
      {
        request: async <T>(path: string, options?: RequestInit) => {
          calls.push({ path, options })
          return {
            data: [{ id: 'demo-1', name: '示例' }],
            meta: { pagination: { page: 2, per_page: 10, total: 1, total_pages: 1 } },
          } as T
        },
      },
      '/api/resources/demo',
    )

    await expect(provider.list({
      page: 2,
      perPage: 10,
      search: '示例',
      filters: { status: 'active' },
      sort: { field: 'name', direction: 'desc' },
    })).resolves.toEqual({
      data: [{ id: 'demo-1', name: '示例' }],
      meta: { pagination: { page: 2, perPage: 10, total: 1, totalPages: 1 } },
    })

    expect(calls[0]?.path).toBe('/api/resources/demo?page=2&per_page=10&search=%E7%A4%BA%E4%BE%8B&sort=name&sort_dir=desc&filter%5Bstatus%5D=active')
  })

  it('uses the REST resource endpoints for mutations', async () => {
    const calls: Array<{ path: string; options?: RequestInit }> = []
    const provider = new HttpResourceDataProvider<{ id: string; name: string }>(
      {
        request: async <T>(path: string, options?: RequestInit) => {
          calls.push({ path, options })
          return { data: { id: 'demo-1', name: '示例' } } as T
        },
      },
      '/api/resources/demo',
    )

    await provider.create({ id: 'demo-1', name: '示例' })
    await provider.update('demo-1', { name: '更新' })
    await provider.delete('demo-1')
    await provider.bulkDelete(['demo-1', 'demo-2'])

    expect(calls.map(({ path, options }) => [path, options?.method, options?.body])).toEqual([
      ['/api/resources/demo', 'POST', JSON.stringify({ id: 'demo-1', name: '示例' })],
      ['/api/resources/demo/demo-1', 'PUT', JSON.stringify({ name: '更新' })],
      ['/api/resources/demo/demo-1', 'DELETE', undefined],
      ['/api/resources/demo/bulk-delete', 'POST', JSON.stringify({ ids: ['demo-1', 'demo-2'] })],
    ])
  })
})
