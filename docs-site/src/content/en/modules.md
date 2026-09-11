# Build business modules

## Module responsibilities

A module owns business domain code, manifest metadata, resources, routes and navigation. It should not depend on admin shell internals or create a second database pool.

## Registration checklist

- Give the module a stable id and version.
- Declare permissions and dependencies.
- Register routes and resources through the runtime bridge.
- Provide translated labels and descriptions.
- Add backend, frontend and browser acceptance tests.
- Remove all registrations when the module is unloaded.

## Menu behavior

Business pages should use a module group when there are multiple entries. A single page can be promoted automatically. The module management page remains the operational entry point; business pages belong in the generated navigation.
