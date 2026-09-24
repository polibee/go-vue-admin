<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useDebounceFn } from '@vueuse/core'
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
import { generatedApi, type RelationOption } from '@/generated/api'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import { localizedFieldLabel, localizedOptionLabel, localizedResourceLabel } from '@/core/resource/resource-i18n'
import { adminRoute } from '@/core/routing/url-namespaces'

interface ResourceDefinition {
  name: string
  label: string
  admin_route: string
  fields: readonly ResourceFormField[]
  relations?: readonly { name: string; kind: 'belongsTo' | 'hasMany'; field: string; selectable: boolean }[]
  form_groups?: readonly { name: string; label: string; columns?: number; fields: readonly string[] }[]
  dependencies?: readonly { field: string; on: string; value: string }[]
}

interface ResourceFormApi {
  show(id: string | number, token: string): Promise<Record<string, unknown>>
  create(payload: Record<string, unknown>, token: string): Promise<unknown>
  update(id: string | number, payload: Record<string, unknown>, token: string): Promise<unknown>
  relationOptions?(relation: string, query: URLSearchParams, token: string): Promise<{ data: RelationOption[]; meta?: Record<string, unknown> }>
}

const props = defineProps<{ resource: ResourceDefinition; api: ResourceFormApi }>()
const { t, te } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const editing = computed(() => Boolean(route.params.id))
const loading = ref(editing.value)
const saving = ref(false)
const error = ref('')
const fieldErrors = ref<Record<string, string>>({})
const form = ref<Record<string, unknown>>({})
const relationOptions = ref<Record<string, RelationOption[]>>({})
const relationMeta = ref<Record<string, { page: number; per_page: number; total: number; last_page: number }>>({})
const relationSearch = ref<Record<string, string>>({})
const relationLoading = ref<Record<string, boolean>>({})
const formFields = computed(() => props.resource.fields.filter((field) => field.visible !== false && field.writable !== false))
const selectableRelations = computed(() => (props.resource.relations || []).filter((relation) => relation.kind === 'belongsTo' && relation.selectable))
const formGroups = computed(() => {
  if (props.resource.form_groups?.length) return props.resource.form_groups.map((group) => ({ ...group, fields: formFields.value.filter((field) => group.fields.includes(field.name)) }))
  return [{ name: 'default', label: resourceLabel.value, columns: 1, fields: formFields.value }]
})
const resourceKey = computed(() => props.resource.admin_route.split('/').filter(Boolean).pop() || '')
const resourceLabel = computed(() => localizedResourceLabel(t, te, resourceKey.value, props.resource.label))
const formTitle = computed(() => t(editing.value ? 'resource.editResource' : 'resource.createResource', { resource: resourceLabel.value }))
function fieldLabel(field: ResourceFormField) { return localizedFieldLabel(t, te, resourceKey.value, field.name, field.label) }
function optionLabel(field: ResourceFormField, value: string, fallback: string) { return localizedOptionLabel(t, te, field.name, value, fallback) }

