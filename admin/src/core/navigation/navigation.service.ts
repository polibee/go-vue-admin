import { apiClient } from '@/core/api/client'

import type { NavigationItem } from './NavigationRegistry'

interface MenuItemPayload {
  id: string
  label: string
  route: string
  permission?: string
}

interface ApiEnvelope<T> {
  data: T
}

export class NavigationService {
  async list(): Promise<NavigationItem[]> {
    const { data } = await apiClient.request<ApiEnvelope<MenuItemPayload[]>>('/api/menu')
    return data
  }
}
