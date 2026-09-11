# Database and production configuration

## Provider selection

Set RESOURCE_PROVIDER=database to use real resource storage. Set DB_CONNECTION=mysql or postgres to select the GORM driver. RESOURCE_PROVIDER=mysql remains a compatibility alias.

The application owns one GORM pool and closes it during shutdown. Modules receive the shared pool and must not create independent connections.

## Example

    RESOURCE_PROVIDER=database
    DB_CONNECTION=postgres
    DB_HOST=127.0.0.1
    DB_PORT=5432
    DB_DATABASE=go_vue_admin
    DB_USERNAME=postgres
    DB_PASSWORD=

Run migrations only after confirming the target database and backup policy. Integration tests are opt-in and must not silently modify a developer database.
