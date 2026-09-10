# Builtin Plugin SDK

Task 13 establishes the v1 plugin contract for statically compiled, trusted
plugins. It does not download packages, execute external processes, load
arbitrary JavaScript, verify signatures, or provide a marketplace.

## Backend contract

Backend plugins implement `plugin.BuiltinPlugin`:

```go
type BuiltinPlugin interface {
    Manifest() PluginManifest
    Enable() error
    Disable() error
}
```

Register plugins in `plugin.Registry`. Every plugin starts as `disabled` and
is listed in registration order. `Enable` invokes the plugin lifecycle once;
`Disable` invokes it once and is rejected while another enabled plugin depends
on it. Dependencies must already be enabled before a dependent plugin can be
enabled.

The v1 runtime only accepts:

- `runtime: "builtin"`
- `uiCompatibility: "shadcn-vue"`
- stable lowercase plugin IDs matching `[a-z][a-z0-9_-]*`

Manifest permissions, menus, and dependency versions are declarative contract
data. Version constraint resolution is intentionally deferred until a package
runtime is introduced.

## Admin SDK

Admin plugins use `defineAdminPlugin` and `PluginRegistry` from
`admin/src/core/extensions/plugin`. The SDK exposes only these registration
operations through `PluginContext`:

- `registerMenu(menu)`
- `registerRoute(route)`
- `registerResource(resource)`

The SDK does not expose Pinia stores, the Vue Router instance, `AdminShell`, or
layout internals. Plugins must use the existing shadcn-vue component
foundation and must not introduce another UI framework or global CSS rules.

Example:

```ts
const plugin = defineAdminPlugin({
  manifest: {
    id: 'example',
    name: 'Example Plugin',
    version: '0.1.0',
    runtime: 'builtin',
    uiCompatibility: 'shadcn-vue',
    permissions: ['example.view'],
    menus: [],
    dependencies: [],
  },
  setup(context) {
    context.registerMenu({
      id: 'example',
      label: 'Example',
      route: '/admin/example',
      permission: 'example.view',
    })
  },
})
```

Enable/disable is deterministic and dependency-aware. Dynamic loading and
external runtime support remain post-v1 work.
