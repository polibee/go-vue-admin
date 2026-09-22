# In-App Notifications Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a secure, PostgreSQL-backed in-app notification center with unread counts, read state, internal navigation, and bilingual AdminShell UI.

**Architecture:** Keep notification persistence and ownership checks in `backend/app/services/notifications/`. Expose authenticated notification endpoints through a focused admin controller, and expose a small typed client/composable under `admin/src/core/notifications/`; `AdminShell.vue` only composes the notification trigger and presentation. The first version is in-app only and uses no new dependency or external delivery channel.

**Tech Stack:** Go 1.27, Goravel HTTP/ORM, PostgreSQL migration, Vue 3 + TypeScript, existing shadcn-vue components, existing vue-i18n and generated API client.

**Spec:** `docs/superpowers/specs/2026-09-22-in-app-notifications-design.md`

## Global Constraints

- Do not add a third-party dependency.
- Notification APIs require authentication and always scope reads/writes to the authenticated user.
- Notification URLs must be internal relative paths; external URLs are rejected.
- Notification text is rendered as plain text, never raw HTML.
- Do not implement email, SMS, WebSocket, plugin notifications, file upload, or import.
- Generate the migration and register it in the migration list, but do not execute it until the migration is reviewed and explicitly approved.
- Keep feature-phase commits local; do not push GitHub during implementation.

## Review Focus

- Cross-user access: a valid notification ID belonging to another user must return 404/403 and must not change state; cover in the controller integration test.
- URL safety: absolute URLs, protocol-relative URLs, and `javascript:` values must be rejected; cover in the service validation test.
- Read idempotency: marking an already-read notification must succeed without changing the original read timestamp; cover in the service test.
- Empty and large result sets: list returns a stable empty page and bounded page size; cover in the controller test.
- UI failure and localization: list, unread count, and mark-read failures show an error state in both locales; cover with composable tests or the existing frontend test boundary.

---

### Task 1: Add the notification schema and domain service

**Files:**
- Create: `backend/app/models/notification.go`
- Create: `backend/database/migrations/20260922000005_create_notifications_table.go`
- Modify: `backend/bootstrap/migrations.go`
- Create: `backend/app/services/notifications/notification_service.go`
- Create: `backend/app/services/notifications/notification_service_test.go`

**Interfaces:**
- Produces `NotificationInput`, `Notification`, `NotificationPage`, `NewNotificationService()`, `Create(userID uint, input NotificationInput)`, `List(userID uint, page, perPage int, unreadOnly bool)`, `UnreadCount(userID uint)`, `MarkRead(userID uint, id uint)`, and `MarkAllRead(userID uint)`.
- `NotificationInput.URL` is validated as an empty string or an internal path beginning with `/` and not `//`.

- [ ] **Step 1: Write failing validation and ownership tests**

Add table-driven tests for valid internal paths, rejected absolute/protocol-relative/script paths, invalid user IDs, and idempotent read semantics. Use the service’s pure URL validator as a small unexported helper test target before connecting database operations.

- [ ] **Step 2: Run the focused service test and verify it fails**

Run `go test ./app/services/notifications -count=1` from `backend/` with `GOCACHE=backend/.gocache`.
Expected: FAIL because the package, validator, and service types do not exist.

- [ ] **Step 3: Create the model and migration**

Create the `notifications` table with `user_id`, `type`, `title`, `body`, nullable `url`, nullable `read_at`, timestamps, an index on `(user_id, read_at)`, and an index on `(user_id, created_at)`. Add the migration to `bootstrap.Migrations()` without running it.

- [ ] **Step 4: Implement the service**

Implement `Create`, `List`, `UnreadCount`, `MarkRead`, and `MarkAllRead`. Every query must include `user_id = ?`; `MarkRead` must update only `user_id = ? AND id = ?`, and a missing row must return a not-found error rather than disclose another user’s row.

- [ ] **Step 5: Run focused and package tests**

Run `go test ./app/services/notifications ./database/migrations ./bootstrap -count=1`.
Expected: PASS without executing the new migration.

- [ ] **Step 6: Commit the domain boundary**

```bash
git add backend/app/models/notification.go backend/database/migrations/20260922000005_create_notifications_table.go backend/bootstrap/migrations.go backend/app/services/notifications
git commit -m "feat: add notification persistence service"
```

### Task 2: Add authenticated notification API and OpenAPI contract

**Files:**
- Create: `backend/app/core/admin/controllers/notification_controller.go`
- Modify: `backend/routes/web.go`
- Modify: `backend/app/openapi/spec.go`
- Modify: `backend/app/openapi/spec_test.go`
- Modify: `backend/app/core/admin/controllers/notification_controller_test.go`

**Interfaces:**
- Routes: `GET /api/v1/notifications`, `GET /api/v1/notifications/unread-count`, `PUT /api/v1/notifications/{id}/read`, and `PUT /api/v1/notifications/read-all`.
- Controller responses: list envelope `{data, meta}`, unread count `{data:{count}}`, and successful mutations with `{data:...}`.

- [ ] **Step 1: Write failing controller contract tests**

Cover authenticated list pagination, unread count, mark-read, mark-all, invalid IDs, and a cross-user ID. Assert that cross-user access cannot mutate another user’s row.

