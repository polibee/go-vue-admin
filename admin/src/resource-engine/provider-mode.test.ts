import { describe, expect, it } from 'vitest'

import type { ResourceDataProvider } from './core/ResourceDataProvider'
import { HttpResourceDataProvider } from './providers/HttpResourceDataProvider'
import { createResourceProvider, resolveResourceProviderMode } from './provider-mode'

interface TestRecord extends Record<string, unknown> {
  id: string
}

const memoryProvider: ResourceDataProvider<TestRecord> = {} as ResourceDataProvider<TestRecord>
const httpProvider = new HttpResourceDataProvider<TestRecord>({ request: async <T>() => ({}) as T }, '/api/resources/demo')

describe('resource provider mode', () => {
  it('defaults to memory and accepts explicit http mode', () => {
    expect(resolveResourceProviderMode()).toBe('memory')
    expect(resolveResourceProviderMode('memory')).toBe('memory')
    expect(resolveResourceProviderMode('http')).toBe('http')
  })

  it('rejects unsupported provider modes', () => {
    expect(() => resolveResourceProviderMode('mysql')).toThrow('Unsupported resource provider mode')
  })

  it('returns the provider selected by the explicit mode', () => {
    expect(createResourceProvider('memory', memoryProvider, httpProvider)).toBe(memoryProvider)
    expect(createResourceProvider('http', memoryProvider, httpProvider)).toBe(httpProvider)
  })
})
