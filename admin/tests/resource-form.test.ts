import assert from 'node:assert/strict'
import test from 'node:test'
import { createResourceForm, serializeResourceForm } from '../src/lib/resource-form.ts'

test('creates typed form values from resource field metadata', () => {
  const form = createResourceForm([
    { name: 'name', label: 'Name', type: 'text' },
    { name: 'age', label: 'Age', type: 'number' },
    { name: 'is_active', label: 'Active', type: 'boolean' },
    { name: 'birthday', label: 'Birthday', type: 'date' },
  ], { name: 'Ada', age: 37, is_active: 1, birthday: null })

  assert.deepEqual(form, { name: 'Ada', age: 37, is_active: true, birthday: '' })
})

test('serializes empty optional values without changing boolean semantics', () => {
  const payload = serializeResourceForm([
    { name: 'name', label: 'Name', type: 'text' },
    { name: 'age', label: 'Age', type: 'number' },
    { name: 'is_active', label: 'Active', type: 'boolean' },
  ], { name: 'Ada', age: '', is_active: false })

  assert.deepEqual(payload, { name: 'Ada', age: null, is_active: false })
})

test('serializes role fields through the same resource form contract', () => {
  const payload = serializeResourceForm([
    { name: 'name', label: 'Name', type: 'text' },
    { name: 'display_name', label: 'Display name', type: 'text' },
  ], { name: 'editor', display_name: 'Editor' })

  assert.deepEqual(payload, { name: 'editor', display_name: 'Editor' })
})

test('defaults a new user resource to active status', () => {
  const form = createResourceForm([
    { name: 'status', label: 'Status', type: 'select', options: [] },
  ])

  assert.equal(form.status, 'active')
})
