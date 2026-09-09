import type { ResourceDataProvider } from './core/ResourceDataProvider'

export type ResourceProviderMode = 'memory' | 'http'

export function resolveResourceProviderMode(value = import.meta.env.VITE_RESOURCE_PROVIDER): ResourceProviderMode {
  const normalized = value?.trim().toLowerCase() || 'memory'
  if (normalized === 'memory' || normalized === 'http') return normalized
  throw new Error(`Unsupported resource provider mode: ${normalized}`)
}

export function createResourceProvider<T extends object>(
  mode: ResourceProviderMode,
  memoryProvider: ResourceDataProvider<T>,
  httpProvider: ResourceDataProvider<T>,
): ResourceDataProvider<T> {
  return mode === 'http' ? httpProvider : memoryProvider
}
