import { mkdir, writeFile } from 'node:fs/promises'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const baseUrl = process.env.OPENAPI_BASE_URL ?? 'http://127.0.0.1:3000'
const response = await fetch(`${baseUrl}/api/openapi.json`)
if (!response.ok) throw new Error(`OpenAPI request failed: ${response.status}`)
const spec = await response.json()
const statuses = spec.components?.schemas?.UserStatus?.enum
if (!Array.isArray(statuses) || statuses.length === 0) throw new Error('UserStatus enum is missing from the OpenAPI contract')
const dataScopes = spec.components?.schemas?.DataScope?.enum
if (!Array.isArray(dataScopes) || dataScopes.length === 0) throw new Error('DataScope enum is missing from the OpenAPI contract')
const paths = Object.keys(spec.paths ?? {})
const outputPath = resolve(dirname(fileURLToPath(import.meta.url)), '../src/generated/api.ts')

const content = `/* eslint-disable */
/* Generated from ${baseUrl}/api/openapi.json. DO NOT EDIT. */
/* Contract paths: ${paths.join(', ')} */

import { ApiError, apiDownload, apiFetch, apiFetchEnvelope } from '@/lib/api'

export type UserStatus = ${statuses.map((status) => JSON.stringify(status)).join(' | ')}
export type DataScope = ${dataScopes.map((scope) => JSON.stringify(scope)).join(' | ')}
export interface AuthUser { id: number; name: string; email: string; status: UserStatus; locale: string; permissions: string[] }
export interface LoginRequest { email: string; password: string }
export interface LoginResponse { access_token: string; token_type: string; user: AuthUser }
export interface RefreshResponse { access_token: string; token_type: string }
export interface ResourceField { name: string; label: string; type: string; required?: boolean; visible: boolean; readable: boolean; writable: boolean; sensitive: boolean; options?: Array<{ value: string; label: string }> }
export interface ResourceManifest { name: string; label: string; route: string; permissions: string[]; data_scope?: DataScope; owner_field?: string; fields: ResourceField[]; columns: Array<{ name: string; label: string; sortable: boolean }>; actions?: Array<{ name: string; label: string; kind: string; permission: string; batch: boolean; payload?: string }>; relations?: ResourceRelation[]; form_groups?: ResourceFormGroup[]; details?: ResourceDetailSection[]; dependencies?: ResourceFieldDependency[] }
export interface ResourceRelation { name: string; kind: 'belongsTo' | 'hasMany'; resource: string; field: string; foreign_field: string; label_field: string; selectable: boolean; multiple: boolean; permission?: string; filter_fields?: string[] }
export interface ResourceFormGroup { name: string; label: string; columns?: number; fields: string[] }
export interface ResourceDetailSection { name: string; label: string; fields: string[] }
export interface ResourceFieldDependency { field: string; on: string; value: string }
export interface RelationOption { value: string; label: string }
export interface RelationOptionList { data: RelationOption[]; meta: ResourceListMeta }
export interface ResourceListMeta { page: number; per_page: number; total: number; last_page: number }
export interface ResourceList<T = Record<string, unknown>> { data: T[]; meta: ResourceListMeta }
export interface ActionRequest { ids: number[]; payload?: Record<string, unknown> }
export interface ActionFailure { id: number; code: string }
export interface ActionResponse { action: string; requested: number; succeeded: number; failed: number; skipped: number; failures: ActionFailure[]; skips: ActionFailure[] }
export interface AdminOverview { users: number; roles: number; permissions: number }
export interface AuditLog { id: number; user_id: number; action: string; metadata: Record<string, unknown> | null; created_at: string }
export interface GlobalSearchResult { resource: string; label: string; id: string | number; title: string; subtitle?: string; route: string }
export interface FieldPermissionOverride { readable: boolean; writable: boolean }
export interface RolePermissionAssignment { id: number; name: string; display_name: string; scope: DataScope; fields?: Record<string, FieldPermissionOverride> }

export const generatedApi = {
  login(request: LoginRequest) { return apiFetch<LoginResponse>('/api/v1/auth/login', { method: 'POST', body: JSON.stringify(request) }) },
  refresh() { return apiFetch<RefreshResponse>('/api/v1/auth/refresh', { method: 'POST' }) },
  currentUser(token: string) { return apiFetch<AuthUser>('/api/v1/auth/me', {}, token) },
  logout(token: string) { return apiFetch<void>('/api/v1/auth/logout', { method: 'POST' }, token) },
  logoutAll(token: string) { return apiFetch<void>('/api/v1/auth/logout-all', { method: 'POST' }, token) },
  resourceRegistry(token: string) {
    return apiFetch<unknown>('/api/v1/admin/registry', {}, token).then((payload) => {
      if (!Array.isArray(payload)) throw new ApiError('Invalid resource registry response', 502, 'INTERNAL_ERROR')
      return payload as ResourceManifest[]
    })
  },
  globalSearch(query: string, token: string) {
    return apiFetch<GlobalSearchResult[]>('/api/v1/admin/search?q=' + encodeURIComponent(query), {}, token)
  },
  overview(token: string) { return apiFetch<AdminOverview>('/api/v1/admin/overview', {}, token) },
  auditLogs(token: string, query: URLSearchParams) { return apiFetchEnvelope<AuditLog[]>('/api/v1/admin/audit-logs?' + query, {}, token) as unknown as Promise<ResourceList<AuditLog>> },
  resourceList<T = Record<string, unknown>>(resource: string, query: URLSearchParams, token: string) { return apiFetchEnvelope<T[]>('/api/v1/admin/' + resource + '?' + query, {}, token) as unknown as Promise<ResourceList<T>> },
  resourceExport(resource: string, query: URLSearchParams, token: string) { return apiDownload('/api/v1/admin/' + resource + '/export?' + query, {}, token) },
  resourceShow<T = Record<string, unknown>>(resource: string, id: string | number, token: string) { return apiFetch<T>('/api/v1/admin/' + resource + '/' + id, {}, token) },
  resourceCreate<T = Record<string, unknown>>(resource: string, payload: Record<string, unknown>, token: string) { return apiFetch<T>('/api/v1/admin/' + resource, { method: 'POST', body: JSON.stringify(payload) }, token) },
  resourceUpdate<T = Record<string, unknown>>(resource: string, id: string | number, payload: Record<string, unknown>, token: string) { return apiFetch<T>('/api/v1/admin/' + resource + '/' + id, { method: 'PUT', body: JSON.stringify(payload) }, token) },
  resourceDelete(resource: string, id: string | number, token: string) { return apiFetch<void>('/api/v1/admin/' + resource + '/' + id, { method: 'DELETE' }, token) },
  resourceAction(resource: string, action: string, request: ActionRequest, token: string) { return apiFetch<ActionResponse>('/api/v1/admin/' + resource + '/actions/' + action, { method: 'POST', body: JSON.stringify(request) }, token).then((payload) => (payload as { data?: ActionResponse }).data || payload as ActionResponse) },
  resourceRelationOptions(resource: string, relation: string, query: URLSearchParams = new URLSearchParams(), token: string) { return apiFetchEnvelope<RelationOption[]>('/api/v1/admin/' + resource + '/relations/' + relation + '/options?' + query, {}, token) },
  resourceRelationRecords(resource: string, id: string | number, relation: string, token: string) { return apiFetchEnvelope<RelationOption[]>('/api/v1/admin/' + resource + '/' + id + '/relations/' + relation, {}, token) },
  replaceRolePermissions(roleID: number, permissionIDs: number[], scopes: Record<string, DataScope>, fields: Record<string, Record<string, FieldPermissionOverride>>, token: string) { return apiFetch<void>('/api/v1/admin/roles/' + roleID + '/permissions', { method: 'PUT', body: JSON.stringify({ permission_ids: permissionIDs, scopes, fields }) }, token) },
}
`

await mkdir(dirname(outputPath), { recursive: true })
await writeFile(outputPath, content, 'utf8')
console.log(`Generated ${outputPath} from ${paths.length} OpenAPI paths.`)
