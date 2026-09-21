<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RefreshCw } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ApiError, errorMessageKey } from '@/lib/api'
import { generatedApi, type AuditLog } from '@/generated/api'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const auth = useAuthStore()
const entries = ref<AuditLog[]>([])
const loading = ref(true)
const error = ref('')

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}

async function load() {
  if (!auth.token) return
  loading.value = true
  error.value = ''
  try {
    entries.value = await generatedApi.auditLogs(auth.token)
  } catch (value) {
    entries.value = []
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ t('auth.auditLogs') }}</h1>
        <p class="text-sm text-muted-foreground">{{ t('auth.auditDescription') }}</p>
      </div>
      <Button variant="outline" :disabled="loading" @click="load"><RefreshCw data-icon="inline-start" />{{ t('resource.refresh') }}</Button>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card>
      <CardHeader><CardTitle>{{ t('auth.auditRecent') }}</CardTitle><CardDescription>{{ t('auth.auditRecentDescription') }}</CardDescription></CardHeader>
      <CardContent>
        <div v-if="loading" class="flex flex-col gap-3"><Skeleton v-for="item in 5" :key="item" class="h-10" /></div>
        <Empty v-else-if="!entries.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('auth.auditEmpty') }}</EmptyDescription></EmptyHeader></Empty>
        <Table v-else><TableHeader><TableRow><TableHead>{{ t('auth.auditAction') }}</TableHead><TableHead>{{ t('auth.auditUser') }}</TableHead><TableHead>{{ t('auth.auditTime') }}</TableHead></TableRow></TableHeader><TableBody><TableRow v-for="entry in entries" :key="entry.id"><TableCell class="font-medium">{{ entry.action }}</TableCell><TableCell>{{ entry.user_id }}</TableCell><TableCell>{{ formatDate(entry.created_at) }}</TableCell></TableRow></TableBody></Table>
      </CardContent>
    </Card>
  </div>
</template>
