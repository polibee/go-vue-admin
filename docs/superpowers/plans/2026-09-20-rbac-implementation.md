# RBAC 与错误码本地化 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 Goravel v1.18 + Vue 3 Admin 中完成稳定错误码本地化、角色管理、角色权限分配、用户角色绑定和后端细粒度权限校验。

**Architecture:** 保留现有 JWT 认证和 Goravel ORM，在 Controller 与 ORM 之间新增 RBAC Service；用全局权限中间件保护写接口，用户权限由角色权限并集计算。前端通过 `ApiError.code` 映射 `errors.*` 语言键，RBAC 页面只组合现有官方 shadcn-vue 组件。

**Tech Stack:** Goravel v1.18、Goravel ORM、PostgreSQL、JWT、Vue 3、TypeScript、Pinia、Vue Router、vue-i18n、shadcn-vue 官方组件、Node 24 内置测试。

**Spec:** `docs/superpowers/specs/2026-09-20-rbac-design.md`

## Global Constraints

- 后端继续使用 Goravel v1.18 和现有 ORM，不引入第二套权限库或 ORM。
- 前端继续使用 `zh-CN/en-US` 目录语言包和官方 shadcn-vue 组件，不修改 `admin/src/components/ui` 源码。
- 前端不直接展示后端 `message`，只根据稳定 `code` 映射本地化文本。
- 所有 RBAC 写接口必须同时经过 JWT 身份认证和权限检查。
- 用户有效权限是其全部角色权限的并集；`super-admin` 只按精确系统角色名处理。
- 角色权限和用户角色更新必须使用完整集合提交并在事务中同步关系表。
- 必须保护系统角色和最后一个有效管理员账号。

## Review Focus

- 未认证、无权限和未知错误码是否分别返回正确 HTTP 状态并显示当前语言文本；测试归属 Task 1 和 Task 3。
- 恶意或未允许的角色、权限、用户 ID 是否不会被越权修改；测试归属 Task 3 和 Task 4。
- 空权限集合、重复 ID、无效 ID 和重复角色名是否返回可本地化错误；测试归属 Task 2 和 Task 3。
- 删除系统角色、移除最后管理员角色和停用最后管理员是否被拒绝；测试归属 Task 3 和 Task 4。
- 浏览器从 Vite 端口调用 Goravel API 时，CORS 预检和实际请求是否都能通过；测试归属 Task 5。

---

### Task 1: 统一错误码与前端错误映射

**Files:**
- Create: `admin/src/locales/zh-CN/errors.json`
- Create: `admin/src/locales/en-US/errors.json`
- Create: `admin/tests/error-code.test.ts`
- Modify: `admin/src/i18n/index.ts`
- Modify: `admin/src/lib/api.ts`
- Modify: `admin/src/views/LoginView.vue`
- Modify: `admin/src/views/RBACView.vue`

**Interfaces:**
- Consumes: 现有 `ApiError { status, code }`、`useI18n()`。
- Produces: `apiFetch` 在 HTTP 错误时保留稳定 `code`；页面通过 `t('errors.' + code)` 获取展示文本；未知 code 使用 `errors.unknown`。

- [ ] **Step 1: Write the failing test**

