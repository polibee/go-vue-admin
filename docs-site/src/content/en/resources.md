# Resources and code generation

## Manifest first

A Resource Manifest is the intermediate structure shared by introspection, code generation and runtime registration. Keep the manifest stable and review it before generating files.

## Single-table CRUD

The Go admin-gen CLI can inspect a configured database table and generate a GORM repository, controller, service, route registration and admin resource definition. Generated code uses the application-owned GORM connection.

Use the generated API client for pages. Do not edit generated files directly; regenerate them when the OpenAPI contract changes.

## Planned extensions

Foreign keys, enums, soft deletes, audit events, batch generation and schema-difference synchronization should be added behind explicit manifest fields and compatibility tests.
