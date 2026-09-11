<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { toast } from 'vue-sonner'

import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type SettingResource } from '@/generated/api'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Spinner } from '@/components/ui/spinner'
import { Switch } from '@/components/ui/switch'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Textarea } from '@/components/ui/textarea'
import { SettingsService, type SettingsByNamespace } from '@/modules/settings/settings.service'
import { useI18n } from 'vue-i18n'

type SettingType = 'string' | 'boolean' | 'integer' | 'number' | 'json'

const service = new SettingsService(createGeneratedApiClient(apiClient))
const { t } = useI18n()
const groups = ref<SettingsByNamespace>({})
const activeNamespace = ref('general')
const loading = ref(true)
const error = ref('')
const savingKey = ref('')
const drafts = reactive<Record<string, string>>({})
const booleanDrafts = reactive<Record<string, boolean>>({})

const namespaces = computed(() => Object.keys(groups.value).sort())

const newForm = reactive<{ namespace: string; key: string; value: string; valueType: SettingType; description: string }>({
  namespace: 'general', key: '', value: '', valueType: 'string', description: '',
})

async function loadSettings() {
  loading.value = true
  error.value = ''
  try {
    groups.value = await service.list()
    if (!groups.value[activeNamespace.value] && namespaces.value[0]) activeNamespace.value = namespaces.value[0]
    for (const item of Object.values(groups.value).flat()) {
      if (item.value_type === 'boolean') booleanDrafts[settingId(item)] = Boolean(item.value)
      else drafts[settingId(item)] = draftValue(item)
    }
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('settings.loadFailed')
  } finally {
    loading.value = false
  }
}

function settingId(item: Pick<SettingResource, 'namespace' | 'key'>) {
  return `${item.namespace}.${item.key}`
}

function draftValue(item: SettingResource): string {
  if (item.value_type === 'boolean') return String(Boolean(item.value))
  if (item.value_type === 'json') return JSON.stringify(item.value, null, 2)
  return String(item.value ?? '')
}

function parseValue(value: string | boolean, type: string): unknown {
  if (type === 'boolean') return Boolean(value)
  if (type === 'integer') return Number.parseInt(String(value), 10)
  if (type === 'number') return Number(String(value))
  if (type === 'json') return JSON.parse(String(value))
  return String(value)
}

async function save(item: SettingResource) {
  const id = settingId(item)
  savingKey.value = id
  try {
    await service.upsert({
      ...item,
      value: parseValue(item.value_type === 'boolean' ? booleanDrafts[id] : (drafts[id] ?? draftValue(item)), item.value_type),
    })
    toast.success(t('settings.saved', { id: `${item.namespace}.${item.key}` }))
    await loadSettings()
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : t('settings.saveFailed'))
  } finally {
    savingKey.value = ''
  }
}

async function createSetting() {
  if (!newForm.namespace.trim() || !newForm.key.trim()) {
    toast.error(t('settings.required'))
    return
  }
  try {
    await service.upsert({
      namespace: newForm.namespace.trim(), key: newForm.key.trim(), value: parseValue(newForm.value, newForm.valueType),
      value_type: newForm.valueType, description: newForm.description.trim(),
    })
    toast.success(t('settings.created'))
    newForm.key = ''
    newForm.value = ''
    newForm.description = ''
    activeNamespace.value = newForm.namespace.trim()
    await loadSettings()
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : t('settings.createFailed'))
  }
}

