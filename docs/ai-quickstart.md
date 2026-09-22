# AI Quickstart

This document is the operating contract for an AI agent working on Go Vue Admin. Read it before inspecting or changing the repository.

## 1. Project identity

- Backend: Goravel/Go in backend.
- Frontend: Vue 3, TypeScript, Vite, and shadcn-vue in admin.
- Database: PostgreSQL managed by Laragon.
- Cache and session acceleration: Redis managed by Laragon.
- Go API port: 3000.
- Go/Vue admin port: 5180.
- Port 5173 belongs to the separate Python admin project.

Never merge the Python and Go/Vue frontends onto one port. Do not replace PostgreSQL or Redis with an in-memory fallback, and do not change backend/.env or Laragon service configuration without explicit approval.

## 2. First checks

Run these checks before making a change:

~~~powershell
git status --short
git log -5 --oneline
Test-NetConnection 127.0.0.1 -Port 5432
Test-NetConnection 127.0.0.1 -Port 6379
~~~

If a service is unavailable, report the exact port and error. Do not silently switch databases, caches, or ports.

## 3. Directory contract

- Platform backend code belongs in backend/app/core.
- Business backend code belongs in backend/app/modules/<module>.
- Backend services stay under backend/app/services/<domain>, grouped by domain such as auth, users, rbac, and audit.
- Console commands are grouped under backend/app/console/<domain>.
- Frontend shared components belong in admin/src/components.
- Frontend infrastructure belongs in admin/src/core.
- Frontend business pages and components belong in admin/src/modules/<module>.

Do not create misc, common, other, generated, or catch-all folders. Do not scatter one domain across unrelated service folders.

## 4. Resource workflow

The public generator entry point is admin:make-resource. It consumes one ResourceSpec and generates the resource backend contract, migration file, permissions, menu metadata, frontend ResourceSpec, tests, and README. Ordinary resources use the core generic list/form/detail pages; only an explicit pageMode custom resource receives dedicated page overrides.

The generator must:

1. Validate the resource name and fields.
2. Preflight conflicts and never overwrite manual files.
3. Generate migrations without executing them.
4. Make the generated resource discoverable by the admin panel and generic pages.
5. Keep custom business logic in the module boundary.

Review the migration, apply it through the project migration process, then verify menu, list, form, detail, permissions, search, export, relations, and batch actions end to end.

## 5. Generic Resource versus custom module

Use a Resource for ordinary CRUD with standard fields, filters, permissions, relations, and batch actions. Use a module for workflows with state machines, external integrations, complex transactions, custom dashboards, or domain-specific UI. A module may reuse Resource components, but it must not bypass authorization or repository/service boundaries.

## 6. Verification

For backend changes:

~~~powershell
cd backend
go test ./...
~~~

For frontend changes:

~~~powershell
cd admin
pnpm exec vue-tsc --noEmit
pnpm run build
~~~

For generator changes, run the golden-file tests and a real Artisan-style smoke command. For UI changes, verify the affected route at port 5180 with PostgreSQL and Redis available.

## 7. Change and Git policy

Prefer small, independently verifiable changes. Keep feature-phase commits local while a phase is incomplete. Push only an independently reviewable milestone after its tests and acceptance checks pass. Never claim a feature is production-ready without fresh verification evidence and an explicit list of remaining environment gates.

## 8. Security boundaries

Preserve authentication, permission checks, data-scope filtering, field filtering, audit redaction, and CSRF/CORS boundaries. Never log passwords, tokens, cookies, authorization headers, or raw secrets. Destructive actions require the existing confirmation and authorization flow.

## 9. Expected AI response

When finishing a task, report:

1. Changed files and the responsibility of each change.
2. Verification commands and their results.
3. Any migration or manual registration step that remains.
4. Any production or environment limitation.
5. Commit and push status, without claiming a push if it was not verified.
