/* eslint-disable */
/* Generated from http://127.0.0.1:3000/api/openapi.json. DO NOT EDIT. */
/* Contract paths: /admin/overview, /admin/resources, /admin/resources/{resource}, /admin/users/status, /auth/login, /auth/me */

import { apiFetch, apiFetchEnvelope } from '@/lib/api'

export type UserStatus = "active" | "disabled" | "locked"
export interface AuthUser { id: number; name: string; email: string; status: UserStatus; locale: string }
export interface LoginRequest { email: string; password: string }
export interface LoginResponse { access_token: string; token_type: string; user: AuthUser }
export interface ResourceManifest { name: string; label: string; route: string; columns: Array<{ name: string; label: string; sortable: boolean }>; actions?: Array<{ name: string; label: string; kind: string; permission: string }> }
export interface ResourceListMeta { page: number; per_page: number; total: number; last_page: number }
export interface ResourceList<T = Record<string, unknown>> { data: T[]; meta: ResourceListMeta }
export interface BulkUserStatusRequest { user_ids: number[]; status: UserStatus }
export interface AdminOverview { users: number; roles: number; permissions: number }

export const generatedApi = {
  login(request: LoginRequest) { return apiFetch<LoginResponse>('/api/v1/auth/login', { method: 'POST', body: JSON.stringify(request) }) },
  currentUser(token: string) { return apiFetch<AuthUser>('/api/v1/auth/me', {}, token) },
  logout(token: string) { return apiFetch<void>('/api/v1/auth/logout', { method: 'POST' }, token) },
  resourceManifests(token: string) { return apiFetch<ResourceManifest[]>('/api/v1/admin/resources', {}, token) },
  overview(token: string) { return apiFetch<AdminOverview>('/api/v1/admin/overview', {}, token) },
  resourceList<T = Record<string, unknown>>(resource: string, query: URLSearchParams, token: string) { return apiFetchEnvelope<T[]>('/api/v1/admin/resources/' + resource + '?' + query, {}, token) as unknown as Promise<ResourceList<T>> },
  bulkSetUserStatus(request: BulkUserStatusRequest, token: string) { return apiFetch<void>('/api/v1/admin/users/status', { method: 'PUT', body: JSON.stringify(request) }, token) },
}
