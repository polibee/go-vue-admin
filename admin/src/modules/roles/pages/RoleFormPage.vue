<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Save } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ApiError, apiFetch, errorMessageKey } from '@/lib/api'
import { generatedApi } from '@/generated/api'
import type { DataScope, RolePermissionAssignment } from '@/generated/api'
import { createResourceForm, serializeResourceForm, type ResourceFormField } from '@/lib/resource-form'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { adminRoute } from '@/core/routing/url-namespaces'
import { isProtectedRole, selectedPermissionIds, type RolePermissionOption } from '../role-permissions'

interface ResourceManifest { name: string; fields: ResourceFormField[] }
interface RoleRecord { permissions?: RolePermissionAssignment[] }

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const editing = computed(() => Boolean(route.params.id))
const protectedRole = computed(() => isProtectedRole(String(form.value.name || '')))
const loading = ref(editing.value)
const saving = ref(false)
const error = ref('')
const fields = ref<ResourceFormField[]>([])
const form = ref<Record<string, any>>({})
const permissionOptions = ref<RolePermissionOption[]>([])
const rolePermissionAssignments = ref<RolePermissionAssignment[]>([])
const selectedPermissionIDs = ref<number[]>([])
const selectedPermissionScopes = ref<Record<string, DataScope>>({})
const permissionLoading = ref(false)
const permissionSaving = ref(false)
const permissionError = ref('')

