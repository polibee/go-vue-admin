# Admin URL Namespace Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Move all Vue admin pages under `/admin/`, reserve root URLs for C-end pages, and keep backend admin APIs under `/api/v1/admin/` with one shared contract.

**Architecture:** The Vue router will use an `/admin` parent route while public/C-end routes remain at the root. Resource manifests will expose separate `adminRoute` and `apiBase` values, and backend route registration, OpenAPI, generated clients, menus, notifications, and tests will consume those distinct namespaces. The migration is a direct cutover with no legacy root admin aliases.

**Tech Stack:** Goravel/Go, Vue 3, TypeScript, Vue Router, Vite, OpenAPI, Vitest/Node tests, Go tests.

**Spec:** `docs/superpowers/specs/2026-09-24-admin-url-namespace-design.md`

## Global Constraints

- Backend admin pages use `/admin/...`.
- C-end pages remain at root paths such as `/` and `/orders`.
- Backend admin APIs use `/api/v1/admin/...`.
- Authentication APIs remain under `/api/v1/auth/...`.
- Resource manifests must distinguish `adminRoute` from `apiBase`.
- Do not preserve old root admin aliases.
- Do not add dependencies or change the API version.
- Do not put `/api/...` paths in Vue Router or `/admin/...` paths in API clients.

## Review Focus

- Direct navigation to `/admin/login`, `/admin/users`, and `/admin/users/:id/edit` must resolve to Vue pages after refresh.
- Unauthenticated admin navigation must redirect to `/admin/login?redirect=...`, while C-end root routes remain independent.
- Generated resource routes and notification links must never lose the `/admin` prefix.
- Backend route registration must not leave a legacy `/users` page/API collision.
- OpenAPI paths, generated TypeScript client paths, Vite proxy behavior, and browser URLs must agree.

---

### Task 1: Establish typed URL namespace helpers

**Files:**
- Create: `admin/src/core/routing/url-namespaces.ts`
- Test: `admin/src/core/routing/url-namespaces.test.mjs`
- Modify: `admin/src/core/resource/manifest-resolver.ts`
- Modify: `admin/src/core/resource/manifest-resolver.test.mjs`

**Interfaces:**
- `adminRoute(resource: string): string` returns `/admin/<resource>`.
- `adminApiBase(resource: string): string` returns `/api/v1/admin/<resource>`.
- `isAdminRoute(path: string): boolean` accepts `/admin` and `/admin/...` only.
- `isAdminApiPath(path: string): boolean` accepts `/api/v1/admin/...` only.

- [ ] **Step 1: Write failing helper tests** for normal resource names, leading slashes, empty names, `/admin` validation, and rejection of `/users` as an admin route.
- [ ] **Step 2: Run the focused tests** with `node --test admin/src/core/routing/url-namespaces.test.mjs`; confirm they fail because the helper module is absent.
- [ ] **Step 3: Implement the helpers** with one normalization function that strips only leading/trailing slashes and rejects empty or path-traversal segments.
- [ ] **Step 4: Update manifest resolution** so a manifest with `name: "departments"` derives `/admin/departments` and `/api/v1/admin/departments` when explicit values are absent.
- [ ] **Step 5: Add manifest assertions** that explicit `adminRoute` and `apiBase` must match the two namespaces and cannot be swapped.
- [ ] **Step 6: Run the focused tests** and commit:

```powershell
git add admin/src/core/routing admin/src/core/resource/manifest-resolver.ts admin/src/core/resource/manifest-resolver.test.mjs
git commit -m "feat: define admin url namespace helpers"
```

### Task 2: Move Vue Router and authentication into `/admin`

**Files:**
- Modify: `admin/src/router/index.ts`
- Modify: `admin/src/modules/auth/pages/LoginPage.vue`
- Modify: `admin/src/core/layouts/AdminShell.vue`
- Modify: `admin/src/core/pages/AdminHomePage.vue`
- Modify: `admin/src/core/pages/ErrorPage.vue`
- Modify: `admin/src/core/pages/ForbiddenPage.vue`
- Modify: `admin/src/core/pages/LoadingPage.vue`
- Test: `admin/tests/admin-routing.test.ts`

**Interfaces:**
- Admin route names use the `admin-` prefix.
- Login route name is `admin-login` and path is `/admin/login`.
- Admin home route name is `admin-home` and path is `/admin`.
- `redirect` query values are restricted to internal `/admin/...` paths before use.

