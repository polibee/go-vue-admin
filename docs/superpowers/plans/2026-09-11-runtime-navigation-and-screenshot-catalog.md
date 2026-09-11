# Runtime Navigation and Screenshot Catalog Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Connect module/plugin declarations to runtime navigation, routes, and resources with safe enable/disable cleanup, then replace the public screenshot catalog with complete high-resolution feature screenshots.

**Architecture:** A frontend `RuntimeRegistrationBridge` owns all dynamic registrations. It mounts module definitions at bootstrap and mounts/unmounts plugin registrations as the plugin management page changes state. Every dynamic navigation item, route, and resource is tagged by owner so teardown is deterministic; backend authorization remains authoritative.

**Tech Stack:** Vue 3, Vue Router 4, TypeScript, shadcn-vue, Vitest, Playwright, Go-generated API client, PNG screenshots at deviceScaleFactor 2.

**Spec:** `docs/superpowers/specs/2026-09-11-runtime-registered-navigation-and-screenshot-catalog-design.md`

## Global Constraints

- Keep `/api/extensions` and `/admin/extensions` compatibility behavior.
- Keep `/api/modules` and `/api/plugins` as management APIs.
- Do not implement network plugin download, external code execution, marketplace loading, or signature verification.
- Frontend hiding is UX only; backend permissions and plugin state remain authoritative.
- Do not commit login screenshots, cookies, storage state, passwords, or tokens.
- Use the existing shadcn-vue components and semantic design tokens.
- Use TDD: each runtime behavior gets a failing test before production code.
- Do not enable GitHub Actions or automatic CI.

## Configuration boundary

- Platform management pages expose lifecycle, dependency, permission, health, and configuration status.
- A module/plugin manifest may declare a dedicated configuration route and schema.
- Domain-specific forms and secret persistence stay in the registered module/plugin page and backend API.
- Secret values are never returned by catalog APIs; this phase only exposes masked-field metadata.

The configuration metadata and entry point are implemented before the runtime navigation bridge so later dynamic routes can mount plugin-owned configuration pages consistently.

## Navigation grouping

- Keep the registry flat and apply grouping only when the sidebar renders.
- Promote one business entry to a first-level menu.
- Group multiple entries from the same owner into a collapsible submenu.
- Preserve separate `业务模块` and `平台插件` management entries.

## File Map

- Create `admin/src/core/extensions/RuntimeRegistrationBridge.ts`: owner-aware dynamic navigation, route, and resource mounting.
- Modify `admin/src/core/navigation/NavigationRegistry.ts`: owner metadata and owner removal.
- Modify `admin/src/core/extensions/ModuleRegistry.ts`: module registration snapshots and idempotent access used by the bridge.
- Modify `admin/src/core/extensions/plugin/PluginRegistry.ts`: expose lifecycle snapshots and make teardown state observable.
- Modify `admin/src/core/extensions/runtime.ts`: initialize the bridge, mount modules, and synchronize plugin state.
- Modify `admin/src/core/router/index.ts`: name the admin parent route and support runtime child routes.
- Modify `admin/src/core/extensions/PluginsPage.vue`: synchronize backend plugin state with the local runtime bridge.
- Create `admin/src/core/extensions/RuntimeRegistrationBridge.test.ts`: unit coverage for mount, unmount, duplicate protection, and owner cleanup.
- Modify `tests/e2e/extensions.spec.ts`: verify runtime business menus appear and disappear with plugin lifecycle.
- Modify `README.md`: document runtime registration and canonical screenshot catalog.
- Replace `output/playwright/*.png`: canonical high-resolution feature screenshots only.

### Task 1: Add owner-aware navigation cleanup

**Files:**
- Modify: `admin/src/core/navigation/NavigationRegistry.ts`
- Create: `admin/src/core/extensions/RuntimeRegistrationBridge.test.ts`

**Interfaces:**
- `NavigationItem` gains optional `owner?: string`.
- `NavigationRegistry.removeOwner(owner: string): void` removes only dynamic items owned by that module/plugin.

- [ ] **Step 1: Write the failing test**

