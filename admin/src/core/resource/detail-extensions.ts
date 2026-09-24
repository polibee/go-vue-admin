import type { Component } from 'vue'

export interface ResourceDetailExtensionProps {
  resource: string
  id: string
  record: Record<string, unknown>
}

export interface ResourceDetailExtension {
  resource: string
  component: Component
}

const extensions = new Map<string, ResourceDetailExtension>()

export function registerResourceDetailExtension(extension: ResourceDetailExtension) {
  if (extensions.has(extension.resource)) {
    throw new Error(`Resource detail extension already registered: ${extension.resource}`)
  }
  extensions.set(extension.resource, extension)
}

export function getResourceDetailExtension(resource: string) {
  return extensions.get(resource)
}
