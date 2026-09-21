<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <main class="flex min-h-svh items-center justify-center bg-background p-6 text-foreground">
    <section class="flex max-w-xl flex-col gap-3 text-center">
      <h1 class="text-2xl font-semibold tracking-tight">{{ t('core.appName') }}</h1>
      <p class="text-muted-foreground">{{ t('core.environmentReady') }}</p>
      <p class="text-sm text-muted-foreground">{{ t('core.environmentDescription') }}</p>
    </section>
  </main>
</template>
