import { fileURLToPath, URL } from 'node:url'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    proxy: {
      '/csrf': 'http://127.0.0.1:3000',
      '/login': 'http://127.0.0.1:3000',
      '/logout': 'http://127.0.0.1:3000',
      '/me': 'http://127.0.0.1:3000',
      '/api': 'http://127.0.0.1:3000',
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
