<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Skeleton } from '@/components/ui/skeleton'
import ResourceFormView from '@/components/resource/ResourceFormView.vue'
import { generatedApi, type ResourceManifest } from '@/generated/api'
import { findResourceManifest, frontendResourceRoute } from '@/core/resource/manifest-resolver'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const props = defineProps<{ resource?: string }>()
const route = useRoute()
const auth = useAuthStore()
const { t } = useI18n()
const manifest = ref<ResourceManifest>()
const loading = ref(true)
const error = ref('')
const resourceName = computed(() => props.resource || String(route.params.resource || ''))
const api = computed(() => ({
  show: (id: string | number, token: string) => generatedApi.resourceShow(resourceName.value, id, token),
  create: (payload: Record<string, unknown>, token: string) => generatedApi.resourceCreate(resourceName.value, payload, token),
  update: (id: string | number, payload: Record<string, unknown>, token: string) => generatedApi.resourceUpdate(resourceName.value, id, payload, token),
  relationOptions: (relation: string, query: URLSearchParams, token: string) => generatedApi.resourceRelationOptions(resourceName.value, relation, query, token),
}))

onMounted(async () => {
  if (!auth.token) {
    error.value = t('errors.unknown')
    loading.value = false
    return
  }
  try {
    const manifests = await generatedApi.resourceRegistry(auth.token)
    const found = findResourceManifest(manifests, resourceName.value)
    manifest.value = found ? { ...found, route: frontendResourceRoute(found, resourceName.value) } : undefined
    if (!manifest.value) error.value = t('resource.resourceNotFound')
  } catch {
    error.value = t('errors.unknown')
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div v-if="loading" class="flex flex-col gap-3"><Skeleton class="h-8 w-48" /><Skeleton class="h-64 w-full" /></div>
  <Alert v-else-if="error" variant="destructive"><AlertTitle>{{ t('states.errorTitle') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
  <ResourceFormView v-else-if="manifest" :resource="manifest" :api="api" />
</template>
