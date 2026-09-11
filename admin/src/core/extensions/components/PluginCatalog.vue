<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import type { ExtensionResource } from '@/generated/api'

defineProps<{ items: ExtensionResource[]; busy: string }>()
const emit = defineEmits<{ toggle: [item: ExtensionResource, state: 'enabled' | 'disabled'] }>()
</script>

<template>
  <div v-if="items.length" class="grid gap-4 md:grid-cols-2">
    <Card v-for="item in items" :key="item.id">
      <CardHeader>
        <CardTitle>{{ item.name }}</CardTitle>
        <CardDescription>平台插件 · {{ item.id }}</CardDescription>
      </CardHeader>
      <CardContent>状态：{{ item.state === 'enabled' ? '已启用' : '已停用' }}</CardContent>
      <CardFooter class="gap-2">
        <Button
          v-if="item.state === 'disabled'"
          :disabled="busy === item.id"
          @click="emit('toggle', item, 'enabled')"
        >
          启用
        </Button>
        <Button
          v-else
          variant="outline"
          :disabled="busy === item.id"
          @click="emit('toggle', item, 'disabled')"
        >
          停用
        </Button>
        <Button variant="outline" as-child>
          <RouterLink :to="`/admin/extensions/${item.id}`">查看</RouterLink>
        </Button>
      </CardFooter>
    </Card>
  </div>
  <p v-else class="text-sm text-muted-foreground">暂无已注册平台插件。</p>
</template>
