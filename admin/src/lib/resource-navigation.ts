import type { ResourceManifest } from '@/generated/api'

export type ResourceNavigationItem = Pick<ResourceManifest, 'name' | 'label' | 'admin_route' | 'permissions'> & {
  navigation?: {
    group?: string
    order?: number
    hidden?: boolean
  }
}

export interface ResourceNavigationGroup {
  name: string
  items: ResourceNavigationItem[]
}

function navigationGroup(resource: ResourceNavigationItem) {
  if (resource.navigation?.group) return resource.navigation.group
  return ['users', 'roles', 'permissions'].includes(resource.name) ? 'system' : 'business'
}

export function visibleResourceNavigation(resources: ResourceNavigationItem[], permissions: string[]) {
  const granted = new Set(permissions)
  return resources
    .filter((resource) => !resource.navigation?.hidden && resource.permissions.some((permission) => granted.has(permission)))
    .sort((left, right) => {
      const order = (left.navigation?.order ?? 100) - (right.navigation?.order ?? 100)
      return order || left.label.localeCompare(right.label)
    })
}

export function groupResourceNavigation(resources: ResourceNavigationItem[], permissions: string[]): ResourceNavigationGroup[] {
  const groups = new Map<string, ResourceNavigationItem[]>()
  for (const resource of visibleResourceNavigation(resources, permissions)) {
    const name = navigationGroup(resource)
    const items = groups.get(name) || []
    items.push(resource)
    groups.set(name, items)
  }
  return [...groups.entries()]
    .sort(([left], [right]) => {
      const rank = (name: string) => name === 'system' ? 0 : name === 'business' ? 10 : 20
      return rank(left) - rank(right) || left.localeCompare(right)
    })
    .map(([name, items]) => ({ name, items }))
}
