import { computed, onMounted, ref } from 'vue'
import { generatedApi, type Notification } from '@/generated/api'
import { useAuthStore } from '@/stores/auth'
import { isAdminRoute } from '@/core/routing/url-namespaces'

export function isInternalNotificationURL(value?: string | null) {
  return Boolean(value && isAdminRoute(value) && !value.startsWith('//') && !/[\r\n]/.test(value))
}

export function useNotifications() {
  const auth = useAuthStore()
  const items = ref<Notification[]>([])
  const loading = ref(false)
  const error = ref(false)
  const unreadCount = ref(0)

  async function refresh() {
    if (!auth.token) return
    loading.value = true
    error.value = false
    try {
      const [list, unread] = await Promise.all([
        generatedApi.notifications(auth.token, new URLSearchParams({ per_page: '10' })),
        generatedApi.notificationUnreadCount(auth.token),
      ])
      items.value = list.data
      unreadCount.value = unread.count
    } catch {
      error.value = true
    } finally {
      loading.value = false
    }
  }

  async function markRead(id: number) {
    if (!auth.token) return
    try {
      await generatedApi.markNotificationRead(id, auth.token)
      const item = items.value.find((entry) => entry.id === id)
      if (item && !item.read_at) {
        item.read_at = new Date().toISOString()
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch {
      error.value = true
    }
  }

  async function markAllRead() {
    if (!auth.token || unreadCount.value === 0) return
    try {
      await generatedApi.markAllNotificationsRead(auth.token)
      const now = new Date().toISOString()
      items.value = items.value.map((item) => ({ ...item, read_at: item.read_at || now }))
      unreadCount.value = 0
    } catch {
      error.value = true
    }
  }

  onMounted(refresh)

  return { items, loading, error, unreadCount: computed(() => unreadCount.value), refresh, markRead, markAllRead }
}
