import assert from 'node:assert/strict'

const { getResourceDetailExtension, registerResourceDetailExtension } = await import('./detail-extensions.ts')

const extension = { resource: 'roles', component: {} }
registerResourceDetailExtension(extension)
assert.equal(getResourceDetailExtension('roles'), extension)
assert.throws(() => registerResourceDetailExtension(extension), /already registered/)
console.log('detail extension tests passed')
