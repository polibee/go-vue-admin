# User Status Model Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the boolean `is_active` user gate with a single `status` field supporting `active`, `disabled`, and `locked` across persistence, authentication, RBAC, resource APIs, and the admin UI.

**Architecture:** Migrate existing PostgreSQL rows before dropping `is_active`, then make `models.User.Status` the only backend source of truth. The service and auth layers validate the same finite set of statuses; the Resource Registry exposes a select field and the Vue list/form consume that contract for editing and filtering.

**Tech Stack:** Goravel v1.18.x, Go, PostgreSQL, Vue 3, TypeScript, Vue Router, vue-i18n, shadcn-vue, Node native tests.

**Spec:** `docs/superpowers/specs/2026-09-21-user-status-design.md`

## Global Constraints

- `status` is the only user status source; do not retain `is_active`.
- Valid statuses are exactly `active`, `disabled`, and `locked`.
- Only `active` users may log in or count as available RBAC administrators.
- Do not implement invitation registration, automatic lockout, status history, audit history, or bulk status changes in this phase.
- Preserve the existing password generation, show/hide, and copy interactions.
- Use the existing shadcn-vue component system and semantic tokens.
- Run PostgreSQL migration steps in the order: add status, copy old values, drop `is_active`.
- Keep commits local; push GitHub only after the complete milestone passes review and verification.

## Review Focus

- Existing rows with `is_active = false` must become `disabled`, not silently become active. Test in Task 1.
- Invalid status strings must be rejected at the service/request boundary. Test in Task 2.
- The last active administrator must not be disabled or locked. Test in Task 2.
- Disabled and locked users must both be denied login without exposing a status-specific authentication detail. Test in Task 2.
- The UI must not submit or render `is_active`, and a blank status must default to `active` only on creation. Test in Task 4.

### Task 1: Migrate the users table and model

**Files:**
- Create: `backend/database/migrations/20260921000001_replace_user_active_with_status.go`
- Modify: `backend/app/models/user.go`
- Modify: `backend/database/migrations/20260920000001_create_users_table.go`
- Test: `backend/app/models/user_test.go`

**Interfaces:**
- Consumes: Existing `users.is_active` rows and the Goravel schema/ORM APIs used by current migrations.
- Produces: `models.User.Status string`, `User.Public()` with `status` and no `is_active`, and a migration that preserves existing state before dropping the old column.

- [ ] **Step 1: Write the failing model/public-contract tests**

  Add a test that constructs a user with `Status: "disabled"` and asserts `Public()` contains `status: "disabled"` and does not contain `is_active`. Add a compile-level assertion that the model exposes `Status`.

- [ ] **Step 2: Run the focused test and verify it fails**

  Run `go test ./app/models -run TestUserPublicStatus -count=1` from `backend/`.

  Expected failure: `User` has no `Status` field and/or `Public()` still returns `is_active`.

- [ ] **Step 3: Implement the model and fresh-schema change**

  Replace `IsActive bool` with `Status string` in `models.User`, return `status` from `Public()`, remove `is_active` from the fresh users migration, and set the fresh-column default to `active`.

- [ ] **Step 4: Add the data migration**

  Create `20260921000001_replace_user_active_with_status.go`. In `Up`, add a nullable/intermediate `status` column if required by the schema driver, update `status` from `is_active` (`true` → `active`, `false` → `disabled`), make the column non-null/defaulted where supported, then drop `is_active`. In `Down`, recreate `is_active` and derive it from `status = active` before dropping `status`.

- [ ] **Step 5: Run focused and backend tests**

  Run `go test ./app/models ./database/migrations -count=1` and then `go test ./...` from `backend/`.

  Expected: PASS; no model output contains `is_active`.

- [ ] **Step 6: Commit the persistence boundary**

  Run `git add backend/app/models backend/database/migrations && git commit -m "feat: migrate users to status field"`.

### Task 2: Update service, authentication, and RBAC rules

**Files:**
- Modify: `backend/app/services/user_service.go`
- Modify: `backend/app/services/rbac_service.go`
- Modify: `backend/app/services/user_role_service.go`
- Modify: `backend/app/http/controllers/auth_controller.go`
- Modify: `backend/app/http/controllers/rbac_controller.go`
- Modify: `backend/routes/web.go`
- Test: `backend/app/services/user_service_test.go`
- Test: `backend/app/services/rbac_service_test.go`
- Test: `backend/tests/feature/auth_test.go`

**Interfaces:**
- Consumes: `models.User.Status` from Task 1.
- Produces: User create/update payloads with `status`, `ValidateUserStatus`, active-only authentication, and active-only administrator counting.

- [ ] **Step 1: Write failing status service tests**

  Add table cases for `active`, `disabled`, `locked`, and `pending`; assert only the first three are accepted. Add a create case with an empty status that resolves to `active`, and update cases that reject an invalid status.

- [ ] **Step 2: Run service tests and verify failure**

  Run `go test ./app/services -run 'Test(User|RBAC).*Status|Test.*Active' -count=1` from `backend/`.

  Expected failure: status validator and status-aware service signatures do not exist.

- [ ] **Step 3: Implement service contracts**

  Add a shared status validation helper. Change `Create` and `Update` to receive status, default blank create status to `active`, reject invalid values, and update all last-admin and user-role availability queries from `is_active = true` to `status = active`.

- [ ] **Step 4: Add authentication regression tests**

  Add feature cases for an active user logging in successfully and disabled/locked users receiving the same account-unavailable response.

