export interface ResourceViewState {
  name: string
  search: string
  filters: Record<string, string>
  sort: string
  direction: 'asc' | 'desc'
  pageSize: string
  hiddenColumns: string[]
}

function storageKey(resource: string) { return `admin.resource.views.${resource}` }

export function loadResourceViews(resource: string): ResourceViewState[] {
  try {
    const value = JSON.parse(localStorage.getItem(storageKey(resource)) || '[]')
    return Array.isArray(value) ? value : []
  } catch {
    return []
  }
}

export function saveResourceView(resource: string, view: ResourceViewState) {
  const views = loadResourceViews(resource).filter((item) => item.name !== view.name)
  localStorage.setItem(storageKey(resource), JSON.stringify([...views, view]))
}

export function deleteResourceView(resource: string, name: string) {
  localStorage.setItem(storageKey(resource), JSON.stringify(loadResourceViews(resource).filter((item) => item.name !== name)))
}
