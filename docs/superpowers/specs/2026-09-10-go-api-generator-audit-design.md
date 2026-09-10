# Unified Go API Generation and Audit Design

## Decision

Use `backend/cmd/admin-gen` as the only project generator. The `api` subcommand will generate the OpenAPI JSON document, schema files, and the TypeScript client from the Go OpenAPI model. The existing Go OpenAPI builder remains the single contract source; `spec-forge` is not added because it would create a second specification source and increase drift risk.

Swagger UI will be added as a read-only documentation surface. It will load the generated `contracts/openapi/openapi.json` and will not generate or mutate the contract. The UI is intended for local development and manual acceptance, not as a production authorization bypass.

## Scope

- Replace shell and Node API-generation entry points with `admin-gen api`.
- Keep Resource Manifest as the only intermediate structure for database-driven CRUD generation.
- Generate TypeScript models and client code from the same OpenAPI document used by Swagger UI.
- Audit every public backend route against its controller, service, repository, OpenAPI metadata, generated client, and frontend consumer.
- Preserve the existing GORM `ResourceRepository` boundary and generated CRUD route registration.
- Verify memory mode first, then verify MySQL using the configured host, port, database, and credentials without changing database permissions.

## Non-goals

- No `spec-forge` dependency in the repository.
- No Bash, PowerShell, or Node program as a generator implementation.
- No automatic GitHub Actions or remote CI trigger.
- No runtime database introspection to generate code.

## Contract flow

```text
Go OpenAPI model
        |
        v
admin-gen api
   |        |        |
   v        v        v
OpenAPI   schemas   TypeScript client
   |
   v
Swagger UI / API audit / frontend services
```

## CLI contract

```text
go -C backend run ./cmd/admin-gen api \
  --output ../contracts/openapi/openapi.json \
  --schema-dir ../contracts/schemas \
  --client-dir ../admin/src/generated/api
```

The command must be deterministic, create parent directories, write stable JSON formatting, and fail on invalid output paths or invalid OpenAPI data. `module` and `resource` remain existing subcommands.

## API audit contract

The audit must report, per route:

```text
method + path -> controller -> service -> repository -> OpenAPI operation -> generated client -> frontend consumer
```

The audit fails for undocumented public routes, missing generated operations, unsafe dynamic resource fields, frontend direct calls that bypass the generated client for generated APIs, or generated-file drift.

## Swagger UI contract

- Serve `/api/docs` from the backend only when `API_DOCS_ENABLED=true`.
- Serve the generated document at `/api/docs/openapi.json`.
- Use a vendored or locally bundled UI asset; do not fetch a CDN at runtime.
- Keep the endpoint unauthenticated only in local development. Production configuration must default to disabled.

## Database and runtime contract

- `ResourceRepository` remains database-independent.
- `GormResourceRepository` receives an explicit `*gorm.DB`.
- The application owns the independent GORM connection pool and closes it during shutdown.
- Generated modules receive the initialized database handle when registering routes.
- MySQL verification is an integration check; SQLite tests remain the deterministic unit/integration substitute.

## Verification

- Go: `go test ./...`, `go vet ./...`, generator golden/contract tests.
- Admin: generated client tests, typecheck, build, and resource provider tests.
- Contract: regenerate and verify a clean git diff.
- Runtime: `/api/health`, `/api/auth/bootstrap`, `/login`, Swagger UI, authenticated resource list/create/update/delete.
- Database: connection ping, migration check, and one real CRUD round trip when MySQL is reachable.
