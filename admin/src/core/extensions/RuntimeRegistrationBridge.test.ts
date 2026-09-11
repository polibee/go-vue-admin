import { describe, expect, it, vi } from 'vitest'
import { createMemoryHistory, createRouter } from 'vue-router'

import { RuntimeRegistrationBridge } from './RuntimeRegistrationBridge'
import { NavigationRegistry } from '@/core/navigation'
import { ResourceRegistry } from '@/resource-engine/core/ResourceRegistry'

const View = { template: '<div />' }

function createBridge() {
  const router = createRouter({ history: createMemoryHistory(), routes: [{ path: '/admin', name: 'admin', component: View }] })
  const navigation = new NavigationRegistry()
  const resources = new ResourceRegistry()
  const bridge = new RuntimeRegistrationBridge({
    router,
    navigationRegistry: navigation,
    resourceRegistry: resources,
    providerFactory: vi.fn(),
  })
  return { bridge, router, navigation, resources }
}

describe('RuntimeRegistrationBridge', () => {
  it('mounts a module navigation item and route once', () => {
    const { bridge, router, navigation } = createBridge()
    const definition = {
      id: 'example',
      routes: [{ path: 'example', name: 'module-example', component: View }],
      navigation: [{ id: 'module-example', label: '示例模块', route: '/admin/example' }],
    }

    bridge.mountModule(definition)
    bridge.mountModule(definition)

    expect(navigation.all().filter((item) => item.id === 'module-example')).toHaveLength(1)
    expect(router.hasRoute('module:example:module-example')).toBe(true)
  })

  it('unmounts plugin navigation, route, and resource together', () => {
    const { bridge, router, navigation, resources } = createBridge()
    const resource = { name: 'plugin-records', label: '插件记录', fields: [] } as never

    bridge.mountPlugin('example-plugin', {
      menus: [{ id: 'plugin-example', label: '示例插件', route: '/admin/example-plugin' }],
      routes: [{ id: 'plugin-example', path: 'example-plugin', view: View }],
      resources: [resource],
    })
    bridge.unmount('plugin:example-plugin')

    expect(navigation.all().some((item) => item.id === 'plugin-example')).toBe(false)
    expect(router.hasRoute('plugin:example-plugin:plugin-example')).toBe(false)
    expect(resources.get('plugin-records')).toBeUndefined()
  })
})
