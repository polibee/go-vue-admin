import assert from 'node:assert/strict'
import test from 'node:test'
import { DEMO_CREDENTIALS, demoCredentialLabel } from '../src/lib/login-demo.ts'

test('exposes the local demo credentials for the login page', () => {
  assert.deepEqual(DEMO_CREDENTIALS, {
    email: 'admin@example.com',
    password: 'Admin123!',
  })
  assert.equal(demoCredentialLabel('email'), 'auth.demoEmail')
  assert.equal(demoCredentialLabel('password'), 'auth.demoPassword')
})
