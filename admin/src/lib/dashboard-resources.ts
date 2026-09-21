export interface DashboardResource {
  name: string
  label: string
  route: string
  permissions: string[]
}

export function visibleDashboardResources(resources: DashboardResource[], permissions: string[]) {
  const granted = new Set(permissions)
  return resources.filter((resource) => resource.permissions.some((permission) => granted.has(permission)))
}

export function dashboardResourceRoute(resource: Pick<DashboardResource, 'name' | 'route'>) {
  if (resource.name === 'roles' || resource.name === 'permissions') return '/rbac'
  return resource.route.replace(/^\/admin(?=\/|$)/, '') || '/' + resource.name
}
