export type AuditMetadata = Record<string, unknown> | string | null

export function formatAuditMetadata(metadata: AuditMetadata) {
  if (!metadata) return ''
  if (typeof metadata !== 'string') return JSON.stringify(metadata, null, 2)
  try {
    return JSON.stringify(JSON.parse(metadata), null, 2)
  } catch {
    return metadata
  }
}
