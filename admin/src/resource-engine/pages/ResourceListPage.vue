<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { FlexRender, tableFeatures, useTable, type ColumnDef } from '@tanstack/vue-table'
import { ArrowDown, ArrowUp, ArrowUpDown } from '@lucide/vue'
import { toast } from 'vue-sonner'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import {
  Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious,
} from '@/components/ui/pagination'
import {
  Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select'
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from '@/components/ui/table'

import { createResourceContext } from '../core/ResourceContext'
import { useAuth } from '@/core/auth'
import type { ResourceListResult, ResourceSort } from '../core/ResourceDataProvider'
import { resourceRegistry } from '../demo'
import { serializeResourceQuery } from '../query/serializeResourceQuery'
import { nextResourceSort } from '../table/state'

const route = useRoute()
const router = useRouter()
const auth = useAuth()
const { t } = useI18n()
const resourceName = computed(() => String(route.params.resource ?? ''))
const definition = computed(() => resourceRegistry.get(resourceName.value))
const provider = computed(() => resourceRegistry.provider(resourceName.value))
const context = computed(() => definition.value && provider.value
  ? createResourceContext(definition.value, provider.value, auth.user?.permissions ?? [])
  : null)
const rows = ref<Record<string, unknown>[]>([])
const search = ref('')
const page = ref(1)
const loading = ref(false)
const error = ref<string | null>(null)
const selectedIds = ref<string[]>([])
const filters = ref<Record<string, string>>({})
const sort = ref<ResourceSort>()
const pagination = ref({ page: 1, perPage: 10, total: 0, totalPages: 0 })
const serializedQuery = ref('')
const deleteTargetId = ref<string | null>(null)
const deleteDialogOpen = ref(false)
const bulkDeleteDialogOpen = ref(false)
const deleting = ref(false)
type ResourceRow = Record<string, unknown>
const features = tableFeatures({})

const columns = ref<ColumnDef<typeof features, ResourceRow>[]>([])
watch(definition, (value) => {
  columns.value = (value?.columns ?? []).map((column) => ({
    accessorKey: String(column.key),
    header: column.label,
    cell: (info) => String(info.getValue() ?? ''),
  }))
}, { immediate: true })
const table = useTable<typeof features, ResourceRow>({ features, columns, data: rows })
const allVisibleSelected = computed(() => rows.value.length > 0 && rows.value.every(isSelected))

async function load() {
  if (!provider.value || !definition.value) return
  loading.value = true
  error.value = null
  serializedQuery.value = serializeResourceQuery({
    page: page.value,
    perPage: 10,
    search: search.value,
    filters: filters.value,
    sort: sort.value,
  })
  try {
    const result = await provider.value.list({
      page: page.value,
      perPage: 10,
      search: search.value,
      filters: filters.value,
      sort: sort.value,
    }) as unknown as ResourceListResult<ResourceRow>
    rows.value = result.data
    pagination.value = result.meta.pagination
    selectedIds.value = []
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('resources.loadFailed')
  } finally {
    loading.value = false
  }
}

function idOf(row: Record<string, unknown>): string {
  return String(row[String(definition.value?.primaryKey ?? 'id')] ?? '')
}

function isSelected(row: Record<string, unknown>): boolean {
  return selectedIds.value.includes(idOf(row))
}

function toggleSelected(row: Record<string, unknown>, checked: boolean) {
  const id = idOf(row)
  selectedIds.value = checked ? [...new Set([...selectedIds.value, id])] : selectedIds.value.filter((value) => value !== id)
}

function toggleAll(checked: boolean) {
  selectedIds.value = checked ? rows.value.map(idOf) : []
}

function setFilter(field: string, value: string) {
  const next = { ...filters.value }
  if (!value || value === 'all') delete next[field]
  else next[field] = value
  filters.value = next
}

function toggleSort(field: string) {
  sort.value = nextResourceSort(sort.value, field)
}

function sortIcon(field: string) {
  if (sort.value?.field !== field) return ArrowUpDown
  return sort.value.direction === 'asc' ? ArrowUp : ArrowDown
}

function openDelete(id: string) {
  deleteTargetId.value = id
  deleteDialogOpen.value = true
}

async function confirmDelete() {
  if (!provider.value || !deleteTargetId.value) return
  deleting.value = true
  try {
    await provider.value.delete(deleteTargetId.value)
    toast.success(t('resources.deleted'))
    deleteDialogOpen.value = false
    await load()
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : t('resources.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

function openBulkDelete() {
  if (selectedIds.value.length) bulkDeleteDialogOpen.value = true
}

async function confirmBulkDelete() {
  if (!provider.value || !selectedIds.value.length) return
  deleting.value = true
  try {
    await provider.value.bulkDelete(selectedIds.value)
    toast.success(t('resources.bulkDeleted'))
    bulkDeleteDialogOpen.value = false
    await load()
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : t('resources.bulkDeleteFailed'))
  } finally {
    deleting.value = false
  }
}

watch([search, filters, sort], () => {
  page.value = 1
  void load()
}, { deep: true })
watch(page, () => void load())
onMounted(() => void load())
</script>

<template>
  <section v-if="definition && provider && context" class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ definition.label }}</h1>
        <p class="text-sm text-muted-foreground">{{ t('resources.records') }} · {{ serializedQuery || t('resources.waiting') }}</p>
      </div>
      <Button v-if="context.can('create')" @click="router.push(`/admin/resources/${resourceName}/create`)">{{ t('resources.create') }}</Button>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('resources.loadFailed') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card>
      <CardHeader>
        <div class="flex flex-wrap items-center justify-between gap-3">
              <div><CardTitle>{{ t('resources.records') }}</CardTitle><CardDescription>{{ t('resources.features') }}</CardDescription></div>
          <div class="flex flex-wrap items-center gap-2">
            <Input v-model="search" class="w-56" :placeholder="t('resources.search')" :aria-label="t('resources.searchAria')" />
            <Select
              v-for="filter in definition.filters ?? []"
              :key="String(filter.field)"
              :model-value="filters[String(filter.field)] ?? 'all'"
              @update:model-value="setFilter(String(filter.field), String($event))"
            >
              <SelectTrigger :aria-label="t('resources.filterAria', { label: filter.label })" class="w-36"><SelectValue :placeholder="filter.label" /></SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem v-for="option in filter.options ?? []" :key="option.value || 'all'" :value="option.value || 'all'">{{ option.label }}</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <Button v-if="selectedIds.length && context.can('bulkDelete')" variant="destructive" @click="openBulkDelete">{{ t('resources.bulkDelete') }}</Button>
          </div>
        </div>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <div class="overflow-x-auto rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-12"><Checkbox :model-value="allVisibleSelected" :aria-label="t('resources.selectPage')" @update:model-value="toggleAll(Boolean($event))" /></TableHead>
                <TableHead v-for="column in definition.columns ?? []" :key="String(column.key)">
                  <Button v-if="column.sortable" variant="ghost" size="sm" :aria-label="t('resources.sortBy', { label: column.label })" @click="toggleSort(String(column.key))">
                    {{ column.label }}
                    <component :is="sortIcon(String(column.key))" data-icon="inline-end" />
                  </Button>
                  <span v-else>{{ column.label }}</span>
                </TableHead>
                <TableHead>{{ t('resources.actions') }}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="loading"><TableCell :colspan="(columns.length || 1) + 2">{{ t('resources.loading') }}</TableCell></TableRow>
              <TableRow v-else-if="!table.getRowModel().rows.length"><TableCell :colspan="(columns.length || 1) + 2"><Empty class="border-0"><EmptyHeader><EmptyTitle>{{ t('resources.noData') }}</EmptyTitle><EmptyDescription>{{ t('resources.features') }}</EmptyDescription></EmptyHeader></Empty></TableCell></TableRow>
              <TableRow v-for="row in table.getRowModel().rows" v-else :key="row.id">
                <TableCell><Checkbox :model-value="isSelected(row.original)" :aria-label="t('resources.select', { id: idOf(row.original) })" @update:model-value="toggleSelected(row.original, Boolean($event))" /></TableCell>
                <TableCell v-for="cell in row.getAllCells()" :key="cell.id"><FlexRender :cell="cell" /></TableCell>
                <TableCell class="flex gap-2">
                  <Button variant="ghost" size="sm" @click="router.push(`/admin/resources/${resourceName}/${idOf(row.original)}`)">{{ t('resources.view') }}</Button>
                  <Button v-if="context.can('update')" variant="ghost" size="sm" @click="router.push(`/admin/resources/${resourceName}/${idOf(row.original)}/edit`)">{{ t('resources.edit') }}</Button>
                  <Button v-if="context.can('delete')" variant="ghost" size="sm" @click="openDelete(idOf(row.original))">{{ t('resources.delete') }}</Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
        <Pagination v-if="pagination.totalPages > 1" v-model:page="page" :items-per-page="pagination.perPage" :total="pagination.total">
          <PaginationContent>
            <PaginationItem :value="Math.max(1, pagination.page - 1)"><PaginationPrevious /></PaginationItem>
            <PaginationItem :value="pagination.page"><span class="px-3 text-sm">{{ t('resources.page', { page: pagination.page, total: pagination.totalPages }) }}</span></PaginationItem>
            <PaginationItem :value="Math.min(pagination.totalPages, pagination.page + 1)"><PaginationNext /></PaginationItem>
          </PaginationContent>
        </Pagination>
      </CardContent>
    </Card>

    <AlertDialog v-model:open="deleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader><AlertDialogTitle>{{ t('resources.confirmDelete') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resources.deleteDescription') }}</AlertDialogDescription></AlertDialogHeader>
        <AlertDialogFooter><AlertDialogCancel>{{ t('resources.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="deleting" variant="destructive" @click="confirmDelete">{{ t('resources.confirm') }}</AlertDialogAction></AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
    <AlertDialog v-model:open="bulkDeleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader><AlertDialogTitle>{{ t('resources.confirmBulkDelete') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resources.bulkDescription', { count: selectedIds.length }) }}</AlertDialogDescription></AlertDialogHeader>
        <AlertDialogFooter><AlertDialogCancel>{{ t('resources.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="deleting" variant="destructive" @click="confirmBulkDelete">{{ t('resources.confirm') }}</AlertDialogAction></AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </section>
  <Alert v-else variant="destructive"><AlertTitle>{{ t('resources.notFound') }}</AlertTitle><AlertDescription>{{ t('resources.notRegistered', { name: resourceName }) }}</AlertDescription></Alert>
</template>
