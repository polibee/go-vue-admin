<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Save } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { ApiError, apiFetch, errorMessageKey } from '@/lib/api'
import { createResourceForm, serializeResourceForm, type ResourceFormField } from '@/lib/resource-form'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

interface ResourceManifest { name: string; fields: ResourceFormField[] }

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const editing = computed(() => Boolean(route.params.id))
const loading = ref(editing.value)
const saving = ref(false)
const error = ref('')
const fields = ref<ResourceFormField[]>([])
const form = ref<Record<string, any>>({})

function localizedError(value: unknown) { return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown') }
function fieldId(field: ResourceFormField) { return `role-field-${field.name}` }

onMounted(async () => {
  if (!auth.token) return
  try {
    const manifests = await apiFetch<ResourceManifest[]>('/api/v1/admin/resources', {}, auth.token)
    const manifest = manifests.find((item) => item.name === 'roles')
    if (!manifest) throw new Error('roles manifest missing')
    fields.value = manifest.fields
    const record = editing.value ? await apiFetch<Record<string, unknown>>(`/api/v1/admin/resources/roles/${route.params.id}`, {}, auth.token) : {}
    form.value = createResourceForm(fields.value, record)
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
})

async function submit() {
  if (!auth.token) return
  saving.value = true
  error.value = ''
  try {
    const path = editing.value ? `/api/v1/admin/roles/${route.params.id}` : '/api/v1/admin/roles'
    const method = editing.value ? 'PUT' : 'POST'
    await apiFetch(path, { method, body: JSON.stringify(serializeResourceForm(fields.value, form.value)) }, auth.token)
    await router.push('/roles')
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center gap-3">
      <Button variant="ghost" size="icon" :aria-label="t('resource.back')" @click="router.push('/roles')"><ArrowLeft /></Button>
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
        <div class="flex gap-2"><Button type="submit" :disabled="saving"><Save data-icon="inline-start" />{{ saving ? t('resource.saving') : t('resource.save') }}</Button><Button type="button" variant="outline" @click="router.push('/roles')">{{ t('resource.cancel') }}</Button></div>
      </form></CardContent>
    </Card>
  </div>
</template>
