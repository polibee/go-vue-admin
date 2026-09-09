export type ResourceRouteAction = 'list' | 'create' | 'show' | 'edit'

export interface ResourceRoute {
  resource: string
  action: ResourceRouteAction
  name: string
  path: string
}

export function resourceRoutes(resource: string): ResourceRoute[] {
  const name = resource.trim()
  return [
    { resource: name, action: 'list', name: `${name}.list`, path: `/admin/resources/${name}` },
    { resource: name, action: 'create', name: `${name}.create`, path: `/admin/resources/${name}/create` },
    { resource: name, action: 'show', name: `${name}.show`, path: `/admin/resources/${name}/:id` },
    { resource: name, action: 'edit', name: `${name}.edit`, path: `/admin/resources/${name}/:id/edit` },
  ]
}
