<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RefreshCw, Trash2 } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious } from '@/components/ui/pagination'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ApiError, errorMessageKey } from '@/lib/api'
import { generatedApi, type AuditLog, type ResourceListMeta } from '@/generated/api'
import { formatAuditMetadata } from '@/lib/audit-log'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const auth = useAuthStore()
const entries = ref<AuditLog[]>([])
const meta = ref<ResourceListMeta>({ page: 1, per_page: 20, total: 0, last_page: 1 })
const pageSize = ref('20')
const actionFilter = ref('all')
const userFilter = ref('')
const selectedEntry = ref<AuditLog | null>(null)
const loading = ref(true)
const error = ref('')
const cleanupOpen = ref(false)
const cleanupMode = ref<'retention' | 'all'>('retention')
const cleanupDays = ref('365')
const cleaning = ref(false)
const cleanupResult = ref<number | null>(null)
const canCleanup = computed(() => auth.can('admin.settings.manage'))

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}

async function load(page = 1) {
  if (!auth.token) return
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams({ page: String(page), per_page: pageSize.value })
    if (actionFilter.value !== 'all') params.set('action', actionFilter.value)
    if (userFilter.value.trim()) params.set('user_id', userFilter.value.trim())
    const response = await generatedApi.auditLogs(auth.token, params)
    entries.value = response.data
    meta.value = response.meta
  } catch (value) {
    entries.value = []
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
}

function changePageSize(value: unknown) {
  pageSize.value = String(value)
  meta.value.per_page = Number(pageSize.value)
  void load(1)
}

function submitFilters() { void load(1) }
function resetFilters() {
  actionFilter.value = 'all'
  userFilter.value = ''
  void load(1)
}

