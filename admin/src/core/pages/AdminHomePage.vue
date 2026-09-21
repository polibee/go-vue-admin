<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { generatedApi, type AdminOverview } from '@/generated/api'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const router = useRouter()
const auth = useAuthStore()
const overview = ref<AdminOverview>()
const loading = ref(true)
const error = ref(false)

onMounted(async () => {
  if (!auth.token) return
  try { overview.value = await generatedApi.overview(auth.token) } catch { error.value = true } finally { loading.value = false }
})
</script>

<template>
  <div class="flex flex-col gap-6">
    <Card>
    <CardHeader>
      <CardTitle>{{ t('auth.shellReady') }}</CardTitle>
    </CardHeader>
    <CardContent class="text-sm text-muted-foreground">
      {{ t('core.environmentDescription') }}
    </CardContent>
    </Card>
    <div class="grid gap-4 md:grid-cols-3">
      <Card v-for="item in [{ key: 'users', route: '/users' }, { key: 'roles', route: '/rbac' }, { key: 'permissions', route: '/rbac' }]" :key="item.key">
        <CardHeader><CardTitle class="text-sm font-medium">{{ t(`auth.overview${item.key[0].toUpperCase()}${item.key.slice(1)}`) }}</CardTitle></CardHeader>
        <CardContent><Skeleton v-if="loading" class="h-8 w-16" /><p v-else-if="error" class="text-sm text-destructive">{{ t('states.errorTitle') }}</p><p v-else class="text-3xl font-semibold">{{ overview?.[item.key as keyof AdminOverview] }}</p><Button v-if="!loading && !error" variant="link" class="mt-2 px-0" @click="router.push(item.route)">{{ t('auth.viewDetails') }}</Button></CardContent>
      </Card>
    </div>
  </div>
</template>
