<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'

const { t } = useI18n()

const milestones = [
  { name: 'Task 1 · API generation foundation', value: 100, variant: 'secondary' as const },
  { name: 'Task 2 · Scalar API reference', value: 100, variant: 'secondary' as const },
  { name: 'Task 3 · API route-to-client audit', value: 100, variant: 'secondary' as const },
  { name: 'Task 4 · GORM CRUD and route registration', value: 100, variant: 'secondary' as const },
  { name: 'Task 5 · Generated API consumer migration', value: 100, variant: 'secondary' as const },
  { name: 'Task 6 · Configured database and browser acceptance', value: 100, variant: 'secondary' as const },
]
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('dashboard.title') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('dashboard.description') }}</p>
    </div>
    <div class="grid gap-4 md:grid-cols-3">
      <Card><CardHeader><CardTitle>{{ t('dashboard.components') }}</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">shadcn-vue</CardContent></Card>
      <Card><CardHeader><CardTitle>{{ t('dashboard.backend') }}</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">Goravel v1.18.x</CardContent></Card>
      <Card><CardHeader><CardTitle>{{ t('dashboard.data') }}</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">GORM · MySQL / PostgreSQL</CardContent></Card>
    </div>
    <Card>
      <CardHeader>
        <div class="flex items-start justify-between gap-4">
          <div>
            <CardTitle>{{ t('dashboard.progress') }}</CardTitle>
            <CardDescription>{{ t('dashboard.progressDescription') }}</CardDescription>
          </div>
          <Badge variant="outline">Foundation 14 / 14</Badge>
        </div>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <div v-for="milestone in milestones" :key="milestone.name" class="flex flex-col gap-2">
          <div class="flex items-center justify-between gap-4 text-sm">
            <span>{{ milestone.name }}</span>
            <Badge :variant="milestone.variant">{{ t('dashboard.done') }}</Badge>
          </div>
          <Progress :model-value="milestone.value" :aria-label="`${milestone.name} ${milestone.value}%`" />
        </div>
      </CardContent>
    </Card>
    <div class="grid gap-4 md:grid-cols-3">
      <Card><CardHeader><CardTitle>{{ t('dashboard.language') }}</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">zh-CN / English</CardContent></Card>
      <Card><CardHeader><CardTitle>{{ t('dashboard.next') }}</CardTitle></CardHeader><CardContent class="text-sm text-muted-foreground">RESOURCE_PROVIDER=database · DB_CONNECTION=mysql|postgres</CardContent></Card>
      <Card><CardHeader><CardTitle>{{ t('dashboard.repository') }}</CardTitle></CardHeader><CardContent class="text-sm"><a class="text-primary underline underline-offset-4" href="https://github.com/polibee/go-vue-admin" target="_blank" rel="noreferrer">polibee/go-vue-admin</a></CardContent></Card>
    </div>
  </section>
</template>
