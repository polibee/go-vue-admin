import { apiFetch } from '@/lib/api'

export type SettingValueType = 'string' | 'boolean' | 'integer' | 'json'

export interface SystemSetting {
  id: number
  key: string
  value: string
  value_type: SettingValueType
  group: string
  description: string
}

export interface SettingInput {
  value: string
  value_type: SettingValueType
  group: string
  description: string
}

export function listSystemSettings(token: string) {
  return apiFetch<SystemSetting[]>('/api/v1/admin/settings', {}, token)
}

export function saveSystemSetting(key: string, input: SettingInput, token: string) {
  return apiFetch<SystemSetting>(`/api/v1/admin/settings/${encodeURIComponent(key)}`, {
    method: 'PUT',
    body: JSON.stringify(input),
  }, token)
}
