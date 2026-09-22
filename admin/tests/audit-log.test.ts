import assert from 'node:assert/strict'
import test from 'node:test'
import { formatAuditMetadata } from '../src/lib/audit-log.ts'

test('formats JSON string audit metadata as readable JSON', () => {
  const formatted = formatAuditMetadata('{"request":{"body":{"password":"[REDACTED]"}}}')
  assert.match(formatted, /"request"/)
  assert.match(formatted, /\[REDACTED\]/)
  assert.doesNotMatch(formatted, /^\"\{/)
})

test('returns empty text for null audit metadata', () => {
  assert.equal(formatAuditMetadata(null), '')
})
