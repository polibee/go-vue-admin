export type Permission = string

export function can(granted: readonly Permission[], required: Permission): boolean {
  const [resource] = required.split('.')
  return granted.some((candidate) => candidate === required || candidate === '*' || candidate === `${resource}.*`)
}

export class PermissionGuard {
  constructor(private readonly granted: readonly Permission[]) {}

  allows(required: Permission): boolean {
    return can(this.granted, required)
  }
}
