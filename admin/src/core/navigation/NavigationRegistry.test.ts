import { describe, expect, it } from 'vitest'

import { NavigationRegistry } from './NavigationRegistry'

describe('NavigationRegistry', () => {
  it('merges registered items and returns only items allowed by permissions', () => {
    const registry = new NavigationRegistry()
    registry.register({ id: 'users', label: '用户', route: '/admin/users', permission: 'users.view' })
    registry.register({ id: 'settings', label: '设置', route: '/admin/settings', permission: 'settings.view' })

    expect(registry.visible(['settings.view']).map((item) => item.id)).toEqual(['settings'])
  })

  it('promotes one owned item and groups multiple owned items', () => {
    const registry = new NavigationRegistry()
    registry.register({ id: 'single', label: '单页', route: '/admin/single', owner: 'plugin:single' })
    registry.register({ id: 'orders', label: '订单', route: '/admin/orders', owner: 'plugin:shop', groupLabel: '商城' })
    registry.register({ id: 'products', label: '商品', route: '/admin/products', owner: 'plugin:shop', groupLabel: '商城' })

    const sections = registry.sections([])
    expect(sections.find((section) => section.kind === 'item' && section.item.id === 'single')).toBeDefined()
    expect(sections.find((section) => section.kind === 'group' && section.group.label === '商城')).toMatchObject({
      group: { items: [{ id: 'orders' }, { id: 'products' }] },
    })
  })
})
