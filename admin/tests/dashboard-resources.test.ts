import assert from 'node:assert/strict'
import test from 'node:test'
import { dashboardResourceRoute, visibleDashboardResources, type DashboardResource } from '../src/lib/dashboard-resources.ts'

const resources: DashboardResource[] = [
  { name: 'users', label: 'Users', admin_route: '/admin/users', permissions: ['admin.users.view'] },
  { name: 'orders', label: 'Orders', admin_route: '/admin/orders', permissions: ['admin.orders.view'] },
]

test('dashboard only exposes resources covered by the current permissions', () => {
  assert.deepEqual(visibleDashboardResources(resources, ['admin.orders.view']).map((resource) => resource.name), ['orders'])
})

test('dashboard maps resource routes to admin frontend routes', () => {
  assert.equal(dashboardResourceRoute(resources[0]), '/admin/users')
  assert.equal(dashboardResourceRoute({ name: 'roles', admin_route: '/admin/roles' }), '/admin/roles')
})
