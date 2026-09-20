import assert from 'node:assert/strict'
import test from 'node:test'
import { errorMessageKey } from '../src/lib/api.ts'

test('maps backend error codes to locale keys and falls back for unknown codes', () => {
  assert.equal(errorMessageKey('AUTH_INVALID_CREDENTIALS'), 'errors.AUTH_INVALID_CREDENTIALS')
  assert.equal(errorMessageKey('RBAC_FORBIDDEN'), 'errors.RBAC_FORBIDDEN')
  assert.equal(errorMessageKey('NOT_DEFINED'), 'errors.unknown')
})
