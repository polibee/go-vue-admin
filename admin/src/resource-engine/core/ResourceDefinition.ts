export type ResourceAction = 'list' | 'get' | 'create' | 'update' | 'delete' | 'bulkDelete'

export type ResourcePermissions = Partial<Record<ResourceAction, string>>

export type ResourceFieldType = 'text' | 'textarea' | 'number' | 'select' | 'checkbox' | 'switch' | 'date' | 'datetime'

export interface ResourceFieldOption {
  label: string
  value: string
}

export interface ResourceField<T extends object = object> {
  name: keyof T & string
  label: string
  type?: ResourceFieldType
  required?: boolean
  placeholder?: string
  options?: readonly ResourceFieldOption[]
}

export interface ResourceColumn<T extends object = object> {
  key: keyof T & string
  label: string
  sortable?: boolean
}

export interface ResourceFilter<T extends object = object> {
  field: keyof T & string
  label: string
  options?: readonly ResourceFieldOption[]
}

export interface ResourceDefinition<T extends object = object> {
  name: string
  label?: string
  endpoint?: string
  primaryKey?: keyof T | string
  permissions?: ResourcePermissions
  columns?: readonly ResourceColumn<T>[]
  fields?: readonly ResourceField<T>[]
  filters?: readonly ResourceFilter<T>[]
  actions?: readonly ResourceAction[]
}

export function defineResource<T extends object>(definition: ResourceDefinition<T>): ResourceDefinition<T> {
  const name = definition.name.trim()
  if (!/^[a-z][a-z0-9_-]*$/.test(name)) {
    throw new Error(`Invalid resource name: ${definition.name}`)
  }

  return {
    ...definition,
    name,
    label: definition.label?.trim() || name,
    endpoint: definition.endpoint?.trim() || name,
    primaryKey: definition.primaryKey ?? 'id',
    permissions: { ...definition.permissions },
  }
}
