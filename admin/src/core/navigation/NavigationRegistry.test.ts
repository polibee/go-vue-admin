import { describe, expect, it } from 'vitest'

import { NavigationRegistry } from './NavigationRegistry'

describe('NavigationRegistry', () => {
  it('merges registered items and returns only items allowed by permissions', () => {
    const registry = new NavigationRegistry()
    registry.register({ id: 'users', label: '用户', route: '/admin/users', permission: 'users.view' })
    registry.register({ id: 'settings', label: '设置', route: '/admin/settings', permission: 'settings.view' })

    expect(registry.visible(['settings.view']).map((item) => item.id)).toEqual(['settings'])
  })
})
