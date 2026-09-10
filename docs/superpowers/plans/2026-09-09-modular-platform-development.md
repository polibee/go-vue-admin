# Goravel + Vue 3 Admin Platform Modular Development Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在 `D:\codex\go-vue-admin` 建立一个可持续演进的 Goravel + Vue 3 Admin Platform 单仓库，使每个模块都能独立开发、独立验收、独立提交，并在验收通过后自动推送到 GitHub。

**Architecture:** 采用 Monorepo，但以 `backend/`、`admin/`、`modules/`、`plugins/`、`contracts/` 为硬边界。Core 只承载平台能力；业务能力放在 `modules/<name>/`；前端 API 只通过 OpenAPI 生成客户端进入 Resource Engine。每个模块通过 manifest、注册入口、测试和验收清单形成闭环，完成后由脚本执行检查、创建提交、创建或确认 GitHub 远程并推送。

**Tech Stack:** Goravel/Go, Vue 3, TypeScript, Vite, Vue Router, Pinia, Vue I18n, Tailwind CSS, shadcn-vue, AI Elements Vue, TanStack Table, vee-validate, Zod, OpenAPI, Vitest, Vue Test Utils, Playwright, GitHub CLI.

**Spec:** `docs/开发文档.md`, `docs/设计.md`, `docs/设计优化.md`

## Global Constraints

- UI 基线只能使用官方 `unovue/shadcn-vue`；不得重新设计颜色、圆角、阴影、间距、字体、组件状态或视觉语言。
- `admin/src/components/ui/` 只放官方 shadcn-vue 源码；`admin/src/components/ai-elements/` 只放官方 AI Elements Vue 源码；两者不得包含业务依赖。
- 自定义扩展只允许放在 `admin/src/components/ui-extensions/`，且只有在官方组件及合理组合都不能满足时才能新增。
- Admin 组合组件放在 `admin/src/components/admin/`；业务逻辑放在 `modules/*`，不得反向污染 Core UI。
- Core 不包含 CMS、Trading、Forum、AI 等具体业务；这些必须作为 Module。
- 前端业务禁止直接使用 `fetch('/api/...')` 或 `axios`；必须使用 Generated API Client 和 `ResourceDataProvider`。
- 后端 API 统一使用 envelope、分页、过滤、排序、搜索和错误代码协议。
- 认证优先使用 Cookie/HttpOnly Session；前端权限仅用于 UX，后端 Middleware/Policy 才是最终授权。
- 所有模块必须具备独立测试、文档、迁移/注册入口和 Definition of Done；未通过门禁不得推送。
- GitHub 仓库默认按 `go-vue-admin` 创建为 private；如已有同名仓库则复用，不覆盖远程历史。
- 第一个远程推送前必须重新执行 `gh auth status`；当前失效令牌不能用于建仓或推送。

## Module Delivery Contract

每个可交付模块都必须包含以下最小契约：

```text
modules/<module>/
├── module.yaml                 # id, version, dependencies, permissions
├── backend/                    # Go module code, routes, migrations, tests
├── admin/                      # routes, resources, locales, tests
├── README.md                   # capability, dependency, verification
└── acceptance.md               # executable acceptance checklist
```

模块完成流程固定为：

```text
定义边界
  ↓
写契约/测试
  ↓
实现模块
  ↓
运行模块级检查
  ↓
运行跨模块回归
  ↓
生成 OpenAPI/SDK 并检查无差异
  ↓
git commit
  ↓
创建或确认 GitHub 仓库
  ↓
git push
```

提交消息采用 Conventional Commits，例如 `feat(auth): add session login`。一个模块一个功能提交；基础设施提交不得夹带无关模块或视觉改动。

---

### Task 1: 建立仓库基础与 GitHub 推送门禁

