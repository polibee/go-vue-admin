<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
import ModuleCatalog from './components/ModuleCatalog.vue'
import PluginCatalog from './components/PluginCatalog.vue'

const client = createGeneratedApiClient(apiClient)
const modules = ref<ExtensionResource[]>([])
const plugins = ref<ExtensionResource[]>([])
const activeTab = ref('modules')
const error = ref('')
const busy = ref('')

async function load() {
  error.value = ''
  try {
    const [moduleResponse, pluginResponse] = await Promise.all([client.listModules(), client.listPlugins()])
    modules.value = moduleResponse.data
    plugins.value = pluginResponse.data
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '模块和插件列表加载失败'
  }
}

async function setPluginState(item: ExtensionResource, state: 'enabled' | 'disabled') {
  busy.value = item.id
  error.value = ''
  try {
    plugins.value = (await client.setPluginState(item.id, state)).data
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '插件状态更新失败'
  } finally {
    busy.value = ''
  }
}

onMounted(load)
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">模块与插件</h1>
      <p class="text-sm text-muted-foreground">业务模块负责产品能力，平台插件负责可插拔基础能力。</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>操作失败</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <Tabs v-model="activeTab" class="gap-6">
      <TabsList>
        <TabsTrigger value="modules">业务模块</TabsTrigger>
        <TabsTrigger value="plugins">平台插件</TabsTrigger>
      </TabsList>
      <TabsContent value="modules">
        <ModuleCatalog :items="modules" />
      </TabsContent>
      <TabsContent value="plugins">
        <PluginCatalog :items="plugins" :busy="busy" @toggle="setPluginState" />
      </TabsContent>
    </Tabs>
  </section>
</template>
