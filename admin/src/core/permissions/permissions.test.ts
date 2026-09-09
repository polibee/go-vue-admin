import { describe, expect, it } from 'vitest'

import { can, PermissionGuard } from './index'

describe('permissions', () => {
  it('checks exact resource actions and supports a resource wildcard', () => {
    expect(can(['users.view'], 'users.view')).toBe(true)
    expect(can(['users.view'], 'users.update')).toBe(false)
    expect(can(['users.*'], 'users.update')).toBe(true)
  })

  it('guards a callback without making frontend authorization authoritative', () => {
    const guard = new PermissionGuard(['users.view'])

    expect(guard.allows('users.view')).toBe(true)
    expect(guard.allows('users.delete')).toBe(false)
  })
})
