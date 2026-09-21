# Data Scope Permissions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a reusable `all`/`own` data-scope contract to the Admin Resource engine and enforce it across resource reads, writes, search, export, and future batch Actions.

**Architecture:** Store the scope selected for a role-permission assignment in `permission_role.scope`, defaulting existing assignments to `all`. Resolve the most permissive scope for the authenticated user and resource permission, then apply the resulting constraint through a single resource-scope service before every database read or mutation. A manifest declares the owner field required by `own`; resources without that field remain `all` until explicitly configured.

**Tech Stack:** Go 1.27, Goravel v1.18, PostgreSQL, Vue 3, TypeScript, shadcn-vue, existing Resource Registry and generated OpenAPI client.

**Spec:** `docs/superpowers/specs/2026-09-22-resource-capabilities-design.md`

## Global Constraints

- Do not add dependencies or change lockfiles.
- PostgreSQL and Redis remain Laragon-managed services; do not change their configuration.
- Migrations are generated and manually reviewed, registered, and executed; code must not execute migrations automatically.
- Backend authorization and data filtering are authoritative; frontend filtering is presentation only.
- `admin:make-resource` remains the only public resource generator entry point.
- Keep Service files under `backend/app/services/<domain>/`; do not place new services in `app/core` or resource modules.
- Keep shared frontend primitives in `admin/src/components/`, infrastructure in `admin/src/core/`, and business pages in `admin/src/modules/<name>/`.
- Keep feature-phase commits local and push only after the complete milestone is independently reviewable.

## Review Focus

- A user with `own` scope must not read, search, export, update, or delete another user's record; test ownership filtering and every mutation path in Task 3.
- A user with multiple roles must receive the least restrictive effective scope (`all` before `own`); test role aggregation in Task 2.
- Existing rows in `permission_role` must remain usable after the migration and resolve to `all`; test migration defaults in Task 1.
- A manifest that requests `own` without a valid owner field must be rejected instead of silently returning all data; test registry validation in Task 1.
- Direct requests with arbitrary owner IDs must not create records outside the authenticated user's scope; test generated create behavior in Task 3.

---

### Task 1: Add the data-scope contract and persistence

**Files:**
- Modify: `backend/app/core/resource/registry.go`
- Modify: `backend/app/modules/announcements/resource/manifest.go`
- Modify: `backend/app/modules/admin/registry/registry.go`
- Create: `backend/app/services/rbac/data_scope_service.go`
- Create: `backend/app/services/rbac/data_scope_service_test.go`
- Create: `backend/database/migrations/20260922000001_add_permission_role_scope.go`
- Modify: `backend/bootstrap/migrations.go`
- Modify: `backend/app/core/resource/registry_test.go`

**Interfaces:**
- `resource.DataScope` is a string enum with values `all` and `own`.
- `resource.Manifest` gains `DataScope resource.DataScope` and `OwnerField string` metadata; an empty `DataScope` normalizes to `all`.
- `rbacservices.DataScopeService.EffectiveScope(userID int64, permission string) (resource.DataScope, error)` returns `all` if any matching role assignment is `all`, otherwise `own` if a matching assignment exists, otherwise an authorization error.
- The migration adds `permission_role.scope` as a non-null string with default `all`, and its `Down` removes the column without changing role or permission rows.

- [ ] **Step 1: Write failing registry tests**

Add tests that register a manifest with `DataScope: "own"` and `OwnerField: "owner_id"`, accept it, reject an unknown scope, and reject `own` without `OwnerField`.

- [ ] **Step 2: Run the focused test and verify it fails**

