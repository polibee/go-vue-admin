<script setup lang="ts">
import { computed } from 'vue'
import { ref } from 'vue'
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
const definition = computed(() => resourceRegistry.get(name.value))
const provider = computed(() => resourceRegistry.provider(name.value))
const context = computed(() => definition.value && provider.value ? createResourceContext(definition.value, provider.value, auth.user?.permissions ?? []) : null)
const error = ref<string | null>(null)
const { t } = useI18n()

async function create(values: Record<string, unknown>) {
  if (!provider.value || !definition.value || !context.value?.can('create')) return
  const id = `${name.value}-${Date.now()}`
  error.value = null
  try {
    await provider.value.create({ id, ...values })
    toast.success(t('resources.created'))
    await router.replace(`/admin/resources/${name.value}/${id}`)
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('resources.createFailed')
  }
}
</script>

<template>
  <section v-if="definition && provider && context" class="flex flex-col gap-6">
    <div><h1 class="text-2xl font-semibold tracking-tight">{{ t('resources.createTitle', { name: definition.label }) }}</h1><p class="text-sm text-muted-foreground">{{ t('resources.createDescription') }}</p></div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('resources.createFailed') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="context.can('create')"><CardHeader><CardTitle>{{ t('resources.basic') }}</CardTitle><CardDescription>{{ t('resources.fieldsDescription') }}</CardDescription></CardHeader><CardContent><ResourceForm :definition="definition" :submit-label="t('resources.create')" @submit="create" /></CardContent></Card>
    <Alert v-else variant="destructive"><AlertTitle>{{ t('resources.noPermission') }}</AlertTitle><AlertDescription>{{ t('resources.createPermission') }}</AlertDescription></Alert>
  </section>
  <Alert v-else variant="destructive"><AlertTitle>{{ t('resources.notFound') }}</AlertTitle><AlertDescription>{{ t('resources.notRegistered', { name }) }}</AlertDescription></Alert>
</template>
