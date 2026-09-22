import assert from 'node:assert/strict'
import { buildActionSelection } from './resource-selection.ts'

assert.deepEqual(buildActionSelection(['2', '1'], false, [], { status: 'active' }), {
  mode: 'ids',
  ids: [2, 1],
})
assert.deepEqual(buildActionSelection([], true, ['4'], { status: 'active', search: 'ada' }), {
  mode: 'query',
  query: { status: 'active', search: 'ada' },
  exclude_ids: [4],
})
