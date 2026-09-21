import assert from 'node:assert/strict'
import test from 'node:test'
import { hasPermission, hasAnyPermission } from '../src/lib/permissions.ts'

test('matches exact permissions and any-permission checks', () => {
  const permissions = ['admin.users.view', 'admin.users.manage']
  assert.equal(hasPermission(permissions, 'admin.users.view'), true)
  assert.equal(hasPermission(permissions, 'admin.roles.manage'), false)
  assert.equal(hasAnyPermission(permissions, ['admin.roles.manage', 'admin.users.manage']), true)
})
