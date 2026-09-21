<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDownUp, Pencil, RefreshCw, Search, Trash2 } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious } from '@/components/ui/pagination'
import { Skeleton } from '@/components/ui/skeleton'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ApiError, apiFetch, errorMessageKey } from '@/lib/api'
import { canDeleteResource, resourceActionPath } from '@/lib/resource-actions'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

interface ResourceColumn { name: string; label: string; sortable: boolean }
interface ResourceManifest { name: string; label: string; route: string; columns: ResourceColumn[] }
interface ResourceMeta { page: number; per_page: number; total: number; last_page: number }
interface ResourceListResponse { data: Record<string, unknown>[]; meta: ResourceMeta }

const { t } = useI18n()
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
const error = ref('')
const deleteDialogOpen = ref(false)
const deleting = ref(false)
const deleteTarget = ref<Record<string, unknown>>()

const resourceName = computed(() => String(route.params.resource || 'users'))
const currentManifest = computed(() => manifests.value.find((item) => item.name === resourceName.value))

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

async function loadManifests() {
  if (!auth.token) return
  manifests.value = await apiFetch<ResourceManifest[]>('/api/v1/admin/resources', {}, auth.token)
  if (!currentManifest.value && manifests.value.length) {
    await router.replace(`/${manifests.value[0].name}`)
  }
}

async function loadRows(page = 1) {
  if (!auth.token || !currentManifest.value) return
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams({ page: String(page), per_page: pageSize.value, sort: sort.value, dir: direction.value })
    if (search.value.trim()) params.set('search', search.value.trim())
    const response = await apiFetch<ResourceListResponse>(`/api/v1/admin/resources/${resourceName.value}?${params}`, {}, auth.token)
    rows.value = response.data
    meta.value = response.meta
  } catch (value) {
    rows.value = []
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
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
function editPath(row: Record<string, unknown>) { return `/${resourceName.value}/${row.id}/edit` }
function openDelete(row: Record<string, unknown>) { deleteTarget.value = row; deleteDialogOpen.value = true }
async function deleteRow() {
  if (!auth.token || !deleteTarget.value?.id) return
  deleting.value = true
  error.value = ''
  try {
    await apiFetch(resourceActionPath(resourceName.value, String(deleteTarget.value.id), 'delete'), { method: 'DELETE' }, auth.token)
    deleteDialogOpen.value = false
    await loadRows(meta.value.page)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    deleting.value = false
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
watch(resourceName, () => { void loadRows(1) })
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ t('resource.title') }}</h1>
        <p class="text-sm text-muted-foreground">{{ t('resource.description') }}</p>
      </div>
      <div class="flex gap-2"><Button v-if="resourceName === 'users'" variant="default" @click="router.push('/users/new')">{{ t('resource.createUser') }}</Button><Button v-if="resourceName === 'roles'" variant="default" @click="router.push('/roles/new')">{{ t('rbac.createRole') }}</Button><Button variant="outline" :disabled="loading" @click="loadRows(meta.page)">
        <RefreshCw data-icon="inline-start" />{{ t('resource.refresh') }}
      </Button></div>
    </div>

    <div v-if="manifests.length" class="flex flex-wrap gap-2">
      <Button v-for="manifest in manifests" :key="manifest.name" :variant="manifest.name === resourceName ? 'default' : 'outline'" size="sm" @click="selectResource(manifest.name)">{{ manifest.label }}</Button>
    </div>

    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card>
      <CardHeader class="gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div><CardTitle>{{ currentManifest?.label || t('resource.resourceNotFound') }}</CardTitle><CardDescription>{{ t('resource.total', { count: meta.total }) }}</CardDescription></div>
        <form class="flex w-full gap-2 sm:w-auto" @submit.prevent="submitSearch"><Input v-model="search" class="sm:w-64" :placeholder="t('resource.searchPlaceholder')" :aria-label="t('resource.search')" /><Button type="submit" size="icon" :aria-label="t('resource.search')"><Search /></Button></form>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="flex flex-col gap-3"><Skeleton v-for="item in 5" :key="item" class="h-10" /></div>
        <Empty v-else-if="!rows.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('resource.noData') }}</EmptyDescription></EmptyHeader></Empty>
        <Table v-else><TableHeader><TableRow><TableHead v-for="column in currentManifest?.columns || []" :key="column.name"><Button v-if="column.sortable" variant="ghost" size="sm" class="-ml-3" @click="sortBy(column)">{{ column.label }}<ArrowDownUp data-icon="inline-end" /></Button><span v-else>{{ column.label }}</span></TableHead><TableHead class="w-36 text-right">{{ t('resource.actions') }}</TableHead></TableRow></TableHeader><TableBody><TableRow v-for="(row, index) in rows" :key="String(row.id || index)" class="cursor-pointer" @click="row.id && router.push(`/${resourceName}/${row.id}`)"><TableCell v-for="column in currentManifest?.columns || []" :key="column.name">{{ displayValue(row[column.name]) }}</TableCell><TableCell class="text-right"><div v-if="row.id && (resourceName === 'users' || resourceName === 'roles')" class="flex justify-end gap-1" @click.stop><Button variant="ghost" size="sm" :aria-label="t('resource.edit')" @click="router.push(editPath(row))"><Pencil data-icon="inline-start" />{{ t('resource.edit') }}</Button><Button variant="ghost" size="sm" :disabled="!canDeleteResource(resourceName, row)" :aria-label="t('resource.delete')" @click="openDelete(row)"><Trash2 data-icon="inline-start" />{{ t('resource.delete') }}</Button></div></TableCell></TableRow></TableBody></Table>
        <div v-if="!loading && meta.total > 0" class="mt-4 flex flex-col items-center gap-3 sm:flex-row sm:justify-between"><p class="text-sm text-muted-foreground">{{ t('resource.page', { page: meta.page }) }}</p><div class="flex items-center gap-2"><Select :model-value="pageSize" @update:model-value="changePageSize"><SelectTrigger class="w-24"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="10">{{ t('resource.perPage', { count: 10 }) }}</SelectItem><SelectItem value="20">{{ t('resource.perPage', { count: 20 }) }}</SelectItem><SelectItem value="50">{{ t('resource.perPage', { count: 50 }) }}</SelectItem></SelectContent></Select></div><Pagination v-model:page="meta.page" :items-per-page="meta.per_page" :total="meta.total" @update:page="loadRows"><PaginationContent v-slot="{ items }"><PaginationPrevious /><template v-for="(item, index) in items" :key="index"><PaginationItem v-if="item.type === 'page'" :value="item.value" :is-active="item.value === meta.page">{{ item.value }}</PaginationItem></template><PaginationNext /></PaginationContent></Pagination></div>
      </CardContent>
    </Card>
    <AlertDialog v-model:open="deleteDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('resource.deleteTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resource.deleteDescription') }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('resource.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="deleting" @click="deleteRow">{{ t('resource.delete') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
  </div>
</template>
