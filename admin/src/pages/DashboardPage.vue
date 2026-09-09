<script setup lang="ts">
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
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
      <Card><CardHeader><CardTitle>开发语言</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">zh-CN 优先，可扩展 en</CardContent></Card>
    </div>
  </section>
</template>
