<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
import ModuleCatalog from './components/ModuleCatalog.vue'

const client = createGeneratedApiClient(apiClient)
const modules = ref<ExtensionResource[]>([])
const error = ref('')

onMounted(async () => {
  try {
    modules.value = (await client.listModules()).data
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '业务模块加载失败'
  }
})
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">业务模块</h1>
      <p class="text-sm text-muted-foreground">业务模块负责产品能力和资源功能。</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>加载失败</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <ModuleCatalog :items="modules" />
  </section>
</template>