onMounted(loadSettings)
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('settings.title') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('settings.description') }}</p>
    </div>

    <Alert v-if="error" variant="destructive">
      <AlertTitle>{{ t('settings.unavailable') }}</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <Card>
      <CardHeader>
        <CardTitle>{{ t('settings.createTitle') }}</CardTitle>
        <CardDescription>{{ t('settings.createDescription') }}</CardDescription>
      </CardHeader>
      <CardContent>
        <FieldGroup class="md:grid md:grid-cols-2">
          <Field>
            <FieldLabel for="setting-namespace">{{ t('settings.namespace') }}</FieldLabel>
            <Input id="setting-namespace" v-model="newForm.namespace" placeholder="general" />
          </Field>
          <Field>
            <FieldLabel for="setting-key">{{ t('settings.key') }}</FieldLabel>
            <Input id="setting-key" v-model="newForm.key" placeholder="site_name" />
          </Field>
          <Field>
            <FieldLabel for="setting-type">{{ t('settings.type') }}</FieldLabel>
            <Select v-model="newForm.valueType">
              <SelectTrigger id="setting-type"><SelectValue :placeholder="t('settings.chooseType')" /></SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="string">{{ t('settings.text') }}</SelectItem>
                  <SelectItem value="boolean">{{ t('settings.boolean') }}</SelectItem>
                  <SelectItem value="integer">{{ t('settings.integer') }}</SelectItem>
                  <SelectItem value="number">{{ t('settings.number') }}</SelectItem>
                  <SelectItem value="json">{{ t('settings.json') }}</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </Field>
          <Field>
            <FieldLabel for="setting-value">{{ t('settings.value') }}</FieldLabel>
            <Textarea v-if="newForm.valueType === 'json'" id="setting-value" v-model="newForm.value" placeholder="{}" />
            <Input v-else id="setting-value" v-model="newForm.value" :placeholder="t('settings.valuePlaceholder')" />
          </Field>
          <Field class="md:col-span-2">
            <FieldLabel for="setting-description">{{ t('settings.note') }}</FieldLabel>
            <Input id="setting-description" v-model="newForm.description" :placeholder="t('settings.notePlaceholder')" />
            <FieldDescription>{{ t('settings.noteHint') }}</FieldDescription>
          </Field>
        </FieldGroup>
      </CardContent>
      <CardFooter class="justify-end">
        <Button @click="createSetting">{{ t('settings.create') }}</Button>
      </CardFooter>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle>{{ t('settings.configured') }}</CardTitle>
        <CardDescription>{{ t('settings.configuredDescription') }}</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="flex items-center gap-2 text-sm text-muted-foreground">
          <Spinner /> {{ t('settings.loading') }}
        </div>
        <Empty v-else-if="!namespaces.length">
          <EmptyHeader>
            <EmptyTitle>{{ t('settings.empty') }}</EmptyTitle>
            <EmptyDescription>{{ t('settings.emptyDescription') }}</EmptyDescription>
          </EmptyHeader>
        </Empty>
        <Tabs v-else v-model="activeNamespace" class="gap-6">
          <TabsList>
            <TabsTrigger v-for="namespace in namespaces" :key="namespace" :value="namespace">{{ namespace }}</TabsTrigger>
          </TabsList>
          <TabsContent v-for="namespace in namespaces" :key="namespace" :value="namespace" class="flex flex-col gap-4">
            <div v-for="item in groups[namespace]" :key="settingId(item)" class="rounded-lg border p-4">
              <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
                <Field class="min-w-0 flex-1">
                  <FieldLabel :for="settingId(item)">{{ item.key }}</FieldLabel>
                  <FieldDescription>{{ item.description || t('settings.missingNote') }} · {{ item.value_type }}</FieldDescription>
                  <Switch v-if="item.value_type === 'boolean'" :id="settingId(item)" v-model="booleanDrafts[settingId(item)]" />
                  <Textarea v-else-if="item.value_type === 'json'" :id="settingId(item)" v-model="drafts[settingId(item)]" />
                  <Input v-else :id="settingId(item)" v-model="drafts[settingId(item)]" />
                </Field>
                <Button variant="outline" :disabled="savingKey === settingId(item)" @click="save(item)">
                  <Spinner v-if="savingKey === settingId(item)" data-icon="inline-start" />
                  {{ savingKey === settingId(item) ? t('settings.saving') : t('settings.save') }}
                </Button>
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  </section>
</template>
