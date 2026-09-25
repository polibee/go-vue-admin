# Admin/App API Boundary Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make the Admin/App boundary explicit without exposing a half-generated App Resource or renaming existing permissions.

**Architecture:** Keep the current Resource generator Admin-only. Extract reusable namespace validation for permission and menu renderers, record the future App contract, and add tests that reject an unsupported App Resource surface before any files are written.

**Tech Stack:** Go, Goravel, table-driven Go tests, Markdown documentation.

**Spec:** `docs/superpowers/specs/2026-09-25-admin-app-api-boundary-design.md`

## Global Constraints

- Existing `admin.*` permission names remain unchanged.
- `admin:make-resource` continues to generate only `/admin/...` pages and `/api/v1/admin/...` APIs.
- Do not add dependencies or migrations.
- Do not create an App API, App registry, or second user table in this phase.
- Keep feature commits local; push only after the stage is independently verified.

## Review Focus

- Empty, malformed, and mixed-case namespaces must be rejected.
- Existing callers that omit a namespace must keep producing `admin.*` permissions and `/admin/...` routes.
- A requested App Resource must fail before writing artifacts.
- Menu and permission renderers must use the same namespace rule.
- Documentation must not claim that App APIs already exist.

---

### Task 1: Namespace contract primitives

**Files:**
- Create: `backend/app/generator/namespace.go`
- Test: `backend/app/generator/namespace_test.go`
- Modify: `backend/app/generator/spec.go`

**Interfaces:**
- `normalizeNamespace(value string) (string, error)` returns `admin` for empty input and accepts only `admin` or `app`.
- `permissionName(namespace, resource, action string) string` returns `<namespace>.<resource>.<action>`.
- `Input.Namespace` is optional and defaults to `admin`.

- [ ] Write failing tests for default admin, accepted app primitive, and invalid namespaces.
- [ ] Run `go test ./app/generator -run 'TestNamespace|TestNormalize'` and observe the missing helper/field failure.
- [ ] Implement the helpers and add `Namespace` to `Input` and `Spec`.
- [ ] Run the focused tests and the existing generator tests.
- [ ] Commit with `feat: define generator API namespaces`.

### Task 2: Namespace-aware renderers

**Files:**
- Modify: `backend/app/generator/permission.go`
- Modify: `backend/app/generator/menu.go`
- Modify: `backend/app/generator/resource_pipeline.go`
- Modify: `backend/app/generator/spec.go`
- Test: `backend/app/generator/permission_test.go`
- Test: `backend/app/generator/menu_test.go`

**Interfaces:**
- Permission and menu inputs accept an optional `Namespace`.
- Existing omitted namespace remains `admin`.
- `Normalize(Input{Namespace: "app"})` returns an explicit unsupported-surface error before rendering.

- [ ] Add failing tests for `app.orders.view` renderer output and App Resource rejection.
- [ ] Run the focused tests and observe failures.
- [ ] Implement namespace-aware permission/menu rendering and the Admin-only Resource guard.
- [ ] Run all generator tests and golden tests.
- [ ] Commit with `feat: guard resource generation by API surface`.

### Task 3: Boundary documentation and verification

**Files:**
- Modify: `docs/architecture.md`
- Modify: `docs/openapi.md`
- Modify: `docs/ai-quickstart.md`
- Test: existing Go generator and OpenAPI tests

- [ ] Document Admin/App route, permission, and OpenAPI boundaries and the current Admin-only generator behavior.
- [ ] Run `go test ./...` with a project-local Go cache.
- [ ] Run `git diff --check` and confirm the working tree contains only this stage.
- [ ] Commit with `docs: document admin and app api boundaries`.

