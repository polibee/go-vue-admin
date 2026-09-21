import assert from 'node:assert/strict'
import test from 'node:test'
import { generatePassword } from '../src/lib/password-generator.ts'

test('generates a strong password with the requested length', () => {
  const password = generatePassword(20)

  assert.equal(password.length, 20)
  assert.match(password, /[A-Z]/)
  assert.match(password, /[a-z]/)
  assert.match(password, /\d/)
  assert.match(password, /[^A-Za-z\d]/)
})
