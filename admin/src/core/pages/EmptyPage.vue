<script setup lang="ts">
import { FolderOpen } from '@lucide/vue'
import { useI18n } from 'vue-i18n'
import { Empty, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@/components/ui/empty'

const { t } = useI18n()
</script>

<template>
  <Empty class="min-h-full">
    <EmptyHeader>
      <EmptyMedia variant="icon"><FolderOpen /></EmptyMedia>
      <EmptyTitle>{{ t('states.emptyTitle') }}</EmptyTitle>
      <EmptyDescription>{{ t('states.emptyDescription') }}</EmptyDescription>
    </EmptyHeader>
  </Empty>
</template>
