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

export const moduleRuntime = new ModuleRegistry()
export const pluginRuntime = new PluginRegistry()

export function registerBuiltinExtensions(): void {
  registerModule(moduleDefinition)
  if (!pluginRuntime.state(examplePlugin.manifest.id)) pluginRuntime.register(examplePlugin)
}

function registerModule(definition: AdminModuleDefinition): void {
  if (!moduleRuntime.get(definition.id)) moduleRuntime.register(definition)
  const mode = resolveResourceProviderMode()
  for (const resource of definition.resources ?? []) {
    if (resourceRegistry.get(resource.name)) continue
    const httpProvider = new HttpResourceDataProvider<Record<string, unknown>>(
      apiClient,
      resource.endpoint ?? `/api/resources/${resource.name}`,
    )
    const memoryProvider = new MemoryResourceDataProvider<Record<string, unknown>>()
    resourceRegistry.register(resource, createResourceProvider(mode, memoryProvider, httpProvider))
  }
}
