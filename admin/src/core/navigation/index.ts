import { LayoutDashboard, Settings } from '@lucide/vue'

import { NavigationRegistry } from './NavigationRegistry'

export const navigationRegistry = new NavigationRegistry([
  { id: 'dashboard', label: '仪表盘', route: '/admin/dashboard', icon: LayoutDashboard, permission: 'dashboard.view' },
  { id: 'settings', label: '设置', route: '/admin/settings', icon: Settings, permission: 'settings.view' },
])

export { NavigationRegistry } from './NavigationRegistry'
export type { NavigationItem } from './NavigationRegistry'
export { NavigationService } from './navigation.service'
