# Resource MySQL Provider Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将示例资源接入可配置的 MySQL 8 持久化 Provider，同时保留内存模式和现有 ResourceDataProvider 前端契约。

**Architecture:** HTTP Controller 负责认证、请求绑定、查询解析和错误映射；`DemoResourceService` 负责领域校验；`DemoResourceRepository` 由内存和 MySQL 两个实现满足。后端读取 `RESOURCE_PROVIDER`，前端读取 `VITE_RESOURCE_PROVIDER`，两者默认保持 memory，显式配置后才使用 HTTP/MySQL。

**Tech Stack:** Go 1.25、Goravel 1.18、Goravel MySQL driver、MySQL 8、Vue 3、TypeScript、Vite、Vitest、shadcn-vue。

**Spec:** `docs/superpowers/specs/2026-09-10-resource-mysql-provider-design.md`

## Global Constraints

- 资源服务层不得导入 HTTP、SQL 或 Goravel facade。
- MySQL Provider 不得把用户输入直接拼入 SQL；排序字段必须白名单化。
- `RESOURCE_PROVIDER` 未配置时使用 memory；未知值必须显式报错。
- 前端 `VITE_RESOURCE_PROVIDER` 未配置时保持现有内存演示数据。
- 所有资源 API 要求当前会话拥有 `dashboard.view`。
- 修改文件使用 `apply_patch`；每个任务完成后运行对应测试并独立提交。

---

### Task 1: Add the persistent resource schema and provider selection contract

**Files:**
- Create: `backend/database/migrations/20260910000002_create_demo_resources_table.go`
- Create: `backend/app/core/resource/provider.go`
- Create: `backend/app/core/resource/provider_test.go`
- Modify: `backend/bootstrap/migrations.go`
- Modify: `backend/config/resource.go`
- Test: `backend/app/core/resource/provider_test.go`

**Interfaces:**
- Produces `ResourceProviderMode`, `ParseResourceProviderMode(value string)`, and migration signature `20260910000002_create_demo_resources_table`.
- The provider mode accepts empty/`memory` as memory and `mysql` as MySQL; all other values return an error.

- [ ] **Step 1: Write the failing provider-mode tests**

```go
func TestParseResourceProviderMode(t *testing.T) {
    require.Equal(t, ResourceProviderMemory, ParseResourceProviderMode(""))
    require.Equal(t, ResourceProviderMemory, ParseResourceProviderMode("memory"))
    require.Equal(t, ResourceProviderMySQL, ParseResourceProviderMode("mysql"))
    require.Error(t, ParseResourceProviderMode("redis"))
}
```

- [ ] **Step 2: Run the focused test and verify it fails because the mode contract is absent**

Run: `go test ./app/core/resource -run TestParseResourceProviderMode -count=1`

Expected: FAIL with undefined `ResourceProviderMode` or `ParseResourceProviderMode`.

- [ ] **Step 3: Implement the mode parser and config entry**

Add `backend/config/resource.go` with `resource.provider` sourced from `RESOURCE_PROVIDER`, and make the parser trim and lowercase input before accepting `memory` or `mysql`.

- [ ] **Step 4: Add the migration**

Create `demo_resources` only when absent, using string primary key `id`, `name`, `status`, `owner`, `created_at`, and `updated_at`; add indexes on `status` and `owner`; use a string status column rather than MySQL ENUM.

- [ ] **Step 5: Register and verify the migration**

Register the migration in `bootstrap.Migrations`, run `gofmt`, and execute `go test ./app/core/resource ./bootstrap ./config`.

- [ ] **Step 6: Commit**

```bash
git add backend/database/migrations/20260910000002_create_demo_resources_table.go backend/app/core/resource/provider.go backend/app/core/resource/provider_test.go backend/bootstrap/migrations.go backend/config/resource.go
git commit -m "feat(resource): add mysql resource schema contract"
```

### Task 2: Implement memory and MySQL repositories behind the domain contract

