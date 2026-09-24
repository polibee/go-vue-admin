<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ApiError, apiFetch, errorMessageKey } from '@/lib/api'
import { generatedApi, type DataScope, type RolePermissionAssignment } from '@/generated/api'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { isProtectedRole, type RolePermissionOption } from '../role-permissions'

const props = defineProps<{
  roleId: number
  roleName: string
  assignments: RolePermissionAssignment[]
}>()

const { t } = useI18n()
const auth = useAuthStore()
const options = ref<RolePermissionOption[]>([])
const selectedIDs = ref<number[]>([])
const scopes = ref<Record<string, DataScope>>({})
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const saved = ref(false)
const protectedRole = computed(() => isProtectedRole(props.roleName))
const groupedOptions = computed(() => {
  const groups = new Map<string, RolePermissionOption[]>()
  for (const permission of options.value) {
    const group = permission.name.split('.')[1] || 'general'
    groups.set(group, [...(groups.get(group) || []), permission])
  }
  return Array.from(groups, ([name, permissions]) => ({ name, permissions }))
})

function syncAssignments(assignments: RolePermissionAssignment[]) {
  selectedIDs.value = assignments.map((assignment) => assignment.id)
  scopes.value = Object.fromEntries(assignments.map((assignment) => [String(assignment.id), assignment.scope || 'all']))
}

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

watch(() => props.assignments, syncAssignments, { immediate: true })

async function loadOptions() {
  if (!auth.token) return
  loading.value = true
  error.value = ''
  try {
    options.value = await apiFetch<RolePermissionOption[]>('/api/v1/admin/permissions', {}, auth.token)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
}

void loadOptions()

function togglePermission(id: number, checked: boolean) {
  selectedIDs.value = checked ? [...new Set([...selectedIDs.value, id])] : selectedIDs.value.filter((item) => item !== id)
  saved.value = false
}

function updateScope(id: number, value: string) {
  scopes.value = { ...scopes.value, [String(id)]: value === 'own' ? 'own' : 'all' }
  saved.value = false
}

async function save() {
  if (!auth.token || protectedRole.value) return
  saving.value = true
  error.value = ''
  saved.value = false
  try {
    await generatedApi.replaceRolePermissions(props.roleId, selectedIDs.value, scopes.value, {}, auth.token)
    saved.value = true
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <Card>
    <CardHeader class="flex flex-row items-start justify-between gap-4">
      <div>
        <CardTitle>{{ t('rbac.assignPermissions') }}</CardTitle>
        <CardDescription>{{ t('rbac.rolePermissionsDescription') }}</CardDescription>
      </div>
      <Badge v-if="protectedRole" variant="secondary">{{ t('rbac.systemRoleReadOnly') }}</Badge>
    </CardHeader>
    <CardContent>
      <Alert v-if="error" class="mb-4" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
      <div v-if="loading" class="py-4 text-sm text-muted-foreground">{{ t('resource.loading') }}</div>
      <div v-else-if="groupedOptions.length" class="grid gap-5 lg:grid-cols-2">
        <section v-for="group in groupedOptions" :key="group.name" class="rounded-lg border bg-muted/20 p-4">
          <div class="mb-3 flex items-center justify-between gap-3"><h3 class="font-medium capitalize">{{ group.name }}</h3><span class="text-xs text-muted-foreground">{{ group.permissions.length }}</span></div>
          <div class="grid gap-2">
            <div v-for="permission in group.permissions" :key="permission.id" class="flex items-center gap-3 rounded-md bg-background p-3">
              <Checkbox :disabled="protectedRole" :model-value="selectedIDs.includes(permission.id)" @update:model-value="togglePermission(permission.id, Boolean($event))" />
              <div class="min-w-0 flex-1"><div class="text-sm font-medium">{{ permission.display_name }}</div><div class="truncate text-xs text-muted-foreground">{{ permission.name }}</div></div>
              <Select v-if="selectedIDs.includes(permission.id)" :disabled="protectedRole" :model-value="scopes[String(permission.id)] || 'all'" @update:model-value="updateScope(permission.id, String($event))"><SelectTrigger class="w-28"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="all">{{ t('rbac.scopeAll') }}</SelectItem><SelectItem value="own">{{ t('rbac.scopeOwn') }}</SelectItem></SelectContent></Select>
            </div>
          </div>
        </section>
      </div>
      <p v-else class="text-sm text-muted-foreground">{{ t('rbac.noAssignablePermissions') }}</p>
      <div v-if="!protectedRole && !loading && groupedOptions.length" class="mt-5 flex items-center justify-end gap-3"><span v-if="saved" class="text-sm text-muted-foreground">{{ t('rbac.saved') }}</span><Button :disabled="saving" @click="save">{{ saving ? t('resource.saving') : t('rbac.savePermissions') }}</Button></div>
    </CardContent>
  </Card>
</template>
