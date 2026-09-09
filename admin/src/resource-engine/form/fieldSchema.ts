import { z } from 'zod'

import type { ResourceField } from '../core/ResourceDefinition'

export function createResourceSchema<T extends object>(fields: readonly ResourceField<T>[]) {
  const shape: Record<string, z.ZodTypeAny> = {}

  for (const field of fields) {
    let schema: z.ZodTypeAny
    switch (field.type) {
      case 'number':
        schema = z.coerce.number()
        break
      case 'checkbox':
      case 'switch':
        schema = z.coerce.boolean()
        break
      default:
        schema = z.string()
    }

    if (field.required && field.type !== 'checkbox' && field.type !== 'switch') {
      schema = schema.refine((value) => String(value).trim().length > 0, `${field.label}不能为空`)
    }

    shape[String(field.name)] = schema
  }

  return z.object(shape)
}
