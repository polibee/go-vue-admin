import { Blocks, FileClock, FileImage, FileCode2, Info, KeyRound, LayoutDashboard, Puzzle, Settings, ShieldCheck, Users } from '@lucide/vue'

import { NavigationRegistry } from './NavigationRegistry'

export const navigationRegistry = new NavigationRegistry([
  { id: 'dashboard', label: '仪表盘', route: '/admin/dashboard', icon: LayoutDashboard, permission: 'dashboard.view' },
  { id: 'settings', label: '设置', route: '/admin/settings', icon: Settings, permission: 'settings.view' },
  { id: 'media', label: '媒体库', route: '/admin/media', icon: FileImage, permission: 'media.view' },
  { id: 'audit', label: '审计日志', route: '/admin/audit', icon: FileClock, permission: 'audit.view' },
  { id: 'api-docs', label: 'API 文档', route: '/admin/api-docs', icon: FileCode2, permission: 'dashboard.view' },
  { id: 'users', label: '用户', route: '/admin/resources/users', icon: Users, permission: 'users.view' },
  { id: 'roles', label: '角色', route: '/admin/resources/roles', icon: ShieldCheck, permission: 'roles.view' },
  { id: 'permissions', label: '权限', route: '/admin/resources/permissions', icon: KeyRound, permission: 'permissions.view' },
  { id: 'modules', label: '业务模块', route: '/admin/modules', icon: Blocks, permission: 'dashboard.view' },
  { id: 'plugins', label: '平台插件', route: '/admin/plugins', icon: Puzzle, permission: 'dashboard.view' },
  { id: 'about', label: '关于项目', route: '/admin/about', icon: Info, permission: 'dashboard.view' },
])

export { NavigationRegistry } from './NavigationRegistry'
export type { NavigationItem } from './NavigationRegistry'
export { NavigationService } from './navigation.service'
