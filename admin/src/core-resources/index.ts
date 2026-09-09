import type { ResourceDefinition } from '@/resource-engine/core/ResourceDefinition'

import { permissionResource } from './permissions/resource'
import { roleResource } from './roles/resource'
import { userResource } from './users/resource'

export const coreResourceDefinitions: ResourceDefinition<object>[] = [
  userResource as ResourceDefinition<object>,
  roleResource as ResourceDefinition<object>,
  permissionResource as ResourceDefinition<object>,
]

export { permissionResource } from './permissions/resource'
export { roleResource } from './roles/resource'
export { userResource } from './users/resource'
