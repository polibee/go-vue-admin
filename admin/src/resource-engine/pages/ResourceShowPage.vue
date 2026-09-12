<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { toast } from 'vue-sonner'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import {
  AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
  AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

import { createResourceContext } from '../core/ResourceContext'
import { useAuth } from '@/core/auth'
import { resourceRegistry } from '../demo'

const route = useRoute()
const router = useRouter()
const auth = useAuth()
const name = computed(() => String(route.params.resource ?? ''))
const id = computed(() => String(route.params.id ?? ''))
const definition = computed(() => resourceRegistry.get(name.value))
const provider = computed(() => resourceRegistry.provider(name.value))
const context = computed(() => definition.value && provider.value ? createResourceContext(definition.value, provider.value, auth.user?.permissions ?? []) : null)
const record = ref<Record<string, unknown> | null>(null)
const error = ref<string | null>(null)
const deleteDialogOpen = ref(false)
const deleting = ref(false)
const { t } = useI18n()

onMounted(async () => {
  try {
    if (provider.value) record.value = await provider.value.get(id.value) as Record<string, unknown>
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('resources.loadFailed')
  }
})

async function deleteRecord() {
  if (!provider.value) return
  deleting.value = true
  try {
    await provider.value.delete(id.value)
    toast.success(t('resources.deleted'))
    deleteDialogOpen.value = false
    await router.replace(`/admin/resources/${name.value}`)
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : t('resources.deleteFailed'))
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <section v-if="definition && provider && context" class="flex flex-col gap-6">
    <div class="flex flex-wrap items-start justify-between gap-4"><div><h1 class="text-2xl font-semibold tracking-tight">{{ definition.label }}{{ t('resources.detail') }}</h1><p class="text-sm text-muted-foreground">{{ id }}</p></div><div class="flex gap-2"><Button v-if="context.can('update')" @click="router.push(`/admin/resources/${name}/${id}/edit`)">{{ t('resources.edit') }}</Button><Button v-if="context.can('delete')" variant="destructive" @click="deleteDialogOpen = true">{{ t('resources.delete') }}</Button></div></div>
    <Alert v-if="error" variant="destructive"><AlertTitle>{{ t('resources.loadFailed') }}</AlertTitle><AlertDescription>{{ error }}</AlertDescription></Alert>
    <Card v-if="record"><CardHeader><CardTitle>{{ record.name }}</CardTitle><CardDescription>{{ t('resources.showDescription') }}</CardDescription></CardHeader><CardContent class="grid gap-4 sm:grid-cols-2"><div v-for="field in definition.fields ?? []" :key="String(field.name)" class="flex flex-col gap-1"><span class="text-sm text-muted-foreground">{{ field.label }}</span><span>{{ record[String(field.name)] }}</span></div></CardContent></Card>
    <AlertDialog v-model:open="deleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader><AlertDialogTitle>{{ t('resources.confirmDelete') }}</AlertDialogTitle><AlertDialogDescription>{{ t('resources.deleteDescription') }}</AlertDialogDescription></AlertDialogHeader>
        <AlertDialogFooter><AlertDialogCancel>{{ t('resources.cancel') }}</AlertDialogCancel><AlertDialogAction :disabled="deleting" variant="destructive" @click="deleteRecord">{{ t('resources.confirm') }}</AlertDialogAction></AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </section>
  <Alert v-else variant="destructive"><AlertTitle>{{ t('resources.notFound') }}</AlertTitle><AlertDescription>{{ t('resources.notRegistered', { name }) }}</AlertDescription></Alert>
</template>
