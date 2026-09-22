export type ResourcePageMode = 'generic' | 'custom'

export interface ResourcePageManifest {
  name: string
  pageMode?: ResourcePageMode
}

export interface ResourcePageOverrides {
  list: unknown
  form: unknown
  detail: unknown
}

export function resolveResourcePageMode(manifest: ResourcePageManifest): ResourcePageMode {
  return manifest.pageMode ?? 'generic'
}

export function validateCustomPageOverrides(overrides: Partial<ResourcePageOverrides>): void {
  for (const key of ['list', 'form', 'detail'] as const) {
    if (!overrides[key]) {
      throw new Error('custom resource page override is missing: ' + key)
    }
  }
}

export function resolveResourcePages<T>(
  manifest: ResourcePageManifest,
  genericPages: ResourcePageOverrides,
  customPages?: Partial<ResourcePageOverrides>,
): ResourcePageOverrides {
  if (resolveResourcePageMode(manifest) === 'generic') {
    return genericPages
  }
  validateCustomPageOverrides(customPages ?? {})
  return customPages as ResourcePageOverrides
}
