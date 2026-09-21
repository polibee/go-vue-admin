import assert from 'node:assert/strict'
import test from 'node:test'
import { USER_STATUSES, userStatusLabelKey } from '../src/lib/user-status.ts'

test('exposes the three user statuses and locale keys', () => {
  assert.deepEqual(USER_STATUSES, ['active', 'disabled', 'locked'])
  assert.equal(userStatusLabelKey('active'), 'resource.statusActive')
  assert.equal(userStatusLabelKey('disabled'), 'resource.statusDisabled')
  assert.equal(userStatusLabelKey('locked'), 'resource.statusLocked')
})
