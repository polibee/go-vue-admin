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
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious } from '@/components/ui/pagination'
import { Skeleton } from '@/components/ui/skeleton'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ApiError, errorMessageKey } from '@/lib/api'
import { generatedApi, type ActionResponse, type ActionRequest, type ResourceFilter, type ResourceManifest as GeneratedResourceManifest, type ResourceListMeta } from '@/generated/api'
import { executableBatchActions } from '@/lib/resource-actions'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

interface ResourceColumn { name: string; label: string; sortable: boolean }
interface ResourceManifest extends GeneratedResourceManifest {}
type ResourceAction = NonNullable<ResourceManifest['actions']>[number]
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
const bulkStatus = ref<'active' | 'disabled' | 'locked'>('disabled')
const bulkStatusDialogOpen = ref(false)
const bulkUpdating = ref(false)
const bulkDeleteDialogOpen = ref(false)
const bulkDeleting = ref(false)
const bulkDeleteAction = ref('bulk-delete')
const bulkUpdateDialogOpen = ref(false)
const bulkUpdateField = ref('')
const bulkUpdateValue = ref('')
const bulkUpdateSaving = ref(false)
const lastActionResult = ref<ActionResponse>()
const filterValues = ref<Record<string, string>>({})
const trashed = ref('default')
const allFilteredSelected = ref(false)
const excludedIds = ref<string[]>([])
const bulkStatusLabel = computed(() => statusLabel('status', bulkStatus.value))

const resourceName = computed(() => props.resource || String(route.params.resource || 'users'))
const currentManifest = computed(() => manifests.value.find((item) => item.name === resourceName.value))
const visibleColumns = computed(() => (currentManifest.value?.columns || []).filter((column) => {
  const field = currentManifest.value?.fields.find((item) => item.name === column.name)
  return !field || (field.visible !== false && field.readable !== false)
}))
const filters = computed(() => currentManifest.value?.filters || [])
const canCreate = computed(() => hasAction('create'))
const batchActions = computed(() => {
  const actions = executableBatchActions(currentManifest.value?.actions, (permission) => auth.can(permission)).slice()
  const deleteAction = currentManifest.value?.actions?.find((action) => action.name === 'delete')
  if (currentManifest.value?.soft_delete && deleteAction && trashed.value === 'only' && auth.can(deleteAction.permission)) {
    actions.push({ ...deleteAction, name: 'restore', label: t('resource.restore'), kind: 'builtin-restore', batch: true, payload: 'trash' })
    actions.push({ ...deleteAction, name: 'force-delete', label: t('resource.forceDelete'), kind: 'builtin-force-delete', batch: true, payload: 'trash' })
  }
  return actions
})
const allVisibleSelected = computed(() => rows.value.length > 0 && rows.value.every((row) => isRowSelected(row.id)))
const headerChecked = computed({
  get: () => allVisibleSelected.value,
  set: (value: boolean | 'indeterminate') => toggleAll(value),
})
const selectedCount = computed(() => allFilteredSelected.value ? Math.max(0, meta.value.total - excludedIds.value.length) : selectedIds.value.length)
const writableFields = computed(() => (currentManifest.value?.fields || []).filter((field) => field.writable !== false && field.visible !== false && !field.sensitive && field.name !== 'password'))

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

