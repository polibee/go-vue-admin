import { describe, expect, it } from 'vitest'

import { defineAdminModule, ModuleRegistry } from './ModuleRegistry'

const moduleDefinition = (id: string) => defineAdminModule({
  id,
  navigation: [{ id: `${id}-nav`, label: id, route: `/admin/${id}` }],
  locales: { 'zh-CN': { label: id }, en: { label: id } },
})

describe('ModuleRegistry', () => {
  it('preserves module registration order', () => {
    const registry = new ModuleRegistry()
    registry.register(moduleDefinition('alpha'))
    registry.register(moduleDefinition('beta'))

    expect(registry.all().map((item) => item.id)).toEqual(['alpha', 'beta'])
  })

  it('rejects duplicate module ids', () => {
    const registry = new ModuleRegistry()
    registry.register(moduleDefinition('alpha'))

    expect(() => registry.register(moduleDefinition('alpha'))).toThrow('Module "alpha" is already registered')
  })

  it('removes one module without changing the others', () => {
    const registry = new ModuleRegistry()
    registry.register(moduleDefinition('alpha'))
    registry.register(moduleDefinition('beta'))

    expect(registry.remove('alpha')).toBe(true)
    expect(registry.all().map((item) => item.id)).toEqual(['beta'])
    expect(registry.remove('alpha')).toBe(false)
  })
})
