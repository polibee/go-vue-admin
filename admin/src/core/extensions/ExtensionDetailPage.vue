<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
const route = useRoute(); const client = createGeneratedApiClient(apiClient); const extension = ref<ExtensionResource | null>(null); const error = ref('')
onMounted(async () => { try { extension.value = (await client.getExtension(String(route.params.id))).data } catch (cause) { error.value = cause instanceof Error ? cause.message : '扩展加载失败' } })
</script>
<template><section class="flex flex-col gap-4"><h1 class="text-2xl font-semibold tracking-tight">扩展详情</h1><Alert v-if="error" variant="destructive"><AlertTitle>无法加载</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert><Card v-else-if="extension"><CardHeader><CardTitle>{{ extension.name }}</CardTitle><CardDescription>{{ extension.kind }} · {{ extension.id }}</CardDescription></CardHeader><CardContent>{{ extension.message }}</CardContent></Card></section></template>
