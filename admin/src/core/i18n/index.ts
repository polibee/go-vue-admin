import { createI18n } from 'vue-i18n'

import en from './locales/en'
import zhCN from './locales/zh-CN'

export const supportedLocales = ['zh-CN', 'en'] as const
export type SupportedLocale = (typeof supportedLocales)[number]

const localeStorageKey = 'go-vue-admin.locale'

function initialLocale(): SupportedLocale {
  const stored = window.localStorage.getItem(localeStorageKey)
  if (stored === 'zh-CN' || stored === 'en') return stored
  // Keep the product default stable and let the user opt into English from the shell.
  return 'zh-CN'
}

export const i18n = createI18n({
  legacy: false,
  locale: initialLocale(),
  fallbackLocale: 'en',
  messages: { 'zh-CN': zhCN, en },
})

export function setLocale(locale: SupportedLocale): void {
  i18n.global.locale.value = locale
  window.localStorage.setItem(localeStorageKey, locale)
  document.documentElement.lang = locale
}
