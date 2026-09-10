<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'

const milestones = [
  { name: 'Task 1 · API generation foundation', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 2 · Scalar API reference', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 3 · API route-to-client audit', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 4 · GORM CRUD and route registration', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 5 · Generated API consumer migration', status: '已完成', value: 100, variant: 'secondary' as const },
  { name: 'Task 6 · Configured database and browser acceptance', status: '已完成', value: 100, variant: 'secondary' as const },
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
      <Card><CardHeader><CardTitle>数据基础</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">GORM · MySQL / PostgreSQL</CardContent></Card>
    </div>
    <Card>
      <CardHeader>
        <div class="flex items-start justify-between gap-4">
          <div>
            <CardTitle>开发进度</CardTitle>
            <CardDescription>基础模块化计划 Task 1–14 与数据库兼容验收已完成；业务资源可按配置切换 MySQL / PostgreSQL。</CardDescription>
          </div>
          <Badge variant="outline">Foundation 14 / 14</Badge>
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
      <Card><CardHeader><CardTitle>下一步</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">配置 RESOURCE_PROVIDER=database，并通过 DB_CONNECTION 选择 MySQL 或 PostgreSQL 后生成业务资源。</CardContent></Card>
      <Card><CardHeader><CardTitle>代码仓库</CardTitle></CardHeader><CardContent class="text-sm"><a class="text-primary underline underline-offset-4" href="https://github.com/polibee/go-vue-admin" target="_blank" rel="noreferrer">polibee/go-vue-admin</a></CardContent></Card>
    </div>
  </section>
</template>
