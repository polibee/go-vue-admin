<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { toast } from 'vue-sonner'
import { useRoute, useRouter } from 'vue-router'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { createResourceContext } from '../core/ResourceContext'
import { useAuth } from '@/core/auth'
import { resourceRegistry } from '../demo'
import ResourceForm from '../form/ResourceForm.vue'
import { useI18n } from 'vue-i18n'

const route = useRoute()
const router = useRouter()
const auth = useAuth()
const name = computed(() => String(route.params.resource ?? ''))
const id = computed(() => String(route.params.id ?? ''))
const definition = computed(() => resourceRegistry.get(name.value))
const provider = computed(() => resourceRegistry.provider(name.value))
const context = computed(() => definition.value && provider.value ? createResourceContext(definition.value, provider.value, auth.user?.permissions ?? []) : null)
const initialValues = ref<Record<string, unknown>>({})
const error = ref<string | null>(null)
const loaded = ref(false)
const { t } = useI18n()

onMounted(async () => {
  try {
    if (!provider.value) return
    initialValues.value = await provider.value.get(id.value) as Record<string, unknown>
    loaded.value = true
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('resources.loadFailed')
  }
})

async function update(values: Record<string, unknown>) {
  if (!provider.value || !context.value?.can('update')) return
  error.value = null
  try {
    await provider.value.update(id.value, values)
    toast.success(t('resources.saved'))
    await router.replace(`/admin/resources/${name.value}/${id.value}`)
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('resources.saveFailed')
  }
}
</script>

<template>
  <section v-if="definition && provider && context" class="flex flex-col gap-6">
    <div><h1 class="text-2xl font-semibold tracking-tight">{{ t('resources.editTitle', { name: definition.label }) }}</h1><p class="text-sm text-muted-foreground">{{ id }}</p></div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('resources.loadFailed') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="context.can('update')"><CardHeader><CardTitle>{{ t('resources.basic') }}</CardTitle><CardDescription>{{ t('resources.editDescription') }}</CardDescription></CardHeader><CardContent><ResourceForm v-if="loaded" :definition="definition" :initial-values="initialValues" :submit-label="t('resources.save')" @submit="update" /></CardContent></Card>
    <Alert v-else variant="destructive"><AlertTitle>{{ t('resources.noPermission') }}</AlertTitle><AlertDescription>{{ t('resources.updatePermission') }}</AlertDescription></Alert>
  </section>
  <Alert v-else variant="destructive"><AlertTitle>{{ t('resources.notFound') }}</AlertTitle><AlertDescription>{{ t('resources.notRegistered', { name }) }}</AlertDescription></Alert>
</template>