Run:

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-scope'); go test ./app/core/resource -run 'Test.*Scope' -count=1; Remove-Item -LiteralPath '.tmp-gocache-scope' -Recurse -Force -ErrorAction SilentlyContinue
```

Expected: FAIL because the scope type, manifest metadata, and validation do not exist.

- [ ] **Step 3: Implement the resource scope metadata**

Add the enum, normalization, and registry validation. Keep `DataScope` empty-compatible by treating it as `all`; only `own` requires a non-empty owner field that is present in the manifest fields.

- [ ] **Step 4: Add the migration and RBAC resolver**

Create the migration with `scope` defaulting to `all`. Implement the resolver with a parameterized query joining `role_user`, `permission_role`, and `permissions`. Keep `super-admin` compatible with the existing protected-role behavior by returning `all` for its permission assignments.

- [ ] **Step 5: Add resolver tests**

Cover no identity, missing assignment, one `own` assignment, and mixed `own` plus `all` assignments. Use the existing RBAC test database conventions and do not mutate a developer database from unit tests.

- [ ] **Step 6: Register the migration and update the generated/sample manifest**

Register the migration in `backend/bootstrap/migrations.go`. Add explicit `DataScope: resource.DataScopeAll` to built-in manifests and keep `announcements` on `all` until an owner field is intentionally added.

- [ ] **Step 7: Run focused verification**

Run:

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-scope'); go test ./app/core/resource ./app/services/rbac ./app/modules/admin/registry -count=1; Remove-Item -LiteralPath '.tmp-gocache-scope' -Recurse -Force -ErrorAction SilentlyContinue
```

Expected: PASS.

- [ ] **Step 8: Commit the contract boundary**

```powershell
git add backend/app/core/resource backend/app/services/rbac backend/app/modules/admin/registry backend/database/migrations backend/bootstrap
git commit -m "feat: add resource data scope contract"
```

### Task 2: Expose scope assignment through RBAC APIs and the role page

**Files:**
- Modify: `backend/app/services/rbac/rbac_role_service.go`
- Modify: `backend/app/core/admin/controllers/rbac_controller.go`
- Modify: `backend/routes/web.go`
- Modify: `backend/app/openapi/spec.go`
- Modify: `backend/app/openapi/spec_test.go`
- Modify: `admin/src/generated/api.ts`
- Modify: `admin/src/modules/rbac/pages/RBACPage.vue`
- Modify: `admin/src/locales/en-US/rbac.json`
- Modify: `admin/src/locales/zh-CN/rbac.json`
- Test: `backend/app/core/admin/controllers/resource_crud_adapter_test.go`
- Test: `admin/tests/permissions.test.ts`

**Interfaces:**
- Extend `PUT /api/v1/admin/roles/{id}/permissions` to accept `scopes`, a map from permission ID string to `all` or `own`, while preserving the existing `permission_ids` input.
- Omitted scope entries resolve to `all`.
- The role permissions response includes each permission's effective assigned scope.
- System roles cannot have their scope assignments changed through the existing protected-role API.

- [ ] **Step 1: Add a failing API contract test**

Submit a role permission replacement with one `own` scope and assert the response exposes that scope; submit an invalid scope and assert `422 VALIDATION_ERROR`.

- [ ] **Step 2: Run the focused backend test and verify it fails**