function localizedError(value: unknown) { return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown') }
function fieldId(field: ResourceFormField) { return `role-field-${field.name}` }

onMounted(async () => {
  if (!auth.token) return
  try {
    const manifests = await generatedApi.resourceRegistry(auth.token)
    const manifest = manifests.find((item) => item.name === 'roles')
    if (!manifest) throw new Error('roles manifest missing')
    fields.value = manifest.fields
    const record = editing.value ? await generatedApi.resourceShow<RoleRecord & Record<string, unknown>>('roles', String(route.params.id), auth.token) : {}
    form.value = createResourceForm(fields.value, record)
    if (editing.value) {
      permissionLoading.value = true
      const [permissions, role] = await Promise.all([
        apiFetch<RolePermissionOption[]>('/api/v1/admin/permissions', {}, auth.token),
        Promise.resolve(record),
      ])
      permissionOptions.value = permissions
      rolePermissionAssignments.value = role.permissions || []
      selectedPermissionIDs.value = selectedPermissionIds(rolePermissionAssignments.value)
      selectedPermissionScopes.value = Object.fromEntries(rolePermissionAssignments.value.map((assignment) => [String(assignment.id), assignment.scope || 'all']))
      permissionLoading.value = false
    }
  } catch (value) {
    error.value = localizedError(value)
    permissionError.value = localizedError(value)
    permissionLoading.value = false
  } finally {
    loading.value = false
  }
})

async function submit() {
  if (!auth.token) return
  saving.value = true
  error.value = ''
  try {
    const payload = serializeResourceForm(fields.value, form.value)
    if (editing.value) await generatedApi.resourceUpdate('roles', String(route.params.id), payload, auth.token)
    else await generatedApi.resourceCreate('roles', payload, auth.token)
    await router.push(adminRoute('roles'))
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    saving.value = false
  }
}

function togglePermission(permissionID: number, checked: boolean) {
  selectedPermissionIDs.value = checked
    ? [...selectedPermissionIDs.value, permissionID]
    : selectedPermissionIDs.value.filter((id) => id !== permissionID)
}

function updatePermissionScope(permissionID: number, value: string) {
  selectedPermissionScopes.value = { ...selectedPermissionScopes.value, [String(permissionID)]: value === 'own' ? 'own' : 'all' }
}

async function savePermissions() {
  if (!auth.token || !editing.value || !route.params.id) return
  permissionSaving.value = true
  permissionError.value = ''
  try {
    await generatedApi.replaceRolePermissions(Number(route.params.id), selectedPermissionIDs.value, selectedPermissionScopes.value, {}, auth.token)
    rolePermissionAssignments.value = permissionOptions.value
      .filter((permission) => selectedPermissionIDs.value.includes(permission.id))
      .map((permission) => ({ ...permission, scope: selectedPermissionScopes.value[String(permission.id)] || 'all' }))
  } catch (value) {
    permissionError.value = localizedError(value)
  } finally {
    permissionSaving.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center gap-3">
      <Button variant="ghost" size="icon" :aria-label="t('resource.back')" @click="router.push(adminRoute('roles'))"><ArrowLeft /></Button>
      <div><h1 class="text-2xl font-semibold tracking-tight">{{ editing ? t('rbac.editRole') : t('rbac.createRole') }}</h1><p class="text-sm text-muted-foreground">{{ t('rbac.roleFormDescription') }}</p></div>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="loading"><CardContent class="py-8">{{ t('resource.loading') }}</CardContent></Card>
    <Card v-else>
      <CardHeader><CardTitle>{{ editing ? t('rbac.editRole') : t('rbac.createRole') }}</CardTitle><CardDescription>{{ t('rbac.roleFormDescription') }}</CardDescription></CardHeader>
      <CardContent><form class="grid gap-5 sm:max-w-xl" @submit.prevent="submit">
        <FieldGroup>
          <Field v-for="field in fields" :key="field.name">
            <FieldLabel :for="fieldId(field)">{{ field.label }}</FieldLabel>
            <Input :id="fieldId(field)" v-model="form[field.name]" required />
          </Field>
        </FieldGroup>
        <div class="flex gap-2"><Button type="submit" :disabled="saving"><Save data-icon="inline-start" />{{ saving ? t('resource.saving') : t('resource.save') }}</Button><Button type="button" variant="outline" @click="router.push(adminRoute('roles'))">{{ t('resource.cancel') }}</Button></div>
      </form></CardContent>
    </Card>
    <Card v-if="editing">
      <CardHeader><CardTitle>{{ t('rbac.assignPermissions') }}</CardTitle><CardDescription>{{ t('rbac.rolePermissionsDescription') }}</CardDescription></CardHeader>
      <CardContent>
        <Alert v-if="permissionError" class="mb-4" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ permissionError }}</AlertDescription></Alert>
        <div v-if="permissionLoading" class="py-4">{{ t('resource.loading') }}</div>
        <div v-else-if="permissionOptions.length" class="grid gap-3 md:grid-cols-2">
          <label v-for="permission in permissionOptions" :key="permission.id" class="flex items-center gap-3 rounded-md border p-3">
            <Checkbox :disabled="protectedRole" :model-value="selectedPermissionIDs.includes(permission.id)" @update:model-value="togglePermission(permission.id, Boolean($event))" />
            <span class="min-w-0 flex-1"><span class="block text-sm font-medium">{{ permission.display_name }}</span><span class="block text-xs text-muted-foreground">{{ permission.name }}</span></span>
            <Select v-if="selectedPermissionIDs.includes(permission.id)" :disabled="protectedRole" :model-value="selectedPermissionScopes[String(permission.id)] || 'all'" @update:model-value="updatePermissionScope(permission.id, String($event))"><SelectTrigger class="w-28"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="all">{{ t('rbac.scopeAll') }}</SelectItem><SelectItem value="own">{{ t('rbac.scopeOwn') }}</SelectItem></SelectContent></Select>
          </label>
        </div>
        <p v-else class="text-sm text-muted-foreground">{{ t('rbac.noAssignablePermissions') }}</p>
        <div class="mt-5 flex justify-end"><Button :disabled="protectedRole || permissionSaving || permissionLoading || permissionOptions.length === 0" @click="savePermissions">{{ permissionSaving ? t('resource.saving') : t('rbac.savePermissions') }}</Button></div>
      </CardContent>
    </Card>
  </div>
</template>