```ts
it('removes only navigation entries owned by a plugin', () => {
  const registry = new NavigationRegistry([{ id: 'dashboard', label: '仪表盘', route: '/admin/dashboard' }])
  registry.register({ id: 'plugin-example', label: '示例插件', route: '/admin/example-plugin', owner: 'plugin:example-plugin' })
  registry.removeOwner('plugin:example-plugin')
  expect(registry.all().map((item) => item.id)).toEqual(['dashboard'])
})
```

- [ ] **Step 2: Run the focused test and confirm it fails**

Run: `pnpm --dir admin exec vitest run src/core/extensions/RuntimeRegistrationBridge.test.ts`

Expected: FAIL because `removeOwner` and `owner` do not yet exist.

- [ ] **Step 3: Implement the smallest registry change**

Store `owner` through the existing merge behavior and filter entries whose `owner` matches the requested owner. Core static entries have no owner and must remain untouched.

- [ ] **Step 4: Run the focused test**

Run: `pnpm --dir admin exec vitest run src/core/extensions/RuntimeRegistrationBridge.test.ts`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add admin/src/core/navigation/NavigationRegistry.ts admin/src/core/extensions/RuntimeRegistrationBridge.test.ts
git commit -m "feat(admin): add owner-aware navigation cleanup"
```

### Task 2: Implement the runtime registration bridge

**Files:**
- Create: `admin/src/core/extensions/RuntimeRegistrationBridge.ts`
- Modify: `admin/src/core/router/index.ts`
- Modify: `admin/src/resource-engine/core/ResourceRegistry.ts` only if owner removal cannot be implemented without breaking existing `remove(name)` behavior.
- Modify: `admin/src/core/extensions/RuntimeRegistrationBridge.test.ts`

**Interfaces:**
- `RuntimeRegistrationBridge.mountModule(definition: AdminModuleDefinition): void`
- `RuntimeRegistrationBridge.mountPlugin(id: string, snapshot: PluginRegistrationSnapshot): void`
- `RuntimeRegistrationBridge.unmount(owner: string): void`
- Constructor: `new RuntimeRegistrationBridge({ router, navigationRegistry, resourceRegistry })`.
- Dynamic route owner format: `module:<id>` or `plugin:<id>`.
- Dynamic route names must be unique and removable with `router.removeRoute(name)`.

- [ ] **Step 1: Add failing bridge tests**

```ts
it('mounts a module navigation item and route once', () => {
  const bridge = createBridgeWithTestRouter()
  const definition = { id: 'example', routes: [{ path: 'example', name: 'module-example', component: TestView }], navigation: [{ id: 'module-example', label: '示例模块', route: '/admin/example' }] }
  bridge.mountModule(definition)
  bridge.mountModule(definition)
  expect(navigationRegistry.all().filter((item) => item.id === 'module-example')).toHaveLength(1)
  expect(router.hasRoute('module:example:module-example')).toBe(true)
})

