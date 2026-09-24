<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Eye, Save } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { ApiError, errorMessageKey } from '@/lib/api'
import { generatedApi } from '@/generated/api'
import { createResourceForm, serializeResourceForm, type ResourceFormField } from '@/lib/resource-form'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { adminRoute } from '@/core/routing/url-namespaces'
import { roleDetailPath } from '../role-routes'

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
function openRoleDetail() { router.push(roleDetailPath(adminRoute('roles'), String(route.params.id))) }

onMounted(async () => {
  if (!auth.token) return
  try {
    const manifests = await generatedApi.resourceRegistry(auth.token)
    const manifest = manifests.find((item) => item.name === 'roles')
    if (!manifest) throw new Error('roles manifest missing')
    fields.value = manifest.fields
    const record = editing.value ? await generatedApi.resourceShow<Record<string, unknown>>('roles', String(route.params.id), auth.token) : {}
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

</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center gap-3">
      <Button variant="ghost" size="icon" :aria-label="t('resource.back')" @click="router.push(adminRoute('roles'))"><ArrowLeft /></Button>
      <div class="min-w-0 flex-1"><h1 class="text-2xl font-semibold tracking-tight">{{ editing ? t('rbac.editRole') : t('rbac.createRole') }}</h1><p class="text-sm text-muted-foreground">{{ t('rbac.roleFormDescription') }}</p></div>
      <Button v-if="editing" type="button" variant="outline" class="shrink-0" @click="openRoleDetail"><Eye data-icon="inline-start" />{{ t('resource.viewDetails') }}</Button>
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
  </div>
</template>
