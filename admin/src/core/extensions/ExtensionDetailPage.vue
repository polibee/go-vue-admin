<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
import { useI18n } from 'vue-i18n'
const route = useRoute(); const client = createGeneratedApiClient(apiClient); const extension = ref<ExtensionResource | null>(null); const error = ref('')
const { t } = useI18n()
onMounted(async () => { try {
  const id = String(route.params.id)
  const kind = String(route.meta.extensionKind ?? route.query.kind ?? '')
  extension.value = (kind === 'plugin' ? await client.getPlugin(id) : await client.getModule(id)).data
} catch (cause) { error.value = cause instanceof Error ? cause.message : t('extensions.loadFailed') } })
</script>
<template><section class="flex flex-col gap-4"><h1 class="text-2xl font-semibold tracking-tight">{{ t('extensions.detailTitle') }}</h1><Alert v-if="error" variant="destructive"><AlertTitle>{{ t('extensions.loadFailed') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert><Card v-else-if="extension"><CardHeader><CardTitle>{{ extension.name }}</CardTitle><CardDescription>{{ extension.kind }} · {{ extension.id }}</CardDescription></CardHeader><CardContent>{{ extension.message }}</CardContent></Card></section></template>
