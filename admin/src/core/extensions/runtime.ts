import { ModuleRegistry } from './ModuleRegistry'
import { PluginRegistry } from './plugin/PluginRegistry'
import { moduleDefinition } from '../../../../modules/example/admin/module'
import { examplePlugin } from '../../../../plugins/example/admin/plugin'
import { apiClient } from '@/core/api/client'
import { HttpResourceDataProvider } from '@/resource-engine/providers/HttpResourceDataProvider'
import { MemoryResourceDataProvider } from '@/resource-engine/core/ResourceDataProvider'
import { resourceRegistry } from '@/resource-engine/demo'
import { createResourceProvider, resolveResourceProviderMode } from '@/resource-engine/provider-mode'
import type { AdminModuleDefinition } from './ModuleRegistry'
import { router } from '@/core/router'
import { navigationRegistry } from '@/core/navigation'
import { RuntimeRegistrationBridge } from './RuntimeRegistrationBridge'

export const moduleRuntime = new ModuleRegistry()
export const pluginRuntime = new PluginRegistry()
export const runtimeBridge = new RuntimeRegistrationBridge({ router, navigationRegistry, resourceRegistry, providerFactory: createProvider })

export function registerBuiltinExtensions(): void {
  registerModule(moduleDefinition)
  if (!pluginRuntime.state(examplePlugin.manifest.id)) pluginRuntime.register(examplePlugin)
}

function registerModule(definition: AdminModuleDefinition): void {
  if (!moduleRuntime.get(definition.id)) moduleRuntime.register(definition)
  runtimeBridge.mountModule(definition)
}

export function mountPlugin(id: string): void {
  const registrations = pluginRuntime.registrations(id)
  if (registrations) runtimeBridge.mountPlugin(id, registrations)
}

export function unmountPlugin(id: string): void {
  runtimeBridge.unmount('plugin:' + id)
}

function createProvider(resource: import('@/resource-engine/core/ResourceDefinition').ResourceDefinition<object>) {
  const mode = resolveResourceProviderMode()
  const httpProvider = new HttpResourceDataProvider<Record<string, unknown>>(
    apiClient,
    resource.endpoint ?? '/api/resources/' + resource.name,
  )
  const memoryProvider = new MemoryResourceDataProvider<Record<string, unknown>>()
  return createResourceProvider(mode, memoryProvider, httpProvider)
}
