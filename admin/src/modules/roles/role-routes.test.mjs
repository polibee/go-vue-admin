import assert from 'node:assert/strict'

const source = await import('./role-routes.ts')

assert.equal(source.roleDetailPath('/admin/roles', 1), '/admin/roles/1')
assert.equal(source.roleDetailPath('/admin/roles', '42'), '/admin/roles/42')
console.log('role route tests passed')
