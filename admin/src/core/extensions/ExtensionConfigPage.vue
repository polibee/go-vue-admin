<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const client = createGeneratedApiClient(apiClient)
const extension = ref<ExtensionResource | null>(null)
const error = ref('')
const { t } = useI18n()

onMounted(async () => {
  try {
    extension.value = (await client.getPlugin(String(route.params.id))).data
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('extensions.loadFailed')
  }
})
</script>

<template>
  <section class="flex flex-col gap-4">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ extension?.config?.label ?? t('extensions.configTitle') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('extensions.configDescription') }}</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>{{ t('extensions.loadFailed') }}</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <Card v-else-if="extension">
      <CardHeader>
        <CardTitle>{{ extension.name }}</CardTitle>
        <CardDescription>{{ t('extensions.configId') }}: {{ extension.config?.schema ?? t('extensions.noSchema') }}</CardDescription>
      </CardHeader>
      <CardContent class="space-y-3 text-sm">
        <p>{{ t('extensions.status') }}: {{ extension.config?.status === 'configured' ? t('extensions.configured') : t('extensions.notConfigured') }}</p>
        <p v-if="extension.config?.secret_fields?.length" class="text-muted-foreground">
          {{ t('extensions.secretFields') }}: {{ extension.config.secret_fields.join(', ') }} ({{ t('extensions.secretHint') }})
        </p>
        <p class="text-muted-foreground">{{ t('extensions.configFormHint') }}</p>
      </CardContent>
    </Card>
  </section>
</template>
