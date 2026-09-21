export function canDeleteResource(resource: string, record: Record<string, unknown>) {
  return Boolean(record.id) && !(resource === 'roles' && record.name === 'super-admin')
}
