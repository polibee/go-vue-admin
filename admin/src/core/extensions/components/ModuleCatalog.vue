<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import type { ExtensionResource } from '@/generated/api'

defineProps<{ items: ExtensionResource[] }>()
const { t } = useI18n()
</script>

<template>
  <div v-if="items.length" class="grid gap-4 md:grid-cols-2">
    <Card v-for="item in items" :key="item.id">
      <CardHeader>
        <CardTitle>{{ item.name }}</CardTitle>
        <CardDescription>{{ t('extensions.moduleKind') }} · {{ item.id }}</CardDescription>
      </CardHeader>
      <CardContent class="flex items-center justify-between gap-4">
        <div class="space-y-1 text-sm">
          <p>{{ t('extensions.status') }}: {{ item.state === 'enabled' ? t('extensions.enabled') : t('extensions.disabled') }}</p>
          <p v-if="item.version" class="text-muted-foreground">{{ t('extensions.version') }}: {{ item.version }}</p>
        </div>
        <div class="flex gap-2">
          <Button v-if="item.config" variant="outline" as-child>
            <RouterLink :to="item.config.route">{{ t('extensions.configure') }}</RouterLink>
          </Button>
          <Button variant="outline" as-child>
            <RouterLink :to="'/admin/modules/' + item.id">{{ t('extensions.details') }}</RouterLink>
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
  <p v-else class="text-sm text-muted-foreground">{{ t('extensions.emptyModules') }}</p>
</template>
