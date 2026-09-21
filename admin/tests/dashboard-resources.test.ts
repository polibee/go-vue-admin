import assert from 'node:assert/strict'
import test from 'node:test'
import { dashboardResourceRoute, visibleDashboardResources, type DashboardResource } from '../src/lib/dashboard-resources.ts'

const resources: DashboardResource[] = [
  { name: 'users', label: 'Users', route: '/admin/users', permissions: ['admin.users.view'] },
  { name: 'orders', label: 'Orders', route: '/admin/orders', permissions: ['admin.orders.view'] },
]

test('dashboard only exposes resources covered by the current permissions', () => {
  assert.deepEqual(visibleDashboardResources(resources, ['admin.orders.view']).map((resource) => resource.name), ['orders'])
})

test('dashboard maps resource routes to admin frontend routes', () => {
  assert.equal(dashboardResourceRoute(resources[0]), '/users')
  assert.equal(dashboardResourceRoute({ name: 'roles', route: '/admin/roles' }), '/rbac')
})
