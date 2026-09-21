import assert from 'node:assert/strict'
import test from 'node:test'
import { canDeleteResource, resourceActionPath } from '../src/lib/resource-actions.ts'

test('builds supported resource action paths', () => {
  assert.equal(resourceActionPath('users', 7, 'delete'), '/api/v1/admin/users/7')
  assert.equal(resourceActionPath('roles', 3, 'delete'), '/api/v1/admin/roles/3')
})

test('protects the built-in super-admin role from list deletion', () => {
  assert.equal(canDeleteResource('roles', { name: 'super-admin' }), false)
  assert.equal(canDeleteResource('roles', { name: 'editor' }), true)
  assert.equal(canDeleteResource('permissions', { name: 'admin.users.view' }), false)
})
