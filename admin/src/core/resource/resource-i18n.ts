type Translate = (key: string, params?: Record<string, unknown>) => string
type HasTranslation = (key: string) => boolean

function localized(t: Translate, te: HasTranslation, key: string, fallback: string) {
  return te(key) ? t(key) : fallback
}

export function localizedResourceLabel(t: Translate, te: HasTranslation, resource: string, fallback: string) {
  return localized(t, te, `resource.labels.${resource}`, fallback)
}

export function localizedFieldLabel(t: Translate, te: HasTranslation, resource: string, field: string, fallback: string) {
  return localized(t, te, `resource.fields.${resource}.${field}`, fallback)
}

export function localizedActionLabel(t: Translate, te: HasTranslation, action: string, fallback: string) {
  return localized(t, te, `resource.actionLabels.${action}`, fallback)
}

export function localizedOptionLabel(t: Translate, te: HasTranslation, field: string, value: string, fallback: string) {
  return localized(t, te, `resource.options.${field}.${value}`, fallback)
}
