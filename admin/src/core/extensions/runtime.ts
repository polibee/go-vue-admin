import { ModuleRegistry } from './ModuleRegistry'
import { PluginRegistry } from './plugin/PluginRegistry'
import { moduleDefinition } from '../../../../modules/example/admin/module'
import { examplePlugin } from '../../../../plugins/example/admin/plugin'

export const moduleRuntime = new ModuleRegistry()
export const pluginRuntime = new PluginRegistry()

export function registerBuiltinExtensions(): void {
  if (!moduleRuntime.get(moduleDefinition.id)) moduleRuntime.register(moduleDefinition)
  if (!pluginRuntime.state(examplePlugin.manifest.id)) pluginRuntime.register(examplePlugin)
}
