export interface ResourceFormField {
  name: string
  label: string
  type: string
  required?: boolean
  options?: readonly { value: string; label: string }[]
  visible?: boolean
  readable?: boolean
  writable?: boolean
  sensitive?: boolean
}

function valueForField(field: ResourceFormField, value: unknown) {
	if (field.type === 'boolean') return Boolean(value)
	if (field.type === 'number' || field.type === 'integer') return value === null || value === undefined || value === '' ? '' : Number(value)
	if (field.name === 'status' && (value === null || value === undefined || value === '')) return 'active'
	return value === null || value === undefined ? '' : String(value)
}

export function createResourceForm(fields: readonly ResourceFormField[], record: Record<string, unknown> = {}) {
  return Object.fromEntries(fields.map((field) => [field.name, valueForField(field, record[field.name])]))
}

export function serializeResourceForm(fields: readonly ResourceFormField[], form: Record<string, unknown>, includeField: (field: ResourceFormField) => boolean = () => true) {
  return Object.fromEntries(fields.filter((field) => field.writable !== false && includeField(field)).map((field) => {
    const value = form[field.name]
    if (field.type === 'number' || field.type === 'integer') return [field.name, value === '' || value === null || value === undefined ? null : Number(value)]
    if (field.type === 'boolean') return [field.name, Boolean(value)]
    return [field.name, value]
  }))
}
