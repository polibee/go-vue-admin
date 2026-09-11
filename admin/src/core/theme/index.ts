export type ThemeMode = 'light' | 'dark' | 'system'
export type ThemePalette = 'shadcn' | 'semi' | 'wechat'

export interface ThemeSettings {
  mode: ThemeMode
  palette: ThemePalette
}

const modeStorageKey = 'go-vue-admin.theme.mode'
const paletteStorageKey = 'go-vue-admin.theme.palette'

function resolvedMode(mode: ThemeMode): 'light' | 'dark' {
  if (mode !== 'system') return mode
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function applyTheme(settings: ThemeSettings): void {
  const root = document.documentElement
  root.classList.toggle('dark', resolvedMode(settings.mode) === 'dark')
  root.dataset.theme = settings.mode
  root.dataset.palette = settings.palette
  window.localStorage.setItem(modeStorageKey, settings.mode)
  window.localStorage.setItem(paletteStorageKey, settings.palette)
}

export function initializeTheme(): ThemeSettings {
  const mode = window.localStorage.getItem(modeStorageKey)
  const palette = window.localStorage.getItem(paletteStorageKey)
  const settings: ThemeSettings = {
    mode: mode === 'light' || mode === 'dark' || mode === 'system' ? mode : 'system',
    palette: palette === 'semi' || palette === 'wechat' || palette === 'shadcn' ? palette : 'shadcn',
  }
  applyTheme(settings)
  return settings
}
