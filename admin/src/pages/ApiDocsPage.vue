<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ApiReference } from '@scalar/api-reference'
import '@scalar/api-reference/style.css'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useI18n } from 'vue-i18n'

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
  { id: 'all', label: 'apiDocs.all', match: () => true },
  { id: 'auth', label: 'apiDocs.auth', match: (path: string) => ['/api/login', '/api/logout', '/api/me', '/api/csrf', '/api/auth/'].some((prefix) => path === prefix || path.startsWith(prefix)) },
  { id: 'resources', label: 'apiDocs.resources', match: (path: string) => path.startsWith('/api/resources/') },
  { id: 'extensions', label: 'apiDocs.extensions', match: (path: string) => ['/api/extensions', '/api/modules', '/api/plugins'].some((prefix) => path === prefix || path.startsWith(prefix + '/')) },
  { id: 'platform', label: 'apiDocs.platform', match: (path: string) => ['/api/health', '/api/menu'].includes(path) },
  { id: 'settings', label: 'apiDocs.settings', match: (path: string) => path.startsWith('/api/settings') },
  { id: 'media', label: 'apiDocs.media', match: (path: string) => path.startsWith('/api/media') },
  { id: 'audit', label: 'apiDocs.audit', match: (path: string) => path.startsWith('/api/audit') },
] as const

const document = ref<OpenApiDocument | null>(null)
const activeCategory = ref('all')
const error = ref('')
const { t } = useI18n()

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
    if (!response.ok) throw new Error(`${t('apiDocs.loadFailed')} (${response.status})`)
    document.value = await response.json() as OpenApiDocument
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('apiDocs.loadFailed')
  }
})
</script>

<template>
  <section class="flex min-h-[calc(100vh-8rem)] flex-col gap-4">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('apiDocs.title') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('apiDocs.description') }}</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>{{ t('apiDocs.unavailable') }}</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <Tabs v-else v-model="activeCategory" class="min-h-0" :key="activeCategory">
      <TabsList class="flex h-auto flex-wrap justify-start gap-1">
        <TabsTrigger v-for="category in categories" :key="category.id" :value="category.id">
          {{ t(category.label) }} ({{ categoryCount(category) }})
        </TabsTrigger>
      </TabsList>
      <div class="min-h-[720px] flex-1 overflow-hidden rounded-lg border bg-background">
        <ApiReference v-if="document" :key="activeCategory" :configuration="configuration" />
        <div v-else class="flex min-h-[720px] items-center justify-center text-sm text-muted-foreground">{{ t('apiDocs.loading') }}</div>
      </div>
    </Tabs>
  </section>
</template>