Run:

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-scope'); go test ./app/core/admin/controllers -run 'Test.*Scope|Test.*Permission' -count=1; Remove-Item -LiteralPath '.tmp-gocache-scope' -Recurse -Force -ErrorAction SilentlyContinue
```

Expected: FAIL because the payload and persistence do not yet accept scopes.

- [ ] **Step 3: Implement scoped role-permission replacement**

Validate every supplied scope, write the permission-role rows in one transaction, and default missing entries to `all`. Keep role deletion, protected system-role rules, and last-admin rules unchanged.

- [ ] **Step 4: Update OpenAPI and generated client**

Add the request property, response shape, enum values, and stable validation error to `spec.go`, update its test, and regenerate `admin/src/generated/api.ts` using the existing repository script without adding dependencies.

- [ ] **Step 5: Update the RBAC page**

Add a compact scope selector beside each permission assignment. Keep the existing permission checkbox behavior; the selector is enabled only when the permission is selected and defaults to `all`. Use existing shadcn-vue components and both locale files.

- [ ] **Step 6: Run frontend and contract verification**

Run:

```powershell
npx vue-tsc -b
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-scope'); go test ./app/openapi ./app/core/admin/controllers -count=1; Remove-Item -LiteralPath '.tmp-gocache-scope' -Recurse -Force -ErrorAction SilentlyContinue
```

Expected: PASS.

- [ ] **Step 7: Commit the RBAC surface**

```powershell
git add backend/app/services/rbac backend/app/core/admin/controllers backend/routes backend/app/openapi admin/src/generated admin/src/modules/rbac admin/src/locales
git commit -m "feat: manage resource data scopes in rbac"
```

### Task 3: Enforce effective scope across Resource reads and writes

**Files:**
- Create: `backend/app/services/rbac/resource_scope_service.go`
- Create: `backend/app/services/rbac/resource_scope_service_test.go`
- Modify: `backend/app/core/admin/controllers/resource_list_controller.go`
- Modify: `backend/app/core/admin/controllers/resource_detail_controller.go`
- Modify: `backend/app/core/admin/controllers/resource_crud_controller.go`
- Modify: `backend/app/core/admin/controllers/resource_export_controller.go`
- Modify: `backend/app/core/admin/controllers/resource_crud_adapter_test.go`
- Modify: `backend/app/core/admin/controllers/resource_list_capabilities_test.go`
- Modify: `backend/app/core/admin/controllers/resource_export_test.go`

**Interfaces:**
- `ResourceScopeService.Apply(ctx http.Context, query orm.Query, manifest resource.Manifest, action string) (orm.Query, error)` applies `all` or `own` to a query.
- For `own`, the query adds a parameterized `manifest.OwnerField = authenticatedUserID` predicate.
- `ResourceScopeService.CanAccess(ctx, manifest, action, id) (bool, error)` provides a single-row guard for update and delete paths that cannot reuse a loaded query.
- Built-in manifests keep `all` behavior; generated resources opt into `own` only when their manifest declares `OwnerField`.

- [ ] **Step 1: Write failing scope tests**

Cover generic list, search, pagination, show, update, delete, and export with two records owned by different users. Assert that an `own` user sees and changes only their record, while `all` sees both.

- [ ] **Step 2: Run the focused tests and verify they fail**

Run:

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-scope'); go test ./app/core/admin/controllers -run 'Test.*Scope|Test.*Export|Test.*Resource' -count=1; Remove-Item -LiteralPath '.tmp-gocache-scope' -Recurse -Force -ErrorAction SilentlyContinue
```

Expected: FAIL because current controllers query by resource and ID without applying ownership.

- [ ] **Step 3: Implement the scope service**

Resolve the authenticated user ID through the existing auth context, resolve the resource permission scope through RBAC, reject `own` manifests without a valid owner field, and use parameter binding for the owner value. Never interpolate user input into SQL identifiers.

- [ ] **Step 4: Apply scope before filters and pagination**

Apply the scope predicate before search, filters, sorting, and pagination in the generic list path. Apply the same predicate before the built-in resource paths are returned, while preserving their existing domain service behavior.

- [ ] **Step 5: Guard detail, update, delete, and export**

Apply `CanAccess` to detail and mutation paths. For export, apply the scope before search and CSV rendering so excluded rows cannot be inferred from totals or output order.

- [ ] **Step 6: Protect generic create ownership**

When an `own` resource is created, ignore a client-supplied owner value and set the owner field from the authenticated user. For `all` resources preserve the current generated create behavior.

- [ ] **Step 7: Run backend verification**