**Files:**
- Create: `backend/go.mod`
- Create: `admin/package.json`
- Create: `package.json`
- Create: `pnpm-workspace.yaml`
- Create: `Makefile`
- Create: `.editorconfig`
- Create: `.gitignore`
- Create: `.github/workflows/ci.yml`
- Create: `scripts/repo-init.sh`
- Create: `scripts/module-check.sh`
- Create: `scripts/module-finish.sh`
- Create: `docs/CONTRIBUTING.md`
- Create: `docs/MODULE_DEVELOPMENT.md`
- Create: `docs/ui/upstream-overrides.md`

**Interfaces:**
- `scripts/module-check.sh <module-path>` exits non-zero on missing contract, failing tests, lint/typecheck/build failure, generated SDK drift, or forbidden UI-layer imports.
- `scripts/module-finish.sh <module-path> <commit-message>` calls the checker, creates a commit, creates/reuses private `polibee/go-vue-admin` through `gh repo create` when no remote exists, adds `origin`, and pushes the current branch. It must refuse to overwrite a non-empty remote.
- CI runs `scripts/module-check.sh` for every changed module and the repository-wide checks on pull requests and pushes.

- [ ] **Step 1: Write the repository contract and module checklist.**

Document the exact directory boundaries, dependency direction, required module files, branch policy, commit policy, and the commands that must pass before a push in `docs/CONTRIBUTING.md` and `docs/MODULE_DEVELOPMENT.md`.

- [ ] **Step 2: Add the root workspace manifests.**

Define `backend`, `admin`, `modules`, `plugins`, and `contracts` as the initial workspaces without adding business dependencies. Add scripts for `lint`, `typecheck`, `test`, `build`, `openapi:generate`, and `module:check`.

- [ ] **Step 3: Implement the module checker.**

The checker must validate `module.yaml`, `README.md`, and `acceptance.md`, detect changed module paths with `git diff --name-only`, run the module's backend/admin tests when present, and fail if `components/ui` or `components/ai-elements` imports application code.

- [ ] **Step 4: Implement the finish-and-push command.**

Use `gh auth status` before any mutation. Use `gh repo view polibee/go-vue-admin` to detect an existing repository; if missing, run `gh repo create polibee/go-vue-admin --private --source . --remote origin --push`. If the remote exists, only add it when the local repository has no remote and then push the current branch. Never force-push and never reset user work.

- [ ] **Step 5: Add CI gates.**

CI must run backend tests, admin lint/typecheck/build/tests, contract generation plus `git diff --exit-code`, module checks, and Playwright smoke tests when the application exists. A generated API mismatch fails CI.

- [ ] **Step 6: Verify and commit.**

Run:

```bash
git init
git add .
git commit -m "chore(repo): establish modular platform foundation"
bash scripts/module-check.sh
```

Expected: local checks pass; GitHub push is reported as blocked until `gh auth status` succeeds.

---

### Task 2: Admin Component Foundation

**Files:**
- Create: `admin/components.json`
- Create: `admin/src/components/ui/README.md`
- Create: `admin/src/components/ai-elements/README.md`
- Create: `admin/src/components/ui-extensions/README.md`
- Create: `admin/component-manifest.json`
- Create: `admin/src/app/bootstrap.ts`
- Create: `admin/src/app/App.vue`
- Create: `admin/src/core/theme/`
- Create: `admin/src/core/i18n/`
- Modify: `admin/package.json`

**Interfaces:**
- `component-manifest.json` records the locked shadcn-vue/AI Elements versions, style, and registry snapshot.
- Bootstrap exposes light/dark/system theme and `zh-CN`/`en` locale without globally registering every component.

- [ ] **Step 1: Initialize shadcn-vue through the official CLI.**
- [ ] **Step 2: Install the locked complete shadcn-vue registry into `admin/src/components/ui/`.**
- [ ] **Step 3: Install the locked AI Elements Vue registry into `admin/src/components/ai-elements/`.**
- [ ] **Step 4: Add manifest and upstream override tracking.**
- [ ] **Step 5: Add smoke tests for imports, dark mode, and tree-local imports.**
- [ ] **Step 6: Run `pnpm lint`, `pnpm typecheck`, `pnpm test`, and `pnpm build`, then commit `feat(admin): add component foundation`.**

