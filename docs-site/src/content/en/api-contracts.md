# API contracts and clients

## Contract source

The backend publishes an OpenAPI document in local development. The admin API docs page renders Scalar and separates contracts by functional area.

Generate contract artifacts from the repository root:

    pnpm run openapi:generate

The generated OpenAPI, schemas and admin TypeScript client are checked for drift. Pages should call the generated client instead of constructing fetch requests directly.

## External integrations

External systems should call the deployed backend API from their own server. A TypeScript SDK is only one consumer option; future generators can target PHP, Go, Python or other languages from the same OpenAPI contract.

Payment, signing, file-processing and other sensitive operations belong in the integrating system's backend, never in browser code.
