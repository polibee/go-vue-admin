<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Save } from '@lucide/vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { ApiError, apiFetch, errorMessageKey } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const editing = computed(() => Boolean(route.params.id))
const loading = ref(editing.value)
const saving = ref(false)
const error = ref('')
const form = reactive({ name: '', email: '', password: '', locale: 'zh-CN', is_active: true })

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

onMounted(async () => {
  if (!editing.value || !auth.token) return
  try {
    const user = await apiFetch<Record<string, unknown>>(`/api/v1/admin/resources/users/${route.params.id}`, {}, auth.token)
    form.name = String(user.name || '')
    form.email = String(user.email || '')
    form.locale = String(user.locale || 'zh-CN')
    form.is_active = Boolean(user.is_active)
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
    await apiFetch(path, { method, body: JSON.stringify(form) }, auth.token)
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
        <div class="grid gap-2"><Label for="user-name">{{ t('resource.userName') }}</Label><Input id="user-name" v-model="form.name" required autocomplete="name" /></div>
        <div class="grid gap-2"><Label for="user-email">{{ t('resource.userEmail') }}</Label><Input id="user-email" v-model="form.email" required type="email" autocomplete="email" /></div>
        <div class="grid gap-2"><Label for="user-password">{{ t('resource.userPassword') }}</Label><Input id="user-password" v-model="form.password" :required="!editing" type="password" autocomplete="new-password" /><p class="text-xs text-muted-foreground">{{ editing ? t('resource.passwordHint') : t('resource.passwordRequired') }}</p></div>
        <div class="grid gap-2"><Label for="user-locale">{{ t('resource.userLocale') }}</Label><Input id="user-locale" v-model="form.locale" /></div>
        <label class="flex items-center gap-3 text-sm"><Switch v-model="form.is_active" /><span>{{ t('resource.userActive') }}</span></label>
        <div class="flex gap-2"><Button type="submit" :disabled="saving"><Save data-icon="inline-start" />{{ saving ? t('resource.saving') : t('resource.save') }}</Button><Button type="button" variant="outline" @click="router.push('/users')">{{ t('resource.cancel') }}</Button></div>
      </form></CardContent>
    </Card>
  </div>
</template>
