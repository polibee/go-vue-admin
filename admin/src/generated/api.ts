/* eslint-disable */
/* Generated from http://127.0.0.1:3000/api/openapi.json. DO NOT EDIT. */
/* Contract paths: /admin/audit-logs, /admin/overview, /admin/registry, /admin/search, /admin/settings, /admin/settings/{key}, /admin/users/status, /admin/{resource}, /admin/{resource}/export, /admin/{resource}/{id}, /auth/login, /auth/logout-all, /auth/me, /auth/refresh */

import { ApiError, apiDownload, apiFetch, apiFetchEnvelope } from '@/lib/api'

export type UserStatus = "active" | "disabled" | "locked"
export type DataScope = "all" | "own"
export interface AuthUser { id: number; name: string; email: string; status: UserStatus; locale: string; permissions: string[] }
export interface LoginRequest { email: string; password: string }
export interface LoginResponse { access_token: string; token_type: string; user: AuthUser }
export interface RefreshResponse { access_token: string; token_type: string }
export interface ResourceManifest { name: string; label: string; route: string; permissions: string[]; data_scope?: DataScope; owner_field?: string; fields: Array<{ name: string; label: string; type: string; options?: Array<{ value: string; label: string }> }>; columns: Array<{ name: string; label: string; sortable: boolean }>; actions?: Array<{ name: string; label: string; kind: string; permission: string }> }
export interface ResourceListMeta { page: number; per_page: number; total: number; last_page: number }
export interface ResourceList<T = Record<string, unknown>> { data: T[]; meta: ResourceListMeta }
export interface BulkUserStatusRequest { user_ids: number[]; status: UserStatus }
export interface AdminOverview { users: number; roles: number; permissions: number }
export interface AuditLog { id: number; user_id: number; action: string; metadata: Record<string, unknown> | null; created_at: string }
export interface GlobalSearchResult { resource: string; label: string; id: string | number; title: string; subtitle?: string; route: string }
export interface RolePermissionAssignment { id: number; name: string; display_name: string; scope: DataScope }

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
  bulkSetUserStatus(request: BulkUserStatusRequest, token: string) { return apiFetch<void>('/api/v1/admin/users/status', { method: 'PUT', body: JSON.stringify(request) }, token) },
  replaceRolePermissions(roleID: number, permissionIDs: number[], scopes: Record<string, DataScope>, token: string) { return apiFetch<void>('/api/v1/admin/roles/' + roleID + '/permissions', { method: 'PUT', body: JSON.stringify({ permission_ids: permissionIDs, scopes }) }, token) },
}
