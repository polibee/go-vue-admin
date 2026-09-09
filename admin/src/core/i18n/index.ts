import { createI18n } from 'vue-i18n'

export const supportedLocales = ['zh-CN', 'en'] as const
export type SupportedLocale = (typeof supportedLocales)[number]

export const i18n = createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'en',
  messages: {
    'zh-CN': {
      app: {
        title: 'Go Vue Admin',
        foundation: '官方 shadcn-vue 基础已就绪',
        description: '业务模块将在此平台边界内独立交付。',
      },
    },
    en: {
      app: {
        title: 'Go Vue Admin',
        foundation: 'Official shadcn-vue foundation is ready',
        description: 'Business modules will be delivered independently within this platform boundary.',
      },
    },
  },
})
