import type { RouteRecordRaw } from 'vue-router'

import type { SupportedLocale } from '@/core/i18n'
import type { NavigationItem } from '@/core/navigation/NavigationRegistry'
import type { ResourceDefinition } from '@/resource-engine/core/ResourceDefinition'

export interface AdminModuleDefinition {
  id: string
  resources?: readonly ResourceDefinition<object>[]
  routes?: readonly RouteRecordRaw[]
  navigation?: readonly NavigationItem[]
  locales?: Partial<Record<SupportedLocale, Record<string, unknown>>>
}

export function defineAdminModule(definition: AdminModuleDefinition): AdminModuleDefinition {
  const id = definition.id.trim()
  if (!/^[a-z][a-z0-9_-]*$/.test(id)) {
    throw new Error(`Invalid module id: ${definition.id}`)
  }

  return {
    id,
    resources: [...(definition.resources ?? [])],
    routes: [...(definition.routes ?? [])],
    navigation: [...(definition.navigation ?? [])],
    locales: { ...definition.locales },
  }
}

export class ModuleRegistry {
  private readonly entries = new Map<string, AdminModuleDefinition>()

  register(module: AdminModuleDefinition): void {
    if (this.entries.has(module.id)) {
      throw new Error(`Module "${module.id}" is already registered`)
    }
    this.entries.set(module.id, module)
  }

  registerMany(modules: readonly AdminModuleDefinition[]): void {
    modules.forEach((module) => this.register(module))
  }

  get(id: string): AdminModuleDefinition | undefined {
    return this.entries.get(id)
  }

  remove(id: string): boolean {
    return this.entries.delete(id)
  }

  all(): AdminModuleDefinition[] {
    return [...this.entries.values()]
  }
}
