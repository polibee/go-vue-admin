import type { ResourceManifest } from '@/generated/api'

export function findResourceManifest(manifests: readonly ResourceManifest[], name: string) {
  return manifests.find((manifest) => manifest.name === name)
}

export function frontendResourceRoute(_manifest: Pick<ResourceManifest, 'route'>, name: string) {
  return `/${name}`
}
