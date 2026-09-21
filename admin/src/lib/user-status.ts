export const USER_STATUSES = ['active', 'disabled', 'locked'] as const
export type UserStatus = typeof USER_STATUSES[number]

export function userStatusLabelKey(status: UserStatus) {
  return `resource.status${status.charAt(0).toUpperCase()}${status.slice(1)}`
}
