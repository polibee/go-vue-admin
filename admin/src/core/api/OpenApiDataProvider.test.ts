import { describe, expect, it, vi } from 'vitest'

import type { GeneratedApiClient } from '@/generated/api'

import { OpenApiDataProvider } from './OpenApiDataProvider'

describe('OpenApiDataProvider', () => {
  it('adapts generated users list/create/update/delete operations', async () => {
    const client = {
      listUsers: vi.fn().mockResolvedValue({
        data: [{ id: 'user-1', email: 'admin@example.com', name: 'Admin', active: true, role_ids: [] }],
        meta: { pagination: { page: 2, per_page: 10, total: 11, total_pages: 2 } },
      }),
      getUser: vi.fn().mockResolvedValue({ data: { id: 'user-1' } }),
      createUser: vi.fn().mockResolvedValue({ data: { id: 'user-2' } }),
      updateUser: vi.fn().mockResolvedValue({ data: { id: 'user-1', name: 'Updated' } }),
      deleteUser: vi.fn().mockResolvedValue({ data: { deleted: true } }),
      bulkDeleteUsers: vi.fn().mockResolvedValue({ data: { deleted: true } }),
    } as unknown as GeneratedApiClient
    const provider = new OpenApiDataProvider(client, 'users')
    const input = { id: 'user-2', email: 'new@example.com', name: 'New', active: true, role_ids: [] }

    await expect(provider.list({ page: 2, perPage: 10, search: 'admin', filters: { active: 'true' } })).resolves.toMatchObject({
      data: [{ id: 'user-1' }],
      meta: { pagination: { page: 2, perPage: 10, total: 11, totalPages: 2 } },
    })
    await expect(provider.get('user-1')).resolves.toMatchObject({ id: 'user-1' })
    await expect(provider.create(input)).resolves.toMatchObject({ id: 'user-2' })
    await expect(provider.update('user-1', { name: 'Updated' })).resolves.toMatchObject({ name: 'Updated' })
    await expect(provider.delete('user-1')).resolves.toBeUndefined()
    await expect(provider.bulkDelete(['user-1', 'user-2'])).resolves.toBeUndefined()

    expect(client.listUsers).toHaveBeenCalledWith({ page: 2, perPage: 10, search: 'admin', filters: { active: 'true' } })
    expect(client.getUser).toHaveBeenCalledWith('user-1')
    expect(client.createUser).toHaveBeenCalledWith(input)
    expect(client.updateUser).toHaveBeenCalledWith('user-1', { name: 'Updated' })
    expect(client.deleteUser).toHaveBeenCalledWith('user-1')
    expect(client.bulkDeleteUsers).toHaveBeenCalledWith({ ids: ['user-1', 'user-2'] })
  })
})
