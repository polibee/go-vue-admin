<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'

const route = useRoute()
const client = createGeneratedApiClient(apiClient)
const extension = ref<ExtensionResource | null>(null)
const error = ref('')

onMounted(async () => {
  try {
    extension.value = (await client.getPlugin(String(route.params.id))).data
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '插件配置加载失败'
  }
})
</script>

<template>
  <section class="flex flex-col gap-4">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ extension?.config?.label ?? '扩展配置' }}</h1>
      <p class="text-sm text-muted-foreground">业务配置由插件页面负责，平台管理页只维护生命周期和配置状态。</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>加载失败</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <Card v-else-if="extension">
      <CardHeader>
        <CardTitle>{{ extension.name }}</CardTitle>
        <CardDescription>配置标识：{{ extension.config?.schema ?? '未声明配置模型' }}</CardDescription>
      </CardHeader>
      <CardContent class="space-y-3 text-sm">
        <p>当前状态：{{ extension.config?.status === 'configured' ? '已配置' : '未配置' }}</p>
        <p v-if="extension.config?.secret_fields?.length" class="text-muted-foreground">
          敏感字段：{{ extension.config.secret_fields.join('、') }}（仅提交，不回显明文）
        </p>
        <p class="text-muted-foreground">具体配置表单和保存接口由对应业务插件实现。</p>
      </CardContent>
    </Card>
  </section>
</template>