- [ ] **Step 5: Update controllers and routes**

  Replace `userPayload.Active *bool` with `Status string`, pass status to the service, update user list/detail serialization, and keep existing permission middleware unchanged.

- [ ] **Step 6: Run the complete backend suite**

  Run `go test ./...` from `backend/`.

  Expected: PASS, including existing last-admin, role assignment, and authentication tests.

- [ ] **Step 7: Commit the backend behavior**

  Run `git add backend/app backend/routes backend/tests && git commit -m "feat: enforce user status in auth and rbac"`.

### Task 3: Expose status through the Resource Registry and list API

**Files:**
- Modify: `backend/app/resources/registry.go`
- Modify: `backend/app/http/controllers/resource_list_controller.go`
- Modify: `backend/app/http/controllers/resource_detail_controller.go`
- Test: `backend/app/resources/registry_test.go`
- Test: `backend/tests/feature/resource_test.go`

**Interfaces:**
- Consumes: The valid status set and `models.User.Status` from Tasks 1–2.
- Produces: A `status` select field with three options, status column metadata, and a list query filter such as `/api/v1/admin/resources/users?status=disabled`.

- [ ] **Step 1: Write failing registry and filter tests**

  Assert the users manifest contains `status` as a select field with `active`, `disabled`, and `locked` options and no `is_active`. Add a feature test that requests `status=disabled` and verifies only disabled records are returned.

- [ ] **Step 2: Run focused tests and verify failure**

  Run `go test ./app/resources ./tests/feature -run 'Resource|Manifest' -count=1` from `backend/`.

  Expected failure: manifest has boolean `is_active` and the list controller ignores the status query.

- [ ] **Step 3: Implement manifest and filter**

  Extend the resource field contract with select options if it does not already support them. Register `status` as the user field/column and apply a validated status predicate in the list controller.

- [ ] **Step 4: Run backend verification**

  Run `go test ./...` from `backend/`.

- [ ] **Step 5: Commit the resource contract**

  Run `git add backend/app/resources backend/app/http/controllers backend/tests && git commit -m "feat: expose user status resource contract"`.

### Task 4: Replace the Active switch with status selection and filtering

**Files:**
- Modify: `admin/src/views/UserFormView.vue`
- Modify: `admin/src/views/ResourceListView.vue`
- Modify: `admin/src/views/RBACView.vue`
- Modify: `admin/src/lib/resource-form.ts`
- Modify: `admin/src/locales/zh-CN/resource.json`
- Modify: `admin/src/locales/en-US/resource.json`
- Modify: `admin/src/locales/zh-CN/rbac.json`
- Modify: `admin/src/locales/en-US/rbac.json`
- Test: `admin/tests/resource-form.test.ts`
- Test: `admin/tests/user-status.test.ts`

**Interfaces:**
- Consumes: The manifest field/options and `status` API contract from Task 3.
- Produces: A status Select in create/edit forms, status Badge and filter in user lists, and no UI reference to `is_active`.

- [ ] **Step 1: Write failing frontend contract tests**

  Assert status form initialization defaults a new user to `active`, serializes `disabled` and `locked`, and excludes `is_active` from rendered resource fields. Add a status-label mapping test for Chinese and English keys.

- [ ] **Step 2: Run the focused tests and verify failure**

  Run `node --experimental-strip-types --test tests/resource-form.test.ts tests/user-status.test.ts` from `admin/`.

  Expected failure: the existing form still renders a boolean `is_active` field and has no status options.

- [ ] **Step 3: Implement the status form**

  Render the manifest select options with the existing shadcn-vue Select components, default creation to `active`, preserve the current password controls, and serialize only `status` for user state.

- [ ] **Step 4: Implement list status display and filtering**

  Replace the Active column with semantic Badge labels and add a status Select filter that reloads page one with the validated `status` query parameter. Update the RBAC user table to read `status`.

- [ ] **Step 5: Run frontend checks**

  Run `node --experimental-strip-types --test tests/*.test.ts`, `.\node_modules\.bin\vue-tsc.cmd -b`, and `npm run build` from `admin/`.

- [ ] **Step 6: Verify the browser flow**

  In the existing right-side preview, open `/users/new`, verify the default active status and no locale/is_active field, edit a user to disabled, return to `/users`, and verify the status filter and Badge. Do not create or delete real data unless explicitly requested.

- [ ] **Step 7: Commit the frontend milestone**

  Run `git add admin/src admin/tests && git commit -m "feat: manage user status in admin"`.

### Task 5: Final integration review and milestone handoff

**Files:**
- Modify: `docs/roadmap.md`
- Review: `docs/superpowers/specs/2026-09-21-user-status-design.md`
- Review: all commits from Tasks 1–4

- [ ] **Step 1: Update the roadmap**

  Mark user status management as complete under Phase 2 and note that automated lockout and status history remain outside this phase.

- [ ] **Step 2: Run the full verification set**

  Run `go test ./...` from `backend/`; run frontend tests, `vue-tsc -b`, and `npm run build` from `admin/`; run `git diff --check` from the repository root.

- [ ] **Step 3: Review the final diff**

  Confirm no production code, API response, migration, or locale references `is_active`; confirm the migration copies old values before dropping the column; confirm no GitHub push was performed.

- [ ] **Step 4: Commit the milestone documentation**

  Run `git add docs/roadmap.md && git commit -m "docs: complete user status milestone"`.

- [ ] **Step 5: Report the milestone**

  Report the migration, status semantics, test/build evidence, local commit range, and the fact that GitHub remains unpushed until the next approved milestone.