- [ ] **Step 1: Add failing route tests** for `/admin`, `/admin/login`, `/admin/users`, `/admin/users/new`, `/admin/users/:id`, and `/admin/audit-logs`; assert `/users` is not registered as an admin route.
- [ ] **Step 2: Run `pnpm exec vitest run admin/tests/admin-routing.test.ts`** and confirm the current root routes fail the new expectations.
- [ ] **Step 3: Add an `/admin` parent route** and move login, dashboard, RBAC, audit, settings, resource, error, forbidden, and loading children under it.
- [ ] **Step 4: Update navigation guards** to redirect unauthenticated admin pages to `/admin/login`, preserving only a validated `/admin/...` redirect.
- [ ] **Step 5: Update login success, logout, forbidden, error, and home links** to use named admin routes or the URL helpers; remove hard-coded root admin paths.
- [ ] **Step 6: Update the routing tests** for direct navigation and refresh-safe paths, then run them and commit:

```powershell
git add admin/src/router admin/src/modules/auth/pages/LoginPage.vue admin/src/core/pages admin/src/core/layouts/AdminShell.vue admin/tests/admin-routing.test.ts
git commit -m "feat: namespace admin pages under admin"
```

### Task 3: Update shared resource pages and generated Resource manifests

**Files:**
- Modify: `admin/src/core/resource/generated.ts`
- Modify: `admin/src/core/resource/pages/ResourceListPage.vue`
- Modify: `admin/src/core/resource/pages/ResourceFormPage.vue`
- Modify: `admin/src/core/resource/pages/ResourceDetailPage.vue`
- Modify: `admin/src/components/resource/ResourceFormView.vue`
- Modify: `admin/src/core/layouts/AdminShell.vue`
- Modify: `admin/src/core/notifications/NotificationMenu.vue`
- Modify: `admin/src/lib/resource-navigation.ts`
- Test: `admin/tests/resource-navigation.test.ts`

**Interfaces:**
- Resource metadata provides `adminRoute` and `apiBase`.
- `resourceListPath`, `resourceCreatePath`, `resourceEditPath`, and `resourceDetailPath` return paths beginning with `/admin/`.
- Notification internal links accept `/admin/...` paths and reject admin-looking root paths such as `/users`.

- [ ] **Step 1: Extend resource navigation tests** to expect `/admin/users`, `/admin/users/new`, `/admin/users/1`, and `/admin/users/1/edit`.
- [ ] **Step 2: Run the focused resource tests** and record the current root-path failures.
- [ ] **Step 3: Update navigation helpers** to use `adminRoute` instead of concatenating `/${resourceName}`.
- [ ] **Step 4: Replace direct router pushes in list, form, detail, shell, and notification components** with the helpers or named admin routes.
- [ ] **Step 5: Update generated resource definitions** for announcements and departments to contain `adminRoute` and `apiBase`.
- [ ] **Step 6: Run resource navigation, form, action, and type checks; commit:

```powershell
git add admin/src/core/resource admin/src/components/resource admin/src/core/layouts/AdminShell.vue admin/src/core/notifications admin/src/lib/resource-navigation.ts admin/src/modules admin/tests
git commit -m "feat: route shared resources under admin namespace"
```

### Task 4: Change generator output and golden contracts

**Files:**
- Modify: `backend/app/generator/spec.go`
- Modify: `backend/app/generator/frontend_render.go`
- Modify: `backend/app/generator/render.go`
- Modify: `backend/app/generator/runtime_registration.go`
- Modify: `backend/app/generator/spec_test.go`
- Modify: `backend/app/generator/frontend_render_test.go`
- Modify: `backend/app/generator/resource_pipeline_test.go`
- Modify: `backend/app/generator/testdata/posts/manifest.go.golden`
- Modify: `backend/app/generator/testdata/posts/README.md.golden`
- Modify: `backend/app/generator/testdata/posts/routes.go.golden`
- Modify: `backend/app/generator/testdata/posts/manifest_test.go.golden`

**Interfaces:**
- Generated `Spec.FrontendRoute` is `/admin/<resource>`.
- Generated frontend ResourceSpec contains `adminRoute` and `apiBase`.
- Generated backend resource routes remain mounted through `/api/v1/admin/<resource>`.

