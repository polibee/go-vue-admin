<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Skeleton } from '@/components/ui/skeleton'

const { t } = useI18n()
</script>

<template>
  <main class="flex min-h-full flex-col gap-4">
    <p class="text-sm text-muted-foreground">{{ t('states.loading') }}</p>
    <Skeleton class="h-24 w-full" />
    <Skeleton class="h-24 w-full" />
  </main>
</template>
