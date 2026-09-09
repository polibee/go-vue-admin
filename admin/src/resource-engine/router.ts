import type { RouteRecordRaw } from 'vue-router'

export function resourceRouteRecords(): RouteRecordRaw[] {
  return [
    { path: 'resources/:resource', component: () => import('./pages/ResourceListPage.vue') },
    { path: 'resources/:resource/create', component: () => import('./pages/ResourceCreatePage.vue') },
    { path: 'resources/:resource/:id/edit', component: () => import('./pages/ResourceEditPage.vue') },
    { path: 'resources/:resource/:id', component: () => import('./pages/ResourceShowPage.vue') },
  ]
}
