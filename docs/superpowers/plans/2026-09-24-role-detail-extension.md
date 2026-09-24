# Role Detail Extensions and Permission UX Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make generic resource details extensible, correctly handle missing roles, and provide scalable permission selection UX with verified role permission persistence.

**Architecture:** The generic detail page will resolve business extensions through a core-owned registry instead of importing role code. The detail API will distinguish a missing record from an empty record, while the role permission extension will own filtering and grouped selection behavior. Existing role permission endpoints and the `super-admin` read-only rule remain unchanged.

**Tech Stack:** Vue 3, TypeScript, vue-i18n, generated API client, Go/Goravel controllers and services, PostgreSQL-backed integration tests.

**Spec:** `docs/superpowers/specs/2026-09-24-role-detail-extension-design.md`

## Global Constraints

- Do not add frontend dependencies or database migrations.
- Keep admin URLs under `/admin/`.
- Keep shared UI in `admin/src/components`, infrastructure in `admin/src/core`, and business UI in `admin/src/modules/<name>`.
- Keep backend services under `backend/app/services/<domain>/`.
- Do not push GitHub during implementation; create local milestone commits only.

## Review Focus

- Missing resource IDs must not render or save a role permission extension: Task 1 adds controller and UI tests.
- Core must not import business modules: Task 2 adds a registry ownership test and source-level import check.
- Search and group selection must preserve checked permissions and scopes: Task 3 adds pure behavior tests.
- Empty permission sets and forbidden role writes must be persisted/rejected: Task 4 adds backend integration coverage.
- The detail link and permission controls must remain usable on narrow screens: Task 5 performs browser verification at a narrow viewport.

---

### Task 1: Make missing resource details explicit

**Files:**
- Create: `admin/src/core/resource/detail-state.ts`
- Create: `admin/src/core/resource/detail-state.test.mjs`
- Modify: `admin/src/core/resource/pages/ResourceDetailPage.vue`
- Test: browser verification of a missing role ID

**Interfaces:**
- The existing detail endpoint already returns `RESOURCE_NOT_FOUND` with HTTP 404.
- The frontend converts that API error into an explicit not-found state, disables edit/delete actions, and does not render extensions.

- [ ] **Step 1: Write the failing controller test**

Add a pure state test for the existing API error contract:

```go
// detail-state.test.mjs
assert.equal(isResourceNotFound(new ApiError('missing', 404, 'RESOURCE_NOT_FOUND')), true)
assert.equal(isResourceNotFound(new ApiError('forbidden', 403, 'RBAC_FORBIDDEN')), false)
```

- [ ] **Step 2: Run the focused test and verify it fails**

Run from `backend`:

```powershell
node --experimental-strip-types src/core/resource/detail-state.test.mjs
```

Expected: FAIL because the not-found state helper does not exist.

- [ ] **Step 3: Implement the not-found contract**

Keep the existing backend 404 contract. Add `isResourceNotFound` and set `notFound` when the detail request raises `RESOURCE_NOT_FOUND`; do not treat a missing record as a generic error.

- [ ] **Step 4: Add the frontend state and translations**

Track a `notFound` state in `ResourceDetailPage.vue`, skip edit/delete and extension rendering when set, and render:

```vue
<Empty v-if="notFound">
  <EmptyHeader><EmptyTitle>{{ t('resource.resourceNotFound') }}</EmptyTitle></EmptyHeader>
</Empty>
```

Reuse the existing localized `resource.resourceNotFound` key and ensure the detail page renders it explicitly.

- [ ] **Step 5: Run focused backend and frontend checks**

```powershell
node --experimental-strip-types src/core/resource/detail-state.test.mjs
npx vue-tsc -b --pretty false
```

Expected: PASS.

- [ ] **Step 6: Commit**

```powershell
git add admin/src/core/resource/detail-state.ts admin/src/core/resource/detail-state.test.mjs admin/src/core/resource/pages/ResourceDetailPage.vue
git commit -m "fix: handle missing resource details"
```

### Task 2: Add the core-owned detail extension registry

**Files:**
- Create: `admin/src/core/resource/detail-extensions.ts`
- Create: `admin/src/core/resource/detail-extensions.test.mjs`
- Create: `admin/src/core/resource/components/ResourceDetailExtensions.vue`
- Modify: `admin/src/core/resource/pages/ResourceDetailPage.vue`
- Modify: `admin/src/modules/roles/components/RolePermissionsPanel.vue`
- Modify: `admin/src/main.ts`
- Create: `admin/src/modules/roles/detail-extension.ts`

