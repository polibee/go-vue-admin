# HTTP Audit Redaction Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Complete database readiness, restore the frontend toolchain, add bounded HTTP request/response audit with recursive redaction, and verify the full generated-resource workflow.

**Architecture:** Keep PostgreSQL and Redis as Laragon-managed services. Reuse the existing `audit_logs.metadata` JSONB field and add a single HTTP middleware after the framework response-origin middleware, so controllers do not each implement request/response capture. Keep redaction as a pure audit-service boundary with unit tests.

**Tech Stack:** Goravel v1.18, Gin adapter, Go 1.25, PostgreSQL, Redis, Vue 3, Vite, TypeScript, shadcn-vue.

**Spec:** `docs/superpowers/specs/2026-09-22-http-audit-redaction-design.md`

## Global Constraints

- Do not add dependencies or change lock files.
- PostgreSQL and Redis remain Laragon-managed services.
- Execute only the two already-reviewed pending migrations; do not generate unrelated migrations.
- Do not audit cookies, authorization headers, tokens, passwords, binary responses, or audit-log queries.
- Keep service code in `backend/app/services/audit/`; middleware remains in `backend/app/http/middleware/`.
- Keep feature-phase commits local; do not push GitHub in this phase.

## Review Focus

- Nested sensitive keys are redacted before JSON encoding; covered by redaction unit tests.
- Oversized request/response data is bounded without breaking the request; covered by size-limit tests.
- Unauthenticated login and failed requests do not break because audit identity is unavailable; covered by middleware/service tests.
- Audit-log reads do not recursively create audit rows; covered by route-scope test.
- Corrupt/non-JSON response bodies do not cause a second response failure; covered by capture helper tests.

### Task 1: Database and frontend preflight

**Files:**
- Modify: `backend/database` only through already-existing pending migrations.
- Modify: `admin/node_modules` only as an environment repair; do not modify tracked lock files.

- [ ] Verify migration status and PostgreSQL/Redis ports.
- [ ] Execute `go run . artisan migrate` from `backend/` with the existing `.env`.
- [ ] Verify both migrations report `Ran` and the `permission_role_field` table exists.
- [ ] Repair the existing Windows native package installation using the repository package manager without changing `pnpm-lock.yaml`.
- [ ] Run the frontend typecheck/build smoke and record the result.

### Task 2: Redaction contract

**Files:**
- Create: `backend/app/services/audit/redaction.go`
- Create: `backend/app/services/audit/redaction_test.go`
- Modify: `backend/app/services/audit/audit_service.go`

- [ ] Write failing tests for recursive sensitive-key replacement, nil handling, and a 32 KiB encoded payload limit.
- [ ] Run the focused Go tests and observe the expected failure.
- [ ] Implement pure redaction and bounded JSON encoding.
- [ ] Update `AuditService.Record` to pass metadata through the same redactor.
- [ ] Run focused tests and the audit service package tests.

### Task 3: HTTP audit middleware

**Files:**
- Create: `backend/app/http/middleware/http_audit.go`
- Create: `backend/app/http/middleware/http_audit_test.go`
- Modify: `backend/bootstrap/app.go`
- Modify: `backend/routes/web.go`

- [ ] Write failing tests for API path filtering, audit-route exclusion, request/response summary shape, and unavailable identity.
- [ ] Run focused tests and observe the expected failure.
- [ ] Implement bounded request capture from `ContextRequest`, response capture from `ResponseOrigin`, and optional authenticated user ID.
- [ ] Register the middleware after CORS without changing route authorization.
- [ ] Run middleware tests and the complete Go suite.

### Task 4: Audit contract and UI details

**Files:**
- Modify: `backend/app/openapi/spec.go` only if the metadata description needs explicit request/response fields.
- Modify: `admin/src/lib/audit-log.ts`
- Modify: `admin/src/modules/audit/pages/AuditLogPage.vue`
- Modify: `admin/src/locales/zh-CN/auth.json`
- Modify: `admin/src/locales/en-US/auth.json`

- [ ] Add safe formatting for structured request/response audit metadata and truncated indicators.
- [ ] Show request/response sections in the existing audit detail dialog without exposing raw secrets.
- [ ] Run frontend typecheck/build.

### Task 5: Full acceptance

**Files:**
- Modify: `backend/scripts/contract-smoke.ps1` if current paths or assertions need to cover audit metadata.
- Modify: `docs/roadmap.md`
- Modify: `docs/testing.md`

- [ ] Start backend and frontend in the background with logs and PIDs.
- [ ] Run contract smoke and verify login, registry, generated resource list, form create/update, detail, permissions, search, export, and relation endpoints.
- [ ] Verify audit rows contain redacted request/response data.
- [ ] Verify the browser pages for menu, list, form, detail, permissions, and search.
- [ ] Run Go tests, frontend typecheck/build, frontend tests, and `git diff --check`.
- [ ] Update roadmap with evidence and create one local stage commit; do not push.
