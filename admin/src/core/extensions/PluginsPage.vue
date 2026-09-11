<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
import PluginCatalog from './components/PluginCatalog.vue'
import { mountPlugin, pluginRuntime, unmountPlugin } from './runtime'

const client = createGeneratedApiClient(apiClient)
const plugins = ref<ExtensionResource[]>([])
const error = ref('')
const busy = ref('')

async function load() {
  try {
    plugins.value = (await client.listPlugins()).data
    for (const item of plugins.value) {
      if (item.state !== 'enabled' || pluginRuntime.state(item.id) === 'enabled') continue
      pluginRuntime.enable(item.id)
      mountPlugin(item.id)
    }
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : '平台插件加载失败'
  }
}

async function setPluginState(item: ExtensionResource, state: 'enabled' | 'disabled') {
  busy.value = item.id
  error.value = ''
  try {
    plugins.value = (await client.setPluginState(item.id, state)).data
    if (state === 'enabled') {
      pluginRuntime.enable(item.id)
      mountPlugin(item.id)
    } else {
      pluginRuntime.disable(item.id)
      unmountPlugin(item.id)
    }
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
      <h1 class="text-2xl font-semibold tracking-tight">平台插件</h1>
      <p class="text-sm text-muted-foreground">平台插件负责可插拔的基础能力和运行时扩展。</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>操作失败</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <PluginCatalog :items="plugins" :busy="busy" @toggle="setPluginState" />
  </section>
</template>
