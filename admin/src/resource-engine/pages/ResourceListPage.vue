<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { FlexRender, tableFeatures, useTable, type ColumnDef } from '@tanstack/vue-table'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import {
  Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious,
} from '@/components/ui/pagination'
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from '@/components/ui/table'

import { createResourceContext } from '../core/ResourceContext'
import type { ResourceListResult } from '../core/ResourceDataProvider'
import { resourceRegistry } from '../demo'
import { serializeResourceQuery } from '../query/serializeResourceQuery'

const route = useRoute()
const router = useRouter()
const resourceName = computed(() => String(route.params.resource ?? ''))
const definition = computed(() => resourceRegistry.get(resourceName.value))
const provider = computed(() => resourceRegistry.provider(resourceName.value))
const context = computed(() => definition.value && provider.value
  ? createResourceContext(definition.value, provider.value, ['dashboard.view'])
  : null)
const rows = ref<Record<string, unknown>[]>([])
const search = ref('')
const page = ref(1)
const loading = ref(false)
const error = ref<string | null>(null)
const selectedIds = ref<string[]>([])
const pagination = ref({ page: 1, perPage: 10, total: 0, totalPages: 0 })
const serializedQuery = ref('')
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

async function load() {
  if (!provider.value || !definition.value) return
  loading.value = true
  error.value = null
  serializedQuery.value = serializeResourceQuery({ page: page.value, perPage: 10, search: search.value })
  try {
    const result = await provider.value.list({ page: page.value, perPage: 10, search: search.value }) as unknown as ResourceListResult<ResourceRow>
    rows.value = result.data
    pagination.value = result.meta.pagination
    selectedIds.value = []
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '资源加载失败'
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

async function bulkDelete() {
  if (!provider.value || !selectedIds.value.length) return
  await provider.value.bulkDelete(selectedIds.value)
  await load()
}

watch(search, () => {
  page.value = 1
  void load()
})
watch(page, () => void load())
onMounted(() => void load())
</script>

<template>
  <section v-if="definition && provider && context" class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ definition.label }}</h1>
        <p class="text-sm text-muted-foreground">通用资源列表 · {{ serializedQuery || '等待查询' }}</p>
      </div>
      <Button v-if="context.can('create')" @click="router.push(`/admin/resources/${resourceName}/create`)">新建</Button>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>加载失败</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card>
      <CardHeader>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div><CardTitle>记录</CardTitle><CardDescription>支持搜索、选择、批量操作和分页。</CardDescription></div>
          <div class="flex items-center gap-2">
            <Input v-model="search" class="w-56" placeholder="搜索记录" aria-label="搜索资源" />
            <Button v-if="selectedIds.length && context.can('bulkDelete')" variant="destructive" @click="bulkDelete">批量删除</Button>
          </div>
        </div>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <div class="overflow-x-auto rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-12"><span class="sr-only">选择</span></TableHead>
                <TableHead v-for="header in table.getHeaderGroups()[0]?.headers" :key="header.id"><FlexRender :header="header" /></TableHead>
                <TableHead>操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="loading"><TableCell :colspan="(columns.length || 1) + 2">加载中…</TableCell></TableRow>
              <TableRow v-else-if="!table.getRowModel().rows.length"><TableCell :colspan="(columns.length || 1) + 2">暂无数据</TableCell></TableRow>
              <TableRow v-for="row in table.getRowModel().rows" v-else :key="row.id">
                <TableCell><Checkbox :checked="isSelected(row.original)" :aria-label="`选择 ${idOf(row.original)}`" @update:checked="toggleSelected(row.original, Boolean($event))" /></TableCell>
                <TableCell v-for="cell in row.getAllCells()" :key="cell.id"><FlexRender :cell="cell" /></TableCell>
                <TableCell class="flex gap-2">
                  <Button variant="ghost" size="sm" @click="router.push(`/admin/resources/${resourceName}/${idOf(row.original)}`)">查看</Button>
                  <Button v-if="context.can('update')" variant="ghost" size="sm" @click="router.push(`/admin/resources/${resourceName}/${idOf(row.original)}/edit`)">编辑</Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
        <Pagination v-if="pagination.totalPages > 1" v-model:page="page" :items-per-page="pagination.perPage" :total="pagination.total">
          <PaginationContent>
            <PaginationItem :value="Math.max(1, pagination.page - 1)"><PaginationPrevious /></PaginationItem>
            <PaginationItem :value="pagination.page"><span class="px-3 text-sm">第 {{ pagination.page }} / {{ pagination.totalPages }} 页</span></PaginationItem>
            <PaginationItem :value="Math.min(pagination.totalPages, pagination.page + 1)"><PaginationNext /></PaginationItem>
          </PaginationContent>
        </Pagination>
      </CardContent>
    </Card>
  </section>
  <Alert v-else variant="destructive"><AlertTitle>资源不存在</AlertTitle><AlertDescription>未注册资源：{{ resourceName }}</AlertDescription></Alert>
</template>
