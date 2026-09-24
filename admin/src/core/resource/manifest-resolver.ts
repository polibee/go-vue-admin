import type { ResourceManifest } from '@/generated/api'
import { adminApiBase, adminRoute } from '../routing/url-namespaces.ts'

export function findResourceManifest(manifests: readonly ResourceManifest[], name: string) {
  return manifests.find((manifest) => manifest.name === name)
}

export function frontendResourceRoute(_manifest: Pick<ResourceManifest, 'route'>, name: string) {
  return adminRoute(name)
}

export function resourceApiBase(name: string) {
  return adminApiBase(name)
}
