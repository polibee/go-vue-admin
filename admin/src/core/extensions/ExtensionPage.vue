<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient } from '@/generated/api'
import { useI18n } from 'vue-i18n'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'

const props = defineProps<{ id: string; title: string; resource: string }>()
const message = ref('')
const error = ref('')
const { t } = useI18n()
onMounted(async () => {
  try {
    const client = createGeneratedApiClient(apiClient)
    const result = props.resource === 'plugins' ? await client.getPlugin(props.id) : await client.getModule(props.id)
    message.value = result.data.message ?? ''
  }
  catch (cause) { error.value = cause instanceof Error ? cause.message : t('extensions.loadFailed') }
})
</script>
<template>
  <section class="flex flex-col gap-4">
    <h1>{{ title }}</h1>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('extensions.unableToLoad') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-else>
      <CardHeader><CardTitle>{{ title }}</CardTitle><CardDescription>{{ t('extensions.registered') }}</CardDescription></CardHeader>
      <CardContent>{{ message || t('extensions.loading') }}</CardContent>
      <CardFooter><Button as-child><RouterLink :to="`/admin/resources/${resource}`">{{ t('extensions.viewResource') }}</RouterLink></Button></CardFooter>
    </Card>
  </section>
</template>
