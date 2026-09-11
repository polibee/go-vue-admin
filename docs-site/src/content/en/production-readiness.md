# Production readiness

## Current assessment

The project is Foundation-ready for controlled development and internal acceptance. It has session authentication, RBAC, generic CRUD, GORM database selection, OpenAPI generation, plugin lifecycle and browser tests.

It is not a blanket production-ready product without business and operational hardening.

## Required before launch

- Use a real identity provider, strong password policy and MFA or SSO.
- Store and rotate secrets through a secret manager.
- Enforce HTTPS, trusted origins and secure cookie attributes.
- Test database backup, restore, migration rollback and high availability.
- Add rate limits, request-size limits, upload scanning and log redaction.
- Add structured logs, metrics, error tracking and alerts.
- Run authorization, tenant-isolation and dependency security tests.
- Define release, rollback and database change approval procedures.
- Disable development bootstrap credentials and verify public documentation contains no operational secrets.

## Local checks

    go -C backend test ./...
    go -C backend vet ./...
    pnpm run typecheck
    pnpm run test
    pnpm run build
    pnpm run openapi:generate
