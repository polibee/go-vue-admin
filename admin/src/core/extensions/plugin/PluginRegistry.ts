import type { Component } from 'vue'

import type { ResourceDefinition } from '@/resource-engine/core/ResourceDefinition'

export type PluginRuntime = 'builtin' | 'external'
export type PluginState = 'disabled' | 'enabled'

export interface PluginDependency {
  id: string
  version?: string
}

export interface PluginMenu {
  id: string
  label: string
  route: string
  permission?: string
}

export interface PluginRoute {
  id: string
  path: string
  view: Component
}

export interface PluginManifest {
  id: string
  name: string
  version: string
  runtime: PluginRuntime
  uiCompatibility: 'shadcn-vue'
  permissions?: readonly string[]
  menus?: readonly PluginMenu[]
  dependencies?: readonly PluginDependency[]
}

export interface PluginRegistrationSnapshot {
  menus: PluginMenu[]
  routes: PluginRoute[]
  resources: ResourceDefinition<object>[]
}

export interface PluginContext {
  registerMenu(menu: PluginMenu): void
  registerRoute(route: PluginRoute): void
  registerResource(resource: ResourceDefinition<object>): void
}

export interface AdminPluginDefinition {
  manifest: PluginManifest
  setup?: (context: PluginContext) => void
  teardown?: (context: PluginContext) => void
}

export interface PluginInfo {
  manifest: PluginManifest
  state: PluginState
}

interface PluginEntry {
  definition: AdminPluginDefinition
  state: PluginState
  registrations: PluginRegistrationSnapshot
}

export function defineAdminPlugin(definition: AdminPluginDefinition): AdminPluginDefinition {
  const id = definition.manifest.id.trim()
  if (!/^[a-z][a-z0-9_-]*$/.test(id)) {
    throw new Error(`Invalid plugin id: ${definition.manifest.id}`)
  }

  return {
    ...definition,
    manifest: {
      ...definition.manifest,
      id,
      permissions: [...(definition.manifest.permissions ?? [])],
      menus: [...(definition.manifest.menus ?? [])],
      dependencies: [...(definition.manifest.dependencies ?? [])],
    },
  }
}

export class PluginRegistry {
  private readonly entries = new Map<string, PluginEntry>()

  register(definition: AdminPluginDefinition): void {
    const manifest = definition.manifest
    if (manifest.runtime !== 'builtin') {
      throw new Error('Only builtin plugins are supported')
    }
    if (manifest.uiCompatibility !== 'shadcn-vue') {
      throw new Error('Only shadcn-vue plugins are supported')
    }
    if (this.entries.has(manifest.id)) {
      throw new Error(`Plugin "${manifest.id}" is already registered`)
    }
    this.entries.set(manifest.id, {
      definition,
      state: 'disabled',
      registrations: { menus: [], routes: [], resources: [] },
    })
  }

  list(): PluginInfo[] {
    return [...this.entries.values()].map((entry) => ({
      manifest: cloneManifest(entry.definition.manifest),
      state: entry.state,
    }))
  }

  state(id: string): PluginState | undefined {
    return this.entries.get(id)?.state
  }

  enable(id: string): void {
    const entry = this.requireEntry(id)
    if (entry.state === 'enabled') return
    for (const dependency of entry.definition.manifest.dependencies ?? []) {
      const dependencyEntry = this.entries.get(dependency.id)
      if (!dependencyEntry) throw new Error(`Plugin dependency "${dependency.id}" is not registered`)
      if (dependencyEntry.state !== 'enabled') throw new Error(`Plugin dependency "${dependency.id}" is not enabled`)
    }

    entry.registrations = { menus: [], routes: [], resources: [] }
    entry.definition.setup?.(createContext(entry.registrations))
    entry.state = 'enabled'
  }

  disable(id: string): void {
    const entry = this.requireEntry(id)
    if (entry.state === 'disabled') return
    for (const [otherId, other] of this.entries) {
      if (other.state !== 'enabled') continue
      if ((other.definition.manifest.dependencies ?? []).some((dependency) => dependency.id === id)) {
        throw new Error(`Plugin dependency "${otherId}" is active`)
      }
    }

    entry.definition.teardown?.(createContext(entry.registrations))
    entry.state = 'disabled'
    entry.registrations = { menus: [], routes: [], resources: [] }
  }

  registrations(id: string): PluginRegistrationSnapshot | undefined {
    const entry = this.entries.get(id)
    if (!entry) return undefined
    return cloneRegistrations(entry.registrations)
  }

  private requireEntry(id: string): PluginEntry {
    const entry = this.entries.get(id)
    if (!entry) throw new Error(`Plugin "${id}" is not registered`)
    return entry
  }
}

function createContext(registrations: PluginRegistrationSnapshot): PluginContext {
  return {
    registerMenu: (menu) => registrations.menus.push({ ...menu }),
    registerRoute: (route) => registrations.routes.push({ ...route }),
    registerResource: (resource) => registrations.resources.push({ ...resource }),
  }
}

function cloneManifest(manifest: PluginManifest): PluginManifest {
  return {
    ...manifest,
    permissions: [...(manifest.permissions ?? [])],
    menus: [...(manifest.menus ?? [])],
    dependencies: [...(manifest.dependencies ?? [])],
  }
}

function cloneRegistrations(registrations: PluginRegistrationSnapshot): PluginRegistrationSnapshot {
  return {
    menus: registrations.menus.map((menu) => ({ ...menu })),
    routes: registrations.routes.map((route) => ({ ...route })),
    resources: registrations.resources.map((resource) => ({ ...resource })),
  }
}
