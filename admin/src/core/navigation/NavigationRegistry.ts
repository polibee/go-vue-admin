import type { Component } from 'vue'

import { can, type Permission } from '@/core/permissions'

export interface NavigationItem {
  id: string
  label: string
  route: string
  icon?: Component
  permission?: Permission
}

export class NavigationRegistry {
  private readonly entries = new Map<string, NavigationItem>()

  constructor(items: NavigationItem[] = []) {
    this.registerMany(items)
  }

  register(item: NavigationItem): void {
    const previous = this.entries.get(item.id)
    this.entries.set(item.id, { ...previous, ...item, icon: item.icon ?? previous?.icon })
  }

  registerMany(items: NavigationItem[]): void {
    items.forEach((item) => this.register(item))
  }

  all(): NavigationItem[] {
    return [...this.entries.values()]
  }

  visible(granted: readonly Permission[]): NavigationItem[] {
    return this.all().filter((item) => !item.permission || can(granted, item.permission))
  }
}
