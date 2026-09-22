# Generic and Custom Resource Pages Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make Manifest plus generic pages the default for ordinary resources while preserving explicit custom-page overrides for complex resources.

**Architecture:** Extend the backend and frontend Resource Manifest with an explicit page mode and optional page overrides. Generated ordinary resources register only their ResourceSpec and reuse the core ResourceListPage, ResourceFormPage, and ResourceDetailPage. The router resolves pages from the manifest and fails validation when a declared custom override is missing.

**Tech Stack:** Go, Goravel, Vue 3, TypeScript, Vue Router, Vite, existing shadcn-vue components, Go tests, existing frontend type/build checks.

**Spec:** docs/superpowers/specs/2026-09-22-generic-resource-page-design.md

## Global Constraints

- Ordinary resources default to generic pages.
- Complex resources use custom pages only through an explicit pageMode declaration.
- The public generator remains admin:make-resource; no second generator workflow is introduced.
- Migrations are generated but never executed automatically.
- Resource URLs remain /<resource>, without /resources in the URL.
- Authentication, resource authorization, data scope, field permissions, audit redaction, and batch-selection validation remain server-side.
- PostgreSQL and Redis remain Laragon-managed services; no in-memory replacement is allowed.
- Existing users, roles, permissions, and announcements behavior must remain available during the migration.

## Review Focus

- Generated ordinary resources must not import missing resource-specific Vue pages.
- Generic forms must render select, boolean, relation, grouped, and dependency fields from the manifest.
- Generic detail pages must respect readable and sensitive field policies.
- Custom page declarations must fail clearly when a component is missing or malformed.
- Existing query selection, cross-page selection, filtered-set selection, bulk actions, soft delete, and resource permissions must not regress.

---

### Task 1: Extend the Resource page-mode contract

**Files:**
- Modify: backend/app/core/resource/registry.go
- Modify: backend/app/generator/spec.go
- Modify: backend/app/generator/spec_test.go
- Modify: backend/app/core/resource/registry_test.go
- Test: backend/app/generator/spec_test.go and backend/app/core/resource/registry_test.go

**Interfaces:**
- Add a page mode value with generic as the default and custom as the only alternative.
- Add optional custom page identifiers to the generated frontend contract without embedding Vue component imports in the Go manifest.
- Reject unknown page modes and reject custom declarations that omit the required override identifier.

- [ ] **Step 1: Write failing tests**

Add tests proving that:

1. A normalized input without a page-mode option becomes generic.
2. An explicit custom page mode is preserved.
3. An unknown page mode is rejected.
4. A manifest cannot register an invalid page mode.

- [ ] **Step 2: Run the focused tests and verify the expected failure**

Run:

~~~powershell
go test ./backend/app/generator ./backend/app/core/resource -run "PageMode|Manifest" -count=1
~~~

Expected: failures because page mode is not yet represented or validated.

- [ ] **Step 3: Implement the smallest contract change**

Add page-mode fields and validation to the Resource Manifest and generator Spec. Keep the existing JSON naming conventions and preserve all existing default values.

- [ ] **Step 4: Run focused tests**

Run the same command and expect all focused tests to pass.

- [ ] **Step 5: Commit**

~~~powershell
git add backend/app/core/resource backend/app/generator
git commit -m "feat: add resource page mode contract"
~~~

### Task 2: Add frontend generic page resolution

**Files:**
- Create: admin/src/core/resource/page-resolver.ts
- Create: admin/src/core/resource/page-resolver.test.ts
- Modify: admin/src/core/resource/generated.ts
- Modify: admin/src/router/index.ts
- Modify: admin/src/core/resource/pages/ResourceListPage.vue
- Create or modify: admin/src/core/resource/pages/ResourceFormPage.vue
- Modify: admin/src/core/resource/pages/ResourceDetailPage.vue
- Test: admin/src/core/resource/page-resolver.test.ts

**Interfaces:**
- page-resolver.ts consumes a generated ResourceDefinition and returns the route component mode for list, form, and detail.
- Generic resources resolve to core pages.
- Custom resources resolve only through explicit component metadata registered by their module.
- The resolver must expose a validation function that returns a deterministic error for a missing custom component.

- [ ] **Step 1: Write failing resolver tests**

Cover:

1. Missing pageMode resolves to generic.
2. Explicit generic resolves to all three core pages.
3. Explicit custom resolves to declared custom components.
4. Missing custom list, form, or detail override throws a descriptive error.

- [ ] **Step 2: Run the frontend test or type-level check and verify failure**

Run the repository’s available frontend test command for the focused file. If no frontend test runner is configured, run the test through the existing project test harness or add only the minimal runner-independent contract assertion used by the generator.

Expected: failure because the resolver does not exist.

- [ ] **Step 3: Implement the resolver and generic route records**

