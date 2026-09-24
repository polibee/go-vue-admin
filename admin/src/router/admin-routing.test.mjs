import assert from 'node:assert/strict'
import test from 'node:test'
import { adminHomePath, adminLoginPath, safeAdminRedirect } from './admin-routing.ts'

test('uses namespaced admin entry paths', () => {
  assert.equal(adminHomePath, '/admin')
  assert.equal(adminLoginPath, '/admin/login')
  assert.equal(safeAdminRedirect('/admin/users'), '/admin/users')
})

test('rejects redirects outside the admin namespace', () => {
  assert.equal(safeAdminRedirect('/users'), '/admin')
  assert.equal(safeAdminRedirect('https://example.com'), '/admin')
  assert.equal(safeAdminRedirect(undefined), '/admin')
})
