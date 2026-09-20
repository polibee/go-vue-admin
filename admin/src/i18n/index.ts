import { createI18n } from 'vue-i18n'
import enUSCore from '@/locales/en-US/core.json'
import enUSAuth from '@/locales/en-US/auth.json'
import zhCNCcore from '@/locales/zh-CN/core.json'
import zhCNAuth from '@/locales/zh-CN/auth.json'
import enUSStates from '@/locales/en-US/states.json'
import enUSRbac from '@/locales/en-US/rbac.json'
import zhCNStates from '@/locales/zh-CN/states.json'
import zhCNRbac from '@/locales/zh-CN/rbac.json'
import enUSErrors from '@/locales/en-US/errors.json'
import zhCNErrors from '@/locales/zh-CN/errors.json'

export const SUPPORTED_LOCALES = ['zh-CN', 'en-US'] as const

const initialLocale = localStorage.getItem('locale') === 'en-US' ? 'en-US' : 'zh-CN'

const i18n = createI18n({
  legacy: false,
  locale: initialLocale,
  fallbackLocale: 'en-US',
  messages: {
    'zh-CN': { core: zhCNCcore, auth: zhCNAuth, errors: zhCNErrors, states: zhCNStates, rbac: zhCNRbac },
    'en-US': { core: enUSCore, auth: enUSAuth, errors: enUSErrors, states: enUSStates, rbac: enUSRbac },
  },
})

export default i18n
