export type ResourceAction = 'list' | 'get' | 'create' | 'update' | 'delete' | 'bulkDelete'

export type ResourcePermissions = Partial<Record<ResourceAction, string>>

export interface ResourceDefinition<T extends object = object> {
  name: string
  label?: string
  endpoint?: string
  primaryKey?: keyof T | string
  permissions?: ResourcePermissions
  columns?: readonly unknown[]
  fields?: readonly unknown[]
  filters?: readonly unknown[]
  actions?: readonly unknown[]
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
