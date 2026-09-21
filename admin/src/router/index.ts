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
        { path: 'rbac', name: 'rbac', component: () => import('@/views/RBACView.vue') },
        { path: ':resource(users|roles|permissions)', name: 'resource-list', component: () => import('@/views/ResourceListView.vue') },
        { path: 'users/new', name: 'user-create', component: () => import('@/views/UserFormView.vue') },
        { path: 'users/:id/edit', name: 'user-edit', component: () => import('@/views/UserFormView.vue') },
        { path: 'roles/new', name: 'role-create', component: () => import('@/views/RoleFormView.vue') },
        { path: 'roles/:id/edit', name: 'role-edit', component: () => import('@/views/RoleFormView.vue') },
        { path: ':resource(users|roles|permissions)/:id', name: 'resource-detail', component: () => import('@/views/ResourceDetailView.vue') },
        { path: 'loading', name: 'loading', component: () => import('@/views/LoadingView.vue') },
        { path: 'empty', name: 'empty', component: () => import('@/views/EmptyView.vue') },
        { path: 'error', name: 'error', component: () => import('@/views/ErrorView.vue') },
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
})

export default router
