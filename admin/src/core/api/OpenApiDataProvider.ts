import type {
  DemoResource,
  GeneratedApiClient,
  PermissionResource,
  ResourceEnvelope,
  ResourceListQuery as GeneratedResourceListQuery,
  RoleResource,
  UserResource,
} from '@/generated/api'

import type {
  ResourceDataProvider,
  ResourceListQuery,
  ResourceListResult,
} from '@/resource-engine/core/ResourceDataProvider'

export type OpenApiResourceName = 'demo' | 'users' | 'roles' | 'permissions'

/**
 * Adapts the generated OpenAPI operations to the resource engine contract.
 * Resource pages depend on this boundary and never construct HTTP requests.
 */
export class OpenApiDataProvider<T extends object> implements ResourceDataProvider<T> {
  constructor(
    private readonly client: GeneratedApiClient,
    private readonly resourceName: OpenApiResourceName,
  ) {}

  async list(query: ResourceListQuery = {}): Promise<ResourceListResult<T>> {
    const envelope = await this.listOperation(query)
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
    return (await this.itemOperation('get', id)).data
  }

  async create(input: T): Promise<T> {
    return (await this.itemOperation('create', input)).data
  }

  async update(id: string, input: Partial<T>): Promise<T> {
    return (await this.itemOperation('update', id, input)).data
  }

  async delete(id: string): Promise<void> {
    await this.itemOperation('delete', id)
  }

  async bulkDelete(ids: string[]): Promise<void> {
    await this.bulkDeleteOperation({ ids })
  }

  private listOperation(query: ResourceListQuery): Promise<ResourceEnvelope<T[]>> {
    const generatedQuery = query as GeneratedResourceListQuery
    switch (this.resourceName) {
      case 'demo': return this.client.listDemoResources(generatedQuery) as unknown as Promise<ResourceEnvelope<T[]>>
      case 'users': return this.client.listUsers(generatedQuery) as unknown as Promise<ResourceEnvelope<T[]>>
      case 'roles': return this.client.listRoles(generatedQuery) as unknown as Promise<ResourceEnvelope<T[]>>
      case 'permissions': return this.client.listPermissions(generatedQuery) as unknown as Promise<ResourceEnvelope<T[]>>
    }
  }

  private itemOperation(
    action: 'get' | 'create' | 'update' | 'delete',
    idOrInput: string | T,
    input?: Partial<T>,
  ): Promise<ResourceEnvelope<T>> {
    switch (this.resourceName) {
      case 'demo':
        return this.demoItemOperation(action, idOrInput, input)
      case 'users':
        return this.userItemOperation(action, idOrInput, input)
      case 'roles':
        return this.roleItemOperation(action, idOrInput, input)
      case 'permissions':
        return this.permissionItemOperation(action, idOrInput, input)
    }
  }

  private demoItemOperation(action: string, idOrInput: string | T, input?: Partial<T>) {
    if (action === 'get') return this.client.getDemoResource(idOrInput as string) as unknown as Promise<ResourceEnvelope<T>>
    if (action === 'create') return this.client.createDemoResource(idOrInput as unknown as DemoResource) as unknown as Promise<ResourceEnvelope<T>>
    if (action === 'update') return this.client.updateDemoResource(idOrInput as string, input as unknown as Partial<DemoResource>) as unknown as Promise<ResourceEnvelope<T>>
    return this.client.deleteDemoResource(idOrInput as string) as unknown as Promise<ResourceEnvelope<T>>
  }

  private userItemOperation(action: string, idOrInput: string | T, input?: Partial<T>) {
    if (action === 'get') return this.client.getUser(idOrInput as string) as unknown as Promise<ResourceEnvelope<T>>
    if (action === 'create') return this.client.createUser(idOrInput as unknown as UserResource) as unknown as Promise<ResourceEnvelope<T>>
    if (action === 'update') return this.client.updateUser(idOrInput as string, input as unknown as Partial<UserResource>) as unknown as Promise<ResourceEnvelope<T>>
    return this.client.deleteUser(idOrInput as string) as unknown as Promise<ResourceEnvelope<T>>
  }

  private roleItemOperation(action: string, idOrInput: string | T, input?: Partial<T>) {
    if (action === 'get') return this.client.getRole(idOrInput as string) as unknown as Promise<ResourceEnvelope<T>>
    if (action === 'create') return this.client.createRole(idOrInput as unknown as RoleResource) as unknown as Promise<ResourceEnvelope<T>>
    if (action === 'update') return this.client.updateRole(idOrInput as string, input as unknown as Partial<RoleResource>) as unknown as Promise<ResourceEnvelope<T>>
    return this.client.deleteRole(idOrInput as string) as unknown as Promise<ResourceEnvelope<T>>
  }

  private permissionItemOperation(action: string, idOrInput: string | T, input?: Partial<T>) {
    if (action === 'get') return this.client.getPermission(idOrInput as string) as unknown as Promise<ResourceEnvelope<T>>
    if (action === 'create') return this.client.createPermission(idOrInput as unknown as PermissionResource) as unknown as Promise<ResourceEnvelope<T>>
    if (action === 'update') return this.client.updatePermission(idOrInput as string, input as unknown as Partial<PermissionResource>) as unknown as Promise<ResourceEnvelope<T>>
    return this.client.deletePermission(idOrInput as string) as unknown as Promise<ResourceEnvelope<T>>
  }

  private bulkDeleteOperation(input: { ids: string[] }): Promise<ResourceEnvelope<{ deleted: boolean }>> {
    switch (this.resourceName) {
      case 'demo': return this.client.bulkDeleteDemoResources(input)
      case 'users': return this.client.bulkDeleteUsers(input)
      case 'roles': return this.client.bulkDeleteRoles(input)
      case 'permissions': return this.client.bulkDeletePermissions(input)
    }
  }
}
