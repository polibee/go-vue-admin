import { defineAdminModule } from '../../../admin/src/core/extensions'
import ExtensionPage from '../../../admin/src/core/extensions/ExtensionPage.vue'

export const moduleDefinition = defineAdminModule({
  id: 'example',
  routes: [{ path: 'example', name: 'module-example', component: ExtensionPage, props: { id: 'example', title: '示例模块', resource: 'example-info' } }],
  navigation: [{ id: 'module-example', label: '示例模块', route: '/admin/example', permission: 'dashboard.view' }],
  locales: { 'zh-CN': {}, en: {} },
})
