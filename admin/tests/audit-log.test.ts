import assert from 'node:assert/strict'
import test from 'node:test'
import { formatAuditMetadata } from '../src/lib/audit-log.ts'

test('formats audit metadata as readable JSON and handles empty metadata', () => {
  assert.equal(formatAuditMetadata({ method: 'password', nested: { ok: true } }), '{\n  "method": "password",\n  "nested": {\n    "ok": true\n  }\n}')
  assert.equal(formatAuditMetadata(null), '')
})
