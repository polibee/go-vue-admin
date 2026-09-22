export function canDeleteResource(resource: string, record: Record<string, unknown>) {
  return Boolean(record.id) && !(resource === 'roles' && record.name === 'super-admin')
}

export interface ResourceActionLike { name: string; label: string; kind: string; permission: string }

export function executableBatchActions(actions: ResourceActionLike[] | undefined, can: (permission: string) => boolean) {
  return (actions || []).filter((action) => Boolean(action.kind) && can(action.permission))
}
