import assert from 'node:assert/strict'
import test from 'node:test'
import { canDeleteResource } from '../src/lib/resource-actions.ts'

test('protects the built-in super-admin role from list deletion', () => {
  assert.equal(canDeleteResource('roles', { name: 'super-admin' }), false)
  assert.equal(canDeleteResource('roles', { name: 'editor' }), true)
  assert.equal(canDeleteResource('permissions', { name: 'admin.users.view' }), false)
})