Expected: all source components are present but runtime imports remain local and on-demand; no second UI framework exists.

---

### Task 3: Admin Semantic Components and Shell

**Files:**
- Create: `admin/src/components/admin/`
- Create: `admin/src/core/router/`
- Create: `admin/src/core/layout/`
- Create: `admin/src/core/navigation/`
- Create: `admin/src/pages/errors/`
- Modify: `admin/src/app/App.vue`

**Interfaces:**
- `AdminPage`, `AdminPageHeader`, `AdminToolbar`, `AdminSearch`, `AdminFilter`, `AdminPagination`, `AdminEmpty`, `AdminLoading`, `AdminError`, `AdminConfirm` compose existing UI primitives.
- `NavigationItem` contains `id`, `parentId`, `label`, `icon`, `route`, `sort`, `permission`, and `enabled`.

- [ ] **Step 1: Add component tests that assert composition and states, not new visual tokens.**
- [ ] **Step 2: Implement semantic components with shadcn-vue composition only.**
- [ ] **Step 3: Implement router, layout, sidebar, header, breadcrumb, user-menu placeholders, and mobile sheet behavior.**
- [ ] **Step 4: Implement 404, 403, 500, loading, and network-error boundaries.**
- [ ] **Step 5: Verify responsive layout and theme switching in a real browser.**
- [ ] **Step 6: Commit `feat(admin): add semantic shell and error states`.**

Expected: `/admin` opens with light/dark/system, Chinese/English, collapsible sidebar, mobile navigation, and error routes.

---

### Task 4: Backend Core Contracts

**Files:**
- Create: `backend/app/core/shared/response/`
- Create: `backend/app/core/shared/errors/`
- Create: `backend/app/core/shared/query/`
- Create: `backend/app/core/shared/pagination/`
- Create: `backend/app/core/openapi/`
- Create: `backend/routes/health.go`
- Create: `backend/tests/core/`

**Interfaces:**
- `PaginationQuery`, `FilterQuery`, `SortQuery`, and `SearchQuery` are the only list-query inputs.
- Success response is `{data, meta}`; errors are `{error:{code,message,details}}`.

- [ ] **Step 1: Write unit tests for envelope, query parsing, and error mapping.**
- [ ] **Step 2: Implement response builders and exception middleware.**
- [ ] **Step 3: Implement health endpoint, logging correlation, and validation error mapping.**
- [ ] **Step 4: Add OpenAPI metadata for the core endpoints.**
- [ ] **Step 5: Run Go unit/feature tests and a contract smoke test.**
- [ ] **Step 6: Commit `feat(core): add API contracts and query protocol`.**

---

### Task 5: Authentication Module

**Files:**
- Create: `backend/app/core/auth/`
- Create: `admin/src/core/auth/`
- Create: `admin/src/pages/auth/LoginPage.vue`
- Create: `tests/e2e/auth.spec.ts`
- Modify: `backend/routes/`

**Interfaces:**
- Backend endpoints: `POST /login`, `POST /logout`, `GET /me`, session refresh/revoke.
- Frontend exports `AuthStore`, `AuthService`, `useAuth`, `AuthGuard`, and `GuestGuard`.

- [ ] **Step 1: Write failing backend tests for login, logout, me, session persistence, and invalid credentials.**
- [ ] **Step 2: Implement HttpOnly Cookie session and CSRF-safe mutation path.**
- [ ] **Step 3: Implement Pinia auth state and route guards.**
- [ ] **Step 4: Build the login page by composing existing shadcn-vue Form/Input/Button components.**
- [ ] **Step 5: Run backend tests, admin tests, and Playwright login/logout flow.**
- [ ] **Step 6: Commit `feat(auth): add session authentication`.**

---

### Task 6: RBAC and Dynamic Navigation Modules

