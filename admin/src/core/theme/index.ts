export type ThemePreference = 'light' | 'dark' | 'system'

const themeStorageKey = 'go-vue-admin.theme'

function resolvedTheme(preference: ThemePreference): 'light' | 'dark' {
  if (preference !== 'system') {
    return preference
  }

  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function applyTheme(preference: ThemePreference): void {
  const root = document.documentElement
  root.classList.toggle('dark', resolvedTheme(preference) === 'dark')
  root.dataset.theme = preference
  window.localStorage.setItem(themeStorageKey, preference)
}

export function initializeTheme(): ThemePreference {
  const stored = window.localStorage.getItem(themeStorageKey)
  const preference: ThemePreference = stored === 'light' || stored === 'dark' || stored === 'system'
    ? stored
    : 'system'

  applyTheme(preference)
  return preference
}
