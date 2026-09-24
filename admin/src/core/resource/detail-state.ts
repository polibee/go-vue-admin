export interface ResourceDetailError {
  status: number
  code?: string
}

export function isResourceNotFound(error: ResourceDetailError) {
  return error.status === 404 && (!error.code || error.code === 'RESOURCE_NOT_FOUND')
}