**Files:**
- Create: `backend/app/core/user/`
- Create: `backend/app/core/role/`
- Create: `backend/app/core/permission/`
- Create: `backend/app/core/menu/`
- Create: `admin/src/core/permissions/`
- Create: `admin/src/core/navigation/NavigationRegistry.ts`
- Create: `tests/e2e/rbac.spec.ts`

**Interfaces:**
- Permissions use `resource.action`, e.g. `users.view`, `roles.update`.
- Frontend exposes `can(permission: string): boolean` and `PermissionGuard`; backend enforces Middleware/Policy.

- [ ] **Step 1: Write authorization tests for API, route, menu, and action visibility.**
- [ ] **Step 2: Implement User/Role/Permission relations and seed data.**
- [ ] **Step 3: Implement backend permission middleware and policies.**
- [ ] **Step 4: Implement menu API and frontend registry merge.**
- [ ] **Step 5: Verify a restricted user cannot see, navigate to, or call forbidden resources.**
- [ ] **Step 6: Commit `feat(rbac): add permissions and dynamic navigation`.**

---

### Task 7: Resource Engine Core

**Files:**
- Create: `admin/src/resource-engine/core/ResourceDefinition.ts`
- Create: `admin/src/resource-engine/core/ResourceRegistry.ts`
- Create: `admin/src/resource-engine/core/ResourceDataProvider.ts`
- Create: `admin/src/resource-engine/core/ResourceContext.ts`
- Create: `admin/src/resource-engine/core/ResourceRoute.ts`
- Create: `admin/src/resource-engine/tests/`

**Interfaces:**
- `defineResource<T>(definition: ResourceDefinition<T>): ResourceDefinition<T>`.
- `ResourceDataProvider` exposes `list`, `get`, `create`, `update`, `delete`, and `bulkDelete` with typed query/mutation results.

- [ ] **Step 1: Write registry and provider contract tests.**
- [ ] **Step 2: Implement resource definition normalization and duplicate-name rejection.**
- [ ] **Step 3: Implement registry lookup, context, route metadata, and permission metadata.**
- [ ] **Step 4: Add a memory provider for deterministic tests.**
- [ ] **Step 5: Run unit tests and typecheck.**
- [ ] **Step 6: Commit `feat(resource): add resource engine core`.**

---

### Task 8: Resource Table, Form, and Generic CRUD

**Files:**
- Create: `admin/src/resource-engine/table/`
- Create: `admin/src/resource-engine/form/`
- Create: `admin/src/resource-engine/pages/ResourceListPage.vue`
- Create: `admin/src/resource-engine/pages/ResourceCreatePage.vue`
- Create: `admin/src/resource-engine/pages/ResourceEditPage.vue`
- Create: `admin/src/resource-engine/pages/ResourceShowPage.vue`
- Create: `admin/src/resource-engine/fields/`
- Create: `admin/src/resource-engine/columns/`
- Create: `admin/src/resource-engine/router.ts`

**Interfaces:**
- Resource routes are `/admin/resources/:resource`, `/create`, `/:id`, and `/:id/edit`.
- Table uses TanStack Table plus shadcn-vue Table for server pagination, sort, search, filter, selection, and bulk actions.
- Form uses vee-validate + Zod + shadcn-vue Form and supports text, textarea, number, select, checkbox, switch, date, and datetime fields.

- [x] **Step 1: Write table/form tests for server state and validation errors.**
- [x] **Step 2: Implement table query serialization using the shared backend protocol.**
- [x] **Step 3: Implement field registry and Zod validation mapping.**
- [x] **Step 4: Implement generic pages and route registration.**
- [x] **Step 5: Add a DemoResource proving CRUD without custom CRUD pages.**
- [x] **Step 6: Run unit, component, browser, and build checks; commit `feat(resource): add generic CRUD engine`.**

---

### Task 9: Core User and RBAC Resources

