import type { ResourceSort } from '../core/ResourceDataProvider'

export function nextResourceSort(current: ResourceSort | undefined, field: string): ResourceSort | undefined {
  if (current?.field !== field) return { field, direction: 'asc' }
  if (current.direction === 'asc') return { field, direction: 'desc' }
  return undefined
}
