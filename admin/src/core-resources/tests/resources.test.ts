import { describe, expect, it } from 'vitest'

import { coreResourceDefinitions } from '../index'

describe('core resource definitions', () => {
  it('declares users, roles, and permissions for the generic CRUD engine', () => {
    expect(coreResourceDefinitions.map((resource) => resource.name)).toEqual(['users', 'roles', 'permissions'])
    expect(coreResourceDefinitions.find((resource) => resource.name === 'users')?.permissions).toMatchObject({
      list: 'users.view',
      create: 'users.create',
      update: 'users.update',
      delete: 'users.delete',
    })
    expect(coreResourceDefinitions.find((resource) => resource.name === 'roles')?.fields?.map((field) => field.name)).toContain('permissions')
    expect(coreResourceDefinitions.find((resource) => resource.name === 'permissions')?.fields?.map((field) => field.name)).toContain('name')
  })
})