**Files:**
- Create: `backend/app/core/user/resources/`
- Create: `backend/app/core/role/resources/`
- Create: `admin/src/core-resources/users/resource.ts`
- Create: `admin/src/core-resources/roles/resource.ts`
- Create: `admin/src/core-resources/permissions/resource.ts`
- Create: `tests/e2e/resources.spec.ts`

- [x] **Step 1: Define resource schemas and permissions from backend contracts.**
- [x] **Step 2: Register Users, Roles, and Permissions through ResourceRegistry.**
- [x] **Step 3: Implement relation/custom fields only where generic fields cannot express the relation.**
- [x] **Step 4: Verify User CRUD and RBAC behavior through generated routes.**
- [x] **Step 5: Commit `feat(core-resources): migrate users and roles to resource engine`.**

---

### Task 10: OpenAPI and Generated Data Provider

**Files:**
- Create: `contracts/openapi/openapi.json`
- Create: `contracts/schemas/`
- Create: `scripts/openapi-generate.sh`
- Create: `admin/src/generated/api/`
- Create: `admin/src/core/api/OpenApiDataProvider.ts`
- Modify: `.github/workflows/ci.yml`

**Interfaces:**
- `OpenApiDataProvider` adapts generated client calls to `ResourceDataProvider`.
- CI runs generation and `git diff --exit-code`; generated artifacts are committed and reviewable.

- [x] **Step 1: Write a contract test for the users list/create/update/delete operations.**
- [x] **Step 2: Generate `openapi.json` from backend annotations/metadata.**
- [x] **Step 3: Generate the TypeScript client into `admin/src/generated/api`.**
- [x] **Step 4: Implement OpenApiDataProvider without direct axios/fetch in resources.**
- [x] **Step 5: Replace DemoResource and core resources with the generated provider.**
- [x] **Step 6: Commit `feat(api): add OpenAPI generated client and provider`.**

---

### Task 11: Settings, Media, and Audit Modules

**Files:**
- Create: `backend/app/core/setting/`
- Create: `backend/app/core/media/`
- Create: `backend/app/core/audit/`
- Create: `admin/src/modules/settings/`
- Create: `admin/src/modules/media/`
- Create: `admin/src/modules/audit/`

- [x] **Step 1: Deliver Settings as a typed namespace/key/value resource without schema changes for new settings.**
- [x] **Step 2: Deliver Media with storage adapters, upload/delete/list/preview, and FileField/ImageField/MediaPicker.**
- [ ] **Step 3: Deliver Audit with actor/action/resource/before/after/ip/user-agent/timestamp and a Diff Viewer.**
- [ ] **Step 4: Add isolated tests and OpenAPI endpoints for each capability.**
- [ ] **Step 5: Commit each capability separately: `feat(settings)`, `feat(media)`, and `feat(audit)`.**

---

### Task 12: Module Runtime and Module Generator

**Files:**
- Create: `backend/app/core/module/Module.go`
- Create: `backend/app/core/module/Registry.go`
- Create: `admin/src/core/extensions/ModuleRegistry.ts`
- Create: `modules/example/module.yaml`
- Create: `modules/example/backend/`
- Create: `modules/example/admin/`
- Create: `scripts/make-module.sh`
- Create: `docs/modules/example/README.md`

**Interfaces:**
- Backend `Module` exposes `Name() string`, `Register(app Application) error`, and `Boot(app Application) error`.
- Frontend `defineAdminModule()` returns module id, resources, routes, navigation, and locales.
- `make module <name>` generates backend/admin/locales/routes/resources and acceptance files.

- [ ] **Step 1: Write registry tests proving registration order, duplicate rejection, and removal isolation.**
- [ ] **Step 2: Implement backend Module interface and registry bootstrap.**
- [ ] **Step 3: Implement frontend module definition and registry.**
- [ ] **Step 4: Add the removable `modules/example` test module.**
- [ ] **Step 5: Implement generator output and generator snapshot tests.**
- [ ] **Step 6: Remove Example Module in a test branch, run Core checks, then commit `feat(module): add module runtime and generator`.**

