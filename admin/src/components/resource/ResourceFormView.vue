<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Save } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Switch } from '@/components/ui/switch'
import { ApiError, errorMessageKey } from '@/lib/api'
import { createResourceForm, serializeResourceForm, type ResourceFormField } from '@/lib/resource-form'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

interface ResourceDefinition {
  label: string
  route: string
  fields: readonly ResourceFormField[]
}

interface ResourceFormApi {
  show(id: string | number, token: string): Promise<Record<string, unknown>>
  create(payload: Record<string, unknown>, token: string): Promise<unknown>
  update(id: string | number, payload: Record<string, unknown>, token: string): Promise<unknown>
}

const props = defineProps<{ resource: ResourceDefinition; api: ResourceFormApi }>()
const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const editing = computed(() => Boolean(route.params.id))
const loading = ref(editing.value)
const saving = ref(false)
const error = ref('')
const form = ref<Record<string, unknown>>({})
const formFields = computed(() => props.resource.fields.filter((field) => field.visible !== false && field.writable !== false))

function fieldId(field: ResourceFormField) { return `resource-field-${field.name}` }
function inputType(field: ResourceFormField) { return field.type === 'email' || field.type === 'password' || field.type === 'number' || field.type === 'date' ? field.type : 'text' }
function isBoolean(field: ResourceFormField) { return field.type === 'boolean' }
function isSelect(field: ResourceFormField) { return field.type === 'select' }
function isNumber(field: ResourceFormField) { return field.type === 'number' }
function fieldRequired(field: ResourceFormField) { return !editing.value && field.name !== 'locale' }
function localizedError(value: unknown) { return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown') }

onMounted(async () => {
  if (!auth.token) return
  try {
    const record = editing.value ? await props.api.show(String(route.params.id), auth.token) : {}
    form.value = createResourceForm(formFields.value, record)
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
    const payload = serializeResourceForm(formFields.value, form.value)
    if (editing.value) await props.api.update(String(route.params.id), payload, auth.token)
    else await props.api.create(payload, auth.token)
    await router.push(props.resource.route)
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
      <Button variant="ghost" size="icon" :aria-label="t('resource.back')" @click="router.back()"><ArrowLeft /></Button>
      <div><h1 class="text-2xl font-semibold tracking-tight">{{ editing ? `Edit ${resource.label}` : `Create ${resource.label}` }}</h1><p class="text-sm text-muted-foreground">{{ resource.label }}</p></div>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="loading"><CardHeader><Skeleton class="h-6 w-40" /></CardHeader><CardContent class="flex flex-col gap-3"><Skeleton v-for="field in formFields" :key="field.name" class="h-10" /></CardContent></Card>
    <Card v-else><CardHeader><CardTitle>{{ editing ? `Edit ${resource.label}` : `Create ${resource.label}` }}</CardTitle><CardDescription>Review the values before saving.</CardDescription></CardHeader><CardContent><form class="grid gap-5 sm:max-w-xl" @submit.prevent="submit">
      <FieldGroup>
        <Field v-for="field in formFields" :key="field.name">
          <FieldLabel :for="fieldId(field)">{{ field.label }}</FieldLabel>
          <Switch v-if="isBoolean(field)" :id="fieldId(field)" v-model="form[field.name] as boolean" />
          <Select v-else-if="isSelect(field)" v-model="form[field.name] as string"><SelectTrigger :id="fieldId(field)"><SelectValue /></SelectTrigger><SelectContent><SelectItem v-for="option in field.options || []" :key="option.value" :value="option.value">{{ option.label }}</SelectItem></SelectContent></Select>
          <Input v-else :id="fieldId(field)" v-model="form[field.name] as string" :type="inputType(field)" :required="fieldRequired(field)" :step="isNumber(field) ? '1' : undefined" />
        </Field>
      </FieldGroup>
      <div class="flex gap-2"><Button type="submit" :disabled="saving"><Save data-icon="inline-start" />{{ saving ? t('resource.saving') : t('resource.save') }}</Button><Button type="button" variant="outline" @click="router.back()">{{ t('resource.cancel') }}</Button></div>
    </form></CardContent></Card>
  </div>
</template>
