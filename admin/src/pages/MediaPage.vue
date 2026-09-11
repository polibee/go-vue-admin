<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Trash2 } from '@lucide/vue'
import { toast } from 'vue-sonner'

import { apiClient } from '@/core/api/client'
import { createGeneratedApiClient, type MediaResource } from '@/generated/api'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Attachment, AttachmentContent, AttachmentDescription, AttachmentMedia, AttachmentTitle } from '@/components/ui/attachment'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { Spinner } from '@/components/ui/spinner'
import { ImageField, MediaPicker } from '@/modules/media/components'
import { MediaService } from '@/modules/media/media.service'
import { useI18n } from 'vue-i18n'

const service = new MediaService(createGeneratedApiClient(apiClient))
const { t } = useI18n()
const items = ref<MediaResource[]>([])
const loading = ref(true)
const uploading = ref(false)
const deleting = ref('')
const error = ref('')

async function loadMedia() {
  loading.value = true
  error.value = ''
  try {
    items.value = await service.list()
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : t('media.unavailable')
  } finally {
    loading.value = false
  }
}

async function upload(file: File) {
  uploading.value = true
  try {
    await service.upload(file)
    toast.success(t('media.uploadSuccess'))
    await loadMedia()
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : t('media.uploadFailed'))
  } finally {
    uploading.value = false
  }
}

async function remove(item: MediaResource) {
  if (!window.confirm(t('media.deleteConfirm', { name: item.original_name }))) return
  deleting.value = item.id
  try {
    await service.remove(item.id)
    items.value = items.value.filter((candidate) => candidate.id !== item.id)
    toast.success(t('media.deleted'))
  } catch (cause) {
    toast.error(cause instanceof Error ? cause.message : t('media.deleteFailed'))
  } finally {
    deleting.value = ''
  }
}

function isImage(item: MediaResource) {
  return item.mime_type.startsWith('image/')
}

function formatSize(size: number) {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / (1024 * 1024)).toFixed(1)} MB`
}

onMounted(loadMedia)
</script>

<template>
  <section class="flex flex-col gap-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('media.title') }}</h1>
      <p class="text-sm text-muted-foreground">{{ t('media.description') }}</p>
    </div>

    <Alert v-if="error" variant="destructive">
      <AlertTitle>{{ t('media.unavailable') }}</AlertTitle>
      <AlertDescription>{{ error }}</AlertDescription>
    </Alert>

    <MediaPicker :busy="uploading" @upload="upload" />

    <Card>
      <CardHeader>
        <CardTitle>{{ t('media.uploaded') }}</CardTitle>
        <CardDescription>{{ t('media.uploadedDescription', { count: items.length }) }}</CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="loading" class="flex items-center gap-2 text-sm text-muted-foreground"><Spinner /> {{ t('media.loading') }}</div>
        <Empty v-else-if="!items.length">
          <EmptyHeader><EmptyTitle>{{ t('media.empty') }}</EmptyTitle></EmptyHeader>
          <EmptyDescription>{{ t('media.emptyDescription') }}</EmptyDescription>
        </Empty>
        <div v-else class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <Attachment v-for="item in items" :key="item.id" orientation="vertical" class="h-full">
            <AttachmentMedia v-if="isImage(item)" class="w-full"><ImageField :src="item.url" :alt="item.original_name" /></AttachmentMedia>
            <AttachmentMedia v-else variant="icon" class="m-4"><ImageField :alt="item.original_name" /></AttachmentMedia>
            <AttachmentContent>
              <AttachmentTitle class="truncate" :title="item.original_name">{{ item.original_name }}</AttachmentTitle>
              <AttachmentDescription>{{ item.mime_type }} · {{ formatSize(item.size) }}</AttachmentDescription>
            </AttachmentContent>
            <Button class="mt-auto self-end" size="icon" variant="ghost" :disabled="deleting === item.id" :aria-label="`${t('resources.delete')} ${item.original_name}`" @click="remove(item)">
              <Trash2 class="size-4" />
            </Button>
          </Attachment>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
