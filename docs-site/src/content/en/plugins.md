# Build platform plugins

## Plugin responsibilities

Plugins extend platform capabilities such as payments, storage, notifications or observability. Their configuration belongs to the plugin configuration page, not inside the module catalog.

## Lifecycle

Declare the plugin id, version, permissions and dependencies. Enable and disable it through the plugin registry. Enabling mounts its routes, resources and navigation; disabling removes those runtime registrations.

## Sensitive configuration

API keys, payment secrets and private certificates must be stored server-side through a secret manager or protected environment. A browser page may edit a masked configuration form, but must never call a sensitive provider directly.
