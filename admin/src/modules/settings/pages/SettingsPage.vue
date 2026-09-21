<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RefreshCw, Save, Settings2 } from '@lucide/vue'
import { useI18n } from 'vue-i18n'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Textarea } from '@/components/ui/textarea'
import { ApiError, errorMessageKey } from '@/lib/api'
import { useAuthStore } from '@/stores/auth'
import { listSystemSettings, saveSystemSetting, type SettingInput, type SettingValueType, type SystemSetting } from '@/modules/settings/api'

const { t } = useI18n()
const auth = useAuthStore()
const settings = ref<SystemSetting[]>([])
const loading = ref(true)
const savingKey = ref('')
const error = ref('')
const success = ref('')
const draft = reactive<Record<string, SettingInput>>({})
const newSetting = reactive<SettingInput & { key: string }>({ key: '', value: '', value_type: 'string', group: 'general', description: '' })

function localizedError(value: unknown) {
  return value instanceof ApiError ? t(errorMessageKey(value.code)) : t('errors.unknown')
}

function toInput(setting: SystemSetting): SettingInput {
  return { value: setting.value, value_type: setting.value_type, group: setting.group, description: setting.description ?? '' }
}

function syncDrafts() {
  for (const setting of settings.value) draft[setting.key] = toInput(setting)
}

async function load() {
  if (!auth.token) return
  loading.value = true
  error.value = ''
  try {
    settings.value = await listSystemSettings(auth.token)
    syncDrafts()
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    loading.value = false
  }
}

async function save(key: string, input: SettingInput) {
  if (!auth.token) return
  savingKey.value = key
  error.value = ''
  success.value = ''
  try {
    const saved = await saveSystemSetting(key, input, auth.token)
    const index = settings.value.findIndex((setting) => setting.key === key)
    if (index >= 0) settings.value[index] = saved
    draft[key] = toInput(saved)
    success.value = t('settings.saved')
  } catch (value) {
    error.value = localizedError(value)
  } finally {
    savingKey.value = ''
  }
}

async function createSetting() {
  const key = newSetting.key.trim()
  if (!key) {
    error.value = t('settings.keyRequired')
    return
  }
  await save(key, newSetting)
  if (!error.value) {
    newSetting.key = ''
    newSetting.value = ''
    newSetting.description = ''
    await load()
  }
}

function updateType(target: SettingInput, value: unknown) {
  target.value_type = String(value) as SettingValueType
}

onMounted(load)
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div><h1 class="text-2xl font-semibold tracking-tight">{{ t('settings.title') }}</h1><p class="text-sm text-muted-foreground">{{ t('settings.description') }}</p></div>
      <Button variant="outline" :disabled="loading" @click="load"><RefreshCw data-icon="inline-start" />{{ t('settings.refresh') }}</Button>
    </div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Alert v-if="success"><AlertTitle>{{ t('settings.saved') }}</AlertTitle><AlertDescription>{{ success }}</AlertDescription></Alert>

    <Card>
      <CardHeader><CardTitle>{{ t('settings.addTitle') }}</CardTitle><CardDescription>{{ t('settings.addDescription') }}</CardDescription></CardHeader>
      <CardContent class="grid gap-4 md:grid-cols-2">
        <div class="flex flex-col gap-2"><Label for="setting-key">{{ t('settings.key') }}</Label><Input id="setting-key" v-model="newSetting.key" placeholder="site.name" /></div>
        <div class="flex flex-col gap-2"><Label for="setting-group">{{ t('settings.group') }}</Label><Input id="setting-group" v-model="newSetting.group" /></div>
        <div class="flex flex-col gap-2"><Label for="setting-type">{{ t('settings.type') }}</Label><Select :model-value="newSetting.value_type" @update:model-value="(value) => updateType(newSetting, value)"><SelectTrigger id="setting-type"><SelectValue /></SelectTrigger><SelectContent><SelectItem value="string">string</SelectItem><SelectItem value="boolean">boolean</SelectItem><SelectItem value="integer">integer</SelectItem><SelectItem value="json">json</SelectItem></SelectContent></Select></div>
        <div class="flex flex-col gap-2"><Label for="setting-value">{{ t('settings.value') }}</Label><Input id="setting-value" v-model="newSetting.value" /></div>
        <div class="flex flex-col gap-2 md:col-span-2"><Label for="setting-description">{{ t('settings.descriptionLabel') }}</Label><Textarea id="setting-description" v-model="newSetting.description" /></div>
        <div class="md:col-span-2"><Button :disabled="savingKey === newSetting.key" @click="createSetting"><Save data-icon="inline-start" />{{ t('settings.add') }}</Button></div>
      </CardContent>
    </Card>

    <div v-if="loading" class="grid gap-4 lg:grid-cols-2"><Skeleton v-for="item in 4" :key="item" class="h-56" /></div>
    <Empty v-else-if="!settings.length"><EmptyHeader><Settings2 class="size-8 text-muted-foreground" /><EmptyTitle>{{ t('settings.emptyTitle') }}</EmptyTitle><EmptyDescription>{{ t('settings.emptyDescription') }}</EmptyDescription></EmptyHeader></Empty>
    <div v-else class="grid gap-4 lg:grid-cols-2">
      <Card v-for="setting in settings" :key="setting.key">
        <CardHeader class="flex-row items-start justify-between gap-4"><div><CardTitle class="font-mono text-base">{{ setting.key }}</CardTitle><CardDescription>{{ setting.description || t('settings.noDescription') }}</CardDescription></div><Badge variant="secondary">{{ setting.value_type }}</Badge></CardHeader>
        <CardContent class="flex flex-col gap-4">
          <div class="grid gap-4 sm:grid-cols-2"><div class="flex flex-col gap-2"><Label>{{ t('settings.group') }}</Label><Input v-model="draft[setting.key].group" /></div><div class="flex flex-col gap-2"><Label>{{ t('settings.value') }}</Label><Input v-model="draft[setting.key].value" /></div></div>
          <div class="flex flex-col gap-2"><Label>{{ t('settings.descriptionLabel') }}</Label><Textarea v-model="draft[setting.key].description" /></div>
          <div class="flex justify-end"><Button :disabled="savingKey === setting.key" @click="save(setting.key, draft[setting.key])"><Save data-icon="inline-start" />{{ t('settings.save') }}</Button></div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