- [ ] **Step 1: Add failing generator assertions** for `adminRoute`, `apiBase`, generated frontend paths, and README examples.
- [ ] **Step 2: Run generator golden and pipeline tests** to confirm the current generator emits root routes.
- [ ] **Step 3: Update the spec normalization** so resource names produce `/admin/<name>` for the frontend and `/api/v1/admin/<name>` for the API.
- [ ] **Step 4: Update frontend and README renderers** to use the two explicit fields and remove ambiguous `route` output.
- [ ] **Step 5: Regenerate or update only the affected golden fixtures** and run `go test ./app/generator/...`.
- [ ] **Step 6: Run a real generator smoke command in a temporary review target**, verify no existing manual file is overwritten, and commit:

```powershell
git add backend/app/generator
git commit -m "feat: generate namespaced admin resource routes"
```

### Task 5: Remove backend page-path collisions and align API documentation

**Files:**
- Modify: `backend/routes/web.go`
- Modify: `backend/app/openapi/spec.go`
- Modify: `backend/app/openapi/spec_test.go`
- Modify: `backend/app/openapi/docs_test.go`
- Modify: `admin/src/generated/api.ts`
- Modify: `docs/openapi.md`
- Modify: `docs/resource-engine.md`
- Modify: `docs/generator.md`
- Modify: `docs/testing.md`

**Interfaces:**
- All admin resource HTTP handlers remain under `/api/v1/admin`.
- No backend admin handler is registered at `/users`, `/roles`, `/permissions`, or another frontend page path.
- Generated API Client paths match the OpenAPI document exactly.

- [ ] **Step 1: Add backend route assertions** that enumerate registered admin paths and reject root page-path handlers.
- [ ] **Step 2: Run the route and OpenAPI tests** to capture the legacy `/users` collision and any stale `/api/v1/users` examples.
- [ ] **Step 3: Remove the legacy backend root page route** and keep only API/OpenAPI/backend infrastructure routes in `backend/routes/web.go`.
- [ ] **Step 4: Update OpenAPI schemas and generated client comments/paths** to use `/api/v1/admin/...`.
- [ ] **Step 5: Update API, resource-engine, generator, and testing docs** with the new three-space contract.
- [ ] **Step 6: Run `go test ./app/openapi/... ./app/core/... ./app/http/...` and commit:

```powershell
git add backend/routes backend/app/openapi admin/src/generated/api.ts docs/openapi.md docs/resource-engine.md docs/generator.md docs/testing.md
git commit -m "fix: remove backend admin page route collisions"
```

### Task 6: Update Vite proxy, browser acceptance, and full verification

**Files:**
- Modify: `admin/vite.config.ts`
- Modify: `admin/tests/resource-navigation.test.ts`
- Modify: `admin/tests/resource-form.test.ts`
- Modify: `admin/tests/dashboard-resources.test.ts`
- Modify: `admin/tests/resource-actions.test.ts`
- Create: `admin/tests/admin-url-namespace.test.ts`
- Modify: `README.md`
- Modify: `README.zh-CN.md`
- Modify: `docs/ai-quickstart.md`

**Interfaces:**
- Vite proxies only `/api` to the Go API.
- `/admin` and root C-end paths are served by the Vue application.
- Browser acceptance uses port 5180 and real PostgreSQL/Redis services without changing their configuration.

- [ ] **Step 1: Add browser acceptance coverage** for `/admin/login`, `/admin`, `/admin/users`, resource form/detail paths, and one existing root C-end route.
- [ ] **Step 2: Run the focused frontend test suite** and fix stale root-path assertions.
- [ ] **Step 3: Inspect `vite.config.ts`** and assert the proxy target is scoped to `/api`, without adding an `/admin` proxy rule.
- [ ] **Step 4: Start the frontend in the required background mode**, read the actual port from its log, and verify Windows access with `wslnet url <port>` before browser testing.
- [ ] **Step 5: Execute the browser flow**: admin login, dashboard, users list, create/edit/detail, roles, permissions, audit logs, notifications, search, export, relation options, and batch action navigation.
- [ ] **Step 6: Run final checks and commit:

```powershell
cd backend
go test ./...

cd ..\admin
pnpm exec vue-tsc --noEmit
pnpm run build
pnpm exec vitest run
```

```powershell
git add admin/vite.config.ts admin/tests README.md README.zh-CN.md docs/ai-quickstart.md
git commit -m "test: verify admin url namespace end to end"
```

## Completion Gate

The implementation is complete only when all six tasks pass, no root admin alias is registered, generated resources use both namespaced values, the browser can refresh every `/admin/...` deep link, and the backend/API/OpenAPI/frontend paths are identical to the spec. Push only after this stage-level verification and review are complete.
