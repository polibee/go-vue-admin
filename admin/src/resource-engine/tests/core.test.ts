import { describe, expect, it } from 'vitest'

import { defineResource } from '../core/ResourceDefinition'
import { MemoryResourceDataProvider } from '../core/ResourceDataProvider'
import { createResourceContext } from '../core/ResourceContext'
import { ResourceRegistry } from '../core/ResourceRegistry'
import { resourceRoutes } from '../core/ResourceRoute'

interface DemoRecord {
  id: string
  name: string
  status: string
}

const users = defineResource<DemoRecord>({
  name: ' users ',
  label: '用户',
  permissions: {
    list: 'users.view',
    update: 'users.update',
  },
})

describe('Resource Engine core', () => {
  it('normalizes resource definitions and rejects duplicate names', () => {
    expect(users.name).toBe('users')
    expect(users.endpoint).toBe('users')
    expect(users.primaryKey).toBe('id')

    const registry = new ResourceRegistry()
    registry.register(users)

    expect(registry.get('users')).toBe(users)
    expect(() => registry.register(defineResource({ name: 'users' }))).toThrow(/already registered/)
  })

  it('provides deterministic list and mutation behavior through the memory provider', async () => {
    const provider = new MemoryResourceDataProvider<DemoRecord>([
      { id: '1', name: 'Beta', status: 'active' },
      { id: '2', name: 'Alpha', status: 'active' },
    ])

    const page = await provider.list({
      page: 1,
      perPage: 1,
      sort: { field: 'name', direction: 'asc' },
      filters: { status: 'active' },
    })

    expect(page.data.map((item) => item.name)).toEqual(['Alpha'])
    expect(page.meta.pagination).toEqual({ page: 1, perPage: 1, total: 2, totalPages: 2 })

    await provider.create({ id: '3', name: 'Gamma', status: 'draft' })
    await provider.update('3', { status: 'active' })
    expect(await provider.get('3')).toMatchObject({ name: 'Gamma', status: 'active' })

    await provider.bulkDelete(['1', '3'])
    await expect(provider.get('1')).rejects.toThrow(/not found/)
  })

  it('keeps context permissions and route metadata close to a resource definition', () => {
    const provider = new MemoryResourceDataProvider<DemoRecord>([])
    const context = createResourceContext(users, provider, ['users.view'])

    expect(context.can('list')).toBe(true)
    expect(context.can('update')).toBe(false)
    expect(resourceRoutes('users').map((route) => route.path)).toEqual([
      '/admin/resources/users',
      '/admin/resources/users/create',
      '/admin/resources/users/:id',
      '/admin/resources/users/:id/edit',
    ])
  })
})
