# Unified Resource URL and CRUD Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove `/resources/` from all public resource URLs and route users, roles, and generated resources through one Resource CRUD contract while preserving domain safety rules.

**Architecture:** The ResourceController becomes the single CRUD HTTP boundary for `/api/v1/admin/{resource}` and `/api/v1/admin/{resource}/{id}`. It resolves standard resources through the Registry, delegates users/roles mutations to their existing domain Services, and uses validated generic table operations for generated resources. The frontend generated Client is the only URL owner; metadata moves to `/api/v1/admin/registry`.

**Tech Stack:** Go 1.27, Goravel, GORM-style ORM, Vue 3, TypeScript, Vue Router, OpenAPI 3.0.3, generated Client templates, Go tests, `vue-tsc`.

**Spec:** `docs/superpowers/specs/2026-09-21-unified-resource-url-design.md`

## Global Constraints

- No user-facing page link or API contract uses `/resources/`.
- No compatibility aliases are retained for old `/resources/` or users/roles CRUD routes.
- Domain protections remain: association cleanup, `super-admin` protection, and last active administrator protection.
- Relationship endpoints remain explicit and are not converted into generic CRUD.
- Do not add dependencies, execute migrations, or push GitHub.
- Keep feature work in independently reviewable local commits.

## Review Focus

- A request to `/api/v1/admin/users/{id}` must use `UserService.Delete` and preserve role cleanup and last-admin errors; covered by Task 2 controller tests.
- A request to `/api/v1/admin/roles/{id}` must use `RoleService.Delete` and reject `super-admin`; covered by Task 2 controller tests.
- `/api/v1/admin/users/status` and role/permission relationship routes must not be captured by the generic `/{resource}/{id}` route; covered by Task 1 route contract tests.
- Generated pages must contain no `/resources/` URL after regeneration; covered by Task 3 generator and Client tests.
- OpenAPI, central Client, and runtime routes must expose the same paths; covered by Task 4 contract tests and smoke checks.

---

### Task 1: Replace the public resource route contract

**Files:**
- Modify: `backend/routes/web.go`
- Modify: `backend/app/openapi/spec.go`
- Modify: `backend/app/openapi/spec_test.go`
- Modify: `backend/scripts/contract-smoke.ps1`
- Test: `backend/routes/web_contract_test.go`

**Interfaces:**
- Consumes: `ResourceController.List`, `Show`, `Create`, `Update`, `Delete` and `RequireResourcePermission`.
- Produces: `GET/POST /api/v1/admin/{resource}`, `GET/PUT/DELETE /api/v1/admin/{resource}/{id}`, and `GET /api/v1/admin/registry`.

- [ ] **Step 1: Write failing route and OpenAPI assertions**

Assert that the new paths exist, old `/admin/resources` paths do not exist, the registry path exists, and relationship routes remain present.

- [ ] **Step 2: Run the focused contract tests and verify the expected failure**

Run:

```text
go test ./app/openapi ./routes
```

Expected: FAIL because the current route and OpenAPI paths still contain `/resources`.

- [ ] **Step 3: Change route registration and OpenAPI paths**

Register relationship-specific routes before the generic resource routes, then replace the generic path prefix with `/api/v1/admin/` and rename the metadata endpoint to `/api/v1/admin/registry`. Update OpenAPI operation paths and the read-only smoke script.

- [ ] **Step 4: Run focused tests and verify they pass**

Run:

```text
go test ./app/openapi ./routes
```

Expected: PASS, with no `/resources/` paths in the contract.

- [ ] **Step 5: Commit**

```text
git add backend/routes backend/app/openapi backend/scripts/contract-smoke.ps1
git commit -m "refactor: unify resource route contract"
```

### Task 2: Route users and roles CRUD through ResourceController adapters

**Files:**
- Modify: `backend/app/core/admin/controllers/resource_crud_controller.go`
- Modify: `backend/app/core/admin/controllers/resource_list_controller.go`
- Modify: `backend/app/core/admin/controllers/resource_detail_controller.go`
- Modify: `backend/app/modules/admin/registry/registry.go`
- Modify: `backend/app/core/admin/controllers/rbac_controller.go`
- Modify: `backend/app/core/admin/controllers/resource_crud_validation_test.go`
- Create or modify: `backend/app/core/admin/controllers/resource_crud_controller_test.go`

**Interfaces:**
- Consumes: `UserService.Create/Update/Delete`, `RoleService.Create/Update/Delete`, `userServiceError`, `roleServiceError`, and Resource Manifest actions.
- Produces: one CRUD controller path that preserves domain-specific Service behavior behind the common route.

- [ ] **Step 1: Write failing adapter tests**

Add controller-level tests for the adapter selection and error mapping. The tests must prove users and roles do not use raw table deletion and that domain errors map to `409`/`404` as before. Add manifest assertions for standard `create`, `update`, and `delete` actions using the manage permission.

- [ ] **Step 2: Run focused tests and verify failure**

Run:

```text
go test ./app/core/admin/controllers ./app/modules/admin/registry
```

Expected: FAIL because the generic controller currently performs raw table mutations and built-in manifests do not expose all standard actions.

- [ ] **Step 3: Implement the adapters**

Dispatch by `manifest.Name` inside ResourceController. Bind the existing user and role payload shapes, call domain Services, reuse existing error mapping, and keep generic validated table mutation for generated resources. Make list/show use the same route resource name without changing their domain-specific response shaping.

