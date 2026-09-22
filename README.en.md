# Go Vue Admin

A modular general-purpose admin platform built with Goravel, Vue 3, TypeScript, and shadcn-vue. Resource Manifests drive resource pages, permissions, menus, search, and batch actions.

[![RackNerd VPS](https://img.shields.io/badge/RackNerd-VPS-2563eb?style=for-the-badge)](https://my.racknerd.com/aff.php?aff=7572)
[![Vast.ai GPU Cloud](https://img.shields.io/badge/Vast.ai-GPU%20Cloud-7c3aed?style=for-the-badge)](https://cloud.vast.ai/?ref_id=91181)

## Features

- JWT login, refresh, logout, and PostgreSQL-authoritative refresh-token storage with Redis as an acceleration layer.
- Users, roles, permissions, data-scope rules, and field permissions.
- Resource Registry-driven list, create, edit, detail, and delete flows.
- Search, filters, sorting, pagination, and CSV export.
- Current-page, cross-page, and filtered-set selection, plus bulk delete, bulk update, and Manifest Actions.
- Resource relations, grouped forms, and field-dependency extension points.
- Global search, menu grouping, and permission filtering.
- Request/response audit logging, sensitive-field redaction, manual cleanup, and scheduled cleanup.
- admin:make-resource generates the backend Resource, migration, permissions, menu, frontend ResourceSpec, tests, and README; ordinary resources reuse generic pages by default and complex resources may explicitly override them.

## Local development

Start PostgreSQL and Redis through Laragon, then run the backend:

~~~powershell
cd backend
go run .
~~~

In another terminal, start the admin panel:

~~~powershell
cd admin
pnpm install
pnpm dev
~~~

Port contract:

- Go API: http://127.0.0.1:3000
- Go/Vue admin panel: http://127.0.0.1:5180
- Python admin panel (separate project): http://127.0.0.1:5173

Review generated migrations before applying them. The generator does not silently execute migrations.

## Resource generation

~~~powershell
cd backend
go run . admin:make-resource announcements --fields="title:text:required,status:select:required:draft=Draft|published=Published"
~~~

Resource generation is the main business-development workflow: generated resources are discovered by the admin panel automatically. The generator avoids overwriting manual files. Complex business flows belong in backend/app/modules/<module> and admin/src/modules/<module> rather than being forced into a generic Resource.

## Verification

~~~powershell
cd backend
go test ./...

cd ..\admin
pnpm exec vue-tsc --noEmit
pnpm run build
~~~

## Directory conventions

- Backend foundation: backend/app/core
- Backend business modules: backend/app/modules/<module>
- Backend services: backend/app/services/<domain>
- Backend console commands: backend/app/console/<domain>
- Frontend shared components: admin/src/components
- Frontend infrastructure: admin/src/core
- Frontend business pages and components: admin/src/modules/<module>

Before production deployment, review secrets, database, Redis, reverse proxy, retention, backups, monitoring, and migration procedures for the target environment. This README describes the development baseline; it is not a substitute for environment acceptance.