**Files:**
- Create: `backend/app/core/resource/memory_repository.go`
- Create: `backend/app/core/resource/mysql_repository.go`
- Create: `backend/app/core/resource/mysql_repository_test.go`
- Modify: `backend/app/core/resource/demo_service.go`
- Test: `backend/app/core/resource/mysql_repository_test.go`

**Interfaces:**
- `NewMemoryDemoResourceRepository(initial ...DemoResource) DemoResourceRepository`.
- `NewMySQLDemoResourceRepository() DemoResourceRepository`.
- MySQL query wrapper methods use the existing Goravel `db.Query` but expose only the operations needed by this resource.

- [ ] **Step 1: Extend service tests for empty update, duplicate bulk IDs, and not-found mapping**

Add tests asserting empty `DemoResourceUpdate` returns `ErrInvalidResource`, duplicate IDs are passed once to the repository, and repository not-found errors remain distinguishable from generic storage errors.

- [ ] **Step 2: Run the focused service tests and verify the new cases fail**

Run: `go test ./app/core/resource -run 'TestDemoResourceService(RejectsEmptyUpdate|BulkDelete|NotFound)' -count=1`

Expected: FAIL because the new assertions and repository error behavior are not implemented.

- [ ] **Step 3: Implement the memory repository**

Use a mutex-protected map keyed by ID, clone records on input/output, apply search/filter/sort/page in memory, return `ErrResourceNotFound` for missing single records, and make bulk delete idempotent for existing IDs.

- [ ] **Step 4: Implement the MySQL query adapter**

Use `facades.DB().WithContext(ctx).Table("demo_resources")`; map only `id`, `name`, `status`, `owner`, `created_at`, and `updated_at`; use `Paginate` for lists; apply search through `WhereAny`; allow only `id`, `name`, `status`, and `owner` as sort fields; use `WhereIn` for bulk delete.

- [ ] **Step 5: Add repository boundary tests**

Use a small fake query interface owned by the resource package to assert that an unknown sort field resolves to `name`, descending sort calls `OrderByDesc`, pagination passes page and per-page values, and update never includes nil fields.

- [ ] **Step 6: Run tests and commit**

Run: `gofmt -w app/core/resource && go test ./app/core/resource -count=1`

```bash
git add backend/app/core/resource
git commit -m "feat(resource): implement memory and mysql repositories"
```

### Task 3: Add authenticated resource HTTP API and provider wiring

**Files:**
- Create: `backend/app/http/controllers/resource_controller.go`
- Create: `backend/app/http/controllers/resource_controller_test.go`
- Modify: `backend/routes/web.go`
- Modify: `backend/app/core/openapi/metadata.go`
- Test: `backend/app/http/controllers/resource_controller_test.go`

**Interfaces:**
- `NewResourceController(auth *AuthController, service *resource.DemoResourceService) *ResourceController`.
- Endpoints are `/api/resources/demo`, `/api/resources/demo/:id`, and `/api/resources/demo/bulk-delete`.
- Controller maps `ErrInvalidResource` to 400, `ErrResourceNotFound` to 404, and other repository errors to 503.

- [ ] **Step 1: Write HTTP contract tests**

Cover unauthenticated list (`401`), authenticated list envelope (`200` with `data` and `meta.pagination`), invalid create (`400`), and missing record (`404`) using a fake service repository.

- [ ] **Step 2: Run the controller tests and verify they fail**

Run: `go test ./app/http/controllers -run TestResourceController -count=1`

Expected: FAIL because the controller and routes do not exist.

- [ ] **Step 3: Implement controller authentication and permission checks**

Call `auth.CurrentUser(ctx)`, return `UNAUTHENTICATED` for missing sessions, use `permission.NewAuthorizer().Require(user.Permissions, "dashboard.view")`, and return the existing error envelope.

- [ ] **Step 4: Implement query and body boundaries**

Convert `ctx.Request().Queries()` to `url.Values`, reuse `query.Parse`, bind typed create/update/bulk DTOs, and pass `requestContext(ctx)` to the service. Do not bind arbitrary maps into domain objects.

- [ ] **Step 5: Wire the configured Provider**

