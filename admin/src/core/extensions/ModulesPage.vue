<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
import ModuleCatalog from './components/ModuleCatalog.vue'

const client = createGeneratedApiClient(apiClient)
const modules = ref<ExtensionResource[]>([])
const error = ref('')
const { t } = useI18n()

onMounted(async () => {
  try {
    modules.value = (await client.listModules()).data
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('extensions.loadFailed')
  }
})
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('extensions.moduleTitle') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('extensions.moduleDescription') }}</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>{{ t('extensions.loadFailed') }}</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <ModuleCatalog :items="modules" />
  </section>
</template>
