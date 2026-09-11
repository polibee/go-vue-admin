# Contributing and releases

## Change boundaries

Keep backend, admin, modules, plugins, contracts and docs-site changes scoped to their responsibility. Update the manifest and OpenAPI contract before changing generated output.

## Verification

Run focused unit tests first, then backend tests, frontend type checks, production build and relevant Playwright flows. Treat browser data as disposable and isolate fixtures between tests.

## Documentation publishing

Build the standalone site with:

    node docs-site/scripts/build.mjs

The Developer Docs Pages workflow is manual-only. It publishes docs-site/dist and does not run as part of normal application deployment.