```ts
import assert from 'node:assert/strict'
import test from 'node:test'
import { errorMessageKey } from '../src/lib/api.ts'

test('maps a backend error code to a locale key and falls back for unknown codes', () => {
  assert.equal(errorMessageKey('AUTH_INVALID_CREDENTIALS'), 'errors.AUTH_INVALID_CREDENTIALS')
  assert.equal(errorMessageKey('NOT_DEFINED'), 'errors.unknown')
})
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `node --experimental-strip-types --test tests/error-code.test.ts`

Expected: FAIL because `errorMessageKey` does not exist.

- [ ] **Step 3: Write the minimal implementation**

Add `errorMessageKey(code?: string)` to `admin/src/lib/api.ts`, returning `errors.${code}` only for the allowlisted backend codes and `errors.unknown` otherwise. Add all codes from the spec to both locale files. Update Login and RBAC catches to use `t(errorMessageKey(error.code))` instead of hard-coded or raw backend text.

- [ ] **Step 4: Run the test and frontend checks**

Run: `node --experimental-strip-types --test tests/error-code.test.ts`

Expected: PASS.

Run: `pnpm exec vue-tsc -b` and `pnpm exec vite build`

Expected: exit code 0 with no TypeScript errors.

- [ ] **Step 5: Commit the isolated change**

```text
git add admin/src/lib/api.ts admin/src/i18n/index.ts admin/src/locales admin/src/views/LoginView.vue admin/src/views/RBACView.vue admin/tests/error-code.test.ts
git commit -m "fix: localize API error codes"
```

### Task 2: RBAC Service 与授权中间件

**Files:**
- Create: `backend/app/services/rbac_service.go`
- Create: `backend/app/http/middleware/permission.go`
- Create: `backend/app/services/rbac_service_test.go`
- Modify: `backend/app/http/controllers/auth_controller.go`
- Modify: `backend/app/http/controllers/rbac_controller.go`
- Modify: `backend/bootstrap/app.go`

**Interfaces:**
- Consumes: `parseAuthToken(ctx)`, Goravel ORM、`users`/`roles`/`permissions`/pivot 表。
- Produces:
  - `NewRBACService() *RBACService`
  - `RBACService.UserHasPermission(ctx http.Context, permission string) (bool, error)`
  - `RBACService.IsLastActiveAdmin(userID int64) (bool, error)`
  - `RequirePermission(permission string) http.Middleware`
  - 403 响应 `{ "code": "RBAC_FORBIDDEN" }`

- [ ] **Step 1: Write failing service tests**

Cover these exact cases with test fixtures or a temporary PostgreSQL test database: a user with a role permission returns true; a user with two roles receives the union; a user without the permission returns false; malformed/missing identity returns an error.

- [ ] **Step 2: Run the focused test to verify failure**

Run: `go test ./app/services ./app/http/middleware`

Expected: FAIL because `RBACService` and `RequirePermission` do not exist.

- [ ] **Step 3: Implement the service**

Resolve the authenticated user from the Goravel Auth context, query `role_user` and `permission_role`, check exact permission names, and handle the exact `super-admin` role name. Keep SQL and relationship logic in the service; controllers must not duplicate it.

- [ ] **Step 4: Implement the middleware**

The middleware must call the existing JWT parser first, call `UserHasPermission`, return `AUTH_UNAUTHORIZED` for invalid identity, `RBAC_FORBIDDEN` for insufficient permission, and call `ctx.Request().Next()` only on success.

- [ ] **Step 5: Run focused and package tests**

Run: `go test ./app/services ./app/http/middleware ./app/http/controllers ./bootstrap`

Expected: PASS.

- [ ] **Step 6: Commit the authorization boundary**

```text
git add backend/app/services backend/app/http/middleware backend/app/http/controllers backend/bootstrap/app.go
git commit -m "feat: add RBAC authorization service"
```

### Task 3: Role CRUD 与角色权限同步

**Files:**
- Create: `backend/app/requests/rbac_requests.go`
- Create: `backend/app/services/rbac_role_service.go`
- Modify: `backend/app/http/controllers/rbac_controller.go`
- Modify: `backend/routes/web.go`
- Modify: `backend/app/models/role.go`
- Create: `backend/app/services/rbac_role_service_test.go`

**Interfaces:**
- Consumes: `RBACService`, `RequirePermission("admin.roles.manage")`、Goravel ORM transactions。
- Produces:
  - `POST /api/v1/admin/roles`
  - `GET /api/v1/admin/roles/:id`
  - `PUT /api/v1/admin/roles/:id`
  - `DELETE /api/v1/admin/roles/:id`
  - `PUT /api/v1/admin/roles/:id/permissions`

- [ ] **Step 1: Write failing service/controller tests**

Test valid role creation, duplicate `name` returning `VALIDATION_ERROR`, missing role returning `RBAC_ROLE_NOT_FOUND`, system role deletion/update protection, and permission replacement in one transaction.

- [ ] **Step 2: Run focused tests to verify failure**

Run: `go test ./app/services ./app/http/controllers`

Expected: FAIL because the service methods and routes do not exist.

- [ ] **Step 3: Implement request validation and service methods**

Validate non-empty `name`, `display_name`, and integer permission IDs. Load all requested permissions before changing pivots; reject any missing ID. Replace `permission_role` rows inside one transaction. Reject deletion or key mutation of `super-admin`.

- [ ] **Step 4: Add protected routes**

Register routes in `backend/routes/web.go`; every write route must use `RequirePermission("admin.roles.manage")`, while role/permission reads use the corresponding view permission.

- [ ] **Step 5: Run tests and an API smoke check**

Run: `go test ./app/... ./bootstrap ./database/... ./routes`

Then authenticate with the local admin and verify role create, permission replacement, duplicate rejection, and system-role protection using PowerShell `Invoke-RestMethod`.

- [ ] **Step 6: Commit role management**

```text
git add backend/app/requests backend/app/services backend/app/http/controllers/rbac_controller.go backend/routes/web.go backend/app/models/role.go
git commit -m "feat: add protected role management"
```

### Task 4: User Role Binding 与管理员保护

**Files:**
- Create: `backend/app/services/user_role_service.go`
- Modify: `backend/app/http/controllers/rbac_controller.go`
- Modify: `backend/routes/web.go`
- Create: `backend/app/services/user_role_service_test.go`
- Modify: `admin/src/views/RBACView.vue`
- Modify: `admin/src/locales/zh-CN/rbac.json`
- Modify: `admin/src/locales/en-US/rbac.json`

**Interfaces:**
- Consumes: `RBACService`, `RequirePermission("admin.roles.manage")`, user and role IDs.
- Produces:
  - `GET /api/v1/admin/users/:id/roles`
  - `PUT /api/v1/admin/users/:id/roles`
  - frontend user-role assignment dialog.

- [ ] **Step 1: Write failing protection tests**

Test user not found, role not found, complete role replacement, removing the last active administrator, and removing the current operator's final management role. Assert exact codes and status codes.

- [ ] **Step 2: Run focused tests to verify failure**

Run: `go test ./app/services ./app/http/controllers`

Expected: FAIL because the binding service and routes do not exist.

- [ ] **Step 3: Implement transactional role binding**

Validate the complete role ID list, load the target user, calculate whether the resulting state leaves at least one active administrator, delete existing `role_user` rows, insert the new set, and commit atomically. Return `RBAC_LAST_ADMIN` for protected removals.

- [ ] **Step 4: Add the user assignment UI**

Use official shadcn-vue `Dialog`, `Checkbox`, `Button`, `Alert`, `Skeleton`, and `Empty`. Load available roles and current user roles, submit the full set, refresh the RBAC view, and render localized `ApiError.code` messages.

- [ ] **Step 5: Run frontend and backend checks**

Run: `node --experimental-strip-types --test tests/error-code.test.ts`, `pnpm exec vue-tsc -b`, `pnpm exec vite build`, and `go test ./app/... ./bootstrap ./database/... ./routes`.

- [ ] **Step 6: Commit user-role binding**

```text
git add backend/app/services backend/app/http/controllers/rbac_controller.go backend/routes/web.go admin/src/views/RBACView.vue admin/src/locales
git commit -m "feat: add user role assignment"
```

### Task 5: 浏览器联调、错误状态与文档收口

**Files:**
- Modify: `docs/authentication.md`
- Modify: `docs/security.md`
- Modify: `docs/roadmap.md`
- Test: `admin/tests/error-code.test.ts`, `backend/app/http/middleware/cors_test.go`

- [ ] **Step 1: Rebuild and restart both services**

Run backend migration/seed verification, build the backend executable, start it on `127.0.0.1:3000`, start Vite preview on an available port, and confirm both return HTTP 200.

- [ ] **Step 2: Verify browser flows**

In the visible browser, verify login, Chinese/English switch, RBAC loading state, empty/error state, role creation, permission assignment, user role binding, and a forbidden request from a user without the required permission.

- [ ] **Step 3: Verify CORS and API error language behavior**

Send an `OPTIONS` request with `Origin: http://127.0.0.1:4174` and expect `204`; trigger a known `RBAC_FORBIDDEN` response and verify the page displays the selected locale rather than the backend message.

- [ ] **Step 4: Update documentation**

Mark completed RBAC capabilities, record error-code rules, document protected endpoints, and keep refresh-token hardening clearly marked as pending.

- [ ] **Step 5: Run the complete verification set**

Run:

```text
node --experimental-strip-types --test tests/*.test.ts
pnpm exec vue-tsc -b
pnpm exec vite build
go test ./app/... ./bootstrap ./config ./database/... ./routes
```

Expected: all commands exit with code 0; browser smoke checks show no mixed-language error text.

- [ ] **Step 6: Commit documentation and verification changes**

```text
git add docs admin/tests backend/app/http/middleware/cors_test.go
git commit -m "docs: finish RBAC implementation notes"
```
