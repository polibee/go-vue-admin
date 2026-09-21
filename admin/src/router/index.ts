import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
    },
    {
      path: '/',
      component: () => import('@/layouts/AdminShell.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('@/views/AdminHomeView.vue'),
        },
        { path: 'rbac', name: 'rbac', meta: { anyPermissions: ['admin.users.view', 'admin.roles.manage', 'admin.permissions.manage'] }, component: () => import('@/views/RBACView.vue') },
        { path: 'audit-logs', name: 'audit-logs', meta: { permission: 'admin.users.view' }, component: () => import('@/views/AuditLogView.vue') },
        { path: ':resource(users|roles|permissions)', name: 'resource-list', component: () => import('@/views/ResourceListView.vue') },
        { path: 'users/new', name: 'user-create', meta: { permission: 'admin.users.manage' }, component: () => import('@/views/UserFormView.vue') },
        { path: 'users/:id/edit', name: 'user-edit', meta: { permission: 'admin.users.manage' }, component: () => import('@/views/UserFormView.vue') },
        { path: 'roles/new', name: 'role-create', meta: { permission: 'admin.roles.manage' }, component: () => import('@/views/RoleFormView.vue') },
        { path: 'roles/:id/edit', name: 'role-edit', meta: { permission: 'admin.roles.manage' }, component: () => import('@/views/RoleFormView.vue') },
        { path: ':resource(users|roles|permissions)/:id', name: 'resource-detail', component: () => import('@/views/ResourceDetailView.vue') },
        { path: 'loading', name: 'loading', component: () => import('@/views/LoadingView.vue') },
        { path: 'empty', name: 'empty', component: () => import('@/views/EmptyView.vue') },
        { path: 'error', name: 'error', component: () => import('@/views/ErrorView.vue') },
        { path: 'forbidden', name: 'forbidden', component: () => import('@/views/ForbiddenView.vue') },
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
