# Resource-Driven Generator Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rework the generator so `admin:make-resource` and `admin:make-crud` are the primary resource-driven entry points that produce backend, permissions, menus, frontend list/form/detail pages, routes, README files, migrations, and tests from one ResourceSpec.

**Architecture:** Parse the input once into `ResourceSpec`; pass that spec through backend, permission, menu, frontend, test, and README renderers; then use one conflict-aware Writer for every artifact. Existing permission/menu commands become thin adapters over the shared renderers, not separate feature pipelines.

**Tech Stack:** Go 1.27, Goravel v1.18 console contracts, `go/format`, Vue 3, Vue Router, existing ResourceList/ResourceForm/ResourceDetail views, standard-library filesystem APIs, Go and Node test runners.

**Spec:** `docs/superpowers/specs/2026-09-21-resource-generator-design.md`

## Global Constraints

- Resource is the primary input and single source of truth for fields, permissions, menu metadata, routes, and pages.
- `admin:make-resource` and `admin:make-crud` are the primary user-facing workflows.
- Permission and menu generation must be internal renderers invoked by the resource pipeline; they must not have standalone Artisan commands.
- Frontend output must reuse existing Admin Shell and generic resource views; do not create duplicate UI primitives.
- All artifacts are preflighted together; any conflict prevents every write.
- Generate migration files but never execute migrations or connect to PostgreSQL.
- Do not modify existing Registry, Sidebar, router, Provider, Seeder, permission rows, role bindings, or route registration automatically.
- No new dependencies or lockfile changes.

## Review Focus

- One field definition must produce matching backend metadata, frontend form metadata, list columns, detail fields, and tests.
- Permission, menu, and route names must be derived identically from the same ResourceSpec.
- A conflict in a frontend artifact must prevent backend and migration artifacts from being written.
- Generated Vue pages must use existing generic components and permission metadata rather than duplicate business UI.
- A generated migration must be present as a file while PostgreSQL remains untouched.

---

### Task 1: Expand the shared ResourceSpec

**Files:**
- Modify: `backend/app/generator/spec.go`
- Test: `backend/app/generator/spec_test.go`

- [ ] Add failing tests for `icon`, frontend route normalization, action defaults, and consistent permission/menu names.
- [ ] Run `go test ./app/generator -run 'Test.*Spec' -count=1` and confirm the new assertions fail.
- [ ] Extend `Spec` with normalized menu and page metadata without introducing a second permission/menu input parser.
- [ ] Run `go test ./app/generator -count=1` and commit `feat: expand resource generator spec`.

### Task 2: Make permission and menu rendering resource-derived

**Files:**
- Modify: `backend/app/generator/permission.go`
- Modify: `backend/app/generator/menu.go`
- Modify: `backend/app/generator/render.go`
- Test: `backend/app/generator/permission_test.go`, `menu_test.go`, `render_test.go`

- [ ] Add failing tests proving a ResourceSpec generates the same permission constants and menu entry as the shared renderers.
- [ ] Run the focused tests and confirm failure before implementation.
- [ ] Add renderer functions that consume ResourceSpec directly; retain standalone command adapters only as compatibility wrappers.
- [ ] Update Golden Files and run `go test ./app/generator -count=1`.
- [ ] Commit `refactor: derive permission and menu artifacts from resources`.

### Task 3: Add frontend resource artifact rendering

**Files:**
- Create: `backend/app/generator/frontend_render.go`
- Test: `backend/app/generator/frontend_render_test.go`
- Create: `admin/src/generated/resources/posts/resource.ts.golden`
- Create: `admin/src/generated/resources/posts/menu.ts.golden`
- Create: `admin/src/generated/resources/posts/routes.ts.golden`
- Create: `admin/src/generated/resources/posts/pages/*.vue.golden`
- Create: `admin/src/generated/resources/posts/posts.test.ts.golden`

- [ ] Write a failing Golden File test for the frontend artifact paths and contents from one `posts` ResourceSpec.
- [ ] Run `go test ./app/generator -run TestRenderFrontend -count=1` and confirm failure.
- [ ] Render typed metadata, menu, route descriptors, thin List/Form/Detail wrappers, and a frontend test. Use existing `ResourceListView`, `UserFormView`/`RoleFormView` patterns only through stable props/contracts; do not duplicate shadcn primitives.
- [ ] Run the generator package tests and inspect the generated Vue syntax.
- [ ] Commit `feat: render resource frontend artifacts`.

### Task 4: Compose the full resource pipeline

**Files:**
- Modify: `backend/app/generator/crud.go`
- Modify: `backend/app/console/resource_generator_command.go`
- Modify: `backend/app/console/crud_generator_command.go`
- Test: `backend/app/generator/crud_test.go`, `backend/app/console/*generator_command_test.go`

- [ ] Add a failing composition test asserting one ResourceSpec produces backend, migration, permission, menu, frontend, README, and test artifacts with no duplicate paths.
- [ ] Run the focused test and confirm failure.
- [ ] Make `admin:make-resource` call the full pipeline; make `admin:make-crud` call the same pipeline without a second parser or renderer set.
- [ ] Remove standalone permission/menu Artisan commands; keep only shared renderers invoked by the resource pipeline.
- [ ] Run command tests and commit `feat: compose resource-driven generator pipeline`.

### Task 5: Add all-artifact conflict safety and smoke coverage

**Files:**
- Modify: `backend/app/generator/writer.go`
- Test: `backend/app/generator/writer_test.go`, `backend/app/generator/crud_test.go`

- [ ] Add a failing test with a conflict in a frontend file and assert no backend, migration, permission, or menu file exists afterward.
- [ ] Implement/verify the existing preflight behavior across the complete artifact set.
- [ ] Run a real temporary-directory smoke for both `admin:make-resource` and `admin:make-crud`; verify generated pages and README exist, migration is not executed, and existing registration files are unchanged.
- [ ] Run `go test ./... -count=1`, frontend `vue-tsc`, build, and tests.
- [ ] Commit `test: verify resource-driven generator safety`.

### Task 6: Update documentation and roadmap

**Files:**
- Modify: `docs/generator.md`
- Modify: `docs/roadmap.md`

- [ ] Document `admin:make-resource` and `admin:make-crud` as the primary workflows and list the full artifact set.
- [ ] Mark standalone permission/menu commands as compatibility helpers.
- [ ] Document manual activation: review generated files, register runtime boundaries, review migration, then execute migration manually.
- [ ] Run `git diff --check` and commit `docs: document resource-driven generator workflow`.
