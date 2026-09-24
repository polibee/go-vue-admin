<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'

const { t } = useI18n()
const router = useRouter()
</script>

<template>
  <main class="flex min-h-full items-center justify-center">
    <Alert variant="destructive" class="max-w-lg">
      <AlertTitle>{{ t('states.errorTitle') }}</AlertTitle>
      <AlertDescription class="flex flex-col gap-4">
        {{ t('states.errorDescription') }}
        <Button variant="outline" class="w-fit" @click="router.replace('/admin')">{{ t('states.backHome') }}</Button>
      </AlertDescription>
    </Alert>
  </main>
</template>
