export function canDeleteResource(resource: string, record: Record<string, unknown>) {
  return Boolean(record.id) && !(resource === 'roles' && record.name === 'super-admin')
}

export interface ResourceActionLike { name: string; label: string; kind: string; permission: string; batch: boolean; payload?: string }

const knownBatchActionKinds = new Set(['user-status'])

export function executableBatchActions(actions: ResourceActionLike[] | undefined, can: (permission: string) => boolean) {
  return (actions || []).filter((action) => action.batch === true && knownBatchActionKinds.has(action.kind) && Boolean(action.payload) && can(action.permission))
}
