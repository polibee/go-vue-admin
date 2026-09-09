<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { toast } from 'vue-sonner'
import { useRoute, useRouter } from 'vue-router'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { createResourceContext } from '../core/ResourceContext'
import { resourceRegistry } from '../demo'
import ResourceForm from '../form/ResourceForm.vue'

const route = useRoute()
const router = useRouter()
const name = computed(() => String(route.params.resource ?? ''))
const id = computed(() => String(route.params.id ?? ''))
const definition = computed(() => resourceRegistry.get(name.value))
const provider = computed(() => resourceRegistry.provider(name.value))
const context = computed(() => definition.value && provider.value ? createResourceContext(definition.value, provider.value, ['dashboard.view']) : null)
const initialValues = ref<Record<string, unknown>>({})
const error = ref<string | null>(null)

onMounted(async () => {
  try {
    if (!provider.value) return
    initialValues.value = await provider.value.get(id.value) as Record<string, unknown>
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '记录加载失败'
  }
})

async function update(values: Record<string, unknown>) {
  if (!provider.value || !context.value?.can('update')) return
  error.value = null
  try {
    await provider.value.update(id.value, values)
    toast.success('记录已保存')
    await router.replace(`/admin/resources/${name.value}/${id.value}`)
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '保存失败'
  }
}
</script>

<template>
  <section v-if="definition && provider && context" class="flex flex-col gap-6">
    <div><h1 class="text-2xl font-semibold tracking-tight">编辑{{ definition.label }}</h1><p class="text-sm text-muted-foreground">{{ id }}</p></div>
    <Alert v-if="error" variant="destructive"><AlertTitle>加载失败</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="context.can('update')"><CardHeader><CardTitle>基本信息</CardTitle><CardDescription>修改后将返回详情页。</CardDescription></CardHeader><CardContent><ResourceForm :definition="definition" :initial-values="initialValues" submit-label="保存" @submit="update" /></CardContent></Card>
    <Alert v-else variant="destructive"><AlertTitle>无权操作</AlertTitle><AlertDescription>当前账号没有编辑此资源的权限。</AlertDescription></Alert>
  </section>
  <Alert v-else variant="destructive"><AlertTitle>资源不存在</AlertTitle><AlertDescription>未注册资源：{{ name }}</AlertDescription></Alert>
</template>