it('unmounts a plugin navigation item, route, and resource together', () => {
  const bridge = createBridgeWithTestRouter()
  bridge.mountPlugin('example-plugin', { menus: [{ id: 'plugin-example', label: '示例插件', route: '/admin/example-plugin' }], routes: [{ id: 'plugin-example', path: 'example-plugin', view: TestView }], resources: [] })
  bridge.unmount('plugin:example-plugin')
  expect(navigationRegistry.all().some((item) => item.id === 'plugin-example')).toBe(false)
  expect(router.hasRoute('plugin:example-plugin:plugin-example')).toBe(false)
})
```

The test file must define `createBridgeWithTestRouter()` with `createRouter(createMemoryHistory(), [{ path: '/admin', name: 'admin' }])`, a fresh `NavigationRegistry`, and a fresh `ResourceRegistry`, then return `{ bridge, router, navigationRegistry, resourceRegistry }`.

- [ ] **Step 2: Run tests and verify the bridge tests fail**

Run: `pnpm --dir admin exec vitest run src/core/extensions/RuntimeRegistrationBridge.test.ts`

Expected: FAIL because the bridge and owner-aware route/resource behavior are absent.

- [ ] **Step 3: Implement `RuntimeRegistrationBridge`**

For modules, convert `navigation` to owned `NavigationItem` entries, add each `RouteRecordRaw` below the named `admin` route, and register resources through the existing provider factory. For plugins, convert `PluginMenu` and `PluginRoute` into owned navigation/route entries and register plugin resources. Before mounting an owner, call `unmount(owner)` so repeated lifecycle calls cannot duplicate registrations.

- [ ] **Step 4: Make the admin parent route addressable**

Set the `/admin` route record name to `admin`; dynamic children must be added with `router.addRoute('admin', childRoute)` and removed with `router.removeRoute(childName)`.

- [ ] **Step 5: Add owner-aware resource cleanup**

Track resource owner names inside the bridge. On unmount, call `resourceRegistry.remove(resource.name)` only for resources owned by that module/plugin. Never remove core resources or resources owned by another extension.

- [ ] **Step 6: Run focused tests and refactor only after green**

Run: `pnpm --dir admin exec vitest run src/core/extensions/RuntimeRegistrationBridge.test.ts`

Expected: PASS with no duplicate entries and complete owner cleanup.

- [ ] **Step 7: Commit**

```bash
git add admin/src/core/extensions/RuntimeRegistrationBridge.ts admin/src/core/extensions/RuntimeRegistrationBridge.test.ts admin/src/core/router/index.ts admin/src/resource-engine/core/ResourceRegistry.ts
git commit -m "feat(admin): add runtime registration bridge"
```

### Task 3: Connect module bootstrap and plugin lifecycle

**Files:**
- Modify: `admin/src/core/extensions/runtime.ts`
- Modify: `admin/src/core/extensions/plugin/PluginRegistry.ts` only for lifecycle snapshot access needed by the bridge.
- Modify: `admin/src/core/extensions/PluginsPage.vue`
- Modify: `admin/src/app/bootstrap.ts` if bridge initialization needs to happen before app mount.
- Modify: `admin/src/core/extensions/RuntimeRegistrationBridge.test.ts`

**Interfaces:**
- `registerBuiltinExtensions()` mounts the example module once and creates the bridge singleton.
- Plugin state synchronization maps API states `enabled|disabled` to local `pluginRuntime.enable|disable` and bridge mount/unmount.

- [ ] **Step 1: Write failing lifecycle tests**

```ts
it('enabling a plugin mounts its declared menu and disabling it removes the menu', () => {
  const builtinManifest = { id: 'example-plugin', name: '示例插件', version: '1.0.0', runtime: 'builtin' as const, uiCompatibility: 'shadcn-vue' as const }
  const menu = { id: 'plugin-example', label: '示例插件', route: '/admin/example-plugin' }
  const plugin = defineAdminPlugin({ manifest: builtinManifest, setup: (context) => context.registerMenu(menu) })
  const runtime = new PluginRegistry()
  runtime.register(plugin)
  const bridge = createBridgeWithTestRouter()
  runtime.enable('example-plugin')
  bridge.mountPlugin('example-plugin', runtime.registrations('example-plugin')!)
  expect(navigationRegistry.all().some((item) => item.id === 'plugin-example')).toBe(true)
  runtime.disable('example-plugin')
  bridge.unmount('plugin:example-plugin')
  expect(navigationRegistry.all().some((item) => item.id === 'plugin-example')).toBe(false)
})
```

- [ ] **Step 2: Run the test and confirm failure**

Run: `pnpm --dir admin exec vitest run src/core/extensions/RuntimeRegistrationBridge.test.ts`

Expected: FAIL because runtime bootstrap and management-page synchronization are not connected.

- [ ] **Step 3: Mount module definitions during bootstrap**

Replace the current resource-only `registerModule` behavior with bridge mounting, while preserving the existing resource-provider selection. The example module’s `navigation` entry must appear after login without manually adding another item to the static navigation list.

- [ ] **Step 4: Synchronize plugin management state**

After `listPlugins()` returns, for every plugin whose backend state is `enabled`, call `pluginRuntime.enable(id)` and mount its registration snapshot. When the user toggles state, call the corresponding local lifecycle method only after the API succeeds; then mount or unmount the bridge owner.

- [ ] **Step 5: Verify tests and page typecheck**

Run: `pnpm --dir admin exec vitest run src/core/extensions/RuntimeRegistrationBridge.test.ts` and `pnpm --dir admin run typecheck`.

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add admin/src/core/extensions/runtime.ts admin/src/core/extensions/plugin/PluginRegistry.ts admin/src/core/extensions/PluginsPage.vue admin/src/app/bootstrap.ts admin/src/core/extensions/RuntimeRegistrationBridge.test.ts
git commit -m "feat(admin): connect extension lifecycle to runtime navigation"
```

