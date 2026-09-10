import { defineAdminModule } from '../../../admin/src/core/extensions'

export const moduleDefinition = defineAdminModule({
  id: 'example',
  resources: [],
  routes: [],
  navigation: [],
  locales: { 'zh-CN': {}, en: {} },
})