**Interfaces:**

```ts
export interface ResourceDetailExtensionProps {
  resource: string
  id: string
  record: Record<string, unknown>
}

export interface ResourceDetailExtension {
  resource: string
  component: Component
}

export function registerResourceDetailExtension(extension: ResourceDetailExtension): void
export function getResourceDetailExtension(resource: string): ResourceDetailExtension | undefined
```

- [ ] **Step 1: Write the failing registry test**

```js
test('registers and resolves one extension per resource', () => {
  const extension = { resource: 'roles', component: {} }
  registerResourceDetailExtension(extension)
  assert.equal(getResourceDetailExtension('roles'), extension)
  assert.throws(() => registerResourceDetailExtension(extension), /already registered/)
})
```

- [ ] **Step 2: Run the test and verify it fails**

```powershell
node --experimental-strip-types src/core/resource/detail-extensions.test.mjs
```

Expected: FAIL because the registry module does not exist.

- [ ] **Step 3: Implement the registry and renderer**

Use a module-local `Map<string, ResourceDetailExtension>`, reject duplicate keys, and render the resolved component with `resource`, `id`, and `record` props. The registry module must not import any file from `admin/src/modules`.

- [ ] **Step 4: Register the role extension at application startup**

Move role-specific registration out of `ResourceDetailPage.vue`. Register `RolePermissionsPanel` for `roles` from the application bootstrap or a modules registration entrypoint. Change the panel props to accept the generic extension contract and derive the current role assignments from `record.permissions`.

- [ ] **Step 5: Verify the core import boundary**

```powershell
rg -n "modules/roles" admin/src/core
```

Expected: no matches.

- [ ] **Step 6: Run registry, type, and build checks**

```powershell
node --experimental-strip-types src/core/resource/detail-extensions.test.mjs
npx vue-tsc -b --pretty false
npm run build
```

- [ ] **Step 7: Commit**

```powershell
git add admin/src/core/resource admin/src/modules/roles admin/src/main.ts
git commit -m "refactor: register resource detail extensions"
```

### Task 3: Add scalable permission filtering and group selection

**Files:**
- Modify: `admin/src/modules/roles/components/RolePermissionsPanel.vue`
- Modify: `admin/src/modules/roles/role-permissions.ts`
- Modify: `admin/src/modules/roles/role-permissions.test.mjs`
- Modify: `admin/src/locales/zh-CN/rbac.json`
- Modify: `admin/src/locales/en-US/rbac.json`

**Interfaces:**

```ts
export function filterPermissionGroups(
  groups: PermissionGroup[],
  query: string,
): PermissionGroup[]

export function togglePermissionGroup(
  selected: number[],
  visiblePermissionIds: number[],
  checked: boolean,
): number[]
```

- [ ] **Step 1: Write failing pure behavior tests**

```js
test('filters by permission name and display name', () => {
  assert.equal(filterPermissionGroups(groups, '用户')[0].permissions.length, 1)
  assert.equal(filterPermissionGroups(groups, 'admin.users.view')[0].permissions.length, 1)
})

test('selects and clears only visible permissions in one group', () => {
  assert.deepEqual(togglePermissionGroup([1, 9], [1, 2], true), [1, 9, 2])
  assert.deepEqual(togglePermissionGroup([1, 2, 9], [1, 2], false), [9])
})
```

- [ ] **Step 2: Run the tests and verify they fail**

```powershell
node src/modules/roles/role-permissions.test.mjs
```

Expected: FAIL because the pure helpers do not exist.

- [ ] **Step 3: Implement filtering and selection helpers**

Normalize query and searchable values with `toLocaleLowerCase()`, preserve original group order, deduplicate selected IDs, and never remove IDs outside the visible group.

- [ ] **Step 4: Wire the panel UI**

Add a localized search input above the groups. Each group gets a select-all checkbox whose state is checked only when every visible permission is selected. Hide empty groups after filtering. Keep all controls disabled for `super-admin`.

- [ ] **Step 5: Verify focused tests and build**

```powershell
node src/modules/roles/role-permissions.test.mjs
npx vue-tsc -b --pretty false
```

- [ ] **Step 6: Commit**

```powershell
git add admin/src/modules/roles admin/src/locales/zh-CN/rbac.json admin/src/locales/en-US/rbac.json
git commit -m "feat: improve role permission selection"
```

