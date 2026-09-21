<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDownUp, Download, Pencil, RefreshCw, Search, Trash2 } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Checkbox } from '@/components/ui/checkbox'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious } from '@/components/ui/pagination'
import { Skeleton } from '@/components/ui/skeleton'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ApiError, errorMessageKey } from '@/lib/api'
import { generatedApi, type ResourceManifest as GeneratedResourceManifest, type ResourceListMeta } from '@/generated/api'
import { canDeleteResource } from '@/lib/resource-actions'
import { useAuthStore } from '@/stores/auth'
import { USER_STATUSES, userStatusLabelKey, type UserStatus } from '@/lib/user-status'
import { useI18n } from 'vue-i18n'

interface ResourceColumn { name: string; label: string; sortable: boolean }
interface ResourceManifest extends GeneratedResourceManifest {}
type ResourceMeta = ResourceListMeta

const { t } = useI18n()
const props = defineProps<{ resource?: string }>()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const manifests = ref<ResourceManifest[]>([])
const rows = ref<Record<string, unknown>[]>([])
const meta = ref<ResourceMeta>({ page: 1, per_page: 10, total: 0, last_page: 1 })
const search = ref('')
const statusFilter = ref('all')
const sort = ref('id')
const direction = ref<'asc' | 'desc'>('desc')
const pageSize = ref('10')
const loading = ref(true)
const exporting = ref(false)
const error = ref('')
const deleteDialogOpen = ref(false)
const deleting = ref(false)
const deleteTarget = ref<Record<string, unknown>>()
const selectedIds = ref<string[]>([])
const bulkStatus = ref<UserStatus>('disabled')
const bulkStatusDialogOpen = ref(false)
const bulkUpdating = ref(false)
const filterValues = ref<Record<string, string>>({})

const resourceName = computed(() => props.resource || String(route.params.resource || 'users'))
const currentManifest = computed(() => manifests.value.find((item) => item.name === resourceName.value))
const filterFields = computed(() => (currentManifest.value?.fields || []).filter((field) => field.type === 'select' || field.type === 'boolean'))
const canManageUsers = computed(() => auth.can('admin.users.manage'))
const canCreate = computed(() => hasAction('create'))
const allVisibleSelected = computed(() => rows.value.length > 0 && rows.value.every((row) => selectedIds.value.includes(String(row.id))))

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

async function loadManifests() {
  if (!auth.token) return
  manifests.value = await generatedApi.resourceRegistry(auth.token)
  if (!currentManifest.value && manifests.value.length) {
    await router.replace(`/${manifests.value[0].name}`)
  }
}