async function loadManifests() {
  if (!auth.token) return
  manifests.value = await generatedApi.resourceRegistry(auth.token)
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
  for (const field of filters.value) {
    const value = filterValues.value[field.name]
    if (value && value !== 'all') params.set(field.name, value)
  }
  if (currentManifest.value?.soft_delete && trashed.value !== 'default') params.set('trashed', trashed.value)
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
function displayValue(value: unknown) { return value === null || value === undefined ? '—' : String(value) }
function statusLabel(fieldName: string, value: unknown) {
  if (value === null || value === undefined) return '—'
  const field = currentManifest.value?.fields.find((item) => item.name === fieldName)
  const option = field?.options?.find((item) => item.value === String(value))
  return option?.label || String(value)
}
function filterOptions(field: ResourceFilter) {
  return field.type === 'boolean' ? [{ value: 'all', label: t('resource.filterAll') }, { value: 'true', label: t('resource.trueValue') }, { value: 'false', label: t('resource.falseValue') }] : [{ value: 'all', label: t('resource.filterAll') }, ...(field.options || [])]
}
function changeResourceFilter(name: string, value: unknown) {
  filterValues.value = { ...filterValues.value, [name]: String(value) }
  void loadRows(1)
}
function toggleRow(id: unknown, checked: boolean | 'indeterminate') {
  const value = String(id)
  if (allFilteredSelected.value) {
    excludedIds.value = checked === true ? excludedIds.value.filter((item) => item !== value) : [...new Set([...excludedIds.value, value])]
    return
  }
  selectedIds.value = checked === true ? [...new Set([...selectedIds.value, value])] : selectedIds.value.filter((item) => item !== value)
}
function isRowSelected(id: unknown) {
  const value = String(id)
  return allFilteredSelected.value ? !excludedIds.value.includes(value) : selectedIds.value.includes(value)
}
function toggleAll(checked: boolean | 'indeterminate') {
  if (allFilteredSelected.value) {
    if (checked === true) excludedIds.value = excludedIds.value.filter((id) => !rows.value.some((row) => String(row.id) === id))
    else excludedIds.value = [...new Set([...excludedIds.value, ...rows.value.map((row) => String(row.id))])]
    return
  }
  if (checked === true) selectedIds.value = [...new Set([...selectedIds.value, ...rows.value.map((row) => String(row.id))])]
  else selectedIds.value = selectedIds.value.filter((id) => !rows.value.some((row) => String(row.id) === id))
}
function selectAllFiltered() { allFilteredSelected.value = true; selectedIds.value = []; excludedIds.value = [] }
function clearSelection() { selectedIds.value = []; allFilteredSelected.value = false; excludedIds.value = [] }
function hasAction(name: string) {
  const action = currentManifest.value?.actions?.find((action) => action.name === name)
  return Boolean(action && auth.can(action.permission))
}
function hasActionKind(kind: string) {
  return batchActions.value.some((action) => action.kind === kind)
}
function openRowStatusAction(row: Record<string, unknown>) {
  selectedIds.value = row.id ? [String(row.id)] : []
  const status = typeof row.status === 'string' ? row.status : 'disabled'
  bulkStatus.value = status === 'active' || status === 'locked' ? status : 'disabled'
  bulkStatusDialogOpen.value = true
}
function openBulkAction(action: ResourceAction) {
  if (action.kind === 'user-status') {
    bulkStatusDialogOpen.value = true
  } else if (action.kind === 'builtin-delete' || action.kind === 'builtin-restore' || action.kind === 'builtin-force-delete') {
    bulkDeleteAction.value = action.name
    bulkDeleteDialogOpen.value = true
  } else if (action.kind === 'builtin-update') {
    bulkUpdateField.value = writableFields.value[0]?.name || ''
    bulkUpdateValue.value = ''
    bulkUpdateDialogOpen.value = true
  }
}
function selectionRequest(payload?: Record<string, unknown>): ActionRequest {
  if (allFilteredSelected.value) {
    const query: Record<string, string> = {}
    for (const [key, value] of buildResourceQuery().entries()) {
      if (!['page', 'per_page', 'sort', 'dir'].includes(key)) query[key] = value
    }
    return { selection: { mode: 'query', query, exclude_ids: excludedIds.value.map(Number) }, payload }
  }
  return { selection: { mode: 'ids', ids: selectedIds.value.map(Number) }, payload }
}
async function applyBulkStatus() {
  if (!auth.token || !selectedIds.value.length) return
  bulkUpdating.value = true
  error.value = ''
  try {
    lastActionResult.value = await generatedApi.resourceAction(resourceName.value, 'set-status', selectionRequest({ status: bulkStatus.value }), auth.token)
    bulkStatusDialogOpen.value = false
    clearSelection()
    await loadRows(meta.value.page)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    bulkUpdating.value = false
  }
}
async function applyBulkDelete() {
  if (!auth.token || !selectedCount.value) return
  bulkDeleting.value = true
  error.value = ''
  try {
    lastActionResult.value = await generatedApi.resourceAction(resourceName.value, bulkDeleteAction.value, selectionRequest(), auth.token)
    bulkDeleteDialogOpen.value = false
    clearSelection()
    await loadRows(1)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    bulkDeleting.value = false
  }
}
async function applyBulkUpdate() {
  if (!auth.token || !selectedCount.value || !bulkUpdateField.value) return
  bulkUpdateSaving.value = true
  error.value = ''
  try {
    lastActionResult.value = await generatedApi.resourceAction(resourceName.value, 'bulk-update', selectionRequest({ [bulkUpdateField.value]: bulkUpdateValue.value }), auth.token)
    bulkUpdateDialogOpen.value = false
    clearSelection()
    await loadRows(meta.value.page)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    bulkUpdateSaving.value = false
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
watch(resourceName, () => { filterValues.value = {}; trashed.value = 'default'; clearSelection(); void loadRows(1) })
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ currentManifest?.label || t('resource.title') }}</h1>
        <p class="text-sm text-muted-foreground">{{ t('resource.description') }}</p>
      </div>
      <div class="flex gap-2"><Button v-if="canCreate" variant="default" @click="router.push(`/${resourceName}/new`)">{{ t('resource.create') }}</Button><Button variant="outline" :disabled="loading || exporting" @click="exportRows"><Download data-icon="inline-start" />{{ t('resource.export') }}</Button><Button variant="outline" :disabled="loading" @click="loadRows(meta.page)">
        <RefreshCw data-icon="inline-start" />{{ t('resource.refresh') }}
      </Button></div>
    </div>

    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Alert v-if="lastActionResult"><AlertTitle>{{ t('resource.actionCompleted') }}</AlertTitle><AlertDescription>{{ t('resource.actionResult', { succeeded: lastActionResult.succeeded, failed: lastActionResult.failed, skipped: lastActionResult.skipped }) }}</AlertDescription></Alert>
    <Card>
      <CardHeader class="gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div><CardTitle>{{ currentManifest?.label || t('resource.resourceNotFound') }}</CardTitle><CardDescription>{{ t('resource.total', { count: meta.total }) }}</CardDescription></div>
        <form class="flex w-full flex-wrap gap-2 sm:w-auto" @submit.prevent="submitSearch"><Select v-if="currentManifest?.soft_delete" :model-value="trashed" @update:model-value="(value) => { trashed = String(value); clearSelection(); void loadRows(1) }"><SelectTrigger class="w-32" :aria-label="t('resource.trashFilter')"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="default">{{ t('resource.activeRecords') }}</SelectItem><SelectItem value="only">{{ t('resource.trashedRecords') }}</SelectItem><SelectItem value="with">{{ t('resource.allRecords') }}</SelectItem></SelectContent></Select><template v-for="field in filters" :key="field.name"><Select v-if="field.type === 'select' || field.type === 'multi-select' || field.type === 'boolean'" :model-value="filterValues[field.name] || 'all'" @update:model-value="changeResourceFilter(field.name, $event)"><SelectTrigger class="w-36" :aria-label="field.label"><SelectValue :placeholder="field.label" /></SelectTrigger><SelectContent><SelectItem v-for="option in filterOptions(field)" :key="option.value" :value="option.value">{{ option.label }}</SelectItem></SelectContent></Select><Input v-else v-model="filterValues[field.name]" class="w-44" :type="field.type === 'date-range' ? 'text' : 'search'" :placeholder="field.type === 'date-range' ? `${field.label} (YYYY-MM-DD..YYYY-MM-DD)` : field.label" /></template><Input v-model="search" class="sm:w-64" :placeholder="t('resource.searchPlaceholder')" :aria-label="t('resource.search')" /><Button type="submit" size="icon" :aria-label="t('resource.search')"><Search /></Button></form>
      </CardHeader>
      <div v-if="selectedCount && batchActions.length" class="overflow-x-auto border-y bg-muted/20 px-4 py-2">
        <div class="flex min-w-max items-center gap-2">
          <span class="text-sm text-muted-foreground">{{ t('resource.selectedCount', { count: selectedCount }) }}</span>
          <Button v-if="!allFilteredSelected && selectedCount === rows.length && meta.total > rows.length" variant="link" size="sm" @click="selectAllFiltered">{{ t('resource.selectAllFiltered', { count: meta.total }) }}</Button>
          <Button v-for="action in batchActions" :key="action.name" size="sm" @click="openBulkAction(action)">{{ action.label }}</Button>
          <Button variant="ghost" size="sm" @click="clearSelection">{{ t('resource.clearSelection') }}</Button>
        </div>
      </div>
      <CardContent>
        <div v-if="loading" class="flex flex-col gap-3"><Skeleton v-for="item in 5" :key="item" class="h-10" /></div>
        <Empty v-else-if="!rows.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('resource.noData') }}</EmptyDescription></EmptyHeader></Empty>
        <Table v-else><TableHeader><TableRow><TableHead v-if="batchActions.length" class="w-10"><Checkbox v-model="headerChecked" :aria-label="t('resource.selectAll')" /></TableHead><TableHead v-for="column in visibleColumns" :key="column.name"><Button v-if="column.sortable" variant="ghost" size="sm" class="-ml-3" @click="sortBy(column)">{{ column.label }}<ArrowDownUp data-icon="inline-end" /></Button><span v-else>{{ column.label }}</span></TableHead><TableHead class="w-36 text-right">{{ t('resource.actions') }}</TableHead></TableRow></TableHeader><TableBody><TableRow v-for="(row, index) in rows" :key="String(row.id || index)" class="cursor-pointer" @click="row.id && router.push(`/${resourceName}/${row.id}`)"><TableCell v-if="batchActions.length" @click.stop><Checkbox :model-value="allFilteredSelected ? !excludedIds.includes(String(row.id)) : selectedIds.includes(String(row.id))" :aria-label="t('resource.selectRow', { name: row.name })" @update:model-value="toggleRow(row.id, $event)" /></TableCell><TableCell v-for="column in visibleColumns" :key="column.name"><Badge v-if="column.name === 'status'" variant="secondary">{{ statusLabel(column.name, row[column.name]) }}</Badge><template v-else>{{ displayValue(row[column.name]) }}</template></TableCell><TableCell class="text-right"><div v-if="row.id && (hasAction('update') || hasAction('delete') || hasActionKind('user-status'))" class="flex justify-end gap-1" @click.stop><Button v-if="hasActionKind('user-status')" variant="ghost" size="sm" :aria-label="t('resource.setStatus')" @click="openRowStatusAction(row)">{{ t('resource.setStatus') }}</Button><Button v-if="hasAction('update')" variant="ghost" size="sm" :aria-label="t('resource.edit')" @click="router.push(editPath(row))"><Pencil data-icon="inline-start" />{{ t('resource.edit') }}</Button><Button v-if="hasAction('delete')" variant="ghost" size="sm" :aria-label="t('resource.delete')" @click="openDelete(row)"><Trash2 data-icon="inline-start" />{{ t('resource.delete') }}</Button></div></TableCell></TableRow></TableBody></Table>
        <div v-if="!loading && meta.total > 0" class="mt-4 flex flex-col items-center gap-3 sm:flex-row sm:justify-between"><p class="text-sm text-muted-foreground">{{ t('resource.page', { page: meta.page }) }}</p><div class="flex items-center gap-2"><Select :model-value="pageSize" @update:model-value="changePageSize"><SelectTrigger class="w-24"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="10">{{ t('resource.perPage', { count: 10 }) }}</SelectItem><SelectItem value="20">{{ t('resource.perPage', { count: 20 }) }}</SelectItem><SelectItem value="50">{{ t('resource.perPage', { count: 50 }) }}</SelectItem></SelectContent></Select></div><Pagination v-model:page="meta.page" :items-per-page="meta.per_page" :total="meta.total" @update:page="loadRows"><PaginationContent v-slot="{ items }"><PaginationPrevious /><template v-for="(item, index) in items" :key="index"><PaginationItem v-if="item.type === 'page'" :value="item.value" :is-active="item.value === meta.page">{{ item.value }}</PaginationItem></template><PaginationNext /></PaginationContent></Pagination></div>
      </CardContent>
    </Card>
    <AlertDialog v-model:open="deleteDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('resource.deleteTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resource.deleteDescription') }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('resource.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="deleting" @click="deleteRow">{{ t('resource.delete') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
    <AlertDialog v-model:open="bulkStatusDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('resource.bulkStatusTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resource.bulkStatusDescription', { count: selectedCount, status: bulkStatusLabel }) }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('resource.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="bulkUpdating" @click="applyBulkStatus">{{ t('resource.applyStatus') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
    <AlertDialog v-model:open="bulkDeleteDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('resource.bulkDeleteTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resource.bulkDeleteDescription', { count: selectedCount }) }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('resource.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="bulkDeleting" @click="applyBulkDelete">{{ t('resource.delete') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
    <Dialog v-model:open="bulkUpdateDialogOpen"><DialogContent><DialogHeader><DialogTitle>{{ t('resource.bulkUpdateTitle') }}</DialogTitle><DialogDescription>{{ t('resource.bulkUpdateDescription', { count: selectedCount }) }}</DialogDescription></DialogHeader><div class="grid gap-3"><Select v-model="bulkUpdateField"><SelectTrigger><SelectValue :placeholder="t('resource.bulkUpdateField')" /></SelectTrigger><SelectContent><SelectItem v-for="field in writableFields" :key="field.name" :value="field.name">{{ field.label }}</SelectItem></SelectContent></Select><Input v-model="bulkUpdateValue" :placeholder="t('resource.bulkUpdateValue')" /></div><DialogFooter><Button variant="outline" @click="bulkUpdateDialogOpen = false">{{ t('resource.cancel') }}</Button><Button :disabled="bulkUpdateSaving || !bulkUpdateField" @click="applyBulkUpdate">{{ t('resource.applyStatus') }}</Button></DialogFooter></DialogContent></Dialog>
  </div>
</template>
