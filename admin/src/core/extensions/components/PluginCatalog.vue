<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import type { ExtensionResource } from '@/generated/api'

defineProps<{ items: ExtensionResource[]; busy: string }>()
const emit = defineEmits<{ toggle: [item: ExtensionResource, state: 'enabled' | 'disabled'] }>()
const { t } = useI18n()
</script>

<template>
  <div v-if="items.length" class="grid gap-4 md:grid-cols-2">
    <Card v-for="item in items" :key="item.id">
      <CardHeader>
        <CardTitle>{{ item.name }}</CardTitle>
        <CardDescription>{{ t('extensions.pluginKind') }} · {{ item.id }}</CardDescription>
      </CardHeader>
      <CardContent class="space-y-1">
        <p>{{ t('extensions.status') }}: {{ item.state === 'enabled' ? t('extensions.enabled') : t('extensions.disabled') }}</p>
        <p v-if="item.version" class="text-sm text-muted-foreground">{{ t('extensions.version') }}: {{ item.version }}</p>
      </CardContent>
      <CardFooter class="gap-2">
        <Button
          v-if="item.state === 'disabled'"
          :disabled="busy === item.id"
          @click="emit('toggle', item, 'enabled')"
        >
          {{ t('extensions.enable') }}
        </Button>
        <Button
          v-else
          variant="outline"
          :disabled="busy === item.id"
          @click="emit('toggle', item, 'disabled')"
        >
          {{ t('extensions.disable') }}
        </Button>
        <Button variant="outline" as-child>
          <RouterLink :to="'/admin/plugins/' + item.id">{{ t('extensions.view') }}</RouterLink>
        </Button>
        <Button v-if="item.config" variant="outline" as-child>
          <RouterLink :to="item.config.route">{{ t('extensions.configure') }}</RouterLink>
        </Button>
      </CardFooter>
    </Card>
  </div>
  <p v-else class="text-sm text-muted-foreground">{{ t('extensions.emptyPlugins') }}</p>
</template>
