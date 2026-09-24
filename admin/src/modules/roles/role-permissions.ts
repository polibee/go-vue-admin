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

export interface RolePermissionGroup {
  name: string
  permissions: RolePermissionOption[]
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

export function filterPermissionGroups(groups: RolePermissionGroup[], query: string) {
  const normalizedQuery = query.trim().toLocaleLowerCase()
  if (!normalizedQuery) return groups
  return groups
    .map((group) => ({
      ...group,
      permissions: group.permissions.filter((permission) => [permission.name, permission.display_name].some((value) => value.toLocaleLowerCase().includes(normalizedQuery))),
    }))
    .filter((group) => group.permissions.length > 0)
}

export function togglePermissionGroup(selected: number[], visiblePermissionIds: number[], checked: boolean) {
  if (checked) return Array.from(new Set([...selected, ...visiblePermissionIds]))
  const visible = new Set(visiblePermissionIds)
  return selected.filter((id) => !visible.has(id))
}