Run:

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-scope'); go test ./app/core/admin/controllers ./app/services/rbac ./app/http/middleware -count=1; Remove-Item -LiteralPath '.tmp-gocache-scope' -Recurse -Force -ErrorAction SilentlyContinue
```

Expected: PASS.

- [ ] **Step 8: Commit enforcement**

```powershell
git add backend/app/services/rbac backend/app/core/admin/controllers backend/app/http/middleware
git commit -m "feat: enforce resource data scopes"
```

### Task 4: Add generator metadata, OpenAPI coverage, and a real acceptance resource

**Files:**
- Modify: `backend/app/generator/spec.go`
- Modify: `backend/app/generator/parse.go`
- Modify: `backend/app/generator/render.go`
- Modify: `backend/app/generator/frontend_render.go`
- Modify: `backend/app/generator/generate_test.go`
- Modify: `backend/app/generator/render_test.go`
- Modify: `backend/app/generator/testdata/posts/manifest.go.golden`
- Modify: `backend/app/generator/testdata/posts/README.md.golden`
- Modify: `backend/app/openapi/spec.go`
- Modify: `backend/app/openapi/spec_test.go`
- Modify: `admin/src/core/resource/pages/ResourceListPage.vue`
- Modify: `admin/src/core/resource/pages/ResourceDetailPage.vue`
- Modify: `admin/src/core/resource/pages/ResourceFormPage.vue`
- Modify: `docs/generator.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- The generator accepts an explicit owner-field option only when the field is declared in the same ResourceSpec.
- Generated Manifest includes `DataScope` and `OwnerField` only when requested; default output remains `all` with no owner field.
- Generated README explains that assigning `own` requires a role scope assignment and that migration execution remains manual.

- [ ] **Step 1: Add failing Golden File coverage**

Add a generator fixture with `owner_id:integer` and an `own` scope option. Assert the manifest, migration, README, and frontend metadata contain the owner field and scope, while the existing default fixture remains `all`.

- [ ] **Step 2: Run generator tests and verify they fail**

Run:

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-scope'); go test ./app/generator -count=1; Remove-Item -LiteralPath '.tmp-gocache-scope' -Recurse -Force -ErrorAction SilentlyContinue
```

Expected: FAIL because parser and renderers do not understand owner-field metadata.

- [ ] **Step 3: Implement parser and renderer changes**

Validate the owner field as a declared integer field, reject inconsistent input before writing any artifact, and render the same metadata to backend, frontend, OpenAPI, tests, and README.

- [ ] **Step 4: Update generic pages**

Use the returned manifest metadata to avoid showing an owner selector for `own` resources and display the current user's ownership behavior in the generated form and detail pages. Do not move authorization logic into Vue.

- [ ] **Step 5: Run the real generator smoke**

Generate a temporary `owned-notices` resource with an owner field, run `admin:check-module`, verify the generated paths and discovery files, and verify no migration was executed automatically.

- [ ] **Step 6: Run full verification**

Run:

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-scope'); go test ./... -count=1; go build -o storage/codex-backend.exe .; Remove-Item -LiteralPath '.tmp-gocache-scope' -Recurse -Force -ErrorAction SilentlyContinue
npx vue-tsc -b
npm run build
```

Expected: all commands pass and the generated resource can be accepted through list, search, detail, export, create, update, and delete under both `all` and `own` role scopes.

- [ ] **Step 7: Update stage documentation and commit**

Update the generator README, resource-engine documentation, and roadmap with the actual scope assignment workflow. Then commit:

```powershell
git add backend/app/generator backend/app/openapi admin/src/core/resource docs/generator.md docs/resource-engine.md docs/roadmap.md
git commit -m "feat: complete resource data scope milestone"
```

### Final milestone checks

- [ ] Confirm `git status --short` is clean.
- [ ] Confirm all temporary Go cache directories are removed.
- [ ] Confirm PostgreSQL and Redis were only inspected or used through existing project configuration.
- [ ] Confirm no GitHub push occurs before the complete milestone is reviewed.
- [ ] Record the final commit hash and verification commands in the handoff.