async function loadRows(page = 1) {
  if (!auth.token || !currentManifest.value) return
  loading.value = true
  error.value = ''
  try {
    const params = buildResourceQuery()
    params.set('page', String(page))
    params.set('per_page', pageSize.value)
    const response = await generatedApi.resourceList<Record<string, unknown>>(resourceName.value, params, auth.token)
    rows.value = response.data
    meta.value = response.meta as unknown as ResourceMeta
    selectedIds.value = []
  } catch (value) {
    rows.value = []
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
}

function buildResourceQuery() {
  const params = new URLSearchParams({ sort: sort.value, dir: direction.value })
  if (search.value.trim()) params.set('search', search.value.trim())
  if (resourceName.value === 'users' && statusFilter.value !== 'all') params.set('status', statusFilter.value)
  for (const field of filterFields.value) {
    const value = filterValues.value[field.name]
    if (value && value !== 'all') params.set(field.name, value)
  }
  return params
}

function submitSearch() { void loadRows(1) }
function changePageSize(value: unknown) {
  pageSize.value = String(value)
  meta.value.per_page = Number(pageSize.value)
  void loadRows(1)
}
function sortBy(column: ResourceColumn) {
  if (!column.sortable) return
  if (sort.value === column.name) direction.value = direction.value === 'asc' ? 'desc' : 'asc'
  else { sort.value = column.name; direction.value = 'asc' }
  void loadRows(1)
}
function selectResource(name: string) { void router.push(`/${name}`) }
function displayValue(value: unknown) { return value === null || value === undefined ? '—' : String(value) }
function statusLabel(value: unknown) { return typeof value === 'string' ? t(userStatusLabelKey(value as UserStatus)) : '—' }
function changeStatusFilter(value: unknown) { statusFilter.value = String(value); void loadRows(1) }
function filterOptions(field: ResourceManifest['fields'][number]) {
  return field.type === 'boolean' ? [{ value: 'all', label: t('resource.filterAll') }, { value: 'true', label: t('resource.trueValue') }, { value: 'false', label: t('resource.falseValue') }] : [{ value: 'all', label: t('resource.filterAll') }, ...(field.options || [])]
}
function changeResourceFilter(name: string, value: unknown) {
  filterValues.value = { ...filterValues.value, [name]: String(value) }
  void loadRows(1)
}
function toggleRow(id: unknown, checked: boolean | 'indeterminate') {
  const value = String(id)
  selectedIds.value = checked === true ? [...new Set([...selectedIds.value, value])] : selectedIds.value.filter((item) => item !== value)
}
function toggleAll(checked: boolean | 'indeterminate') {
  selectedIds.value = checked === true ? rows.value.map((row) => String(row.id)) : []
}
function clearSelection() { selectedIds.value = [] }
function hasAction(name: string) {
  const action = currentManifest.value?.actions?.find((action) => action.name === name)
  return Boolean(action && auth.can(action.permission))
}
function openRowStatusAction(row: Record<string, unknown>) {
  selectedIds.value = row.id ? [String(row.id)] : []
  bulkStatus.value = (typeof row.status === 'string' ? row.status : 'disabled') as UserStatus
  bulkStatusDialogOpen.value = true
}
async function applyBulkStatus() {
  if (!auth.token || !selectedIds.value.length) return
  bulkUpdating.value = true
  error.value = ''
  try {
    await generatedApi.bulkSetUserStatus({ user_ids: selectedIds.value.map(Number), status: bulkStatus.value }, auth.token)
    bulkStatusDialogOpen.value = false
    clearSelection()
    await loadRows(meta.value.page)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    bulkUpdating.value = false
  }
}
function editPath(row: Record<string, unknown>) { return `/${resourceName.value}/${row.id}/edit` }
function openDelete(row: Record<string, unknown>) { deleteTarget.value = row; deleteDialogOpen.value = true }
async function deleteRow() {
  if (!auth.token || !deleteTarget.value?.id) return
  deleting.value = true
  error.value = ''
  try {
    await generatedApi.resourceDelete(resourceName.value, String(deleteTarget.value.id), auth.token)
    deleteDialogOpen.value = false
    await loadRows(meta.value.page)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    deleting.value = false
  }
}

async function exportRows() {
  if (!auth.token) return
  exporting.value = true
  error.value = ''
  try {
    const blob = await generatedApi.resourceExport(resourceName.value, buildResourceQuery(), auth.token)
    const url = URL.createObjectURL(blob)
    const anchor = document.createElement('a')
    anchor.href = url
    anchor.download = `${resourceName.value}-export.csv`
    anchor.click()
    URL.revokeObjectURL(url)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    exporting.value = false
  }
}

onMounted(async () => {
  try {
    await loadManifests()
    await loadRows()
  } catch (value) {
    error.value = localizedError(value)
    loading.value = false
  }
})
watch(resourceName, () => { statusFilter.value = 'all'; filterValues.value = {}; void loadRows(1) })
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ currentManifest?.label || t('resource.title') }}</h1>
        <p class="text-sm text-muted-foreground">{{ t('resource.description') }}</p>
      </div>
      <div class="flex gap-2"><Button v-if="canCreate" variant="default" @click="router.push(`/${resourceName}/new`)">{{ resourceName === 'users' ? t('resource.createUser') : resourceName === 'roles' ? t('rbac.createRole') : t('resource.create') }}</Button><Button variant="outline" :disabled="loading || exporting" @click="exportRows"><Download data-icon="inline-start" />{{ t('resource.export') }}</Button><Button variant="outline" :disabled="loading" @click="loadRows(meta.page)">
        <RefreshCw data-icon="inline-start" />{{ t('resource.refresh') }}
      </Button></div>
    </div>

    <div v-if="manifests.length" class="flex flex-wrap gap-2">
      <Button v-for="manifest in manifests" :key="manifest.name" :variant="manifest.name === resourceName ? 'default' : 'outline'" size="sm" @click="selectResource(manifest.name)">{{ manifest.label }}</Button>
    </div>
    <div v-if="resourceName === 'users' && selectedIds.length && canManageUsers" class="flex flex-wrap items-center gap-2 rounded-lg border bg-muted/30 p-3">
      <span class="text-sm text-muted-foreground">{{ t('resource.selectedCount', { count: selectedIds.length }) }}</span>
      <Select v-model="bulkStatus"><SelectTrigger class="w-36" :aria-label="t('resource.bulkStatus')"><SelectValue /></SelectTrigger><SelectContent><SelectItem v-for="status in USER_STATUSES" :key="status" :value="status">{{ t(userStatusLabelKey(status)) }}</SelectItem></SelectContent></Select>
      <Button size="sm" @click="bulkStatusDialogOpen = true">{{ t('resource.applyStatus') }}</Button>
      <Button variant="ghost" size="sm" @click="clearSelection">{{ t('resource.clearSelection') }}</Button>
    </div>

    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card>
      <CardHeader class="gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div><CardTitle>{{ currentManifest?.label || t('resource.resourceNotFound') }}</CardTitle><CardDescription>{{ t('resource.total', { count: meta.total }) }}</CardDescription></div>
        <form class="flex w-full flex-wrap gap-2 sm:w-auto" @submit.prevent="submitSearch"><Select v-if="resourceName === 'users'" :model-value="statusFilter" @update:model-value="changeStatusFilter"><SelectTrigger class="w-32" :aria-label="t('resource.statusFilter')"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="all">{{ t('resource.statusAll') }}</SelectItem><SelectItem v-for="status in USER_STATUSES" :key="status" :value="status">{{ t(userStatusLabelKey(status)) }}</SelectItem></SelectContent></Select><Select v-for="field in filterFields" :key="field.name" :model-value="filterValues[field.name] || 'all'" @update:model-value="changeResourceFilter(field.name, $event)"><SelectTrigger class="w-36" :aria-label="field.label"><SelectValue :placeholder="field.label" /></SelectTrigger><SelectContent><SelectItem v-for="option in filterOptions(field)" :key="option.value" :value="option.value">{{ option.label }}</SelectItem></SelectContent></Select><Input v-model="search" class="sm:w-64" :placeholder="t('resource.searchPlaceholder')" :aria-label="t('resource.search')" /><Button type="submit" size="icon" :aria-label="t('resource.search')"><Search /></Button></form>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="flex flex-col gap-3"><Skeleton v-for="item in 5" :key="item" class="h-10" /></div>
        <Empty v-else-if="!rows.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('resource.noData') }}</EmptyDescription></EmptyHeader></Empty>
        <Table v-else><TableHeader><TableRow><TableHead v-if="resourceName === 'users'" class="w-10"><Checkbox :checked="allVisibleSelected" :aria-label="t('resource.selectAll')" @click="toggleAll(!allVisibleSelected)" /></TableHead><TableHead v-for="column in currentManifest?.columns || []" :key="column.name"><Button v-if="column.sortable" variant="ghost" size="sm" class="-ml-3" @click="sortBy(column)">{{ column.label }}<ArrowDownUp data-icon="inline-end" /></Button><span v-else>{{ column.label }}</span></TableHead><TableHead class="w-36 text-right">{{ t('resource.actions') }}</TableHead></TableRow></TableHeader><TableBody><TableRow v-for="(row, index) in rows" :key="String(row.id || index)" class="cursor-pointer" @click="row.id && router.push(`/${resourceName}/${row.id}`)"><TableCell v-if="resourceName === 'users'" @click.stop><Checkbox :checked="selectedIds.includes(String(row.id))" :aria-label="t('resource.selectRow', { name: row.name })" @click="toggleRow(row.id, !selectedIds.includes(String(row.id)))" /></TableCell><TableCell v-for="column in currentManifest?.columns || []" :key="column.name"><Badge v-if="column.name === 'status'" variant="secondary">{{ statusLabel(row[column.name]) }}</Badge><template v-else>{{ displayValue(row[column.name]) }}</template></TableCell><TableCell class="text-right"><div v-if="row.id && (hasAction('update') || hasAction('delete') || hasAction('set-status'))" class="flex justify-end gap-1" @click.stop><Button v-if="hasAction('set-status')" variant="ghost" size="sm" :aria-label="t('resource.setStatus')" @click="openRowStatusAction(row)">{{ t('resource.setStatus') }}</Button><Button v-if="hasAction('update')" variant="ghost" size="sm" :aria-label="t('resource.edit')" @click="router.push(editPath(row))"><Pencil data-icon="inline-start" />{{ t('resource.edit') }}</Button><Button v-if="hasAction('delete')" variant="ghost" size="sm" :disabled="!canDeleteResource(resourceName, row)" :aria-label="t('resource.delete')" @click="openDelete(row)"><Trash2 data-icon="inline-start" />{{ t('resource.delete') }}</Button></div></TableCell></TableRow></TableBody></Table>
        <div v-if="!loading && meta.total > 0" class="mt-4 flex flex-col items-center gap-3 sm:flex-row sm:justify-between"><p class="text-sm text-muted-foreground">{{ t('resource.page', { page: meta.page }) }}</p><div class="flex items-center gap-2"><Select :model-value="pageSize" @update:model-value="changePageSize"><SelectTrigger class="w-24"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="10">{{ t('resource.perPage', { count: 10 }) }}</SelectItem><SelectItem value="20">{{ t('resource.perPage', { count: 20 }) }}</SelectItem><SelectItem value="50">{{ t('resource.perPage', { count: 50 }) }}</SelectItem></SelectContent></Select></div><Pagination v-model:page="meta.page" :items-per-page="meta.per_page" :total="meta.total" @update:page="loadRows"><PaginationContent v-slot="{ items }"><PaginationPrevious /><template v-for="(item, index) in items" :key="index"><PaginationItem v-if="item.type === 'page'" :value="item.value" :is-active="item.value === meta.page">{{ item.value }}</PaginationItem></template><PaginationNext /></PaginationContent></Pagination></div>
      </CardContent>
    </Card>
    <AlertDialog v-model:open="deleteDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('resource.deleteTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resource.deleteDescription') }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('resource.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="deleting" @click="deleteRow">{{ t('resource.delete') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
    <AlertDialog v-model:open="bulkStatusDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('resource.bulkStatusTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resource.bulkStatusDescription', { count: selectedIds.length, status: t(userStatusLabelKey(bulkStatus)) }) }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('resource.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="bulkUpdating" @click="applyBulkStatus">{{ t('resource.applyStatus') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
  </div>
</template>