Replace resource-specific imports in generated route output with generic core-page imports for generic resources. Keep custom component imports isolated to the custom resource module.

- [ ] **Step 4: Run frontend checks**

Run:

~~~powershell
cd admin
pnpm exec vue-tsc --noEmit
pnpm run build
~~~

Expected: pass with the existing users, roles, permissions, and announcements routes intact.

- [ ] **Step 5: Commit**

~~~powershell
git add admin/src/core/resource admin/src/router/index.ts
git commit -m "feat: resolve generic resource pages from manifests"
~~~

### Task 3: Change generator output to ordinary-resource mode

**Files:**
- Modify: backend/app/generator/frontend_render.go
- Modify: backend/app/generator/runtime_registration.go
- Modify: backend/app/generator/render_test.go
- Modify: backend/app/generator/frontend_render_test.go
- Modify: backend/app/generator/runtime_registration_test.go
- Modify: backend/app/generator/resource_pipeline_test.go
- Modify: backend/app/generator/testdata/posts/*

**Interfaces:**
- RenderFrontend returns resource.ts, api.ts, the generated resource contract test, and any explicitly requested custom-page artifacts.
- Generic generation must not create ListPage.vue, FormPage.vue, or DetailPage.vue.
- Runtime registration must register generic routes without importing resource-specific Vue pages.
- Custom mode must keep the existing module page directory and generate or validate its explicit override metadata.

- [ ] **Step 1: Update Golden File expectations first**

Change the expected ordinary-resource artifact list so it contains no resource-specific Vue pages and expects pageMode generic in the generated resource.ts.

- [ ] **Step 2: Run generator tests and verify failure**

Run:

~~~powershell
go test ./backend/app/generator -count=1
~~~

Expected: failures showing the current renderer still emits resource-specific page files and imports.

- [ ] **Step 3: Implement generic output**

Remove the default three page artifacts from RenderFrontend. Update runtime registration to generate generic route records with manifest-provided resource names and permissions. Keep custom-page output behind the explicit custom mode only.

- [ ] **Step 4: Run the complete generator suite**

Run:

~~~powershell
go test ./backend/app/generator -count=1
~~~

Expected: all generator and Golden File tests pass.

- [ ] **Step 5: Commit**

~~~powershell
git add backend/app/generator
git commit -m "feat: generate generic resources without dedicated pages"
~~~

### Task 4: Complete generic form and detail rendering

**Files:**
- Modify: admin/src/core/resource/pages/ResourceFormPage.vue
- Modify: admin/src/core/resource/pages/ResourceDetailPage.vue
- Modify: admin/src/components/resource/ResourceFormView.vue
- Modify: admin/src/core/resource/pages/ResourceListPage.vue
- Create: admin/src/core/resource/types.ts
- Test: existing resource contract tests plus new focused generic-page tests

**Interfaces:**
- Generic pages consume only Resource Manifest data and resource API adapters.
- Form rendering must honor required, readonly, writable, relation, form-group, and field-dependency metadata.
- Detail rendering must honor readable, visible, sensitive, relation, and detail-section metadata.
- List rendering must preserve current filters, pagination, export, soft-delete, and batch-selection behavior.

- [ ] **Step 1: Add failing contract assertions**

Add tests for:

1. A generated ordinary resource can render its form from fields alone.
2. A readonly field is not submitted.
3. An unreadable or sensitive field is not shown in detail output.
4. A relation field uses the relation metadata instead of a hardcoded resource branch.

- [ ] **Step 2: Run the focused frontend checks and verify failure**

Run the focused test/check command and confirm the failures identify missing generic rendering behavior.

- [ ] **Step 3: Implement manifest-driven rendering**

Move any remaining announcements-specific or resource-name-specific rendering into manifest adapters or explicit custom modules. Keep shared UI primitives in admin/src/components and generic orchestration in admin/src/core/resource.

- [ ] **Step 4: Run frontend validation**

Run:

~~~powershell
cd admin
pnpm exec vue-tsc --noEmit
pnpm run build
~~~

Expected: pass without changing the existing port contract.

- [ ] **Step 5: Commit**

~~~powershell
git add admin/src/core/resource admin/src/components/resource admin/src/core/resource/types.ts
git commit -m "feat: render generic forms and details from manifests"
~~~

### Task 5: Add the departments ordinary-resource acceptance fixture

**Files:**
- Create: backend/app/modules/departments/resource/manifest.go
- Create: backend/app/modules/departments/resource/manifest_test.go
- Create: backend/app/modules/departments/repository.go
- Create: backend/app/modules/departments/service.go
- Create: backend/database/migrations/*_create_departments_table.go
- Create: admin/src/modules/departments/resource.ts
- Modify: backend/app/modules/admin/registry/generated_resources.go
- Modify: admin/src/core/resource/generated.ts
- Create: backend/app/modules/departments/README.md
- Test: backend/app/modules/departments/resource/manifest_test.go and generator/runtime registration tests

**Interfaces:**
- Departments must be a normal generated resource with no dedicated Vue ListPage, FormPage, or DetailPage.
- Its Manifest must include searchable, sortable, filterable fields and standard CRUD permissions.
- Its registry entry must be discovered by the existing generated registration mechanism.

- [ ] **Step 1: Add the failing acceptance assertions**

Assert that departments is discoverable in the backend registry and frontend resource registry, that its route is /departments, and that no departments-specific page artifact exists.

- [ ] **Step 2: Run the acceptance tests and verify failure**

Run:

~~~powershell
go test ./backend/app/modules/departments/... ./backend/app/modules/admin/registry ./backend/app/generator -count=1
~~~

Expected: failure because departments does not yet exist.

- [ ] **Step 3: Generate and implement the departments fixture**

Use the resource generator output as the source of truth. Review the migration without executing it automatically. Keep the fixture small enough to test the generic path rather than introduce business-specific UI.

- [ ] **Step 4: Run backend and frontend checks**

Run:

~~~powershell
go test ./...
cd admin
pnpm exec vue-tsc --noEmit
pnpm run build
~~~

Expected: all checks pass. If the local PostgreSQL/Redis services are not available, report the exact port failure and do not substitute another service.

- [ ] **Step 5: Commit**

~~~powershell
git add backend/app/modules/departments backend/database/migrations admin/src/modules/departments backend/app/modules/admin/registry/generated_resources.go admin/src/core/resource/generated.ts
git commit -m "test: add generic departments resource acceptance"
~~~

### Task 6: Verify custom-page override behavior

**Files:**
- Modify: backend/app/generator/spec.go
- Modify: backend/app/generator/frontend_render.go
- Modify: admin/src/core/resource/page-resolver.ts
- Create: admin/src/modules/orders/pages/OrderListPage.vue
- Create: admin/src/modules/orders/pages/OrderFormPage.vue
- Create: admin/src/modules/orders/pages/OrderDetailPage.vue
- Create: admin/src/modules/orders/resource.ts
- Test: generator custom-mode tests and page-resolver tests

**Interfaces:**
- Custom mode must reuse the same manifest permissions, API, Action, and audit boundaries.
- Custom pages are selected only by explicit pageMode custom metadata.
- Missing custom components fail validation rather than silently falling back.

- [ ] **Step 1: Add failing custom-mode tests**

Test successful resolution with all declared overrides and failure when one override is missing.

- [ ] **Step 2: Run focused tests and verify failure**

Run the generator and page-resolver focused suites.

- [ ] **Step 3: Implement the explicit override path**

Keep the order of resolution deterministic: generic default, explicit generic, explicit custom. Do not add resource-name conditionals to the generic pages.

- [ ] **Step 4: Run full validation**

Run:

~~~powershell
go test ./...
cd admin
pnpm exec vue-tsc --noEmit
pnpm run build
~~~

- [ ] **Step 5: Commit**

~~~powershell
git add backend/app/generator admin/src/core/resource admin/src/modules/orders
git commit -m "feat: support explicit custom resource page overrides"
~~~

### Task 7: Update documentation and perform the milestone acceptance

**Files:**
- Modify: README.zh-CN.md
- Modify: README.md
- Modify: docs/ai-quickstart.md
- Modify: docs/generator.md
- Modify: docs/resource-engine.md
- Modify: docs/roadmap.md
- Modify: backend/app/generator/testdata/posts/README.md.golden

**Interfaces:**
- Documentation must state that generic mode is the default.
- Documentation must explain when to use custom mode.
- Documentation must show the departments acceptance workflow and migration review boundary.

- [ ] **Step 1: Add documentation assertions or checklist entries**

Record the exact generated file list and the absence of ordinary-resource page files.

- [ ] **Step 2: Run the complete verification set**

Run:

~~~powershell
git diff --check
go test ./...
cd admin
pnpm exec vue-tsc --noEmit
pnpm run build
~~~

Then manually verify departments at http://127.0.0.1:5180/departments after the reviewed migration is applied in the local acceptance database.

- [ ] **Step 3: Commit the stage milestone**

~~~powershell
git add README.md README.zh-CN.md docs/ai-quickstart.md docs/generator.md docs/resource-engine.md docs/roadmap.md backend/app/generator/testdata/posts
git commit -m "docs: document generic and custom resource modes"
~~~

- [ ] **Step 4: Push only after the complete milestone passes**

Confirm the working tree, tests, frontend build, migration review, and browser acceptance before pushing the complete stage. Do not push intermediate task commits individually.
