<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { ApiError, apiFetch, errorMessageKey } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import { userStatusLabelKey, type UserStatus } from '@/lib/user-status'

interface RBACUser { id: number; name: string; email: string; status: UserStatus }
interface RBACRole { id: number; name: string; display_name: string }
interface RBACPermission { id: number; name: string; display_name: string }
interface RBACRoleDetail extends RBACRole { permissions: RBACPermission[] }

const { t } = useI18n()
const auth = useAuthStore()
const users = ref<RBACUser[]>([])
const roles = ref<RBACRole[]>([])
const permissions = ref<RBACPermission[]>([])
const loading = ref(true)
const error = ref('')
const actionError = ref('')
const roleDialogOpen = ref(false)
const permissionDialogOpen = ref(false)
const userRoleDialogOpen = ref(false)
const deleteRoleDialogOpen = ref(false)
const editingRoleId = ref<number | null>(null)
const selectedPermissionIDs = ref<number[]>([])
const selectedRoleIDs = ref<number[]>([])
const selectedRole = ref<RBACRole>()
const roleToDelete = ref<RBACRole>()
const selectedUser = ref<RBACUser>()
const roleForm = reactive({ name: '', display_name: '' })
const saving = ref(false)

function localizedError(errorValue: unknown) {
  return errorValue instanceof ApiError ? t(errorMessageKey(errorValue.code)) : t('errors.unknown')
}

async function loadRBAC() {
  if (!auth.token) return
  loading.value = true
  error.value = ''
  try {
    const [userData, roleData, permissionData] = await Promise.all([
      auth.can('admin.users.view') ? apiFetch<RBACUser[]>('/api/v1/admin/users', {}, auth.token) : Promise.resolve([]),
      auth.canAny(['admin.roles.manage', 'admin.permissions.manage']) ? apiFetch<RBACRole[]>('/api/v1/admin/roles', {}, auth.token) : Promise.resolve([]),
      auth.canAny(['admin.roles.manage', 'admin.permissions.manage']) ? apiFetch<RBACPermission[]>('/api/v1/admin/permissions', {}, auth.token) : Promise.resolve([]),
    ])
    users.value = userData
    roles.value = roleData
    permissions.value = permissionData
  } catch (loadError) {
    error.value = localizedError(loadError)
  } finally {
    loading.value = false
  }
}

function openCreateRole() {
  editingRoleId.value = null
  roleForm.name = ''
  roleForm.display_name = ''
  actionError.value = ''
  roleDialogOpen.value = true
}

function openEditRole(role: RBACRole) {
  editingRoleId.value = role.id
  roleForm.name = role.name
  roleForm.display_name = role.display_name
  actionError.value = ''
  roleDialogOpen.value = true
}

async function saveRole() {
  if (!auth.token) return
  saving.value = true
  actionError.value = ''
  try {
    const path = editingRoleId.value ? `/api/v1/admin/roles/${editingRoleId.value}` : '/api/v1/admin/roles'
    await apiFetch(path, { method: editingRoleId.value ? 'PUT' : 'POST', body: JSON.stringify(roleForm) }, auth.token)
    roleDialogOpen.value = false
    await loadRBAC()
  } catch (saveError) {
    actionError.value = localizedError(saveError)
  } finally {
    saving.value = false
  }
}

function openDeleteRole(role: RBACRole) {
  roleToDelete.value = role
  actionError.value = ''
  deleteRoleDialogOpen.value = true
}

async function deleteRole() {
  if (!auth.token || !roleToDelete.value) return
  saving.value = true
  actionError.value = ''
  try {
    await apiFetch(`/api/v1/admin/roles/${roleToDelete.value.id}`, { method: 'DELETE' }, auth.token)
    deleteRoleDialogOpen.value = false
    await loadRBAC()
  } catch (deleteError) {
    actionError.value = localizedError(deleteError)
  } finally {
    saving.value = false
  }
}

async function openRolePermissions(role: RBACRole) {
  if (!auth.token) return
  selectedRole.value = role
  actionError.value = ''
  try {
    const current = await apiFetch<RBACRoleDetail>(`/api/v1/admin/roles/${role.id}`, {}, auth.token)
    selectedPermissionIDs.value = current.permissions.map((permission) => permission.id)
  } catch (loadError) {
    actionError.value = localizedError(loadError)
    selectedPermissionIDs.value = []
  }
  permissionDialogOpen.value = true
}

async function saveRolePermissions() {
  if (!auth.token || !selectedRole.value) return
  saving.value = true
  actionError.value = ''
  try {
    await apiFetch(`/api/v1/admin/roles/${selectedRole.value.id}/permissions`, { method: 'PUT', body: JSON.stringify({ permission_ids: selectedPermissionIDs.value }) }, auth.token)
    permissionDialogOpen.value = false
  } catch (saveError) {
    actionError.value = localizedError(saveError)
  } finally {
    saving.value = false
  }
}

