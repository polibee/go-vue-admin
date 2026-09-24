import assert from 'node:assert/strict'
import test from 'node:test'
import { isProtectedRole, mergeRolePermissions, selectedPermissionIds } from './role-permissions.ts'

test('builds an editable permission view from the role assignments', () => {
  const permissions = [
    { id: 1, name: 'admin.users.view', display_name: 'View users' },
    { id: 2, name: 'admin.roles.manage', display_name: 'Manage roles' },
  ]
  const assignments = [{ id: 2, name: 'admin.roles.manage', display_name: 'Manage roles', scope: 'own' }]

  assert.deepEqual(mergeRolePermissions(permissions, assignments), [
    { ...permissions[0], assigned: false, scope: 'all' },
    { ...permissions[1], assigned: true, scope: 'own' },
  ])
  assert.deepEqual(selectedPermissionIds(assignments), [2])
})

test('protects the built-in super-admin role from permission edits', () => {
  assert.equal(isProtectedRole('super-admin'), true)
  assert.equal(isProtectedRole('content-admin'), false)
})
