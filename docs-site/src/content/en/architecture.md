# Architecture and boundaries

## Repository layout

- backend contains HTTP routes, domain services, repositories and migrations.
- admin contains the Vue application, shell, resource engine and generated client.
- modules contains business modules with manifests and registration entry points.
- plugins contains platform extensions with lifecycle and configuration boundaries.
- contracts contains OpenAPI and schema artifacts.
- docs-site is a standalone static developer documentation site.

## Request flow

The browser calls the generated API client. The backend authenticates the HttpOnly session, checks permissions, then calls a service and repository. The repository may be memory-backed for tests or GORM-backed for configured databases.

The application owns one GORM resource connection pool. Generated modules receive that connection and must not open or close their own pools.

## Extension rule

Modules and plugins declare navigation, routes and resources through runtime registration. A single owned page may be promoted to a menu item; multiple pages are grouped to avoid filling the primary navigation.