async function openUserRoles(user: RBACUser) {
  if (!auth.token) return
  selectedUser.value = user
  actionError.value = ''
  try {
    const current = await apiFetch<RBACRole[]>(`/api/v1/admin/users/${user.id}/roles`, {}, auth.token)
    selectedRoleIDs.value = current.map((role) => role.id)
  } catch (loadError) {
    actionError.value = localizedError(loadError)
    selectedRoleIDs.value = []
  }
  userRoleDialogOpen.value = true
}

async function saveUserRoles() {
  if (!auth.token || !selectedUser.value) return
  saving.value = true
  actionError.value = ''
  try {
    await apiFetch(`/api/v1/admin/users/${selectedUser.value.id}/roles`, { method: 'PUT', body: JSON.stringify({ role_ids: selectedRoleIDs.value }) }, auth.token)
    userRoleDialogOpen.value = false
  } catch (saveError) {
    actionError.value = localizedError(saveError)
  } finally {
    saving.value = false
  }
}

function toggleID(target: number[], id: number, checked: boolean) {
  if (checked && !target.includes(id)) target.push(id)
  if (!checked) {
    const index = target.indexOf(id)
    if (index >= 0) target.splice(index, 1)
  }
}

onMounted(loadRBAC)
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div><h1 class="text-2xl font-semibold tracking-tight">{{ t('rbac.title') }}</h1><p class="text-sm text-muted-foreground">{{ t('rbac.description') }}</p></div>
      <Button v-if="auth.can('admin.roles.manage')" @click="openCreateRole">{{ t('rbac.createRole') }}</Button>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Alert v-if="actionError" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ actionError }}</AlertDescription></Alert>
    <div v-if="loading" class="grid gap-4 lg:grid-cols-3"><Skeleton v-for="item in 3" :key="item" class="h-64" /></div>
    <div v-else class="grid gap-4 lg:grid-cols-3">
      <Card v-if="auth.can('admin.users.view')">
        <CardHeader><CardTitle>{{ t('rbac.users') }}</CardTitle><CardDescription>{{ users.length }}</CardDescription></CardHeader>
        <CardContent>
          <Empty v-if="!users.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('rbac.noUsers') }}</EmptyDescription></EmptyHeader></Empty>
          <Table v-else><TableHeader><TableRow><TableHead>{{ t('rbac.name') }}</TableHead><TableHead>{{ t('rbac.status') }}</TableHead><TableHead /></TableRow></TableHeader><TableBody><TableRow v-for="user in users" :key="user.id"><TableCell><div class="font-medium">{{ user.name }}</div><div class="text-xs text-muted-foreground">{{ user.email }}</div></TableCell><TableCell><Badge variant="secondary">{{ t(userStatusLabelKey(user.status)) }}</Badge></TableCell><TableCell><Button v-if="auth.can('admin.roles.manage')" variant="ghost" size="sm" @click="openUserRoles(user)">{{ t('rbac.assignRoles') }}</Button></TableCell></TableRow></TableBody></Table>
        </CardContent>
      </Card>
      <Card v-if="auth.can('admin.roles.manage')" class="lg:col-span-2">
        <CardHeader><CardTitle>{{ t('rbac.roles') }}</CardTitle><CardDescription>{{ roles.length }}</CardDescription></CardHeader>
        <CardContent>
          <Empty v-if="!roles.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('rbac.noRoles') }}</EmptyDescription></EmptyHeader></Empty>
          <Table v-else><TableHeader><TableRow><TableHead>{{ t('rbac.name') }}</TableHead><TableHead>{{ t('rbac.displayName') }}</TableHead><TableHead /></TableRow></TableHeader><TableBody><TableRow v-for="role in roles" :key="role.id"><TableCell><div class="font-medium">{{ role.name }}</div></TableCell><TableCell>{{ role.display_name }}</TableCell><TableCell class="flex justify-end gap-1"><Button variant="ghost" size="sm" @click="openRolePermissions(role)">{{ t('rbac.assignPermissions') }}</Button><Button variant="ghost" size="sm" @click="openEditRole(role)">{{ t('rbac.editRole') }}</Button><Button variant="ghost" size="sm" :disabled="role.name === 'super-admin'" @click="openDeleteRole(role)">{{ t('rbac.deleteRole') }}</Button></TableCell></TableRow></TableBody></Table>
        </CardContent>
      </Card>
      <Card v-if="auth.canAny(['admin.roles.manage', 'admin.permissions.manage'])" class="lg:col-span-3">
        <CardHeader><CardTitle>{{ t('rbac.permissions') }}</CardTitle><CardDescription>{{ permissions.length }}</CardDescription></CardHeader>
        <CardContent><Empty v-if="!permissions.length"><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('rbac.noPermissions') }}</EmptyDescription></EmptyHeader></Empty><ul v-else class="grid gap-3 md:grid-cols-2 lg:grid-cols-3"><li v-for="permission in permissions" :key="permission.id" class="rounded-md border p-3"><div class="font-medium">{{ permission.display_name }}</div><div class="text-xs text-muted-foreground">{{ permission.name }}</div></li></ul></CardContent>
      </Card>
    </div>

    <Dialog v-model:open="roleDialogOpen"><DialogContent><DialogHeader><DialogTitle>{{ editingRoleId ? t('rbac.editRole') : t('rbac.createRole') }}</DialogTitle><DialogDescription>{{ t('rbac.description') }}</DialogDescription></DialogHeader><div class="flex flex-col gap-4"><div class="flex flex-col gap-2"><Label for="role-name">{{ t('rbac.roleName') }}</Label><Input id="role-name" v-model="roleForm.name" /></div><div class="flex flex-col gap-2"><Label for="role-display-name">{{ t('rbac.displayName') }}</Label><Input id="role-display-name" v-model="roleForm.display_name" /></div></div><DialogFooter><Button variant="outline" @click="roleDialogOpen = false">{{ t('rbac.cancel') }}</Button><Button :disabled="saving" @click="saveRole">{{ t('rbac.save') }}</Button></DialogFooter></DialogContent></Dialog>

    <Dialog v-model:open="permissionDialogOpen"><DialogContent><DialogHeader><DialogTitle>{{ t('rbac.assignPermissions') }}</DialogTitle><DialogDescription>{{ selectedRole?.display_name }}</DialogDescription></DialogHeader><div v-if="permissions.length" class="flex max-h-80 flex-col gap-3 overflow-y-auto"><label v-for="permission in permissions" :key="permission.id" class="flex items-start gap-3 rounded-md border p-3"><Checkbox :model-value="selectedPermissionIDs.includes(permission.id)" @update:model-value="toggleID(selectedPermissionIDs, permission.id, Boolean($event))" /><span><span class="block text-sm font-medium">{{ permission.display_name }}</span><span class="block text-xs text-muted-foreground">{{ permission.name }}</span></span></label></div><Empty v-else><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('rbac.noAssignablePermissions') }}</EmptyDescription></EmptyHeader></Empty><DialogFooter><Button variant="outline" @click="permissionDialogOpen = false">{{ t('rbac.cancel') }}</Button><Button :disabled="saving || selectedRole?.name === 'super-admin'" @click="saveRolePermissions">{{ t('rbac.save') }}</Button></DialogFooter></DialogContent></Dialog>

    <Dialog v-model:open="userRoleDialogOpen"><DialogContent><DialogHeader><DialogTitle>{{ t('rbac.assignRoles') }}</DialogTitle><DialogDescription>{{ selectedUser?.email }}</DialogDescription></DialogHeader><div v-if="roles.length" class="flex max-h-80 flex-col gap-3 overflow-y-auto"><label v-for="role in roles" :key="role.id" class="flex items-center gap-3 rounded-md border p-3"><Checkbox :model-value="selectedRoleIDs.includes(role.id)" @update:model-value="toggleID(selectedRoleIDs, role.id, Boolean($event))" /><span><span class="block text-sm font-medium">{{ role.display_name }}</span><span class="block text-xs text-muted-foreground">{{ role.name }}</span></span></label></div><Empty v-else><EmptyHeader><EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('rbac.noAssignableRoles') }}</EmptyDescription></EmptyHeader></Empty><DialogFooter><Button variant="outline" @click="userRoleDialogOpen = false">{{ t('rbac.cancel') }}</Button><Button :disabled="saving" @click="saveUserRoles">{{ t('rbac.save') }}</Button></DialogFooter></DialogContent></Dialog>
    <AlertDialog v-model:open="deleteRoleDialogOpen"><AlertDialogContent><AlertDialogHeader><AlertDialogTitle>{{ t('rbac.deleteRoleTitle') }}</AlertDialogTitle><AlertDialogDescription>{{ t('rbac.deleteRoleDescription') }} {{ roleToDelete?.display_name }}</AlertDialogDescription></AlertDialogHeader><AlertDialogFooter><AlertDialogCancel>{{ t('rbac.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="saving" @click="deleteRole">{{ t('rbac.delete') }}</AlertDialogAction></AlertDialogFooter></AlertDialogContent></AlertDialog>
  </div>
</template>
