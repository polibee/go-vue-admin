export interface ResourceSelection {
  mode: 'ids' | 'query'
  ids?: number[]
  query?: Record<string, string>
  exclude_ids?: number[]
}

export function buildActionSelection(selectedIds: string[], allFilteredSelected: boolean, excludedIds: string[], query: Record<string, string>): ResourceSelection {
  if (allFilteredSelected) {
    return { mode: 'query', query, exclude_ids: excludedIds.map(Number) }
  }
  return { mode: 'ids', ids: selectedIds.map(Number) }
}
