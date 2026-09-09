export interface ResourceSort {
  field: string
  direction: 'asc' | 'desc'
}

export interface ResourceListQuery {
  page?: number
  perPage?: number
  search?: string
  filters?: Record<string, string>
  sort?: ResourceSort
}

export interface ResourcePagination {
  page: number
  perPage: number
  total: number
  totalPages: number
}

export interface ResourceListResult<T> {
  data: T[]
  meta: {
    pagination: ResourcePagination
  }
}

export interface ResourceDataProvider<T extends object> {
  list(query?: ResourceListQuery): Promise<ResourceListResult<T>>
  get(id: string): Promise<T>
  create(input: T): Promise<T>
  update(id: string, input: Partial<T>): Promise<T>
  delete(id: string): Promise<void>
  bulkDelete(ids: string[]): Promise<void>
}

export class MemoryResourceDataProvider<T extends object> implements ResourceDataProvider<T> {
  private rows: T[]

  constructor(initial: readonly T[] = [], private readonly primaryKey = 'id') {
    this.rows = initial.map((item) => ({ ...item }))
  }

  async list(query: ResourceListQuery = {}): Promise<ResourceListResult<T>> {
    const page = Math.max(1, query.page ?? 1)
    const perPage = Math.max(1, query.perPage ?? 20)
    let rows = this.rows.map((item) => ({ ...item }))

    if (query.search?.trim()) {
      const search = query.search.trim().toLowerCase()
      rows = rows.filter((item) => Object.values(asRecord(item)).some((value) => String(value).toLowerCase().includes(search)))
    }

    for (const [field, expected] of Object.entries(query.filters ?? {})) {
      rows = rows.filter((item) => String(asRecord(item)[field] ?? '') === expected)
    }

    if (query.sort) {
      const { field, direction } = query.sort
      rows.sort((left, right) => {
        const comparison = String(asRecord(left)[field] ?? '').localeCompare(String(asRecord(right)[field] ?? ''))
        return direction === 'desc' ? -comparison : comparison
      })
    }

    const total = rows.length
    const start = (page - 1) * perPage
    const data = rows.slice(start, start + perPage)
    return {
      data,
      meta: {
        pagination: {
          page,
          perPage,
          total,
          totalPages: total === 0 ? 0 : Math.ceil(total / perPage),
        },
      },
    }
  }

  async get(id: string): Promise<T> {
    const item = this.rows.find((row) => String(asRecord(row)[this.primaryKey] ?? '') === id)
    if (!item) throw new Error(`Resource "${id}" not found`)
    return { ...item }
  }

  async create(input: T): Promise<T> {
    const id = String(asRecord(input)[this.primaryKey] ?? '')
    if (!id || this.rows.some((row) => String(asRecord(row)[this.primaryKey] ?? '') === id)) {
      throw new Error(`Resource "${id || 'unknown'}" already exists`)
    }
    const item = { ...input }
    this.rows.push(item)
    return { ...item }
  }

  async update(id: string, input: Partial<T>): Promise<T> {
    const index = this.rows.findIndex((row) => String(asRecord(row)[this.primaryKey] ?? '') === id)
    if (index < 0) throw new Error(`Resource "${id}" not found`)
    this.rows[index] = { ...this.rows[index], ...input }
    return { ...this.rows[index] }
  }

  async delete(id: string): Promise<void> {
    const originalLength = this.rows.length
    this.rows = this.rows.filter((row) => String(asRecord(row)[this.primaryKey] ?? '') !== id)
    if (this.rows.length === originalLength) throw new Error(`Resource "${id}" not found`)
  }

  async bulkDelete(ids: string[]): Promise<void> {
    const wanted = new Set(ids)
    this.rows = this.rows.filter((row) => !wanted.has(String(asRecord(row)[this.primaryKey] ?? '')))
  }
}

function asRecord(value: object): Record<string, unknown> {
  return value as Record<string, unknown>
}