- [ ] **Step 4: Remove CRUD-specific controller route dependencies**

Keep relation operations in `RBACController`, but remove its standalone CRUD delete dependency from route registration. Ensure the generic controller is the only delete boundary.

- [ ] **Step 5: Run focused tests and verify pass**

Run:

```text
go test ./app/core/admin/controllers ./app/modules/admin/registry ./app/services/users ./app/services/rbac
```

Expected: PASS, including last-admin and `super-admin` protections.

- [ ] **Step 6: Commit**

```text
git add backend/app/core/admin/controllers backend/app/modules/admin/registry
git commit -m "feat: adapt builtin resources to unified crud"
```

### Task 3: Make the central TypeScript Client the only public URL owner

**Files:**
- Modify: `admin/scripts/generate-api-client.mjs`
- Modify: `admin/src/generated/api.ts`
- Modify: `admin/src/core/resource/pages/ResourceListPage.vue`
- Modify: `admin/src/core/resource/pages/ResourceDetailPage.vue`
- Modify: `admin/src/lib/resource-actions.ts`
- Modify: `admin/src/modules/users/pages/UserFormPage.vue`
- Modify: `admin/src/modules/roles/pages/RoleFormPage.vue`
- Modify: `admin/src/modules/rbac/pages/RBACPage.vue`
- Test: `backend/app/generator/frontend_render_test.go`

**Interfaces:**
- Consumes: unified `generatedApi.resourceList/resourceShow/resourceCreate/resourceUpdate/resourceDelete` methods.
- Produces: frontend calls using `/api/v1/admin/{resource}` only; no special users/roles delete path.

- [ ] **Step 1: Write failing generated-client assertions**

Assert generated client strings use `/api/v1/admin/` without `/resources/`, and generated resource module APIs delegate to `generatedApi`.

- [ ] **Step 2: Run the focused generator test and verify failure**

Run:

```text
go test ./app/generator -run TestRenderFrontendArtifactsFromResourceSpec
```

Expected: FAIL because generated and central Client paths still use `/resources/`.

- [ ] **Step 3: Update central Client and generated wrappers**

Change all generic resource methods to `/api/v1/admin/` and make metadata call `/api/v1/admin/registry`. Remove users/roles branching from delete URL construction; retain only policy checks such as disabling deletion of `super-admin`.

- [ ] **Step 4: Update specialized existing pages**

Change user and role forms and the RBAC role delete action to call central resource Client methods. Keep role permission assignment and user-role binding on their explicit relation endpoints.

- [ ] **Step 5: Run frontend verification**

Run:

```text
npx vue-tsc -b
rg -n "api/v1/admin/resources|resourceActionPath\(" admin/src admin/scripts
```

Expected: type check passes and the search returns no public CRUD URL usage.

- [ ] **Step 6: Commit**

```text
git add admin/src admin/scripts backend/app/generator
git commit -m "refactor: remove resources segment from admin client"
```

### Task 4: Align OpenAPI, generator documentation, and smoke acceptance

**Files:**
- Modify: `backend/app/openapi/spec_test.go`
- Modify: `admin/scripts/generate-api-client.mjs`
- Modify: `docs/openapi.md`
- Modify: `docs/generator.md`
- Modify: `docs/resource-engine.md`
- Modify: `docs/roadmap.md`
- Modify: `backend/app/generator/generate_test.go`
- Modify: `backend/app/generator/module_check_test.go`

**Interfaces:**
- Consumes: final runtime route and central Client paths from Tasks 1–3.
- Produces: documented, repeatable acceptance checks for the unified contract.

- [ ] **Step 1: Add no-legacy-path assertions**

Assert OpenAPI paths, generated Client output, generator output, and discovery output contain no `/resources/` URL.

- [ ] **Step 2: Run contract and generator tests to verify any stale references fail**

Run:

```text
go test ./app/openapi ./app/generator
```

- [ ] **Step 3: Update documentation and read-only smoke checks**

Document the new `/api/v1/admin/{resource}` contract, `/api/v1/admin/registry`, preserved relationship routes, and domain Service protections. Keep migration execution out of all checks.

- [ ] **Step 4: Run the complete verification suite**

Run:

```text
go test ./...
npx vue-tsc -b
git diff --check
```

Expected: all commands pass and the worktree has no untracked generated artifacts.

- [ ] **Step 5: Commit**

```text
git add backend admin docs
git commit -m "test: verify unified resource url architecture"
```

### Task 5: Final local review and handoff

**Files:**
- Review: `docs/superpowers/specs/2026-09-21-unified-resource-url-design.md`
- Review: all commits from Tasks 1–4

- [ ] **Step 1: Inspect the final diff and search for legacy paths**

Run:

```text
git diff d54716b..HEAD --stat
rg -n "api/v1/admin/resources|/resources/|DELETE.*admin/users|DELETE.*admin/roles" backend admin docs
```

Expected: no legacy CRUD route or public Client URL remains; documentation references only the intentional historical migration statement.

- [ ] **Step 2: Run final tests again**

Run the complete verification suite from Task 4.

- [ ] **Step 3: Report the local commits and remaining external action**

Report that no migration was executed and no GitHub push was performed. If the user later wants the stage pushed, handle that as a separate explicit action.