function fieldId(field: ResourceFormField) { return `resource-field-${field.name}` }
function inputType(field: ResourceFormField) { return field.type === 'email' || field.type === 'password' || field.type === 'number' || field.type === 'date' ? field.type : 'text' }
function isBoolean(field: ResourceFormField) { return field.type === 'boolean' }
function isSelect(field: ResourceFormField) { return field.type === 'select' }
function isNumber(field: ResourceFormField) { return field.type === 'number' }
// Required is enforced by the API on update as well, so the field hint must not
// depend on create mode.
function fieldRequired(field: ResourceFormField) { return field.required === true }
function relationForField(field: ResourceFormField) { return selectableRelations.value.find((relation) => relation.field === field.name) }
function fieldVisible(field: ResourceFormField) {
  const dependency = props.resource.dependencies?.find((item) => item.field === field.name)
  return !dependency || String(form.value[dependency.on] ?? '') === dependency.value
}
function groupClass(columns = 1) { return columns > 1 ? 'sm:grid-cols-2' : 'grid-cols-1' }
async function loadRelationOptions(relationFilter?: { name: string; field: string }, append = false) {
  if (!auth.token) return
  const relations = relationFilter ? [relationFilter] : selectableRelations.value
  for (const relation of relations) {
    relationLoading.value = { ...relationLoading.value, [relation.name]: true }
    try {
      const currentPage = append ? (relationMeta.value[relation.name]?.page || 1) + 1 : 1
      const selected = form.value[relation.field]
      const query = new URLSearchParams({ search: relationSearch.value[relation.name] || '', selected: selected ? String(selected) : '', page: String(currentPage), per_page: '20' })
      const loader = props.api.relationOptions || ((name: string, request: URLSearchParams, token: string) => generatedApi.resourceRelationOptions(resourceKey.value, name, request, token))
      const response = await loader(relation.name, query, auth.token)
      const previous = append ? relationOptions.value[relation.name] || [] : []
      relationOptions.value = { ...relationOptions.value, [relation.name]: [...previous, ...response.data.filter((item) => !previous.some((existing) => existing.value === item.value))] }
      if (response.meta) {
        relationMeta.value = { ...relationMeta.value, [relation.name]: {
          page: Number(response.meta.page || currentPage), per_page: Number(response.meta.per_page || 20),
          total: Number(response.meta.total || response.data.length), last_page: Number(response.meta.last_page || currentPage),
        } }
      }
    } catch (value) {
      error.value = localizedError(value)
    } finally {
      relationLoading.value = { ...relationLoading.value, [relation.name]: false }
    }
  }
}
function localizedError(value: unknown) { return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown') }
const debouncedRelationSearch = useDebounceFn((relation?: { name: string; field: string }) => { if (relation) void loadRelationOptions(relation, false) }, 300)
function captureFieldError(value: unknown) {
  if (!(value instanceof ApiError)) return
  // Prefer the structured field key; fall back to the legacy message format.
  const field = value.field ?? value.message.match(/field "([a-zA-Z0-9_]+)"/)?.[1]
  if (!field) return
  fieldErrors.value = { [field]: localizedError(value) }
  requestAnimationFrame(() => document.getElementById(fieldId({ name: field } as ResourceFormField))?.focus())
}

onMounted(async () => {
  if (!auth.token) return
  try {
    const record = editing.value ? await props.api.show(String(route.params.id), auth.token) : {}
    form.value = createResourceForm(formFields.value, record)
    await loadRelationOptions()
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
  fieldErrors.value = {}
  try {
    const payload = serializeResourceForm(formFields.value, form.value, fieldVisible)
    if (editing.value) await props.api.update(String(route.params.id), payload, auth.token)
    else await props.api.create(payload, auth.token)
    await router.push(adminRoute(props.resource.name))
  } catch (value) {
    captureFieldError(value)
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
      <div><h1 class="text-2xl font-semibold tracking-tight">{{ formTitle }}</h1><p class="text-sm text-muted-foreground">{{ resourceLabel }}</p></div>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="loading"><CardHeader><Skeleton class="h-6 w-40" /></CardHeader><CardContent class="flex flex-col gap-3"><Skeleton v-for="field in formFields" :key="field.name" class="h-10" /></CardContent></Card>
    <Card v-else><CardHeader><CardTitle>{{ formTitle }}</CardTitle><CardDescription>{{ t('resource.formDescription') }}</CardDescription></CardHeader><CardContent><form class="grid gap-5 sm:max-w-xl" @submit.prevent="submit">
      <FieldGroup>
        <div v-for="group in formGroups" :key="group.name" class="rounded-lg border p-4">
          <h3 class="mb-4 text-sm font-medium">{{ group.label }}</h3>
          <FieldGroup :class="['grid gap-5', groupClass(group.columns)]">
        <Field v-for="field in group.fields" v-show="fieldVisible(field)" :key="field.name">
          <FieldLabel :for="fieldId(field)">{{ fieldLabel(field) }}</FieldLabel>
          <Switch v-if="isBoolean(field)" :id="fieldId(field)" v-model="form[field.name] as boolean" />
          <div v-else-if="relationForField(field)" class="grid gap-2"><Input v-model="relationSearch[relationForField(field)?.name || '']" :placeholder="t('resource.relationSearch')" @input="debouncedRelationSearch(relationForField(field))" @keydown.enter.prevent="loadRelationOptions(relationForField(field), false)" /><Select v-model="form[field.name] as string" :disabled="relationLoading[relationForField(field)?.name || '']"><SelectTrigger :id="fieldId(field)"><SelectValue /></SelectTrigger><SelectContent><SelectItem v-for="option in relationOptions[relationForField(field)?.name || ''] || []" :key="option.value" :value="option.value">{{ option.label }}</SelectItem></SelectContent></Select><Button v-if="(relationMeta[relationForField(field)?.name || '']?.last_page || 1) > (relationMeta[relationForField(field)?.name || '']?.page || 1)" type="button" variant="outline" size="sm" :disabled="relationLoading[relationForField(field)?.name || '']" @click="loadRelationOptions(relationForField(field), true)">{{ t('resource.loadMore') }}</Button></div>
          <Select v-else-if="isSelect(field)" v-model="form[field.name] as string"><SelectTrigger :id="fieldId(field)"><SelectValue /></SelectTrigger><SelectContent><SelectItem v-for="option in field.options || []" :key="option.value" :value="option.value">{{ optionLabel(field, option.value, option.label) }}</SelectItem></SelectContent></Select>
          <Input v-else :id="fieldId(field)" v-model="form[field.name] as string" :type="inputType(field)" :required="fieldRequired(field)" :step="isNumber(field) ? '1' : undefined" :aria-invalid="Boolean(fieldErrors[field.name])" />
          <p v-if="fieldErrors[field.name]" class="text-sm text-destructive">{{ fieldErrors[field.name] }}</p>
        </Field>
          </FieldGroup>
        </div>
      </FieldGroup>
      <div class="flex gap-2"><Button type="submit" :disabled="saving || Boolean(error)"><Save data-icon="inline-start" />{{ saving ? t('resource.saving') : t('resource.save') }}</Button><Button type="button" variant="outline" @click="router.back()">{{ t('resource.cancel') }}</Button></div>
    </form></CardContent></Card>
  </div>
</template>
