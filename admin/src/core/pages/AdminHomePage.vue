<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { generatedApi, type AdminOverview, type ResourceManifest } from '@/generated/api'
import { dashboardResourceRoute, visibleDashboardResources } from '@/lib/dashboard-resources'
import { useAuthStore } from '@/stores/auth'
import { localizedResourceLabel } from '@/core/resource/resource-i18n'

const { t, te } = useI18n()
const router = useRouter()
const auth = useAuthStore()
const overview = ref<AdminOverview>()
const manifests = ref<ResourceManifest[]>([])
const loading = ref(true)
const error = ref(false)
const overviewKeys: Record<string, keyof AdminOverview> = { users: 'users', roles: 'roles', permissions: 'permissions' }
const resources = computed(() => visibleDashboardResources(manifests.value, auth.user?.permissions || []).map((resource) => ({
  ...resource,
  label: localizedResourceLabel(t, te, resource.name, resource.label),
  route: dashboardResourceRoute(resource),
  overviewKey: overviewKeys[resource.name],
})))

onMounted(async () => {
  if (!auth.token) return
  try {
    const [overviewData, registry] = await Promise.all([generatedApi.overview(auth.token), generatedApi.resourceRegistry(auth.token)])
    overview.value = overviewData
    manifests.value = registry
  } catch { error.value = true } finally { loading.value = false }
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
    <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <Card v-for="item in resources" :key="item.name">
        <CardHeader><CardTitle class="text-sm font-medium">{{ item.label }}</CardTitle></CardHeader>
        <CardContent><Skeleton v-if="loading" class="h-8 w-16" /><p v-else-if="error" class="text-sm text-destructive">{{ t('states.errorTitle') }}</p><p v-else class="text-3xl font-semibold">{{ item.overviewKey ? overview?.[item.overviewKey] : '—' }}</p><Button v-if="!loading && !error" variant="link" class="mt-2 px-0" @click="router.push(item.route)">{{ t('auth.viewDetails') }}</Button></CardContent>
      </Card>
    </div>
  </div>
</template>
