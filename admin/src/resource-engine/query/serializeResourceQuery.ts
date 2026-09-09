import type { ResourceListQuery } from '../core/ResourceDataProvider'

export function serializeResourceQuery(query: ResourceListQuery = {}): string {
  const params = new URLSearchParams()
  if (query.page !== undefined) params.set('page', String(query.page))
  if (query.perPage !== undefined) params.set('per_page', String(query.perPage))
  if (query.search?.trim()) params.set('search', query.search.trim())
  if (query.sort) {
    params.set('sort', query.sort.field)
    params.set('sort_dir', query.sort.direction)
  }

  for (const [field, value] of Object.entries(query.filters ?? {}).sort(([left], [right]) => left.localeCompare(right))) {
    params.set(`filter[${field}]`, value)
  }

  return params.toString()
}
