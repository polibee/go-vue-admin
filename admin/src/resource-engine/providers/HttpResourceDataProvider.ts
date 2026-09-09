import type { ApiClient } from '@/core/api/client'

import type {
  ResourceDataProvider,
  ResourceListQuery,
  ResourceListResult,
} from '../core/ResourceDataProvider'
import { serializeResourceQuery } from '../query/serializeResourceQuery'

interface ResourceEnvelope<T> {
  data: T
  meta?: {
    pagination?: {
      page: number
      per_page: number
      total: number
      total_pages: number
    }
  }
}

export class HttpResourceDataProvider<T extends object> implements ResourceDataProvider<T> {
  constructor(
    private readonly client: ApiClient,
    private readonly resourcePath: string,
  ) {}

  async list(query: ResourceListQuery = {}): Promise<ResourceListResult<T>> {
    const serialized = serializeResourceQuery(query)
    const path = serialized ? `${this.resourcePath}?${serialized}` : this.resourcePath
    const envelope = await this.client.request<ResourceEnvelope<T[]>>(path)
    const pagination = envelope.meta?.pagination
    return {
      data: envelope.data,
      meta: {
        pagination: {
          page: pagination?.page ?? query.page ?? 1,
          perPage: pagination?.per_page ?? query.perPage ?? 20,
          total: pagination?.total ?? envelope.data.length,
          totalPages: pagination?.total_pages ?? 1,
        },
      },
    }
  }

  async get(id: string): Promise<T> {
    const envelope = await this.client.request<ResourceEnvelope<T>>(this.itemPath(id))
    return envelope.data
  }

  async create(input: T): Promise<T> {
    const envelope = await this.client.request<ResourceEnvelope<T>>(this.resourcePath, this.jsonOptions('POST', input))
    return envelope.data
  }

  async update(id: string, input: Partial<T>): Promise<T> {
    const envelope = await this.client.request<ResourceEnvelope<T>>(this.itemPath(id), this.jsonOptions('PUT', input))
    return envelope.data
  }

  async delete(id: string): Promise<void> {
    await this.client.request(this.itemPath(id), { method: 'DELETE' })
  }

  async bulkDelete(ids: string[]): Promise<void> {
    await this.client.request(`${this.resourcePath}/bulk-delete`, this.jsonOptions('POST', { ids }))
  }

  private itemPath(id: string): string {
    return `${this.resourcePath}/${encodeURIComponent(id)}`
  }

  private jsonOptions(method: string, body: unknown): RequestInit {
    return {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    }
  }
}
