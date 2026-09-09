import { can } from '@/core/permissions'

import type { ResourceDataProvider } from './ResourceDataProvider'
import type { ResourceAction, ResourceDefinition } from './ResourceDefinition'

export interface ResourceContext<T extends object> {
  definition: ResourceDefinition<T>
  provider: ResourceDataProvider<T>
  can(action: ResourceAction): boolean
}

export function createResourceContext<T extends object>(
  definition: ResourceDefinition<T>,
  provider: ResourceDataProvider<T>,
  grantedPermissions: readonly string[] = [],
): ResourceContext<T> {
  return {
    definition,
    provider,
    can: (action) => {
      const required = definition.permissions?.[action]
      return !required || can(grantedPermissions, required)
    },
  }
}
