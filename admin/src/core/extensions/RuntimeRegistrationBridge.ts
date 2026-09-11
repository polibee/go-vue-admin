import type { Router, RouteRecordRaw } from 'vue-router'

import type { AdminModuleDefinition } from './ModuleRegistry'
import type { PluginRegistrationSnapshot } from './plugin/PluginRegistry'
import { navigationRegistry, type NavigationItem } from '@/core/navigation'
import type { ResourceDataProvider } from '@/resource-engine/core/ResourceDataProvider'
import type { ResourceDefinition } from '@/resource-engine/core/ResourceDefinition'
import type { ResourceRegistry } from '@/resource-engine/core/ResourceRegistry'

interface RuntimeRegistrationBridgeOptions {
  router: Router
  navigationRegistry: typeof navigationRegistry
  resourceRegistry: ResourceRegistry
  providerFactory?: (resource: ResourceDefinition<object>) => ResourceDataProvider<object> | undefined
}

export class RuntimeRegistrationBridge {
  private readonly routeNames = new Map<string, string[]>()
  private readonly resourceOwners = new Map<string, string>()

  constructor(private readonly options: RuntimeRegistrationBridgeOptions) {}

  mountModule(definition: AdminModuleDefinition): void {
    const owner = `module:${definition.id}`
    this.unmount(owner)
    this.mountNavigation(owner, definition.navigation ?? [])
    this.mountRoutes(owner, definition.routes ?? [])
    this.mountResources(owner, (definition.resources ?? []) as readonly ResourceDefinition<object>[])
  }

  mountPlugin(id: string, snapshot: PluginRegistrationSnapshot): void {
    const owner = `plugin:${id}`
    this.unmount(owner)
    this.mountNavigation(owner, snapshot.menus)
    this.mountRoutes(owner, snapshot.routes.map((route) => ({
      path: route.path,
      name: route.id,
      component: route.view,
    })))
    this.mountResources(owner, snapshot.resources)
  }

  unmount(owner: string): void {
    this.options.navigationRegistry.removeOwner(owner)
    for (const routeName of this.routeNames.get(owner) ?? []) {
      if (this.options.router.hasRoute(routeName)) this.options.router.removeRoute(routeName)
    }
    this.routeNames.delete(owner)
    for (const [resourceName, resourceOwner] of this.resourceOwners) {
      if (resourceOwner !== owner) continue
      if (this.options.resourceRegistry.get(resourceName)) this.options.resourceRegistry.remove(resourceName)
      this.resourceOwners.delete(resourceName)
    }
  }

  private mountNavigation(owner: string, items: readonly NavigationItem[]): void {
    this.options.navigationRegistry.registerMany(items.map((item) => ({ ...item, owner })))
  }

  private mountRoutes(owner: string, routes: readonly RouteRecordRaw[]): void {
    const routeNames: string[] = []
    routes.forEach((route, index) => {
      const name = `${owner}:${String(route.name ?? index)}`
      this.options.router.addRoute('admin', { ...route, name })
      routeNames.push(name)
    })
    this.routeNames.set(owner, routeNames)
  }

  private mountResources(owner: string, resources: readonly ResourceDefinition<object>[]): void {
    resources.forEach((resource) => {
      if (this.options.resourceRegistry.get(resource.name)) return
      this.options.resourceRegistry.register(resource, this.options.providerFactory?.(resource))
      this.resourceOwners.set(resource.name, owner)
    })
  }
}