Create a resource service from `RESOURCE_PROVIDER`; memory mode seeds `demo-1` and `demo-2`, MySQL mode constructs `NewMySQLDemoResourceRepository`. An invalid mode produces a controller configuration error rather than falling back.

- [ ] **Step 6: Register routes and OpenAPI metadata**

Register bulk-delete before `/:id`, add all six endpoint descriptions, and keep the existing auth and menu routes unchanged.

- [ ] **Step 7: Run controller and backend tests, then commit**

Run: `go test ./app/http/controllers ./tests/core ./tests/feature -count=1 && go vet ./...`

```bash
git add backend/app/http/controllers/resource_controller.go backend/app/http/controllers/resource_controller_test.go backend/routes/web.go backend/app/core/openapi/metadata.go
git commit -m "feat(resource): expose authenticated resource api"
```

### Task 4: Switch the frontend resource registry by explicit mode

**Files:**
- Modify: `admin/src/resource-engine/demo.ts`
- Modify: `admin/src/resource-engine/providers/HttpResourceDataProvider.ts`
- Create: `admin/src/resource-engine/provider-mode.ts`
- Create: `admin/src/resource-engine/provider-mode.test.ts`
- Modify: `admin/src/resource-engine/providers/HttpResourceDataProvider.test.ts`
- Test: `admin/src/resource-engine/provider-mode.test.ts`

**Interfaces:**
- `resolveResourceProviderMode(value?: string): 'memory' | 'http'`.
- `createDemoProvider(mode, apiClient)` returns a `ResourceDataProvider<DemoResourceRecord>`.

- [ ] **Step 1: Write provider-mode tests**

Assert empty/`memory` selects the existing seeded memory provider, `http` creates an `HttpResourceDataProvider` targeting `/api/resources/demo`, and unsupported values throw a clear error.

- [ ] **Step 2: Run the frontend focused test and verify it fails**

Run: `pnpm --dir admin test -- src/resource-engine/provider-mode.test.ts`

Expected: FAIL because provider mode resolution is absent.

- [ ] **Step 3: Implement explicit provider selection**

Read `import.meta.env.VITE_RESOURCE_PROVIDER`, default to `memory`, and keep seed data unchanged. Do not expose a silent HTTP-to-memory fallback.

- [ ] **Step 4: Preserve Provider contract behavior**

Ensure HTTP list maps `per_page` and `total_pages`, mutation bodies include JSON content type, and IDs remain URL encoded; add tests for the backend pagination envelope and bulk delete.

- [ ] **Step 5: Run frontend verification and commit**

Run: `pnpm --dir admin test && pnpm --dir admin typecheck && pnpm --dir admin build`

```bash
git add admin/src/resource-engine
git commit -m "feat(admin): select resource provider by environment"
```

### Task 5: End-to-end verification and GitHub push

**Files:**
- Modify: `.env.example` or the existing environment example files with `RESOURCE_PROVIDER=memory` and `VITE_RESOURCE_PROVIDER=memory` comments.
- Modify: `docs/superpowers/specs/2026-09-10-resource-mysql-provider-design.md` only if verification discovers a contract correction.

- [ ] **Step 1: Run all backend checks from the backend module directory**

Run: `go test ./... && go vet ./...`

- [ ] **Step 2: Run all frontend checks**

Run: `pnpm test && pnpm typecheck && pnpm build`

- [ ] **Step 3: Verify runtime memory mode**

With `RESOURCE_PROVIDER=memory`, verify `GET /api/resources/demo` returns `200`, the seeded records are listed, and the panel resource page remains available.

- [ ] **Step 4: Verify runtime MySQL mode without changing credentials**

With `RESOURCE_PROVIDER=mysql`, run migrations and request the list endpoint. If the configured MySQL credentials reject the connection, record the exact access error and leave the repository unmodified; do not change user permissions.

- [ ] **Step 5: Review the diff and push**

Run `git diff --check`, `git status --short`, and `git log --oneline -5`; confirm only intended files are changed.

```bash
git push origin HEAD
```
