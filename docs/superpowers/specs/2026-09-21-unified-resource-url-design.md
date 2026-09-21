# Unified Resource URL and CRUD Design

## Goal

Unify standard admin resources behind one CRUD architecture while keeping browser URLs and API URLs free of the `/resources/` segment.

The public resource identity is the resource name itself. Domain-specific safety rules remain inside domain services and are not expressed as separate CRUD routes.

## URL contract

Frontend routes:

```text
/users
/roles
/permissions
/<generated-resource>
/<generated-resource>/new
/<generated-resource>/:id
/<generated-resource>/:id/edit
```

Admin CRUD API:

```text
GET    /api/v1/admin/{resource}
POST   /api/v1/admin/{resource}
GET    /api/v1/admin/{resource}/{id}
PUT    /api/v1/admin/{resource}/{id}
DELETE /api/v1/admin/{resource}/{id}
```

The resource registry metadata endpoint is:

```text
GET /api/v1/admin/registry
```

No user-facing page link or API contract uses `/resources/`.

Relationship-specific endpoints remain explicit because they are not CRUD operations:

```text
PUT /api/v1/admin/roles/{id}/permissions
GET /api/v1/admin/users/{id}/roles
PUT /api/v1/admin/users/{id}/roles
PUT /api/v1/admin/users/status
```

## Backend boundary

The generic `ResourceController` owns the standard CRUD route and resolves the resource through the registry.

For standard resources:

- `users` delete delegates to `UserService.Delete`;
- `roles` delete delegates to `RoleService.Delete`;
- other generated resources use the validated generic table operation;
- create and update may delegate to domain services when a resource has domain invariants;
- the route and permission model remains the same regardless of the internal adapter.

The following protections must remain unchanged:

- user role-association cleanup;
- role permission/user-association cleanup;
- `super-admin` protection;
- last active administrator protection;
- resource action permission checks;
- generated-field and select-option validation.

The old routes `/api/v1/admin/users/{id}` and `/api/v1/admin/roles/{id}` must not remain as CRUD delete aliases after migration.

## Frontend boundary

The central generated API client owns the URL templates. Shared list, detail, and generated form pages call that client and never construct a `/resources/` path.

`resourceActionPath` is removed or reduced to non-URL action policy helpers. Delete buttons use the same generic client for every resource.

## OpenAPI and generator

- OpenAPI paths use `/admin/{resource}` and `/admin/{resource}/{id}`.
- The registry metadata path becomes `/admin/registry`.
- The generated TypeScript client follows the same paths.
- `admin:make-resource` continues to generate resource modules and discovery files; no new per-resource route registration is required.
- Existing generated modules are updated by the normal generator/discovery flow, not by a database migration.

## Compatibility and migration

This is an intentional URL contract change. No compatibility aliases are retained for the old `/resources/` API paths or the old users/roles CRUD delete paths.

The implementation must update routes, controllers, OpenAPI, generated Client, shared pages, generator templates, tests, and documentation as one stage. Database schema and migrations are out of scope.

## Acceptance criteria

1. Browser navigation contains no `/resources/` segment.
2. Generic CRUD works for generated resources and for users/roles through the same API shape.
3. Users and roles preserve their domain deletion protections.
4. Old specialized CRUD delete routes are absent.
5. OpenAPI, generated Client, generator output, and runtime routes agree.
6. Backend tests, frontend type checks, and generator smoke tests pass.
