import type { DataScope, RolePermissionAssignment } from '@/generated/api'

export interface RolePermissionOption {
  id: number
  name: string
  display_name: string
}

export interface RolePermissionView extends RolePermissionOption {
  assigned: boolean
  scope: DataScope
}

export function isProtectedRole(name: string) {
  return name === 'super-admin'
}

export function selectedPermissionIds(assignments: RolePermissionAssignment[]) {
  return assignments.map((assignment) => assignment.id)
}

export function mergeRolePermissions(
  permissions: RolePermissionOption[],
  assignments: RolePermissionAssignment[],
): RolePermissionView[] {
  const assignmentByID = new Map(assignments.map((assignment) => [assignment.id, assignment]))
  return permissions.map((permission) => {
    const assignment = assignmentByID.get(permission.id)
    return {
      ...permission,
      assigned: Boolean(assignment),
      scope: assignment?.scope || 'all',
    }
  })
}
