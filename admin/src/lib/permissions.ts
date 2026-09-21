export function hasPermission(permissions: string[] | undefined, permission: string) {
  return Boolean(permissions?.includes(permission))
}

export function hasAnyPermission(permissions: string[] | undefined, required: string[]) {
  return required.some((permission) => hasPermission(permissions, permission))
}
