<script setup lang="ts">
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import type { ExtensionResource } from '@/generated/api'

defineProps<{ items: ExtensionResource[] }>()
</script>

<template>
  <div v-if="items.length" class="grid gap-4 md:grid-cols-2">
    <Card v-for="item in items" :key="item.id">
      <CardHeader>
        <CardTitle>{{ item.name }}</CardTitle>
        <CardDescription>业务模块 · {{ item.id }}</CardDescription>
      </CardHeader>
      <CardContent class="flex items-center justify-between gap-4">
        <div class="space-y-1 text-sm">
          <p>状态：{{ item.state === 'enabled' ? '已启用' : '已停用' }}</p>
          <p v-if="item.version" class="text-muted-foreground">版本：{{ item.version }}</p>
        </div>
        <div class="flex gap-2">
          <Button v-if="item.config" variant="outline" as-child>
            <RouterLink :to="item.config.route">配置</RouterLink>
          </Button>
          <Button variant="outline" as-child>
            <RouterLink :to="`/admin/modules/${item.id}`">详情</RouterLink>
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
  <p v-else class="text-sm text-muted-foreground">暂无已注册业务模块。</p>
</template>
