<script setup lang="ts">
import { Bell, CheckCheck } from '@lucide/vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { useNotifications, isInternalNotificationURL } from './useNotifications'

const { t } = useI18n()
const router = useRouter()
const { items, loading, error, unreadCount, markRead, markAllRead } = useNotifications()

function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, { dateStyle: 'short', timeStyle: 'short' }).format(new Date(value))
}

async function openNotification(id: number, url?: string | null) {
  await markRead(id)
  if (url && isInternalNotificationURL(url)) void router.push(url)
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" size="icon" class="relative" :aria-label="t('core.notifications')">
        <Bell />
        <span v-if="unreadCount" class="absolute right-0.5 top-0.5 min-w-4 rounded-full bg-destructive px-1 text-[10px] leading-4 text-destructive-foreground">{{ unreadCount > 99 ? '99+' : unreadCount }}</span>
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end" class="w-80">
      <DropdownMenuLabel class="flex items-center justify-between gap-3"><span>{{ t('core.notifications') }}</span><Button v-if="unreadCount" variant="ghost" size="sm" class="h-7 px-2" @click.stop="markAllRead"><CheckCheck data-icon="inline-start" />{{ t('core.markAllRead') }}</Button></DropdownMenuLabel>
      <DropdownMenuSeparator />
      <DropdownMenuItem v-if="loading" disabled>{{ t('core.notificationsLoading') }}</DropdownMenuItem>
      <DropdownMenuItem v-else-if="error" disabled>{{ t('core.notificationsError') }}</DropdownMenuItem>
      <DropdownMenuItem v-else-if="!items.length" disabled>{{ t('core.notificationsEmpty') }}</DropdownMenuItem>
      <template v-else>
        <DropdownMenuItem v-for="item in items" :key="item.id" class="items-start gap-3 whitespace-normal" @click="openNotification(item.id, item.url)">
          <span class="mt-1 size-2 shrink-0 rounded-full" :class="item.read_at ? 'bg-muted' : 'bg-primary'" />
          <span class="min-w-0 flex-1"><span class="block font-medium">{{ item.title }}</span><span class="block text-xs text-muted-foreground">{{ item.body }}</span><span class="mt-1 block text-[11px] text-muted-foreground">{{ formatDate(item.created_at) }}</span></span>
        </DropdownMenuItem>
      </template>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
