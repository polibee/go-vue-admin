import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { generatedResourceRoutes } from '@/core/resource/generated'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/modules/auth/pages/LoginPage.vue'),
    },
    {
      path: '/',
      component: () => import('@/core/layouts/AdminShell.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('@/core/pages/AdminHomePage.vue'),
        },
        { path: 'rbac', name: 'rbac', meta: { anyPermissions: ['admin.users.view', 'admin.roles.manage', 'admin.permissions.manage'] }, component: () => import('@/modules/rbac/pages/RBACPage.vue') },
        { path: 'audit-logs', name: 'audit-logs', meta: { permission: 'admin.users.view' }, component: () => import('@/modules/audit/pages/AuditLogPage.vue') },
        ...generatedResourceRoutes,
        { path: ':resource(users|roles|permissions)', name: 'resource-list', component: () => import('@/core/resource/pages/ResourceListPage.vue') },
        { path: 'users/new', name: 'user-create', meta: { permission: 'admin.users.manage' }, component: () => import('@/modules/users/pages/UserFormPage.vue') },
        { path: 'users/:id/edit', name: 'user-edit', meta: { permission: 'admin.users.manage' }, component: () => import('@/modules/users/pages/UserFormPage.vue') },
        { path: 'roles/new', name: 'role-create', meta: { permission: 'admin.roles.manage' }, component: () => import('@/modules/roles/pages/RoleFormPage.vue') },
        { path: 'roles/:id/edit', name: 'role-edit', meta: { permission: 'admin.roles.manage' }, component: () => import('@/modules/roles/pages/RoleFormPage.vue') },
        { path: ':resource(users|roles|permissions)/:id', name: 'resource-detail', component: () => import('@/core/resource/pages/ResourceDetailPage.vue') },
        { path: 'loading', name: 'loading', component: () => import('@/core/pages/LoadingPage.vue') },
        { path: 'empty', name: 'empty', component: () => import('@/core/pages/EmptyPage.vue') },
        { path: 'error', name: 'error', component: () => import('@/core/pages/ErrorPage.vue') },
        { path: 'forbidden', name: 'forbidden', component: () => import('@/core/pages/ForbiddenPage.vue') },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  await auth.restore()
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.isAuthenticated) {
    return { name: 'home' }
  }
  if (to.meta.permission && !auth.can(String(to.meta.permission))) {
    return { name: 'forbidden' }
  }
  if (to.meta.anyPermissions && !auth.canAny(to.meta.anyPermissions as string[])) {
    return { name: 'forbidden' }
  }
  const resource = typeof to.params.resource === 'string' ? to.params.resource : ''
  const resourcePermission: Record<string, string> = { users: 'admin.users.view', roles: 'admin.roles.manage', permissions: 'admin.permissions.manage' }
  if (resource && resourcePermission[resource] && !auth.can(resourcePermission[resource])) {
    return { name: 'forbidden' }
  }
})

export default router
