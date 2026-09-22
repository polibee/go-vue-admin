import assert from 'node:assert/strict'
import test from 'node:test'
import { executableBatchActions } from '../src/lib/resource-actions.ts'

const actions = [
  { name: 'view', label: 'View', kind: '', permission: 'admin.users.view' },
  { name: 'set-status', label: 'Set status', kind: 'user-status', permission: 'admin.users.manage' },
]

test('only permitted manifest actions with a kind are batch actions', () => {
  assert.deepEqual(executableBatchActions(actions, (permission) => permission === 'admin.users.manage').map((action) => action.name), ['set-status'])
  assert.deepEqual(executableBatchActions(actions, () => false), [])
})
