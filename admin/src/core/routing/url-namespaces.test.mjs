import assert from 'node:assert/strict'
import test from 'node:test'
import { adminApiBase, adminRoute, isAdminApiPath, isAdminRoute } from './url-namespaces.ts'

test('builds separate admin page and API namespaces', () => {
  assert.equal(adminRoute('departments'), '/admin/departments')
  assert.equal(adminApiBase('departments'), '/api/v1/admin/departments')
})

test('normalizes resource names and rejects unsafe names', () => {
  assert.equal(adminRoute('/departments/'), '/admin/departments')
  assert.throws(() => adminRoute(''))
  assert.throws(() => adminRoute('../users'))
  assert.throws(() => adminApiBase('users/roles'))
})

test('validates admin page and API paths independently', () => {
  assert.equal(isAdminRoute('/admin'), true)
  assert.equal(isAdminRoute('/admin/users'), true)
  assert.equal(isAdminRoute('/users'), false)
  assert.equal(isAdminRoute('/api/v1/admin/users'), false)
  assert.equal(isAdminApiPath('/api/v1/admin/users'), true)
  assert.equal(isAdminApiPath('/admin/users'), false)
  assert.equal(isAdminApiPath('/api/v1/users'), false)
})
