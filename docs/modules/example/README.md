# Example module

This fixture demonstrates the compiled module contract used by Task 12.

- Manifest: `modules/example/module.yaml`
- Backend module: `modules/example/backend/module.go`
- Admin definition: `modules/example/admin/module.ts`
- Contract check: `bash scripts/module-check.sh modules/example`

Use `bash scripts/make-module.sh <module-name>` to create another module with
the same backend, admin, route, resource, locale, and acceptance-test layout.
