export function canDeleteResource(resource: string, record: Record<string, unknown>) {
  return Boolean(record.id) && !(resource === 'roles' && record.name === 'super-admin')
}

export interface ResourceActionLike { name: string; label: string; kind: string; permission: string; batch: boolean; payload?: string; payload_fields?: Array<{ name: string; label: string; type: 'text' | 'number' | 'boolean' | 'select'; required?: boolean; options?: Array<{ value: string; label: string }> }> }

// Kinds the admin shell renders its own batch dialog for. A kind outside this
// set is only executable when the action declares its payload contract, because
// otherwise there is no dialog to collect parameters and the request would fail
// on the server with an unregistered handler.
const BUILTIN_BATCH_KINDS = new Set([
  'builtin-delete',
  'builtin-update',
  'builtin-restore',
  'builtin-force-delete',
  'user-status',
])

export function executableBatchActions(actions: ResourceActionLike[] | undefined, can: (permission: string) => boolean) {
  return (actions || []).filter((action) =>
    action.batch === true
    && Boolean(action.payload)
    && can(action.permission)
    && (BUILTIN_BATCH_KINDS.has(action.kind) || (action.payload_fields?.length ?? 0) > 0),
  )
}
