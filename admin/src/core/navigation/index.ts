import { KeyRound, LayoutDashboard, Settings, ShieldCheck, Users } from '@lucide/vue'

import { NavigationRegistry } from './NavigationRegistry'

export const navigationRegistry = new NavigationRegistry([
  { id: 'dashboard', label: '仪表盘', route: '/admin/dashboard', icon: LayoutDashboard, permission: 'dashboard.view' },
  { id: 'settings', label: '设置', route: '/admin/settings', icon: Settings, permission: 'settings.view' },
  { id: 'users', label: '用户', route: '/admin/resources/users', icon: Users, permission: 'users.view' },
  { id: 'roles', label: '角色', route: '/admin/resources/roles', icon: ShieldCheck, permission: 'roles.view' },
  { id: 'permissions', label: '权限', route: '/admin/resources/permissions', icon: KeyRound, permission: 'permissions.view' },
])

export { NavigationRegistry } from './NavigationRegistry'
export type { NavigationItem } from './NavigationRegistry'
export { NavigationService } from './navigation.service'
