<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
import PluginCatalog from './components/PluginCatalog.vue'
import { mountPlugin, pluginRuntime, unmountPlugin } from './runtime'

const client = createGeneratedApiClient(apiClient)
const plugins = ref<ExtensionResource[]>([])
const error = ref('')
const { t } = useI18n()
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
    error.value = cause instanceof Error ? cause.message : t('extensions.loadFailed')
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
    error.value = cause instanceof Error ? cause.message : t('extensions.actionFailed')
  } finally {
    busy.value = ''
  }
}

onMounted(load)
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('extensions.pluginTitle') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('extensions.pluginDescription') }}</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>{{ t('extensions.actionFailed') }}</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <PluginCatalog :items="plugins" :busy="busy" @toggle="setPluginState" />
  </section>
</template>
