# Contributing

## Repository boundaries

The repository is a modular platform, not a shared application folder.

| Path | Responsibility | Must not contain |
| --- | --- | --- |
| `backend/` | Goravel platform core, shared infrastructure, HTTP/API entry points | CMS, trading, forum, or other product-specific domains |
| `admin/` | Vue admin shell, resource engine, generated API client, shared admin composition | direct API calls from product screens |
| `modules/<name>/` | independently deliverable product capability | changes that make another module a hidden dependency |
| `plugins/<name>/` | runtime plugin packages | platform-core business rules |
| `contracts/` | versioned OpenAPI inputs, schemas, and generated SDK artifacts | hand-edited generated output |

Dependency direction is `modules/plugins -> admin resource engine or backend platform -> contracts`. The backend produces OpenAPI; `contracts` produces the generated TypeScript client; frontend business code consumes that client through `ResourceDataProvider`. Product code must not call `fetch('/api/...')` or `axios` directly.

The admin UI layers are fixed: `admin/src/components/ui/` contains upstream shadcn-vue code; `admin/src/components/ai-elements/` contains upstream AI Elements Vue code; `admin/src/components/admin/` composes business-neutral admin behavior; and `admin/src/components/ui-extensions/` contains the rare, justified extension. Upstream layers must never import application, resource, provider, permission, router, store, or module code. Do not redesign upstream colors, radius, shadows, typography, spacing, or component states.

`docs/shadcn-vue-dev/` is a local reference copy and is intentionally ignored. Never add, modify, or delete it as part of repository work.

## Branches and commits

Keep `main` releasable. Develop one module or one coherent infrastructure change per branch and per commit. Use Conventional Commits, for example `feat(auth): add session login` or `chore(repo): establish modular platform foundation`. Do not mix unrelated visual cleanup, generated SDK drift, or multiple modules in one commit.

Never force-push, reset shared history, or overwrite a non-empty remote. The target GitHub repository is private: `polibee/go-vue-admin`.

## Required checks before push

Run these from the repository root:

```bash
pnpm run lint
pnpm run typecheck
pnpm run test
pnpm run build
pnpm run openapi:generate
git diff --exit-code -- contracts/openapi contracts/schemas admin/src/generated
bash scripts/module-check.sh modules/<module>
```

Use `bash scripts/module-finish.sh modules/<module> "feat(scope): message"` only after those checks pass. It performs an authentication gate before any mutation, runs the module gate, commits the selected module, then safely creates or reuses the private remote and pushes without force.

## GitHub Actions policy

The repository workflow is manual-only (`workflow_dispatch`). Run the checks
locally before pushing; normal pushes and pull requests must not start a GitHub
Actions job.
