export type ResourceAction = 'delete'

export function resourceActionPath(resource: string, id: number | string, action: ResourceAction) {
  if (action === 'delete' && (resource === 'users' || resource === 'roles')) return `/api/v1/admin/${resource}/${id}`
  throw new Error(`Unsupported resource action: ${resource}/${action}`)
}

export function canDeleteResource(resource: string, record: Record<string, unknown>) {
  return (resource === 'users' || resource === 'roles') && !(resource === 'roles' && record.name === 'super-admin')
}
