<script setup lang="ts">
import { ApiReference } from '@scalar/api-reference'
import '@scalar/api-reference/style.css'

const configuration = {
  url: '/api/docs/openapi.json',
  hideClientButton: true,
  hideModels: false,
  hideDownloadButton: true,
  hideTestRequestButton: true,
  hideSearch: false,
  hideDarkModeToggle: true,
} as const
</script>

<template>
  <section class="flex min-h-[calc(100vh-8rem)] flex-col gap-4">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">API 文档</h1>
      <p class="text-sm text-muted-foreground">只读浏览当前后台平台的 OpenAPI 契约。</p>
    </div>
    <div class="min-h-[720px] flex-1 overflow-hidden rounded-lg border bg-background">
      <ApiReference
        :configuration="configuration"
      />
    </div>
  </section>
</template>
