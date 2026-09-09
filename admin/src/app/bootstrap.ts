import { createApp } from 'vue'
import App from './App.vue'
import { i18n } from '@/core/i18n'
import { router } from '@/core/router'
import { initializeTheme } from '@/core/theme'
import '../assets/index.css'

export function bootstrap(): void {
  initializeTheme()
  createApp(App).use(i18n).use(router).mount('#app')
}
