import { fileURLToPath, URL } from 'node:url'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig, loadEnv } from 'vite'

import { resolveBackendUrl } from './vite-backend'

export default defineConfig(async ({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const backendUrl = await resolveBackendUrl(env)
  console.info(`[vite] API proxy target: ${backendUrl}`)

  return {
    plugins: [vue(), tailwindcss()],
    server: {
      proxy: {
        '/api': backendUrl,
      },
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
  }
})
