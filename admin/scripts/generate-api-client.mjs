import { mkdir, writeFile } from 'node:fs/promises'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const baseUrl = process.env.OPENAPI_BASE_URL ?? 'http://127.0.0.1:3000'
const response = await fetch(`${baseUrl}/api/openapi.json`)
if (!response.ok) throw new Error(`OpenAPI request failed: ${response.status}`)
const spec = await response.json()
const statuses = spec.components?.schemas?.UserStatus?.enum
if (!Array.isArray(statuses) || statuses.length === 0) throw new Error('UserStatus enum is missing from the OpenAPI contract')
const paths = Object.keys(spec.paths ?? {})
const outputPath = resolve(dirname(fileURLToPath(import.meta.url)), '../src/generated/api.ts')

const content = `/* eslint-disable */
/* Generated from ${baseUrl}/api/openapi.json. DO NOT EDIT. */
/* Contract paths: ${paths.join(', ')} */

import { apiFetch, apiFetchEnvelope } from '@/lib/api'

export type UserStatus = ${statuses.map((status) => JSON.stringify(status)).join(' | ')}
export interface AuthUser { id: number; name: string; email: string; status: UserStatus; locale: string }
export interface LoginRequest { email: string; password: string }
export interface LoginResponse { access_token: string; token_type: string; user: AuthUser }
export interface RefreshResponse { access_token: string; token_type: string }
export interface ResourceManifest { name: string; label: string; route: string; columns: Array<{ name: string; label: string; sortable: boolean }>; actions?: Array<{ name: string; label: string; kind: string; permission: string }> }
export interface ResourceListMeta { page: number; per_page: number; total: number; last_page: number }
export interface ResourceList<T = Record<string, unknown>> { data: T[]; meta: ResourceListMeta }
export interface BulkUserStatusRequest { user_ids: number[]; status: UserStatus }
export interface AdminOverview { users: number; roles: number; permissions: number }
export interface AuditLog { id: number; user_id: number; action: string; metadata: Record<string, unknown> | null; created_at: string }

export const generatedApi = {
  login(request: LoginRequest) { return apiFetch<LoginResponse>('/api/v1/auth/login', { method: 'POST', body: JSON.stringify(request) }) },
  refresh() { return apiFetch<RefreshResponse>('/api/v1/auth/refresh', { method: 'POST' }) },
  currentUser(token: string) { return apiFetch<AuthUser>('/api/v1/auth/me', {}, token) },
  logout(token: string) { return apiFetch<void>('/api/v1/auth/logout', { method: 'POST' }, token) },
  logoutAll(token: string) { return apiFetch<void>('/api/v1/auth/logout-all', { method: 'POST' }, token) },
  resourceManifests(token: string) { return apiFetch<ResourceManifest[]>('/api/v1/admin/resources', {}, token) },
  overview(token: string) { return apiFetch<AdminOverview>('/api/v1/admin/overview', {}, token) },
  auditLogs(token: string) { return apiFetch<AuditLog[]>('/api/v1/admin/audit-logs', {}, token) },
  resourceList<T = Record<string, unknown>>(resource: string, query: URLSearchParams, token: string) { return apiFetchEnvelope<T[]>('/api/v1/admin/resources/' + resource + '?' + query, {}, token) as unknown as Promise<ResourceList<T>> },
  bulkSetUserStatus(request: BulkUserStatusRequest, token: string) { return apiFetch<void>('/api/v1/admin/users/status', { method: 'PUT', body: JSON.stringify(request) }, token) },
}
`

await mkdir(dirname(outputPath), { recursive: true })
await writeFile(outputPath, content, 'utf8')
console.log(`Generated ${outputPath} from ${paths.length} OpenAPI paths.`)
