import { describe, expect, it } from 'vitest'
import type { Component } from 'vue'

import { defineAdminPlugin, PluginRegistry } from './PluginRegistry'

const plugin = (id: string, dependencies: string[] = []) => defineAdminPlugin({
  manifest: {
    id,
    name: `${id} plugin`,
    version: '1.0.0',
    runtime: 'builtin',
    uiCompatibility: 'shadcn-vue',
    permissions: [`${id}.view`],
    menus: [{ id, label: id, route: `/admin/${id}`, permission: `${id}.view` }],
    dependencies: dependencies.map((dependency) => ({ id: dependency })),
  },
})

describe('PluginRegistry', () => {
  it('keeps plugin order and starts disabled', () => {
    const registry = new PluginRegistry()
    registry.register(plugin('alpha'))
    registry.register(plugin('beta'))

    expect(registry.list().map((item) => item.manifest.id)).toEqual(['alpha', 'beta'])
    expect(registry.list().every((item) => item.state === 'disabled')).toBe(true)
  })

  it('enables and disables a builtin plugin once per transition', () => {
    const calls: string[] = []
    const registry = new PluginRegistry()
    registry.register(defineAdminPlugin({
      manifest: plugin('alpha').manifest,
      setup: () => calls.push('enable'),
      teardown: () => calls.push('disable'),
    }))

    registry.enable('alpha')
    registry.enable('alpha')
    registry.disable('alpha')
    registry.disable('alpha')

    expect(calls).toEqual(['enable', 'disable'])
    expect(registry.state('alpha')).toBe('disabled')
  })

  it('requires enabled dependencies and rejects external runtimes', () => {
    const registry = new PluginRegistry()
    registry.register(plugin('alpha'))
    registry.register(plugin('beta', ['alpha']))

    expect(() => registry.enable('beta')).toThrow('Plugin dependency "alpha" is not enabled')
    registry.enable('alpha')
    registry.enable('beta')
    expect(() => registry.disable('alpha')).toThrow('Plugin dependency "beta" is active')

    expect(() => registry.register(defineAdminPlugin({
      manifest: { ...plugin('external').manifest, runtime: 'external' },
    }))).toThrow('Only builtin plugins are supported')
  })

  it('provides a context registration surface without framework internals', () => {
    const registry = new PluginRegistry()
    registry.register(defineAdminPlugin({
      manifest: plugin('alpha').manifest,
      setup: (context) => {
        context.registerMenu({ id: 'alpha-menu', label: 'Alpha', route: '/admin/alpha' })
        context.registerRoute({ id: 'alpha-route', path: '/alpha', view: {} as Component })
      },
    }))

    registry.enable('alpha')

    expect(registry.registrations('alpha')).toEqual({
      menus: [{ id: 'alpha-menu', label: 'Alpha', route: '/admin/alpha' }],
      routes: [{ id: 'alpha-route', path: '/alpha', view: {} }],
      resources: [],
    })
  })
})
