const RESOURCE_NAME = /^[a-z][a-z0-9-]*$/

function normalizeResourceName(resource: string) {
  const normalized = resource.replace(/^\/+|\/+$/g, '')
  if (!RESOURCE_NAME.test(normalized)) {
    throw new Error(`Invalid resource name: ${resource}`)
  }
  return normalized
}

export function adminRoute(resource: string) {
  return `/admin/${normalizeResourceName(resource)}`
}

export function adminApiBase(resource: string) {
  return `/api/v1/admin/${normalizeResourceName(resource)}`
}

export function isAdminRoute(path: string) {
  return path === '/admin' || path.startsWith('/admin/')
}

export function isAdminApiPath(path: string) {
  return path.startsWith('/api/v1/admin/')
}
