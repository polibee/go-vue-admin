import { defineAdminPlugin } from '../../../admin/src/core/extensions/plugin'
import ExtensionPage from '../../../admin/src/core/extensions/ExtensionPage.vue'

export const examplePlugin = defineAdminPlugin({
  manifest: {
    id: 'example-plugin', name: '示例插件', version: '1.0.0', runtime: 'builtin', uiCompatibility: 'shadcn-vue',
  },
  setup(context) {
    context.registerMenu({ id: 'plugin-example', label: '示例插件', route: '/admin/example-plugin', permission: 'dashboard.view' })
    context.registerRoute({ id: 'plugin-example', path: 'example-plugin', view: ExtensionPage })
  },
})
