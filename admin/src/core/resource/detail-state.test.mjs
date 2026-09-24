import assert from 'node:assert/strict'

const { isResourceNotFound } = await import('./detail-state.ts')

assert.equal(isResourceNotFound({ status: 404, code: 'RESOURCE_NOT_FOUND' }), true)
assert.equal(isResourceNotFound({ status: 404, code: undefined }), true)
assert.equal(isResourceNotFound({ status: 403, code: 'RBAC_FORBIDDEN' }), false)
assert.equal(isResourceNotFound({ status: 500, code: 'RESOURCE_NOT_FOUND' }), false)
console.log('detail state tests passed')
