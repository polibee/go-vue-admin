<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { FileClock } from '@lucide/vue'
import { toast } from 'vue-sonner'

import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type AuditResource } from '@/generated/api'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Separator } from '@/components/ui/separator'
import { Spinner } from '@/components/ui/spinner'
import { DiffViewer } from '@/modules/audit/components'
import { AuditService } from '@/modules/audit/audit.service'

const service = new AuditService(createGeneratedApiClient(apiClient))
const entries = ref<AuditResource[]>([])
const selected = ref<AuditResource | null>(null)
const loading = ref(true)
const detailLoading = ref(false)
const error = ref('')

async function loadAudit() {
  loading.value = true
  error.value = ''
  try {
    entries.value = await service.list()
    if (entries.value[0]) await selectEntry(entries.value[0])
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '审计日志加载失败，请稍后重试。'
  } finally {
    loading.value = false
  }
}

async function selectEntry(entry: AuditResource) {
  detailLoading.value = true
  try {
    selected.value = await service.get(entry.id)
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : '审计详情加载失败，请稍后重试。')
  } finally {
    detailLoading.value = false
  }
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat('zh-CN', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}

onMounted(loadAudit)
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">审计日志</h1>
      <p class="text-sm text-muted-foreground">查看关键设置和媒体操作的操作者、请求来源与前后变更快照。</p>
    </div>

    <Alert v-if="error" variant="destructive">
      <AlertTitle>审计服务暂时不可用</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <div v-if="loading" class="flex items-center gap-2 text-sm text-muted-foreground"><Spinner /> 正在加载审计日志…</div>
    <Empty v-else-if="!entries.length">
      <EmptyHeader><FileClock class="size-10 text-muted-foreground" /><EmptyTitle>暂无审计记录</EmptyTitle></EmptyHeader>
      <EmptyDescription>当设置或媒体发生变更后，记录会显示在这里。</EmptyDescription>
    </Empty>
    <div v-else class="grid gap-6 xl:grid-cols-[minmax(18rem,0.8fr)_minmax(0,2fr)]">
      <Card class="h-fit">
        <CardHeader>
          <CardTitle>操作记录</CardTitle>
          <CardDescription>共 {{ entries.length }} 条，最近的记录在前。</CardDescription>
        </CardHeader>
        <CardContent class="grid gap-2">
          <button
            v-for="entry in entries"
            :key="entry.id"
            type="button"
            class="rounded-lg border p-3 text-left transition-colors hover:bg-muted/60"
            :class="selected?.id === entry.id ? 'border-primary bg-muted' : 'border-border'"
            @click="selectEntry(entry)"
          >
            <div class="flex items-center justify-between gap-2">
              <span class="truncate text-sm font-medium">{{ entry.action }}</span>
              <Badge variant="secondary">{{ entry.resource_type }}</Badge>
            </div>
            <p class="mt-1 truncate text-xs text-muted-foreground">{{ entry.resource_id }}</p>
            <p class="mt-2 text-xs text-muted-foreground">{{ formatDate(entry.created_at) }}</p>
          </button>
        </CardContent>
      </Card>

      <Card v-if="selected">
        <CardHeader>
          <div class="flex items-start justify-between gap-4">
            <div>
              <CardTitle>{{ selected.action }}</CardTitle>
              <CardDescription>{{ selected.resource_type }} / {{ selected.resource_id }}</CardDescription>
            </div>
            <Spinner v-if="detailLoading" />
          </div>
        </CardHeader>
        <CardContent class="grid gap-5">
          <div class="grid gap-3 text-sm sm:grid-cols-3">
            <div><span class="text-muted-foreground">操作者</span><p class="font-medium">{{ selected.actor_email || selected.actor_id || '系统' }}</p></div>
            <div><span class="text-muted-foreground">来源 IP</span><p class="font-medium">{{ selected.ip || '—' }}</p></div>
            <div><span class="text-muted-foreground">时间</span><p class="font-medium">{{ formatDate(selected.created_at) }}</p></div>
          </div>
          <Separator />
          <DiffViewer :before="selected.before" :after="selected.after" />
        </CardContent>
      </Card>
    </div>
  </section>
</template>