---

### Task 13: Builtin Plugin Registry and SDK Contract

**Files:**
- Create: `backend/app/core/plugin/`
- Create: `admin/src/core/extensions/plugin/`
- Create: `plugins/example/manifest.json`
- Create: `docs/plugins/PLUGIN_SDK.md`

- [ ] **Step 1: Define `PluginManifest`, `PluginState`, runtime values, permissions, menus, and dependencies.**
- [ ] **Step 2: Implement list/enable/disable for builtin plugins only.**
- [ ] **Step 3: Implement the frontend SDK registration surface without exposing Pinia/router/layout internals.**
- [ ] **Step 4: Add plugin enable/disable tests and UI compatibility checks.**
- [ ] **Step 5: Commit `feat(plugin): add builtin registry and SDK contract`.**

External process plugins, package installation, signature verification, sandboxing, marketplace, and dynamic JS loading remain v2/v3 work and are not part of v1.

---

### Task 14: Repository-Level Verification and v1 Foundation Release

**Files:**
- Modify: `.github/workflows/ci.yml`
- Create: `docs/releases/v1-foundation.md`
- Create: `scripts/release-check.sh`

- [ ] **Step 1: Run backend unit, feature, integration, permission, and API contract tests.**
- [ ] **Step 2: Run admin lint, typecheck, component tests, Resource Engine tests, and production build.**
- [ ] **Step 3: Run Playwright flows for login, users CRUD, RBAC, resource CRUD, module registration, and plugin enable/disable.**
- [ ] **Step 4: Run generated OpenAPI/SDK drift check and module contract check.**
- [ ] **Step 5: Verify no forbidden UI redesign or second UI framework was introduced.**
- [ ] **Step 6: Commit `chore(release): verify v1 foundation` and invoke `scripts/module-finish.sh` only after GitHub authentication is repaired.**

## Phase-to-Module Delivery Order

| Delivery | Modules | Exit result |
| --- | --- | --- |
| M0 | Repository Foundation + push gates | Reproducible monorepo and CI |
| M1 | shadcn-vue Foundation | Locked official UI source |
| M2 | AI Elements Foundation | Locked AI UI source, no AI business coupling |
| M3 | Admin Semantic + Shell | Usable `/admin`, theme, i18n, errors |
| M4 | Backend Core Contracts | Stable envelope/query/error/OpenAPI base |
| M5 | Authentication | Session login/logout/me |
| M6 | RBAC + Navigation | Enforced authorization and dynamic sidebar |
| M7 | Resource Engine Core | Declarative resources and provider contract |
| M8 | Resource CRUD | Generic table/form/routes |
| M9 | Core Resources | Users/Roles/Permissions migrated |
| M10 | OpenAPI Provider | Generated client is the only API path |
| M11 | Settings/Media/Audit | Platform foundation capabilities |
| M12 | Module Runtime/Generator | Business modules are removable and generatable |
| M13 | Builtin Plugin/SDK | v1 plugin contract without external runtime |
| M14 | Release Verification | v1 Foundation Complete |

Each row is a separate delivery boundary. A row is not pushed until its acceptance checklist and CI checks pass. M11 capabilities should remain separate commits even though they share the same milestone.

## Self-Review

- Spec coverage: Core boundaries, Admin Foundation, shadcn-vue/AI Elements constraints, auth, RBAC, navigation, Resource Engine, OpenAPI, settings, media, audit, modules, plugins, testing, CI, and v1 exclusions are mapped to Tasks 1-14.
- Placeholder scan: no implementation step is left as TBD; all later interfaces name concrete files, commands, or exported contracts.
- Type consistency: `ResourceDataProvider`, `defineResource`, `defineAdminModule`, backend `Module`, and `PluginManifest` are defined before their consumers.
- Scope guard: external runtime plugins, marketplace, sandbox, GraphQL, workflow engine, and schema builder are intentionally excluded from v1 as required by the source docs.
