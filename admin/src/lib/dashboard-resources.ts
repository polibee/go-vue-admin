export interface DashboardResource {
  name: string
  label: string
  admin_route: string
  permissions: string[]
}

export function visibleDashboardResources(resources: DashboardResource[], permissions: string[]) {
  const granted = new Set(permissions)
  return resources.filter((resource) => resource.permissions.some((permission) => granted.has(permission)))
}

export function dashboardResourceRoute(resource: Pick<DashboardResource, 'name' | 'admin_route'>) {
  return resource.admin_route.startsWith('/admin/') ? resource.admin_route : adminRoute(resource.name)
}
import { adminRoute } from '../core/routing/url-namespaces.ts'
