<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { supportedLocales } from '@/core/i18n'

const { locale, t } = useI18n()

function switchLocale(): void {
  locale.value = locale.value === supportedLocales[0] ? supportedLocales[1] : supportedLocales[0]
}
</script>

<template>
  <main class="flex min-h-svh items-center justify-center bg-background p-6">
    <Card class="w-full max-w-xl">
      <CardHeader>
        <div class="flex items-center justify-between gap-4">
          <CardTitle>{{ t('app.title') }}</CardTitle>
          <Badge variant="secondary">shadcn-vue</Badge>
        </div>
        <CardDescription>{{ t('app.foundation') }}</CardDescription>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <p class="text-sm text-muted-foreground">{{ t('app.description') }}</p>
        <Button class="w-fit" variant="outline" @click="switchLocale">
          {{ locale === 'zh-CN' ? 'English' : '中文' }}
        </Button>
      </CardContent>
    </Card>
  </main>
</template>
