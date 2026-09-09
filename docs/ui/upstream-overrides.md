# Upstream UI overrides

`admin/src/components/ui/` is reserved for source imported from the official `unovue/shadcn-vue` registry. `admin/src/components/ai-elements/` is reserved for source imported from the official AI Elements Vue registry. Both layers keep upstream visual language and stay free of application dependencies.

There are no overrides in the repository foundation.

## Recording an override

An override is exceptional. Before changing an upstream file, verify that composition, a documented prop, or a component in `admin/src/components/ui-extensions/` cannot solve the need. Add one entry per override using this template:

| Date | Layer/component | Upstream version or snapshot | Reason | Behavior impact | Upgrade/removal plan |
| --- | --- | --- | --- | --- | --- |
| YYYY-MM-DD | `ui/button` | version/commit | required functional gap | user-visible effect | how it will be reconciled with upstream |

Do not use this document to justify visual redesigns. Colors, typography, radius, shadows, spacing, accessibility behavior, and component states remain upstream unless an explicit product requirement says otherwise.
