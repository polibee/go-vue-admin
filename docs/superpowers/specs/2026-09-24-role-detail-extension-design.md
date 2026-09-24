# Role Detail Extension Design

## Goal

Keep the generic resource detail page reusable while allowing business modules,
starting with roles, to add domain-specific detail sections without importing
business components directly from `core`.

## Scope

- Show a localized not-found state for missing resource records.
- Load detail extensions only after the base record is confirmed to exist.
- Register the role permission panel from the roles module through a core-owned
  extension registry.
- Add permission search and per-resource select-all controls.
- Keep `super-admin` read-only and preserve the existing permission API.
- Make the role edit-page detail link wrap cleanly on narrow screens.
- Add automated coverage for normal-role persistence, revocation, and
  forbidden writes.

## Architecture

`core/resource/pages/ResourceDetailPage.vue` owns generic loading, not-found,
error, and extension rendering. `core/resource/detail-extensions.ts` exposes a
small registry keyed by resource name. Business modules register extension
descriptors during application setup; the core page resolves descriptors and
passes the loaded record and route identity to them. The core layer never
imports `modules/roles`.

An extension descriptor contains a resource key, a Vue component, and a stable
renderer contract. The registry rejects duplicate resource keys so registration
cannot silently replace an existing extension.

## Behavior

1. A detail request returning no record renders `resource.resourceNotFound`,
   does not render business extensions, and disables edit/delete actions.
2. Permission search filters by localized display name or permission name.
3. A resource group can select or clear all visible permissions in that group.
4. Saving a normal role reloads its detail record; the selected state remains
   after refresh.
5. Removing all permissions persists an empty assignment set.
6. `super-admin` cannot edit permissions and keeps all assignments checked and
   disabled.
7. A user without role-management permission cannot save role assignments; the
   UI reports the API authorization error.
8. The edit-page detail link uses a responsive layout and remains accessible on
   narrow viewports.

## Non-goals

- No new permission model or database migration.
- No changes to the generic resource URL contract.
- No separate role-specific detail page.
- No new frontend dependency.

## Acceptance Criteria

- `core/resource/pages/ResourceDetailPage.vue` has no import from
  `modules/roles`.
- Missing role IDs show a not-found state and cannot submit permissions.
- Role permissions are rendered through the extension registry.
- Search and group select-all work for normal roles and do not change
  `super-admin` read-only behavior.
- Backend tests cover permission save, revoke, and forbidden writes.
- Frontend type check, production build, and browser verification pass.
