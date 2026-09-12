<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type ExtensionResource } from '@/generated/api'
import { useI18n } from 'vue-i18n'
import ModuleCatalog from './components/ModuleCatalog.vue'
import PluginCatalog from './components/PluginCatalog.vue'

const client = createGeneratedApiClient(apiClient)
const modules = ref<ExtensionResource[]>([])
const plugins = ref<ExtensionResource[]>([])
const activeTab = ref('modules')
const error = ref('')
const busy = ref('')
const { t } = useI18n()

async function load() {
  error.value = ''
  try {
    const [moduleResponse, pluginResponse] = await Promise.all([client.listModules(), client.listPlugins()])
    modules.value = moduleResponse.data
    plugins.value = pluginResponse.data
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('extensions.loadFailed')
  }
}

async function setPluginState(item: ExtensionResource, state: 'enabled' | 'disabled') {
  busy.value = item.id
  error.value = ''
  try {
    plugins.value = (await client.setPluginState(item.id, state)).data
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
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('extensions.moduleTitle') }} & {{ t('extensions.pluginTitle') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('extensions.moduleDescription') }} {{ t('extensions.pluginDescription') }}</p>
    </div>
    <Alert v-if="error" variant="destructive">
      <AlertTitle>{{ t('extensions.actionFailed') }}</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>
    <Tabs v-model="activeTab" class="gap-6">
      <TabsList>
        <TabsTrigger value="modules">{{ t('extensions.moduleTitle') }}</TabsTrigger>
        <TabsTrigger value="plugins">{{ t('extensions.pluginTitle') }}</TabsTrigger>
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
