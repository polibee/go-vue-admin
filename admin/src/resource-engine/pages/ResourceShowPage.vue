<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { createResourceContext } from '../core/ResourceContext'
import { resourceRegistry } from '../demo'

const route = useRoute()
const router = useRouter()
const name = computed(() => String(route.params.resource ?? ''))
const id = computed(() => String(route.params.id ?? ''))
const definition = computed(() => resourceRegistry.get(name.value))
const provider = computed(() => resourceRegistry.provider(name.value))
const context = computed(() => definition.value && provider.value ? createResourceContext(definition.value, provider.value, ['dashboard.view']) : null)
const record = ref<Record<string, unknown> | null>(null)
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    if (provider.value) record.value = await provider.value.get(id.value) as Record<string, unknown>
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '记录加载失败'
  }
})
</script>

<template>
  <section v-if="definition && provider && context" class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4"><div><h1 class="text-2xl font-semibold tracking-tight">{{ definition.label }}详情</h1><p class="text-sm text-muted-foreground">{{ id }}</p></div><Button v-if="context.can('update')" @click="router.push(`/admin/resources/${name}/${id}/edit`)">编辑</Button></div>
    <Alert v-if="error" variant="destructive"><AlertTitle>加载失败</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="record"><CardHeader><CardTitle>{{ record.name }}</CardTitle><CardDescription>通用 ResourceShowPage</CardDescription></CardHeader><CardContent class="grid gap-4 sm:grid-cols-2"><div v-for="field in definition.fields ?? []" :key="String(field.name)" class="flex flex-col gap-1"><span class="text-sm text-muted-foreground">{{ field.label }}</span><span>{{ record[String(field.name)] }}</span></div></CardContent></Card>
  </section>
  <Alert v-else variant="destructive"><AlertTitle>资源不存在</AlertTitle><AlertDescription>未注册资源：{{ name }}</AlertDescription></Alert>
</template>
