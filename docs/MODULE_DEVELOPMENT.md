# Module development checklist

Each product capability belongs in `modules/<module>/` and must be independently understandable, testable, and deliverable.

```text
modules/<module>/
├── module.yaml       # id, version, dependencies, permissions
├── backend/          # Go routes, services, migrations, and tests
├── admin/            # routes, resources, locales, and tests
├── README.md         # capability, dependencies, and verification
└── acceptance.md     # executable acceptance checklist
```

`module.yaml`, `README.md`, and `acceptance.md` are mandatory even when a module has no backend or admin implementation yet. Module dependencies must be explicit in `module.yaml`; no module may reach into another module's private implementation.

## Delivery sequence

1. Define the module boundary, manifest, permissions, and acceptance criteria.
2. Write module tests before implementation and keep backend/admin tests under the module.
3. Regenerate OpenAPI and the SDK when a backend API changes; generated output is committed and never edited by hand.
4. Run `bash scripts/module-check.sh modules/<module>`.
5. Review the staged diff, then use `bash scripts/module-finish.sh modules/<module> "feat(scope): message"`.

The checker always runs repository lint, typecheck, test, build, OpenAPI generation, generated-SDK drift detection, and UI-boundary validation. It then validates the module contract and runs the module backend/admin checks when present. A missing contract, failing command, generated diff, or forbidden import fails the gate.

## UI checklist

Before adding UI, inspect `components/ui`, then `components/ai-elements`, then `components/admin`. Compose existing primitives before considering `components/ui-extensions`. Modules may consume UI primitives and admin components, but neither upstream UI directory may import module, page, provider, store, router, resource-engine, or permission code.

For every change to an upstream component, add an entry to `docs/ui/upstream-overrides.md` explaining the upstream version, reason, affected behavior, and upgrade plan. Application-specific behavior belongs outside the upstream layers.
