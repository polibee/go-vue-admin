import { createI18n } from 'vue-i18n'

import zhApp from './locales/zh-CN/app'
import zhAuth from './locales/zh-CN/auth'
import zhDashboard from './locales/zh-CN/dashboard'
import zhAbout from './locales/zh-CN/about'
import zhExtensions from './locales/zh-CN/extensions'
import enApp from './locales/en/app'
import enAuth from './locales/en/auth'
import enDashboard from './locales/en/dashboard'
import enAbout from './locales/en/about'
import enExtensions from './locales/en/extensions'

export const supportedLocales = ['zh-CN', 'en'] as const
export type SupportedLocale = (typeof supportedLocales)[number]
const localeStorageKey = 'go-vue-admin.locale'

function initialLocale(): SupportedLocale {
  const stored = window.localStorage.getItem(localeStorageKey)
  return stored === 'en' || stored === 'zh-CN' ? stored : 'zh-CN'
}

export const i18n = createI18n({
  legacy: false,
  locale: initialLocale(),
  fallbackLocale: 'en',
  messages: {
    'zh-CN': { ...zhApp, ...zhAuth, ...zhDashboard, ...zhAbout, ...zhExtensions },
    en: { ...enApp, ...enAuth, ...enDashboard, ...enAbout, ...enExtensions },
  },
})

export function setLocale(locale: SupportedLocale): void {
  i18n.global.locale.value = locale
  document.documentElement.lang = locale
  window.localStorage.setItem(localeStorageKey, locale)
}
