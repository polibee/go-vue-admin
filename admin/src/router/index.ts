import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { generatedResourceRoutes } from '@/core/resource/generated'
import ResourceListPage from '@/core/resource/pages/ResourceListPage.vue'
import ResourceDetailPage from '@/core/resource/pages/ResourceDetailPage.vue'
import { adminHomePath, adminLoginPath, safeAdminRedirect } from './admin-routing'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: adminLoginPath,
      name: 'admin-login',
      component: () => import('@/modules/auth/pages/LoginPage.vue'),
    },
    {
      path: adminHomePath,
      component: () => import('@/core/layouts/AdminShell.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'admin-home',
          component: () => import('@/core/pages/AdminHomePage.vue'),
        },
        { path: 'rbac', name: 'admin-rbac', meta: { anyPermissions: ['admin.users.view', 'admin.roles.manage', 'admin.permissions.manage'] }, component: () => import('@/modules/rbac/pages/RBACPage.vue') },
        { path: 'audit-logs', name: 'admin-audit-logs', meta: { permission: 'admin.users.view' }, component: () => import('@/modules/audit/pages/AuditLogPage.vue') },
        ...generatedResourceRoutes,
        { path: ':resource(users|roles|permissions)', name: 'admin-resource-list', component: ResourceListPage },
        { path: 'users/new', name: 'admin-user-create', meta: { permission: 'admin.users.manage' }, component: () => import('@/modules/users/pages/UserFormPage.vue') },
        { path: 'users/:id/edit', name: 'admin-user-edit', meta: { permission: 'admin.users.manage' }, component: () => import('@/modules/users/pages/UserFormPage.vue') },
        { path: 'roles/new', name: 'admin-role-create', meta: { permission: 'admin.roles.manage' }, component: () => import('@/modules/roles/pages/RoleFormPage.vue') },
        { path: 'roles/:id/edit', name: 'admin-role-edit', meta: { permission: 'admin.roles.manage' }, component: () => import('@/modules/roles/pages/RoleFormPage.vue') },
        { path: ':resource(users|roles|permissions)/:id', name: 'admin-resource-detail', component: ResourceDetailPage },
        { path: 'loading', name: 'admin-loading', component: () => import('@/core/pages/LoadingPage.vue') },
        { path: 'empty', name: 'admin-empty', component: () => import('@/core/pages/EmptyPage.vue') },
        { path: 'error', name: 'admin-error', component: () => import('@/core/pages/ErrorPage.vue') },
        { path: 'forbidden', name: 'admin-forbidden', component: () => import('@/core/pages/ForbiddenPage.vue') },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  await auth.restore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'admin-login', query: { redirect: safeAdminRedirect(to.fullPath) } }
  }
  if (to.name === 'admin-login' && auth.isAuthenticated) {
    return { name: 'admin-home' }
  }
  if (to.meta.permission && !auth.can(String(to.meta.permission))) {
    return { name: 'admin-forbidden' }
  }
  if (to.meta.anyPermissions && !auth.canAny(to.meta.anyPermissions as string[])) {
    return { name: 'admin-forbidden' }
  }
  const resource = typeof to.params.resource === 'string' ? to.params.resource : ''
  const resourcePermission: Record<string, string> = { users: 'admin.users.view', roles: 'admin.roles.manage', permissions: 'admin.permissions.manage' }
  if (resource && resourcePermission[resource] && !auth.can(resourcePermission[resource])) {
    return { name: 'forbidden' }
  }
})

export default router
