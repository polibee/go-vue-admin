export interface ResourceFormField {
  name: string
  label: string
  type: string
  options?: readonly { value: string; label: string }[]
  visible?: boolean
  readable?: boolean
  writable?: boolean
  sensitive?: boolean
}

function valueForField(field: ResourceFormField, value: unknown) {
	if (field.type === 'boolean') return Boolean(value)
	if (field.type === 'number') return value === null || value === undefined ? '' : Number(value)
	if (field.name === 'status' && (value === null || value === undefined || value === '')) return 'active'
	return value === null || value === undefined ? '' : String(value)
}

export function createResourceForm(fields: readonly ResourceFormField[], record: Record<string, unknown> = {}) {
  return Object.fromEntries(fields.map((field) => [field.name, valueForField(field, record[field.name])]))
}

export function serializeResourceForm(fields: readonly ResourceFormField[], form: Record<string, unknown>) {
  return Object.fromEntries(fields.map((field) => {
    const value = form[field.name]
    if (field.type === 'number') return [field.name, value === '' || value === null || value === undefined ? null : Number(value)]
    if (field.type === 'boolean') return [field.name, Boolean(value)]
    return [field.name, value]
  }))
}
