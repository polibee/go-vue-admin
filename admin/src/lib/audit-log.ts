export function formatAuditMetadata(metadata: Record<string, unknown> | null) {
  return metadata ? JSON.stringify(metadata, null, 2) : ''
}
