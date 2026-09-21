<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Copy, Dices, Eye, EyeOff, Save } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { InputGroup, InputGroupAddon, InputGroupButton, InputGroupInput } from '@/components/ui/input-group'
import { Switch } from '@/components/ui/switch'
import { ApiError, apiFetch, errorMessageKey } from '@/lib/api'
import { createResourceForm, serializeResourceForm, type ResourceFormField } from '@/lib/resource-form'
import { generatePassword } from '@/lib/password-generator'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

interface ResourceManifest { name: string; label: string; fields: ResourceFormField[] }

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
const showPassword = ref(false)
const copiedPassword = ref(false)
const visibleFields = computed(() => fields.value.filter((field) => field.name !== 'locale'))

function localizedError(value: unknown) { return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown') }
function fieldId(field: ResourceFormField) { return `resource-field-${field.name}` }
function isPassword(field: ResourceFormField) { return field.type === 'password' }
function isBoolean(field: ResourceFormField) { return field.type === 'boolean' }
function isNumber(field: ResourceFormField) { return field.type === 'number' }
function inputType(field: ResourceFormField) { return field.type === 'email' || field.type === 'password' || field.type === 'number' || field.type === 'date' ? field.type : 'text' }
function fieldRequired(field: ResourceFormField) { return !editing.value && field.name !== 'is_active' }
function fieldDescription(field: ResourceFormField) { return isPassword(field) ? (editing.value ? t('resource.passwordHint') : t('resource.passwordRequired')) : '' }
function generateUserPassword() {
  form.value.password = generatePassword()
  showPassword.value = true
  copiedPassword.value = false
}
async function copyUserPassword() {
  if (!form.value.password) return
  await navigator.clipboard.writeText(String(form.value.password))
  copiedPassword.value = true
}

onMounted(async () => {
  if (!auth.token) return
  try {
    const manifests = await apiFetch<ResourceManifest[]>('/api/v1/admin/resources', {}, auth.token)
    const manifest = manifests.find((item) => item.name === 'users')
    if (!manifest) throw new Error('users manifest missing')
    fields.value = manifest.fields
    const record = editing.value ? await apiFetch<Record<string, unknown>>(`/api/v1/admin/resources/users/${route.params.id}`, {}, auth.token) : {}
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
    const path = editing.value ? `/api/v1/admin/users/${route.params.id}` : '/api/v1/admin/users'
    const method = editing.value ? 'PUT' : 'POST'
    await apiFetch(path, { method, body: JSON.stringify(serializeResourceForm(fields.value, form.value)) }, auth.token)
    await router.push('/users')
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
      <Button variant="ghost" size="icon" :aria-label="t('resource.back')" @click="router.push('/users')"><ArrowLeft /></Button>
      <div><h1 class="text-2xl font-semibold tracking-tight">{{ editing ? t('resource.editUser') : t('resource.createUser') }}</h1><p class="text-sm text-muted-foreground">{{ t('resource.userFormDescription') }}</p></div>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="loading"><CardContent class="py-8">{{ t('resource.loading') }}</CardContent></Card>
    <Card v-else>
      <CardHeader><CardTitle>{{ editing ? t('resource.editUser') : t('resource.createUser') }}</CardTitle><CardDescription>{{ t('resource.userFormDescription') }}</CardDescription></CardHeader>
      <CardContent><form class="grid gap-5 sm:max-w-xl" @submit.prevent="submit">
        <FieldGroup>
          <Field v-for="field in visibleFields" :key="field.name">
            <FieldLabel :for="fieldId(field)">{{ field.label }}</FieldLabel>
            <Switch v-if="isBoolean(field)" :id="fieldId(field)" v-model="form[field.name]" />
            <InputGroup v-else-if="isPassword(field)">
              <InputGroupInput :id="fieldId(field)" v-model="form[field.name]" :type="showPassword ? 'text' : 'password'" :required="fieldRequired(field)" autocomplete="new-password" />
              <InputGroupAddon align="inline-end">
                <InputGroupButton :aria-label="showPassword ? t('resource.hidePassword') : t('resource.showPassword')" :title="showPassword ? t('resource.hidePassword') : t('resource.showPassword')" @click="showPassword = !showPassword"><EyeOff v-if="showPassword" /><Eye v-else /></InputGroupButton>
                <InputGroupButton :aria-label="t('resource.generatePassword')" :title="t('resource.generatePassword')" @click="generateUserPassword"><Dices /></InputGroupButton>
                <InputGroupButton :aria-label="copiedPassword ? t('resource.copiedPassword') : t('resource.copyPassword')" :title="copiedPassword ? t('resource.copiedPassword') : t('resource.copyPassword')" @click="copyUserPassword"><Copy /></InputGroupButton>
              </InputGroupAddon>
            </InputGroup>
            <Input v-else :id="fieldId(field)" v-model="form[field.name]" :type="inputType(field)" :required="fieldRequired(field)" :autocomplete="isPassword(field) ? 'new-password' : undefined" :step="isNumber(field) ? '1' : undefined" />
            <p v-if="fieldDescription(field)" class="text-xs text-muted-foreground">{{ fieldDescription(field) }}</p>
          </Field>
        </FieldGroup>
        <div class="flex gap-2"><Button type="submit" :disabled="saving"><Save data-icon="inline-start" />{{ saving ? t('resource.saving') : t('resource.save') }}</Button><Button type="button" variant="outline" @click="router.push('/users')">{{ t('resource.cancel') }}</Button></div>
      </form></CardContent>
    </Card>
  </div>
</template>
