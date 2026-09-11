<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ApiReference } from '@scalar/api-reference'
import '@scalar/api-reference/style.css'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'

interface OpenApiDocument {
  openapi?: string
  info?: Record<string, unknown>
  servers?: unknown[]
  paths?: Record<string, unknown>
  components?: Record<string, unknown>
  tags?: unknown[]
  [key: string]: unknown
}

const categories = [
  { id: 'all', label: '全部 API', match: () => true },
  { id: 'auth', label: '认证与会话', match: (path: string) => ['/api/login', '/api/logout', '/api/me', '/api/csrf', '/api/auth/'].some((prefix) => path === prefix || path.startsWith(prefix)) },
  { id: 'resources', label: '业务资源', match: (path: string) => path.startsWith('/api/resources/') },
  { id: 'extensions', label: '模块与插件', match: (path: string) => ['/api/extensions', '/api/modules', '/api/plugins'].some((prefix) => path === prefix || path.startsWith(prefix + '/')) },
  { id: 'platform', label: '平台能力', match: (path: string) => ['/api/health', '/api/menu'].includes(path) },
  { id: 'settings', label: '设置', match: (path: string) => path.startsWith('/api/settings') },
  { id: 'media', label: '媒体', match: (path: string) => path.startsWith('/api/media') },
  { id: 'audit', label: '审计', match: (path: string) => path.startsWith('/api/audit') },
] as const

const document = ref<OpenApiDocument | null>(null)
const activeCategory = ref('all')
const error = ref('')

const activeDefinition = computed(() => categories.find((category) => category.id === activeCategory.value) ?? categories[0])
const filteredDocument = computed(() => {
  if (!document.value) return null
  const paths = Object.fromEntries(Object.entries(document.value.paths ?? {}).filter(([path]) => activeDefinition.value.match(path)))
  return { ...document.value, paths }
})

const categoryCount = (category: (typeof categories)[number]) => Object.keys(document.value?.paths ?? {}).filter((path) => category.match(path)).length

const configuration = computed(() => ({
  spec: { content: filteredDocument.value ?? {} },
  hideClientButton: true,
  hideModels: false,
  hideDownloadButton: true,
  hideTestRequestButton: true,
  hideSearch: false,
  hideDarkModeToggle: true,
}))

onMounted(async () => {
  try {
    const response = await fetch('/api/docs/openapi.json', { credentials: 'include' })
    if (!response.ok) throw new Error(`API 文档加载失败（${response.status}）`)
    document.value = await response.json() as OpenApiDocument
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : 'API 文档加载失败'
  }
})
</script>

<template>
  <section class="flex min-h-[calc(100vh-8rem)] flex-col gap-4">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">API 文档</h1>
      <p class="text-sm text-muted-foreground">按功能分页浏览 OpenAPI 契约；“全部 API”保留完整契约视图。</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>加载失败</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <Tabs v-else v-model="activeCategory" class="min-h-0" :key="activeCategory">
      <TabsList class="flex h-auto flex-wrap justify-start gap-1">
        <TabsTrigger v-for="category in categories" :key="category.id" :value="category.id">
          {{ category.label }}（{{ categoryCount(category) }}）
        </TabsTrigger>
      </TabsList>
      <div class="min-h-[720px] flex-1 overflow-hidden rounded-lg border bg-background">
        <ApiReference v-if="document" :key="activeCategory" :configuration="configuration" />
        <div v-else class="flex min-h-[720px] items-center justify-center text-sm text-muted-foreground">正在加载 API 契约…</div>
      </div>
    </Tabs>
  </section>
</template>
