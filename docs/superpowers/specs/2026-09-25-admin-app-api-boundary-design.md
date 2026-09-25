# Admin/App API Boundary Design

## Decision

The repository has two future API surfaces, but only the Admin surface is currently implemented:

- Admin pages: `/admin/...`
- Admin API: `/api/v1/admin/...`
- Admin permission namespace: `admin.<resource>.<action>`
- Future end-user pages: root paths such as `/orders`
- Future end-user API: `/api/v1/app/...`
- Future end-user permission namespace: `app.<resource>.<action>`

The same domain Service may be reused by both surfaces. Controllers, route middleware, OpenAPI documents, and authorization policies remain surface-specific. An App endpoint must not reuse an Admin permission merely because the business operation looks similar.

## Current boundary

`admin:make-resource` is an Admin Resource generator. It creates an Admin route, Admin API base, Admin menu metadata, and `admin.*` permissions. It must reject an App namespace until the repository has an App registry, App route registration, App authentication policy, and App OpenAPI contract. This prevents a partially generated App resource from accidentally being exposed through the Admin runtime.

The namespace validation and permission/menu renderers are reusable primitives. Existing `admin.*` names are not renamed or migrated.

## Authorization model

- Admin authorization uses RBAC, data scope, and field permissions.
- App authorization uses authenticated ownership or tenant membership, plus a separate entitlement/policy check when required.
- Shared domain services enforce domain invariants; controllers do not write authorization tables directly.
- A second user table and polymorphic subject model are not introduced until separate identity stores or multi-tenancy are required.

## OpenAPI contract

The canonical Admin contract is `/api/openapi/admin.json`; `/api/openapi.json` remains a compatibility alias. It is not a complete contract for future App APIs or module-specific routes. When App APIs are introduced, publish separate Admin and App documents and add route-coverage tests rather than silently mixing both surfaces.

## Non-goals

- No permission rename.
- No database migration.
- No second user table.
- No App business API implementation in this phase.
