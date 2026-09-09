import { LayoutDashboard, Settings } from '@lucide/vue'

export const navigationItems = [
  { id: 'dashboard', label: '仪表盘', route: '/admin/dashboard', icon: LayoutDashboard },
  { id: 'settings', label: '设置', route: '/admin/settings', icon: Settings },
] as const
