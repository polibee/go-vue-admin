<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'

const milestones = [
  { name: 'Task 1 · Repository foundation', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 2 · shadcn-vue UI foundation', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 3 · Admin shell and routing', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 4 · API contracts', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 5 · Cookie session authentication', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 6 · RBAC and dynamic navigation', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 7 · Resource Engine Core', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 8 · Generic CRUD pages', status: '待开始', value: 0, variant: 'outline' as const },
]
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">仪表盘</h1>
      <p class="text-sm text-muted-foreground">平台基础能力已准备好，业务模块可按契约接入。</p>
    </div>
    <div class="grid gap-4 md:grid-cols-3">
      <Card><CardHeader><CardTitle>组件基础</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">官方 shadcn-vue</CardContent></Card>
      <Card><CardHeader><CardTitle>后端基础</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">Goravel v1.18.x</CardContent></Card>
      <Card><CardHeader><CardTitle>数据基础</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">MySQL 8.x · utf8mb4</CardContent></Card>
    </div>
    <Card>
      <CardHeader>
        <div class="flex items-start justify-between gap-4">
          <div>
            <CardTitle>开发进度</CardTitle>
            <CardDescription>模块按阶段完成后独立提交并推送 GitHub。</CardDescription>
          </div>
          <Badge variant="outline">Task 7 / 14</Badge>
        </div>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <div v-for="milestone in milestones" :key="milestone.name" class="flex flex-col gap-2">
          <div class="flex items-center justify-between gap-4 text-sm">
            <span>{{ milestone.name }}</span>
            <Badge :variant="milestone.variant">{{ milestone.status }}</Badge>
          </div>
          <Progress :model-value="milestone.value" :aria-label="`${milestone.name} ${milestone.value}%`" />
        </div>
      </CardContent>
    </Card>
    <div class="grid gap-4 md:grid-cols-3">
      <Card><CardHeader><CardTitle>开发语言</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">zh-CN 优先，可扩展 en</CardContent></Card>
      <Card><CardHeader><CardTitle>下一步</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">完成资源引擎并接入通用 CRUD。</CardContent></Card>
      <Card><CardHeader><CardTitle>代码仓库</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">polibee/go-vue-admin</CardContent></Card>
    </div>
  </section>
</template>