### Task 4: Add browser regression coverage for dynamic business entries

**Files:**
- Modify: `tests/e2e/extensions.spec.ts`
- Modify: `tests/e2e/rbac.spec.ts` if dynamic permission visibility belongs in the existing RBAC flow.

- [ ] **Step 1: Add failing browser assertions**

After login, assert the module’s declared menu `示例模块` exists and opens `/admin/example`. On `/admin/plugins`, enable the example plugin, assert `示例插件` appears in the sidebar, open `/admin/example-plugin`, then disable the plugin and assert the menu disappears and direct navigation is rejected or redirected.

- [ ] **Step 2: Run the focused E2E and confirm failure**

Run: `pnpm exec playwright test tests/e2e/extensions.spec.ts --reporter=line`

Expected: FAIL before the bridge is connected because declared business menus are not present.

- [ ] **Step 3: Update the E2E flow after the bridge is implemented**

Use stable accessible labels and explicit routes. Do not assert on internal registry implementation details.

- [ ] **Step 4: Run the focused E2E**

Run: `pnpm exec playwright test tests/e2e/extensions.spec.ts --reporter=line`

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add tests/e2e/extensions.spec.ts tests/e2e/rbac.spec.ts
git commit -m "test(e2e): verify runtime extension navigation"
```

### Task 5: Replace the screenshot catalog and README references

**Files:**
- Delete only obsolete screenshot files under `output/playwright/` after enumerating them with `Get-ChildItem`.
- Create: `output/playwright/users.png`
- Create: `output/playwright/permissions.png`
- Create: `output/playwright/media.png`
- Create: `output/playwright/api-docs.png`
- Replace: `output/playwright/dashboard.png`, `modules.png`, `plugins.png`, `roles.png`, `settings.png`, `audit.png`
- Modify: `README.md`

- [ ] **Step 1: Write the screenshot inventory check**

Add a small local verification command or test assertion that the README references exactly these ten files:

```text
dashboard.png modules.png plugins.png users.png roles.png
permissions.png settings.png media.png audit.png api-docs.png
```

- [ ] **Step 2: Run the inventory check and confirm it fails**

Run: `rg -o "output/playwright/[a-z-]+\.png" README.md | Sort-Object -Unique` and compare it with the canonical list. Expected: FAIL because the four new feature screenshots are absent.

- [ ] **Step 3: Capture screenshots with a temporary authenticated browser context**

Use a single Playwright context with viewport `{ width: 1440, height: 900 }` and `deviceScaleFactor: 2`. Capture each canonical route after `networkidle`; never save storage state or login screenshots. Include the dynamic module and enabled-plugin business pages if they are part of the README narrative, but do not add extra filenames outside the canonical list.

- [ ] **Step 4: Remove obsolete files and update README**

Delete only files outside the canonical list, update image captions, and document that screenshots are generated at 2× pixel density. Do not delete source code or test artifacts.

- [ ] **Step 5: Verify dimensions and references**

Run a read-only image-dimension check. Every screenshot must have width `2880`; normal pages should have height `1800` or their full-page height. Confirm README has no references to deleted files and no login credential image.

- [ ] **Step 6: Commit**

```bash
git add README.md output/playwright
git commit -m "docs: refresh high resolution feature screenshots"
```

### Task 6: Full verification and release handoff

**Files:**
- No production files; inspect the complete diff and generated artifacts.

- [ ] **Step 1: Run backend verification**

Run: `go -C backend test ./...` and `go -C backend vet ./...`.

- [ ] **Step 2: Run frontend verification**

Run: `pnpm --dir admin run typecheck`, `pnpm --dir admin run test`, and `pnpm --dir admin run build`.

- [ ] **Step 3: Run browser verification**

Run: `pnpm exec playwright test tests/e2e --reporter=line`.

- [ ] **Step 4: Check the final diff**

Run: `git diff --check`, `git status --short`, and `rg -n "TBD|/admin/extensions" README.md tests/e2e admin/src/core/extensions`.

Expected: no whitespace errors, only intentional compatibility references to `/admin/extensions`, and no untracked credential artifacts.

- [ ] **Step 5: Push the completed module**

```bash
git push origin main
```

Do not enable GitHub Actions or add automatic CI triggers.
