<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

const { t } = useI18n()
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>{{ t('auth.shellReady') }}</CardTitle>
    </CardHeader>
    <CardContent class="text-sm text-muted-foreground">
      {{ t('core.environmentDescription') }}
    </CardContent>
  </Card>
</template>
