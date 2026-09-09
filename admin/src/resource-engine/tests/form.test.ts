import { describe, expect, it } from 'vitest'

import { createResourceSchema } from '../form/fieldSchema'

describe('resource form schema', () => {
  it('maps required and numeric fields to validation rules', () => {
    const schema = createResourceSchema([
      { name: 'name', label: '名称', required: true },
      { name: 'count', label: '数量', type: 'number', required: true },
    ])

    expect(schema.safeParse({ name: '', count: 'x' }).success).toBe(false)
    expect(schema.parse({ name: '记录', count: '2' })).toEqual({ name: '记录', count: 2 })
  })
})
