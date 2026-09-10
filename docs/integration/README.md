# Public integration documentation

This directory is the future public-facing guide for consumers of the Go Vue
Admin platform API. It is intentionally separate from `admin/`: the admin app
contains an authenticated, internal Scalar viewer, while this directory will
describe the stable external contract.

The machine-readable source is the versioned OpenAPI document under
`contracts/openapi`. Do not document endpoints here that are not present in
that contract.

Planned sections:

- authentication and session/API-key boundaries;
- response envelopes, errors, pagination, filtering, and rate limits;
- webhook and event delivery when a business product introduces them;
- versioning, deprecation, changelog, and language-specific SDK guidance.

Payment, signing, secret management, and WordPress adapter behavior belong to
future business-product documentation, not the generic admin platform.
