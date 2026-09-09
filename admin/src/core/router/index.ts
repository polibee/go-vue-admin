import { createRouter, createWebHistory } from 'vue-router'
import { AuthGuard, GuestGuard } from '@/core/auth'
import { resourceRouteRecords } from '@/resource-engine/router'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/admin/dashboard' },
    { path: '/login', component: () => import('@/pages/auth/LoginPage.vue'), beforeEnter: GuestGuard },
    { path: '/admin', redirect: '/admin/dashboard' },
    {
      path: '/admin',
      component: () => import('@/components/admin/AdminShell.vue'),
      beforeEnter: AuthGuard,
      children: [
        { path: 'dashboard', component: () => import('@/pages/DashboardPage.vue') },
        { path: 'settings', component: () => import('@/pages/SettingsPage.vue') },
        ...resourceRouteRecords(),
      ],
    },
    { path: '/:pathMatch(.*)*', component: () => import('@/pages/errors/NotFoundPage.vue') },
  ],
})
