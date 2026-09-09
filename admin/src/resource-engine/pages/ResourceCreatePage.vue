<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { createResourceContext } from '../core/ResourceContext'
import { resourceRegistry } from '../demo'
import ResourceForm from '../form/ResourceForm.vue'

const route = useRoute()
const router = useRouter()
const name = computed(() => String(route.params.resource ?? ''))
const definition = computed(() => resourceRegistry.get(name.value))
const provider = computed(() => resourceRegistry.provider(name.value))
const context = computed(() => definition.value && provider.value ? createResourceContext(definition.value, provider.value, ['dashboard.view']) : null)

async function create(values: Record<string, unknown>) {
  if (!provider.value || !definition.value || !context.value?.can('create')) return
  const id = `${name.value}-${Date.now()}`
  await provider.value.create({ id, ...values })
  await router.replace(`/admin/resources/${name.value}/${id}`)
}
</script>

<template>
  <section v-if="definition && provider && context" class="flex flex-col gap-6">
    <div><h1 class="text-2xl font-semibold tracking-tight">新建{{ definition.label }}</h1><p class="text-sm text-muted-foreground">使用通用资源表单创建记录。</p></div>
    <Card v-if="context.can('create')"><CardHeader><CardTitle>基本信息</CardTitle><CardDescription>字段定义由 ResourceDefinition 提供。</CardDescription></CardHeader><CardContent><ResourceForm :definition="definition" submit-label="创建" @submit="create" /></CardContent></Card>
    <Alert v-else variant="destructive"><AlertTitle>无权操作</AlertTitle><AlertDescription>当前账号没有创建此资源的权限。</AlertDescription></Alert>
  </section>
  <Alert v-else variant="destructive"><AlertTitle>资源不存在</AlertTitle><AlertDescription>未注册资源：{{ name }}</AlertDescription></Alert>
</template>
