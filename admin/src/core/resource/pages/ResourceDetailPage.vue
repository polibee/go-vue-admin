<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Pencil, Trash2 } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Skeleton } from '@/components/ui/skeleton'
import { ApiError, errorMessageKey } from '@/lib/api'
import { generatedApi, type RelationOption, type ResourceDetailSection, type ResourceManifest, type RolePermissionAssignment } from '@/generated/api'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { localizedFieldLabel, localizedResourceLabel } from '@/core/resource/resource-i18n'
import { adminRoute } from '@/core/routing/url-namespaces'
import RolePermissionsPanel from '@/modules/roles/components/RolePermissionsPanel.vue'

const { t, te } = useI18n()
const props = defineProps<{ resource?: string }>()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
interface ResourceAction { name: string; permission: string }
const data = ref<Record<string, unknown>>({})
const manifest = ref<ResourceManifest>()
const relationRecords = ref<Record<string, RelationOption[]>>({})
const loading = ref(true)
const error = ref('')
const deleteDialogOpen = ref(false)
const deleting = ref(false)
const resourceName = computed(() => props.resource || String(route.params.resource || 'users'))
const displayData = computed(() => Object.fromEntries(Object.entries(data.value).filter(([key]) => key === 'id' || manifest.value?.fields.some((field) => field.name === key && field.visible !== false && field.readable !== false))))
const hasAction = (name: string) => {
  const action = manifest.value?.actions?.find((item) => item.name === name)
  return Boolean(action && auth.can(action.permission))
}
const canEdit = computed(() => hasAction('update'))
const canDelete = computed(() => hasAction('delete'))
const detailSections = computed<ResourceDetailSection[]>(() => manifest.value?.details || [])
const hasManyRelations = computed(() => (manifest.value?.relations || []).filter((relation) => relation.kind === 'hasMany'))
const resourceLabel = computed(() => localizedResourceLabel(t, te, resourceName.value, manifest.value?.label || resourceName.value))
const isRoleDetail = computed(() => resourceName.value === 'roles')
const rolePermissions = computed(() => Array.isArray(data.value.permissions) ? data.value.permissions as RolePermissionAssignment[] : [])

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

onMounted(async () => {
  if (!auth.token) return
  try {
    const [manifests, record] = await Promise.all([
      generatedApi.resourceRegistry(auth.token),
      generatedApi.resourceShow<Record<string, unknown>>(resourceName.value, String(route.params.id), auth.token),
    ])
    manifest.value = manifests.find((item) => item.name === resourceName.value)
    data.value = record
    for (const relation of hasManyRelations.value) {
      const response = await generatedApi.resourceRelationRecords(resourceName.value, String(route.params.id), relation.name, auth.token)
      relationRecords.value = { ...relationRecords.value, [relation.name]: response.data }
    }
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
})

function fieldLabel(name: string) {
  const field = manifest.value?.fields.find((item) => item.name === name)
  return localizedFieldLabel(t, te, resourceName.value, name, field?.label || name)
}

async function deleteRecord() {
  if (!auth.token || !data.value.id || !canDelete.value) return
  deleting.value = true
  error.value = ''
  try {
    await generatedApi.resourceDelete(resourceName.value, String(route.params.id), auth.token)
    await router.push(adminRoute(resourceName.value))
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    deleting.value = false
    deleteDialogOpen.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center gap-3">
      <Button variant="ghost" size="icon" :aria-label="t('resource.back')" @click="router.back()"><ArrowLeft /></Button>
      <div class="flex-1"><h1 class="text-2xl font-semibold tracking-tight">{{ t('resource.detail') }}</h1><p class="text-sm text-muted-foreground">{{ resourceLabel }} #{{ route.params.id }}</p></div><div v-if="canEdit || canDelete" class="flex gap-2"><Button v-if="canEdit" variant="outline" @click="router.push(`${adminRoute(resourceName)}/${route.params.id}/edit`)"><Pencil data-icon="inline-start" />{{ t('resource.edit') }}</Button><Button v-if="canDelete" variant="destructive" @click="deleteDialogOpen = true"><Trash2 data-icon="inline-start" />{{ t('resource.delete') }}</Button></div>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="loading"><CardHeader><Skeleton class="h-6 w-40" /><Skeleton class="h-4 w-64" /></CardHeader><CardContent class="flex flex-col gap-3"><Skeleton v-for="item in 4" :key="item" class="h-10" /></CardContent></Card>
    <Empty v-else-if="error"><EmptyHeader><EmptyTitle>{{ t('states.errorTitle') }}</EmptyTitle><EmptyDescription>{{ error }}</EmptyDescription></EmptyHeader></Empty>
    <Card v-else-if="Object.keys(displayData).length"><CardHeader><CardTitle>{{ String(data.display_name || data.name || data.email || route.params.id) }}</CardTitle><CardDescription>{{ t('resource.detailDescription') }}</CardDescription></CardHeader><CardContent><dl class="grid gap-4 sm:grid-cols-2"> <div v-for="(value, key) in displayData" :key="key" class="rounded-md border p-3"><dt class="text-xs text-muted-foreground">{{ fieldLabel(String(key)) }}</dt><dd class="mt-1 break-words text-sm font-medium">{{ value === null || value === undefined ? '—' : String(value) }}</dd></div></dl></CardContent></Card>
    <Card v-for="section in detailSections" :key="section.name"><CardHeader><CardTitle>{{ section.label }}</CardTitle></CardHeader><CardContent><dl class="grid gap-4 sm:grid-cols-2"><div v-for="field in section.fields" v-show="displayData[field] !== undefined" :key="field" class="rounded-md border p-3"><dt class="text-xs text-muted-foreground">{{ fieldLabel(field) }}</dt><dd class="mt-1 break-words text-sm font-medium">{{ displayData[field] === null || displayData[field] === undefined ? '—' : String(displayData[field]) }}</dd></div></dl></CardContent></Card>
    <Card v-for="relation in hasManyRelations" :key="relation.name"><CardHeader><CardTitle>{{ relation.name }}</CardTitle><CardDescription>{{ t('resource.relatedRecords') }}</CardDescription></CardHeader><CardContent><div v-if="relationRecords[relation.name]?.length" class="flex flex-wrap gap-2"><Button v-for="item in relationRecords[relation.name]" :key="item.value" variant="outline" size="sm">{{ item.label }}</Button></div><p v-else class="text-sm text-muted-foreground">{{ t('resource.noRelatedRecords') }}</p></CardContent></Card>
    <RolePermissionsPanel v-if="isRoleDetail" :role-id="Number(route.params.id)" :role-name="String(data.name || '')" :assignments="rolePermissions" />
    <Empty v-if="!loading && !error && !Object.keys(displayData).length && !detailSections.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('resource.noData') }}</EmptyDescription></EmptyHeader></Empty>
    <AlertDialog v-model:open="deleteDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('resource.deleteTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resource.deleteDescription') }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('resource.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="deleting" @click="deleteRecord">{{ t('resource.delete') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
  </div>
</template>
