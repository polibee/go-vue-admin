import type { ResourceDefinition } from './ResourceDefinition'

export class ResourceRegistry {
  private readonly resources = new Map<string, ResourceDefinition<object>>()

  register<T extends object>(definition: ResourceDefinition<T>): void {
    if (this.resources.has(definition.name)) {
      throw new Error(`Resource "${definition.name}" is already registered`)
    }
    this.resources.set(definition.name, definition as ResourceDefinition<object>)
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

  all(): ResourceDefinition<object>[] {
    return [...this.resources.values()]
  }
}
