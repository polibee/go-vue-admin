import assert from 'node:assert/strict'
import { existsSync } from 'node:fs'

assert.equal(existsSync(new URL('./pages/', import.meta.url)), false)
assert.equal(existsSync(new URL('./resource.ts', import.meta.url)), true)
