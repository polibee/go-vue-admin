import { shallowReactive, type Component } from 'vue'

import { can, type Permission } from '@/core/permissions'

export interface NavigationItem {
  id: string
  label: string
  route: string
  icon?: Component
  permission?: Permission
  owner?: string
  group?: string
  groupLabel?: string
  order?: number
}

export interface NavigationGroup {
  id: string
  label: string
  icon?: Component
  items: NavigationItem[]
}

export type NavigationSection =
  | { kind: 'item'; item: NavigationItem }
  | { kind: 'group'; group: NavigationGroup }

export class NavigationRegistry {
  private readonly entries = shallowReactive(new Map<string, NavigationItem>())

  remove(id: string): void { this.entries.delete(id) }

  removeOwner(owner: string): void {
    for (const [id, item] of this.entries) {
      if (item.owner === owner) this.entries.delete(id)
    }
  }

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

  sections(granted: readonly Permission[]): NavigationSection[] {
    const items = this.visible(granted)
    const owned = new Map<string, NavigationItem[]>()
    const standalone: NavigationItem[] = []

    for (const item of items) {
      if (!item.owner) {
        standalone.push(item)
        continue
      }
      const groupKey = item.group ? `group:${item.group}` : `owner:${item.owner}`
      const groupItems = owned.get(groupKey) ?? []
      groupItems.push(item)
      owned.set(groupKey, groupItems)
    }

    const sections: NavigationSection[] = standalone.map((item) => ({ kind: 'item', item }))
    for (const [groupKey, groupItems] of owned) {
      const sorted = [...groupItems].sort((left, right) => (left.order ?? 0) - (right.order ?? 0))
      if (sorted.length === 1) {
        sections.push({ kind: 'item', item: sorted[0] })
        continue
      }
      const first = sorted[0]
      sections.push({
        kind: 'group',
        group: {
          id: groupKey,
          label: first.groupLabel ?? first.owner?.replace(/^(module|plugin):/, '') ?? groupKey,
          icon: first.icon,
          items: sorted,
        },
      })
    }

    return sections
  }
}
