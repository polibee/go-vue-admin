import { createI18n } from 'vue-i18n'

import zhApp from './locales/zh-CN/app'
import zhAuth from './locales/zh-CN/auth'
import zhDashboard from './locales/zh-CN/dashboard'
import zhAbout from './locales/zh-CN/about'
import zhExtensions from './locales/zh-CN/extensions'
import zhSettings from './locales/zh-CN/settings'
import zhMedia from './locales/zh-CN/media'
import zhAudit from './locales/zh-CN/audit'
import zhApi from './locales/zh-CN/api'
import zhResources from './locales/zh-CN/resources'
import zhErrors from './locales/zh-CN/errors'
import enApp from './locales/en/app'
import enAuth from './locales/en/auth'
import enDashboard from './locales/en/dashboard'
import enAbout from './locales/en/about'
import enExtensions from './locales/en/extensions'
import enSettings from './locales/en/settings'
import enMedia from './locales/en/media'
import enAudit from './locales/en/audit'
import enApi from './locales/en/api'
import enResources from './locales/en/resources'
import enErrors from './locales/en/errors'

export const supportedLocales = ['zh-CN', 'en'] as const
export type SupportedLocale = (typeof supportedLocales)[number]
const localeStorageKey = 'go-vue-admin.locale'

function initialLocale(): SupportedLocale {
  const stored = typeof window !== 'undefined' ? window.localStorage.getItem(localeStorageKey) : null
  return stored === 'en' || stored === 'zh-CN' ? stored : 'zh-CN'
}

export const i18n = createI18n({
  legacy: false,
  locale: initialLocale(),
  fallbackLocale: 'en',
  messages: {
    'zh-CN': { ...zhApp, ...zhAuth, ...zhDashboard, ...zhAbout, ...zhExtensions, ...zhSettings, ...zhMedia, ...zhAudit, ...zhApi, ...zhResources, ...zhErrors },
    en: { ...enApp, ...enAuth, ...enDashboard, ...enAbout, ...enExtensions, ...enSettings, ...enMedia, ...enAudit, ...enApi, ...enResources, ...enErrors },
  },
})

export function setLocale(locale: SupportedLocale): void {
  i18n.global.locale.value = locale
  if (typeof document !== 'undefined') document.documentElement.lang = locale
  if (typeof window !== 'undefined') window.localStorage.setItem(localeStorageKey, locale)
}
