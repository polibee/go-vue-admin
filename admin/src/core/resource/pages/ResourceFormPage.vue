<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import ResourceFormView from '@/components/resource/ResourceFormView.vue'
import { generatedApi, type ResourceManifest } from '@/generated/api'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{ resource?: string }>()
const route = useRoute()
const auth = useAuthStore()
const manifest = ref<ResourceManifest>()
const resourceName = computed(() => props.resource || String(route.params.resource || ''))
const api = computed(() => ({
  show: (id: string | number, token: string) => generatedApi.resourceShow(resourceName.value, id, token),
  create: (payload: Record<string, unknown>, token: string) => generatedApi.resourceCreate(resourceName.value, payload, token),
  update: (id: string | number, payload: Record<string, unknown>, token: string) => generatedApi.resourceUpdate(resourceName.value, id, payload, token),
  relationOptions: (relation: string, token: string) => generatedApi.resourceRelationOptions(resourceName.value, relation, token),
}))

onMounted(async () => {
  if (!auth.token) return
  const manifests = await generatedApi.resourceRegistry(auth.token)
  manifest.value = manifests.find((item) => item.name === resourceName.value)
})
</script>

<template>
  <ResourceFormView v-if="manifest" :resource="manifest" :api="api" />
</template>
