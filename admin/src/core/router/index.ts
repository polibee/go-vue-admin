import { createRouter, createWebHistory } from 'vue-router'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/admin/dashboard' },
    { path: '/admin', redirect: '/admin/dashboard' },
    {
      path: '/admin',
      component: () => import('@/components/admin/AdminShell.vue'),
      children: [
        { path: 'dashboard', component: () => import('@/pages/DashboardPage.vue') },
        { path: 'settings', component: () => import('@/pages/SettingsPage.vue') },
      ],
    },
    { path: '/:pathMatch(.*)*', component: () => import('@/pages/errors/NotFoundPage.vue') },
  ],
})
