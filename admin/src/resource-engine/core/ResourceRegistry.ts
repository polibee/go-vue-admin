import type { ResourceDefinition } from './ResourceDefinition'
import { shallowReactive } from 'vue'
import type { ResourceDataProvider } from './ResourceDataProvider'

export class ResourceRegistry {
  private readonly resources = shallowReactive(new Map<string, ResourceDefinition<object>>())
  private readonly providers = shallowReactive(new Map<string, ResourceDataProvider<object>>())

  remove(name: string): void {
    this.resources.delete(name)
    this.providers.delete(name)
  }

  register<T extends object>(definition: ResourceDefinition<T>, provider?: ResourceDataProvider<T>): void {
    if (this.resources.has(definition.name)) {
      throw new Error(`Resource "${definition.name}" is already registered`)
    }
    this.resources.set(definition.name, definition as ResourceDefinition<object>)
    if (provider) this.providers.set(definition.name, provider as ResourceDataProvider<object>)
  }

  get<T extends object = object>(name: string): ResourceDefinition<T> | undefined {
    return this.resources.get(name) as ResourceDefinition<T> | undefined
  }

  require<T extends object = object>(name: string): ResourceDefinition<T> {
    const definition = this.get<T>(name)
    if (!definition) {
      throw new Error(`Resource "${name}" is not registered`)
    }
    return definition
  }

  provider<T extends object = object>(name: string): ResourceDataProvider<T> | undefined {
    return this.providers.get(name) as ResourceDataProvider<T> | undefined
  }

  requireProvider<T extends object = object>(name: string): ResourceDataProvider<T> {
    const provider = this.provider<T>(name)
    if (!provider) throw new Error(`Resource provider "${name}" is not registered`)
    return provider
  }

  all(): ResourceDefinition<object>[] {
    return [...this.resources.values()]
  }
}
