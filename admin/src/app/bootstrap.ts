import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { i18n } from '@/core/i18n'
import { router } from '@/core/router'
import { initializeTheme } from '@/core/theme'
import '../assets/index.css'

export function bootstrap(): void {
  initializeTheme()
  const app = createApp(App)
  app.use(createPinia()).use(i18n).use(router).mount('#app')
}