async function cleanupLogs() {
  if (!auth.token) return
  const days = Number(cleanupDays.value)
  if (cleanupMode.value === 'retention' && (!Number.isInteger(days) || days < 1 || days > 3650)) {
    error.value = t('auth.auditRetentionInvalid')
    return
  }
  cleaning.value = true
  error.value = ''
  try {
    const result = await generatedApi.cleanupAuditLogs(cleanupMode.value === 'all' ? { mode: 'all' } : { mode: 'retention', retention_days: days }, auth.token)
    cleanupResult.value = result.deleted
    cleanupOpen.value = false
    await load(1)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    cleaning.value = false
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
      <div class="flex gap-2"><Button v-if="canCleanup" variant="outline" :disabled="loading" @click="cleanupOpen = true"><Trash2 data-icon="inline-start" />{{ t('auth.auditCleanup') }}</Button><Button variant="outline" :disabled="loading" @click="load"><RefreshCw data-icon="inline-start" />{{ t('resource.refresh') }}</Button></div>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card>
      <CardHeader class="gap-4 sm:flex-row sm:items-center sm:justify-between"><div><CardTitle>{{ t('auth.auditRecent') }}</CardTitle><CardDescription>{{ t('auth.auditRecentDescription', { count: meta.total }) }}</CardDescription></div><form class="flex w-full flex-wrap gap-2 sm:w-auto" @submit.prevent="submitFilters"><Select :model-value="actionFilter" @update:model-value="(value) => { actionFilter = String(value); submitFilters() }"><SelectTrigger class="w-40" :aria-label="t('auth.auditActionFilter')"><SelectValue :placeholder="t('auth.auditActionAll')" /></SelectTrigger><SelectContent><SelectItem value="all">{{ t('auth.auditActionAll') }}</SelectItem><SelectItem v-for="action in ['auth.login', 'auth.refresh', 'auth.logout', 'auth.logout_all', 'user.create', 'user.update', 'user.delete', 'user.status.bulk', 'user.roles.replace', 'role.create', 'role.update', 'role.delete', 'role.permissions.replace']" :key="action" :value="action">{{ action }}</SelectItem></SelectContent></Select><Input v-model="userFilter" class="w-28" :placeholder="t('auth.auditUserFilter')" :aria-label="t('auth.auditUserFilter')" /><Button type="submit" size="sm">{{ t('resource.search') }}</Button><Button type="button" variant="ghost" size="sm" @click="resetFilters">{{ t('resource.clearSelection') }}</Button></form></CardHeader>
      <CardContent>
        <div v-if="loading" class="flex flex-col gap-3"><Skeleton v-for="item in 5" :key="item" class="h-10" /></div>
        <Empty v-else-if="!entries.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('auth.auditEmpty') }}</EmptyDescription></EmptyHeader></Empty>
        <Table v-else><TableHeader><TableRow><TableHead>{{ t('auth.auditAction') }}</TableHead><TableHead>{{ t('auth.auditUser') }}</TableHead><TableHead>{{ t('auth.auditTime') }}</TableHead></TableRow></TableHeader><TableBody><TableRow v-for="entry in entries" :key="entry.id" class="cursor-pointer" @click="selectedEntry = entry"><TableCell class="font-medium">{{ entry.action }}</TableCell><TableCell>{{ entry.user_id }}</TableCell><TableCell>{{ formatDate(entry.created_at) }}</TableCell></TableRow></TableBody></Table>
        <div v-if="!loading && meta.total > 0" class="mt-4 flex flex-col items-center gap-3 sm:flex-row sm:justify-between"><p class="text-sm text-muted-foreground">{{ t('resource.page', { page: meta.page }) }}</p><div class="flex items-center gap-2"><Select :model-value="pageSize" @update:model-value="changePageSize"><SelectTrigger class="w-24"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="10">{{ t('resource.perPage', { count: 10 }) }}</SelectItem><SelectItem value="20">{{ t('resource.perPage', { count: 20 }) }}</SelectItem><SelectItem value="50">{{ t('resource.perPage', { count: 50 }) }}</SelectItem></SelectContent></Select></div><Pagination v-model:page="meta.page" :items-per-page="meta.per_page" :total="meta.total" @update:page="load"><PaginationContent v-slot="{ items }"><PaginationPrevious /><template v-for="(item, index) in items" :key="index"><PaginationItem v-if="item.type === 'page'" :value="item.value" :is-active="item.value === meta.page">{{ item.value }}</PaginationItem></template><PaginationNext /></PaginationContent></Pagination></div>
      </CardContent>
    </Card>
    <Dialog :open="selectedEntry !== null" @update:open="(open) => { if (!open) selectedEntry = null }"><DialogContent class="max-h-[calc(100vh-2rem)] overflow-y-auto"><DialogHeader><DialogTitle>{{ t('auth.auditDetail') }}</DialogTitle><DialogDescription>{{ selectedEntry?.action }} · {{ selectedEntry ? formatDate(selectedEntry.created_at) : '' }}</DialogDescription></DialogHeader><dl v-if="selectedEntry" class="grid gap-3 text-sm"><div class="flex justify-between gap-4"><dt class="text-muted-foreground">{{ t('auth.auditUser') }}</dt><dd>{{ selectedEntry.user_id }}</dd></div><div><dt class="mb-2 text-muted-foreground">{{ t('auth.auditMetadata') }}</dt><dd v-if="formatAuditMetadata(selectedEntry.metadata)" class="max-h-[50vh] overflow-auto rounded-md bg-muted p-3"><pre class="whitespace-pre-wrap break-words text-xs">{{ formatAuditMetadata(selectedEntry.metadata) }}</pre></dd><dd v-else class="text-muted-foreground">{{ t('auth.auditNoMetadata') }}</dd></div></dl></DialogContent></Dialog>
    <Dialog v-model:open="cleanupOpen"><DialogContent class="max-h-[calc(100vh-2rem)] overflow-y-auto"><DialogHeader><DialogTitle>{{ t('auth.auditCleanup') }}</DialogTitle><DialogDescription>{{ t(cleanupMode === 'all' ? 'auth.auditCleanupAllDescription' : 'auth.auditCleanupDescription') }}</DialogDescription></DialogHeader><div class="grid gap-3"><div class="grid gap-2"><label class="text-sm font-medium">{{ t('auth.auditCleanupMode') }}</label><Select v-model="cleanupMode"><SelectTrigger><SelectValue /></SelectTrigger><SelectContent><SelectItem value="retention">{{ t('auth.auditCleanupRetentionMode') }}</SelectItem><SelectItem value="all">{{ t('auth.auditCleanupAllMode') }}</SelectItem></SelectContent></Select></div><div v-if="cleanupMode === 'retention'" class="grid gap-2"><label for="audit-retention-days" class="text-sm font-medium">{{ t('auth.auditRetentionDays') }}</label><Input id="audit-retention-days" v-model="cleanupDays" type="number" min="1" max="3650" /></div><Alert v-else variant="destructive"><AlertTitle>{{ t('auth.auditCleanupAllWarningTitle') }}</AlertTitle><AlertDescription>{{ t('auth.auditCleanupAllWarning') }}</AlertDescription></Alert></div><DialogFooter><Button variant="outline" @click="cleanupOpen = false">{{ t('resource.cancel') }}</Button><Button :disabled="cleaning" @click="cleanupLogs">{{ t(cleanupMode === 'all' ? 'auth.auditCleanupAllConfirm' : 'auth.auditCleanupConfirm') }}</Button></DialogFooter></DialogContent></Dialog>
    <Alert v-if="cleanupResult !== null"><AlertTitle>{{ t('auth.auditCleanupDone') }}</AlertTitle><AlertDescription>{{ t('auth.auditCleanupResult', { count: cleanupResult }) }}</AlertDescription></Alert>
  </div>
</template>
