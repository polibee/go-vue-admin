import assert from 'node:assert/strict'
import { executableBatchActions } from './resource-actions.ts'

const actions = executableBatchActions([
  { name: 'archive', label: 'Archive', kind: 'archive', permission: 'admin.posts.archive', batch: true, payload: 'archive' },
], (permission) => permission === 'admin.posts.archive')

assert.equal(actions.length, 1)
assert.equal(actions[0].kind, 'archive')
