# Unified Go API Generator and Audit Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace platform-specific API generation with one deterministic Go CLI, add Swagger UI documentation, and audit the complete backend/frontend API contract while preserving generated GORM CRUD and automatic route registration.

**Architecture:** `admin-gen api` calls the existing Go OpenAPI builder and a Go TypeScript renderer. OpenAPI JSON is the single contract artifact consumed by Swagger UI, generated schemas, TypeScript client code, and audit checks. Resource generation continues to use Manifest as its intermediate structure and injects an explicit GORM database handle into generated modules.

**Tech Stack:** Go 1.25, Goravel, GORM, OpenAPI JSON, Vue 3, TypeScript, pnpm, Swagger UI assets.

**Spec:** `docs/superpowers/specs/2026-09-10-go-api-generator-audit-design.md`

## Global Constraints

- `admin-gen` is the only generator implementation.
- No new Bash, PowerShell, or Node API generation implementation.
- `spec-forge` is not added; the Go OpenAPI model remains the source of truth.
- Swagger UI is disabled by default outside local development and never generates contracts.
- Resource SQL identifiers and writable fields remain Manifest-whitelisted.
- Generated outputs must be deterministic and checked for drift.
- GitHub Actions remain manual-only.

### Task 1: Consolidate API generation under `admin-gen api`

**Files:**
- Create: `backend/app/core/generator/api.go`
- Create: `backend/app/core/generator/api_test.go`
- Modify: `backend/cmd/admin-gen/main.go`
- Modify: `backend/app/core/openapi/document.go`
- Delete: `backend/cmd/openapi/main.go`
- Delete: `scripts/openapi-generate.sh`
- Delete: `scripts/generate-ts-client.mjs`
- Modify: `package.json`

**Interfaces:**
- `GenerateAPI(options APIOptions) error`
- `APIOptions{RootDir, Output, SchemaDir, ClientDir string}`
- `admin-gen api --output --schema-dir --client-dir`

- [ ] Write a failing test for deterministic OpenAPI and TypeScript output in a temporary root.
- [ ] Run `go test ./app/core/generator -run TestGenerateAPI -count=1` and confirm the command is absent.
- [ ] Implement the Go renderer using the existing OpenAPI document and stable JSON encoding.
- [ ] Add generated client rendering for models, request methods, query serialization, and JSON/FormData bodies.
- [ ] Wire the CLI and replace package scripts with `go -C backend run ./cmd/admin-gen api`.
- [ ] Delete the old shell, Node, and duplicate OpenAPI command entry points.
- [ ] Run generator tests and verify generated output is byte-stable.
- [ ] Commit `feat(generator): unify api generation in go cli`.

### Task 2: Add Swagger UI as a read-only backend documentation surface

**Files:**
- Create: `backend/app/http/controllers/api_docs_controller.go`
- Create: `backend/app/http/controllers/api_docs_controller_test.go`
- Create: `backend/resources/views/api-docs.tmpl`
- Modify: `backend/routes/web.go`
- Modify: `backend/config/app.go` or the project API docs configuration file
- Modify: `backend/app/core/openapi/metadata.go`

**Interfaces:**
- `NewAPIDocsController(documentPath string, enabled bool) *APIDocsController`
- `Index(ctx http.Context) http.Response`
- `Document(ctx http.Context) http.Response`

- [ ] Test disabled docs return 404 and enabled docs return the Swagger HTML shell.
- [ ] Test the document endpoint returns the generated JSON with the documented media type.
- [ ] Add local bundled Swagger UI assets or a pinned project-local package; do not use a runtime CDN.
- [ ] Register `/api/docs` and `/api/docs/openapi.json` only when enabled.
- [ ] Run controller tests and verify production defaults remain disabled.
- [ ] Commit `feat(docs): add opt-in swagger ui`.

### Task 3: Implement API route-to-client audit checks

**Files:**
- Create: `backend/app/core/generator/audit.go`
- Create: `backend/app/core/generator/audit_test.go`
- Modify: `backend/cmd/admin-gen/main.go`
- Modify: `scripts/release-check.sh` or replace its generator call with `admin-gen audit`
- Create: `docs/API_AUDIT.md`

**Interfaces:**
- `AuditAPI(options APIAuditOptions) ([]APIAuditFinding, error)`
- `admin-gen audit --root .`
- Finding fields: `method`, `path`, `layer`, `message`, `severity`.

- [ ] Write fixtures for an undocumented route, missing generated operation, and clean route.
- [ ] Run the audit test and verify each fixture fails for the expected reason.
- [ ] Implement route and OpenAPI inventory checks without executing the server.
- [ ] Add generated-client drift and direct-call checks for generated resource endpoints.
- [ ] Make release-check invoke the Go audit command and fail with actionable findings.
- [ ] Document the route/controller/service/repository/client evidence chain.
- [ ] Commit `feat(audit): verify api contract coverage`.

### Task 4: Harden generated GORM CRUD and route registration

**Files:**
- Modify: `backend/app/core/generator/resource.go`
- Modify: `backend/app/core/generator/resource_test.go`
- Modify: `backend/app/core/resource/gorm_database.go`
- Modify: `backend/routes/web.go`
- Modify: generated module templates in `backend/app/core/generator/module.go`

- [ ] Add a failing multi-resource generation test proving imports and `RegisterRoutes` do not collide.
- [ ] Refactor module registration to accept multiple generated resources.
- [ ] Move resource DB creation to application startup ownership and close it on shutdown.
- [ ] Make provider selection choose the configured GORM driver without silently falling back to memory.
- [ ] Run generated module compile tests and all Go tests.
- [ ] Commit `fix(generator): support multiple gorm resources`.

### Task 5: Regenerate and audit frontend API consumers

**Files:**
- Regenerate: `contracts/openapi/openapi.json`
- Regenerate: `contracts/schemas/*.json`
- Regenerate: `admin/src/generated/api/*`
- Modify: `admin/src/core/api/OpenApiDataProvider.ts`
- Modify: `admin/src/modules/*/*.service.ts` only where audit finds direct bypasses
- Create/modify: frontend contract tests

- [ ] Run `go -C backend run ./cmd/admin-gen api`.
- [ ] Verify generated files are deterministic and type-safe.
- [ ] Update consumers to use generated methods instead of handwritten API paths.
- [ ] Run frontend unit tests, typecheck, and build.
- [ ] Commit `refactor(admin): consume unified generated api client`.

### Task 6: MySQL integration and browser acceptance

**Files:**
- Modify: `.env.example` and `backend/.env.example`
- Create: `backend/app/core/resource/mysql_integration_test.go` with an opt-in environment guard
- Modify: `docs/RESOURCE_GENERATION.md`
- Modify: `docs/API_AUDIT.md`

- [ ] Check the configured MySQL host and port without changing credentials.
- [ ] Run GORM ping and migration checks when the database is reachable.
- [ ] Run one generated resource CRUD round trip against MySQL.
- [ ] Start backend and Vite in the background with logs and PID files.
- [ ] Verify `/api/health`, `/api/auth/bootstrap`, `/login`, Swagger UI, and one resource page.
- [ ] Use Playwright snapshots for login and resource CRUD acceptance.
- [ ] Run local release-check with E2E enabled, without invoking GitHub CI.
- [ ] Commit `test(resource): verify mysql and browser contracts`.
