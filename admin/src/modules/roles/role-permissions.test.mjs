import assert from 'node:assert/strict'
import test from 'node:test'
import { filterPermissionGroups, isProtectedRole, mergeRolePermissions, selectedPermissionIds, togglePermissionGroup } from './role-permissions.ts'

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

test('filters permission groups by display name and permission name', () => {
  const groups = [
    { name: 'users', permissions: [{ id: 1, name: 'admin.users.view', display_name: 'View users' }] },
    { name: 'roles', permissions: [{ id: 2, name: 'admin.roles.manage', display_name: 'Manage roles' }] },
  ]

  assert.equal(filterPermissionGroups(groups, 'users')[0].permissions.length, 1)
  assert.equal(filterPermissionGroups(groups, 'admin.roles.manage')[0].permissions.length, 1)
  assert.equal(filterPermissionGroups(groups, 'missing').length, 0)
})

test('selects and clears only visible permissions in one group', () => {
  assert.deepEqual(togglePermissionGroup([1, 9], [1, 2], true), [1, 9, 2])
  assert.deepEqual(togglePermissionGroup([1, 2, 9], [1, 2], false), [9])
})
