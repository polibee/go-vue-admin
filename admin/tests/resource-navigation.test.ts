import assert from 'node:assert/strict'
import test from 'node:test'
import { groupResourceNavigation, visibleResourceNavigation, type ResourceNavigationItem } from '../src/lib/resource-navigation.ts'

const resources: ResourceNavigationItem[] = [
  { name: 'announcements', label: 'Announcements', admin_route: '/admin/announcements', permissions: ['admin.announcements.view'], navigation: { group: 'business', order: 20 } },
  { name: 'users', label: 'Users', admin_route: '/admin/users', permissions: ['admin.users.view'], navigation: { group: 'system', order: 10 } },
  { name: 'hidden', label: 'Hidden', admin_route: '/admin/hidden', permissions: ['admin.hidden.view'], navigation: { group: 'business', order: 1, hidden: true } },
  { name: 'orders', label: 'Orders', admin_route: '/admin/orders', permissions: ['admin.orders.view'], navigation: { group: 'business', order: 10 } },
]

test('resource navigation filters hidden and unauthorized resources', () => {
  assert.deepEqual(visibleResourceNavigation(resources, ['admin.users.view', 'admin.orders.view']).map((item) => item.name), ['orders', 'users'])
})

test('resource navigation groups and orders resources without a shared list page', () => {
  assert.deepEqual(groupResourceNavigation(resources, ['admin.announcements.view', 'admin.users.view', 'admin.orders.view']), [
    { name: 'system', items: [resources[1]] },
    { name: 'business', items: [resources[3], resources[0]] },
  ])
})

test('resource navigation keeps built-in resources in system group during mixed-version rollout', () => {
  const users = { name: 'users', label: 'Users', admin_route: '/admin/users', permissions: ['admin.users.view'] }
  assert.deepEqual(groupResourceNavigation([users], ['admin.users.view']), [{ name: 'system', items: [users] }])
})
