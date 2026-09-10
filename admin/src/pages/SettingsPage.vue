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

type SettingType = 'string' | 'boolean' | 'integer' | 'number' | 'json'

const service = new SettingsService(createGeneratedApiClient(apiClient))
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
    error.value = cause instanceof Error ? cause.message : '设置加载失败，请稍后重试。'
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
    toast.success(`已保存 ${item.namespace}.${item.key}`)
    await loadSettings()
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : '设置保存失败，请检查输入。')
  } finally {
    savingKey.value = ''
  }
}

async function createSetting() {
  if (!newForm.namespace.trim() || !newForm.key.trim()) {
    toast.error('请填写命名空间和键名。')
    return
  }
  try {
    await service.upsert({
      namespace: newForm.namespace.trim(), key: newForm.key.trim(), value: parseValue(newForm.value, newForm.valueType),
      value_type: newForm.valueType, description: newForm.description.trim(),
    })
    toast.success('设置已创建')
    newForm.key = ''
    newForm.value = ''
    newForm.description = ''
    activeNamespace.value = newForm.namespace.trim()
    await loadSettings()
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : '设置创建失败，请检查输入。')
  }
}

onMounted(loadSettings)
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">设置</h1>
      <p class="text-sm text-muted-foreground">按命名空间管理平台配置，值类型由契约校验。</p>
    </div>

    <Alert v-if="error" variant="destructive">
      <AlertTitle>设置暂时不可用</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <Card>
      <CardHeader>
        <CardTitle>新增设置</CardTitle>
        <CardDescription>使用 namespace + key 定位设置，修改不会改变数据库结构。</CardDescription>
      </CardHeader>
      <CardContent>
        <FieldGroup class="md:grid md:grid-cols-2">
          <Field>
            <FieldLabel for="setting-namespace">命名空间</FieldLabel>
            <Input id="setting-namespace" v-model="newForm.namespace" placeholder="general" />
          </Field>
          <Field>
            <FieldLabel for="setting-key">键名</FieldLabel>
            <Input id="setting-key" v-model="newForm.key" placeholder="site_name" />
          </Field>
          <Field>
            <FieldLabel for="setting-type">值类型</FieldLabel>
            <Select v-model="newForm.valueType">
              <SelectTrigger id="setting-type"><SelectValue placeholder="选择值类型" /></SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="string">文本</SelectItem>
                  <SelectItem value="boolean">布尔值</SelectItem>
                  <SelectItem value="integer">整数</SelectItem>
                  <SelectItem value="number">数字</SelectItem>
                  <SelectItem value="json">JSON</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </Field>
          <Field>
            <FieldLabel for="setting-value">值</FieldLabel>
            <Textarea v-if="newForm.valueType === 'json'" id="setting-value" v-model="newForm.value" placeholder="{}" />
            <Input v-else id="setting-value" v-model="newForm.value" placeholder="输入设置值" />
          </Field>
          <Field class="md:col-span-2">
            <FieldLabel for="setting-description">说明</FieldLabel>
            <Input id="setting-description" v-model="newForm.description" placeholder="帮助管理员理解这个设置" />
            <FieldDescription>建议使用稳定、可读的键名，敏感凭据不应直接写入此处。</FieldDescription>
          </Field>
        </FieldGroup>
      </CardContent>
      <CardFooter class="justify-end">
        <Button @click="createSetting">创建设置</Button>
      </CardFooter>
    </Card>

    <Card>
      <CardHeader>
        <CardTitle>已配置设置</CardTitle>
        <CardDescription>当前显示 Memory 或 MySQL Provider 返回的设置。</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="flex items-center gap-2 text-sm text-muted-foreground">
          <Spinner /> 正在加载设置…
        </div>
        <Empty v-else-if="!namespaces.length">
          <EmptyHeader>
            <EmptyTitle>暂无设置</EmptyTitle>
            <EmptyDescription>请先创建一个命名空间设置。</EmptyDescription>
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
                  <FieldDescription>{{ item.description || '未填写说明' }} · {{ item.value_type }}</FieldDescription>
                  <Switch v-if="item.value_type === 'boolean'" :id="settingId(item)" v-model="booleanDrafts[settingId(item)]" />
                  <Textarea v-else-if="item.value_type === 'json'" :id="settingId(item)" v-model="drafts[settingId(item)]" />
                  <Input v-else :id="settingId(item)" v-model="drafts[settingId(item)]" />
                </Field>
                <Button variant="outline" :disabled="savingKey === settingId(item)" @click="save(item)">
                  <Spinner v-if="savingKey === settingId(item)" data-icon="inline-start" />
                  {{ savingKey === settingId(item) ? '保存中…' : '保存' }}
                </Button>
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  </section>
</template>
