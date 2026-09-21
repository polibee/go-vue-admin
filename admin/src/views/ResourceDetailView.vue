<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Pencil, Trash2 } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Skeleton } from '@/components/ui/skeleton'
import { ApiError, apiFetch, errorMessageKey } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
interface ResourceField { name: string; label: string; type: string }
interface ResourceManifest { name: string; fields: ResourceField[] }
const data = ref<Record<string, unknown>>({})
const manifest = ref<ResourceManifest>()
const loading = ref(true)
const error = ref('')
const deleteDialogOpen = ref(false)
const deleting = ref(false)

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

onMounted(async () => {
  if (!auth.token) return
  try {
    const [manifests, record] = await Promise.all([
      apiFetch<ResourceManifest[]>('/api/v1/admin/resources', {}, auth.token),
      apiFetch<Record<string, unknown>>(`/api/v1/admin/resources/${route.params.resource}/${route.params.id}`, {}, auth.token),
    ])
    manifest.value = manifests.find((item) => item.name === route.params.resource)
    data.value = record
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
})

function fieldLabel(name: string) {
  return manifest.value?.fields.find((field) => field.name === name)?.label || name
}

async function deleteUser() {
  if (!auth.token) return
  deleting.value = true
  error.value = ''
  try {
    await apiFetch(`/api/v1/admin/users/${route.params.id}`, { method: 'DELETE' }, auth.token)
    await router.push('/users')
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
      <div class="flex-1"><h1 class="text-2xl font-semibold tracking-tight">{{ t('resource.detail') }}</h1><p class="text-sm text-muted-foreground">{{ route.params.resource }} #{{ route.params.id }}</p></div><div v-if="route.params.resource === 'users'" class="flex gap-2"><Button variant="outline" @click="router.push(`/users/${route.params.id}/edit`)"><Pencil data-icon="inline-start" />{{ t('resource.editUser') }}</Button><Button variant="destructive" @click="deleteDialogOpen = true"><Trash2 data-icon="inline-start" />{{ t('resource.deleteUser') }}</Button></div>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="loading"><CardHeader><Skeleton class="h-6 w-40" /><Skeleton class="h-4 w-64" /></CardHeader><CardContent class="flex flex-col gap-3"><Skeleton v-for="item in 4" :key="item" class="h-10" /></CardContent></Card>
    <Empty v-else-if="error"><EmptyHeader><EmptyTitle>{{ t('states.errorTitle') }}</EmptyTitle><EmptyDescription>{{ error }}</EmptyDescription></EmptyHeader></Empty>
    <Card v-else-if="Object.keys(data).length"><CardHeader><CardTitle>{{ String(data.display_name || data.name || data.email || route.params.id) }}</CardTitle><CardDescription>{{ t('resource.detailDescription') }}</CardDescription></CardHeader><CardContent><dl class="grid gap-4 sm:grid-cols-2"> <div v-for="(value, key) in data" :key="key" class="rounded-md border p-3"><dt class="text-xs text-muted-foreground">{{ fieldLabel(String(key)) }}</dt><dd class="mt-1 break-words text-sm font-medium">{{ value === null || value === undefined ? '—' : String(value) }}</dd></div></dl></CardContent></Card>
    <Empty v-else><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('resource.noData') }}</EmptyDescription></EmptyHeader></Empty>
    <AlertDialog v-model:open="deleteDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('resource.deleteUserTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resource.deleteUserDescription') }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('resource.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="deleting" @click="deleteUser">{{ t('resource.deleteUser') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
  </div>
</template>