### Task 4: Verify persistence, revocation, and forbidden writes

**Files:**
- Modify: `backend/app/core/admin/controllers/rbac_controller.go`
- Modify: `backend/app/services/rbac/rbac_role_service.go` only if the tests expose a contract defect
- Create or modify: `backend/app/core/admin/controllers/rbac_controller_test.go`
- Modify: `backend/app/openapi/spec_test.go` if the API contract changes

**Interfaces:**
- `PUT /api/v1/admin/roles/{id}/permissions` accepts an empty `permission_ids` array and persists revocation.
- Requests without `admin.roles.manage` return the existing authorization error and do not mutate `permission_role`.

- [ ] **Step 1: Write the failing integration tests**

Cover these exact cases:

The tests must use the existing controller test helpers and assert these concrete outcomes:

```go
func TestReplaceRolePermissionsPersistsAndRevokesAssignments(t *testing.T) {
	// PUT one permission, query PermissionAssignments, PUT an empty list,
	// query again, and assert the final assignment count is zero.
}

func TestReplaceRolePermissionsRejectsUserWithoutRoleManagement(t *testing.T) {
	// Snapshot PermissionAssignments, issue the PUT as a user without
	// admin.roles.manage, assert http.StatusForbidden, then compare the snapshot.
}

func TestRoleDetailReturnsAssignmentsAfterReplacement(t *testing.T) {
	// Replace one permission, call the generic role detail endpoint, and assert
	// the response data.permissions contains the selected permission ID.
}
```

The comments describe the exact assertions and existing helpers to use; do not add a new test harness or change the permission schema.

- [ ] **Step 2: Run focused tests and verify any failure is contractual**

```powershell
go test ./app/core/admin/controllers ./app/services/rbac -run 'Test(ReplaceRolePermissions|RoleDetail)' -count=1
```

- [ ] **Step 3: Make the smallest backend correction required**

Do not change the permission schema. Preserve transaction boundaries and the built-in `super-admin` protection. Only adjust validation or response handling if a test demonstrates a defect.

- [ ] **Step 4: Run backend RBAC and OpenAPI tests**

```powershell
go test ./app/core/admin/controllers ./app/services/rbac ./app/openapi -count=1
```

- [ ] **Step 5: Commit**

```powershell
git add backend/app/core/admin/controllers backend/app/services/rbac backend/app/openapi
git commit -m "test: verify role permission persistence and authorization"
```

### Task 5: Make the edit link responsive and perform end-to-end verification

**Files:**
- Modify: `admin/src/modules/roles/pages/RoleFormPage.vue`
- Modify: `admin/src/locales/zh-CN/resource.json`
- Modify: `admin/src/locales/en-US/resource.json`
- Test: browser flow against the running local services

- [ ] **Step 1: Implement the responsive header**

Use a wrapping header layout so the detail action moves below the title on narrow screens:

```vue
<div class="flex flex-wrap items-start gap-3">
  <Button ... />
  <div class="min-w-0 flex-1">...</div>
  <Button v-if="editing" class="w-full sm:w-auto" ...>
    {{ t('resource.viewDetails') }}
  </Button>
</div>
```

- [ ] **Step 2: Run production checks**

```powershell
npx vue-tsc -b --pretty false
npm run build
```

- [ ] **Step 3: Run browser acceptance**

Verify at desktop and narrow viewport:

1. `/admin/roles/999999` shows “Resource not found” and no permission save button.
2. `/admin/roles/1` shows all permissions checked and disabled for `super-admin`.
3. A normal role detail page allows search, group select-all, save, reload, and revoke.
4. A user without `admin.roles.manage` receives a forbidden response and the assignment set remains unchanged.
5. `/admin/roles/1/edit` exposes the detail link and wraps it on a narrow viewport.

- [ ] **Step 4: Run the complete verification suite**

```powershell
cd ..\backend
$env:GOCACHE = "$pwd\.gocache-check"
go test ./...
cd ..\admin
npx vue-tsc -b --pretty false
npm run build
node src/core/resource/detail-extensions.test.mjs
node src/modules/roles/role-permissions.test.mjs
node src/modules/roles/role-routes.test.mjs
```

Remove the temporary `.gocache-check` directory after verification and confirm `git status --short` is clean.

- [ ] **Step 5: Commit the verified milestone**

```powershell
git add admin backend
git commit -m "feat: harden role detail permissions workflow"
```
