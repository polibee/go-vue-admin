import assert from 'node:assert/strict'
import { findResourceManifest, frontendResourceRoute, resourceApiBase } from './manifest-resolver.ts'

const manifests = [{ name: 'departments', label: 'Departments' }]
assert.deepEqual(findResourceManifest(manifests, 'departments'), manifests[0])
assert.equal(findResourceManifest(manifests, 'missing'), undefined)
assert.equal(frontendResourceRoute({ route: '/admin/departments' }, 'departments'), '/admin/departments')
assert.equal(resourceApiBase('departments'), '/api/v1/admin/departments')
