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
      name: 'admin',
      component: () => import('@/components/admin/AdminShell.vue'),
      beforeEnter: AuthGuard,
      children: [
        { path: 'dashboard', component: () => import('@/pages/DashboardPage.vue') },
        { path: 'settings', component: () => import('@/pages/SettingsPage.vue') },
        { path: 'media', component: () => import('@/pages/MediaPage.vue') },
        { path: 'audit', component: () => import('@/pages/AuditPage.vue') },
        { path: 'api-docs', component: () => import('@/pages/ApiDocsPage.vue') },
        { path: 'about', component: () => import('@/pages/AboutPage.vue') },
        { path: 'modules', component: () => import('@/core/extensions/ModulesPage.vue') },
        { path: 'modules/:id', component: () => import('@/core/extensions/ExtensionDetailPage.vue'), meta: { extensionKind: 'module' } },
        { path: 'plugins', component: () => import('@/core/extensions/PluginsPage.vue') },
        { path: 'plugins/:id', component: () => import('@/core/extensions/ExtensionDetailPage.vue'), meta: { extensionKind: 'plugin' } },
        { path: 'plugins/:id/config', component: () => import('@/core/extensions/ExtensionConfigPage.vue'), meta: { extensionKind: 'plugin' } },
        { path: 'extensions', redirect: '/admin/modules' },
        { path: 'extensions/:id', component: () => import('@/core/extensions/ExtensionDetailPage.vue') },
        ...resourceRouteRecords(),
      ],
    },
    { path: '/:pathMatch(.*)*', component: () => import('@/pages/errors/NotFoundPage.vue') },
  ],
})
