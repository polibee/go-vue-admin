import { describe, expect, it } from 'vitest'

import { serializeResourceQuery } from '../query/serializeResourceQuery'
import { nextResourceSort } from '../table/state'

describe('resource query serialization', () => {
  it('serializes the shared backend list-query protocol deterministically', () => {
    expect(serializeResourceQuery({
      page: 2,
      perPage: 10,
      search: ' alpha ',
      sort: { field: 'name', direction: 'desc' },
      filters: { status: 'active' },
    })).toBe('page=2&per_page=10&search=alpha&sort=name&sort_dir=desc&filter%5Bstatus%5D=active')
  })

  it('cycles a sortable column through ascending, descending, and reset states', () => {
    expect(nextResourceSort(undefined, 'name')).toEqual({ field: 'name', direction: 'asc' })
    expect(nextResourceSort({ field: 'name', direction: 'asc' }, 'name')).toEqual({ field: 'name', direction: 'desc' })
    expect(nextResourceSort({ field: 'name', direction: 'desc' }, 'name')).toBeUndefined()
    expect(nextResourceSort({ field: 'status', direction: 'asc' }, 'name')).toEqual({ field: 'name', direction: 'asc' })
  })
})
