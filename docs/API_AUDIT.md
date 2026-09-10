# API audit

`admin-gen audit` statically verifies the evidence chain for the generic admin
API:

```text
route -> controller/service/repository -> OpenAPI operation -> generated client -> frontend consumer
```

Run it locally from the repository root:

```bash
go -C backend run ./cmd/admin-gen audit --root ..
```

The audit reports undocumented `/api` routes, missing generated operations,
stale OpenAPI/client artifacts, and direct frontend calls to generated
resource endpoints. Errors fail the command; direct frontend calls are
warnings so existing migration work remains visible without hiding release
blockers.

The check does not start the server or connect to MySQL. Runtime behavior,
authorization, and browser acceptance remain part of the local release-check.