- [ ] **Step 2: Run the controller tests and verify the contract fails**

Run `go test ./app/core/admin/controllers -run Notification -count=1`.
Expected: FAIL because the controller and routes do not exist.

- [ ] **Step 3: Implement the controller and routes**

Extract the authenticated user ID from the existing auth facade, call the notification service, normalize page/per-page with the existing list limits, and attach `RequireAuthentication()` middleware to every route. Do not add a sidebar permission or resource manifest.

- [ ] **Step 4: Extend OpenAPI and generated-client source**

Add `Notification`, `NotificationList`, `NotificationUnreadCount`, and request/response schemas to `spec.go`, add the four paths, and update `admin/scripts/generate-api-client.mjs` plus `admin/src/generated/api.ts` with typed notification methods.

- [ ] **Step 5: Run API and OpenAPI tests**

Run `go test ./app/core/admin/controllers ./app/openapi ./tests/feature -count=1` and `go build ./...`.
Expected: PASS with the new paths visible in the generated OpenAPI document.

- [ ] **Step 6: Commit the API boundary**

```bash
git add backend/app/core/admin/controllers backend/routes/web.go backend/app/openapi admin/scripts/generate-api-client.mjs admin/src/generated/api.ts
git commit -m "feat: expose authenticated notification api"
```

### Task 3: Add the frontend notification core and AdminShell entry point

**Files:**
- Create: `admin/src/core/notifications/useNotifications.ts`
- Create: `admin/src/core/notifications/NotificationMenu.vue`
- Modify: `admin/src/core/layouts/AdminShell.vue`
- Modify: `admin/src/locales/zh-CN/core.json`
- Modify: `admin/src/locales/en-US/core.json`

**Interfaces:**
- `useNotifications()` exposes `items`, `unreadCount`, `loading`, `error`, `refresh()`, `markRead(id)`, and `markAllRead()`.
- `NotificationMenu` consumes the composable and emits no business-specific events; navigation uses the router only after validating the API-provided internal path.

- [ ] **Step 1: Write the frontend behavior tests or pure helper tests**

Pin unread count refresh, mark-read removal from the unread count, empty state, and rejected external URL behavior. If the repository has no frontend test runner, extract URL validation and state transitions into a small pure module and test it with the existing TypeScript test boundary rather than introducing a new framework.

- [ ] **Step 2: Implement the composable**

Fetch the first notification page and unread count on mount, guard requests with the auth token, expose localized error state, and avoid polling or WebSocket behavior in this first version.

- [ ] **Step 3: Implement the menu**

Use existing dropdown/popover/button/badge primitives. Render plain text title/body, show unread styling, provide single and all-read actions, and show loading, empty, and failure states. Do not use `v-html`.

- [ ] **Step 4: Compose it into AdminShell**

Add a notification trigger in the desktop and compact header layouts without adding a sidebar menu item. Clicking an internal URL uses `router.push`; empty URLs do not navigate.

- [ ] **Step 5: Add bilingual copy and run frontend checks**

Run `npx vue-tsc -b` and `npm run build` from `admin/`. Expected: PASS in both locales.

- [ ] **Step 6: Commit the UI boundary**

```bash
git add admin/src/core/notifications admin/src/core/layouts/AdminShell.vue admin/src/locales
git commit -m "feat: add in-app notification center"
```

### Task 4: Add a safe business event and end-to-end verification

**Files:**
- Modify: `backend/app/services/users/user_service.go` or the existing user status mutation boundary identified during implementation
- Modify: `backend/app/services/notifications/notification_service.go`
- Create or modify: `backend/tests/feature/notifications_test.go`
- Modify: `docs/testing.md`
- Modify: `docs/roadmap.md`

**Interfaces:**
- A real existing business mutation creates a notification through `NotificationService.Create`; no controller writes the table directly.
- The feature test verifies create -> list -> unread count -> mark read -> zero unread.

- [ ] **Step 1: Choose the narrowest existing event**

Use the existing user status or role-assignment service boundary only if it can identify a target user without changing unrelated business behavior. If no safe event exists, use a test fixture inserted by the feature test and document that notification producers are ready for module integration; do not add a fake production endpoint.

- [ ] **Step 2: Add the feature test**

Run the notification API sequence for two users and assert that user A never sees or changes user B’s notification. Keep test data isolated and do not require a production migration to run automatically.

- [ ] **Step 3: Run browser verification after explicit migration approval**

After the new migration is reviewed and executed, start the isolated backend/frontend ports, log in, verify the header badge, open the notification menu, mark one notification and all notifications read, and verify an internal navigation target.

- [ ] **Step 4: Run the full verification suite**

Run `go test ./... -count=1`, `go build ./...`, `npx vue-tsc -b`, and `npm run build`. Run `git diff --check` and inspect the browser console for notification errors.

- [ ] **Step 5: Commit the acceptance milestone**

```bash
git add backend docs/testing.md docs/roadmap.md
git commit -m "test: verify in-app notification flow"
```

## Migration and Execution Gate

The migration file is generated and registered but is not executed as part of implementation. Before browser acceptance, review the SQL/schema change and explicitly approve running the migration against the Laragon-managed PostgreSQL instance.
