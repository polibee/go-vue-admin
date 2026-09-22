import assert from 'node:assert/strict'
import { createResourceForm, serializeResourceForm } from './resource-form.ts'

const fields = [
  { name: 'name', label: 'Name', type: 'text', writable: true },
  { name: 'status', label: 'Status', type: 'select', writable: true, options: [{ value: 'active', label: 'Active' }] },
  { name: 'id', label: 'ID', type: 'integer', writable: false, visible: false, readable: true },
]

assert.deepEqual(createResourceForm(fields, { name: 'Ada', status: 'active', id: 7 }), {
  name: 'Ada',
  status: 'active',
  id: '7',
})
assert.deepEqual(serializeResourceForm(fields, { name: 'Ada', status: 'active', id: '8' }), {
  name: 'Ada',
  status: 'active',
})
assert.deepEqual(serializeResourceForm(fields, { name: 'Ada', status: 'active' }, (field) => field.name !== 'status'), {
  name: 'Ada',
})
