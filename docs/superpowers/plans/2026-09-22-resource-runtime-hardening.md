# Resource runtime hardening implementation plan

## Scope

Extend the existing Resource Manifest runtime without creating a second page system or adding dependencies. The batch Action registry remains the server-side execution boundary; relation options remain registry-resolved and permission-filtered; generic pages remain the default UI.

## Tasks

### 1. Unified Action Handler contract

- Add a resource-aware handler registration boundary around the existing `actions.Registry`.
- Validate Manifest `kind`, batch flag, payload contract, and handler availability before execution.
- Keep built-in delete, restore, status, and update actions on the same registry path.
- Add Go tests for duplicate registration, unknown handlers, payload mismatch, permission rejection, and successful execution.

### 2. Remote relation options

- Extend the relation options endpoint with search, page, per-page, and nullable selection semantics.
- Return bounded option results plus pagination metadata; include a selected value when editing.
- Preserve target resource permission, readable label, data-scope, and relation-write validation.
- Add backend contract tests and generated client types.

### 3. Generic complex-form validation

- Add field-level error mapping to the shared form component.
- Validate dependency visibility before serialization and omit hidden fields.
- Show group-aware validation errors and focus the first invalid field after a failed submit.
- Add frontend tests for dependencies, relation-load failures, field errors, and retry behavior.

### 4. Batch selection and Action regression coverage

- Extract/test selection request construction for current-page IDs and query-wide selection with exclusions.
- Test batch Action availability by permission and handler kind.
- Test result summaries and failure/skip rendering without changing the existing toolbar layout.

### 5. Real-resource browser acceptance

- Use `departments` as the generic resource baseline and a relation-enabled fixture already present in the repository where possible.
- Verify login, menu, list, filters, search, create, edit, detail, relation options, batch selection, batch Action, permission filtering, and error recovery.
- Run the browser checks against the existing backend/frontend processes only; do not execute migrations or alter Laragon service configuration.

## Verification

```text
go test ./... -count=1
go build ./...
npm run build
node --test tests/resource-actions.test.ts
git diff --check
```

The work ends with one local stage commit. GitHub push remains deferred until this milestone is independently reviewable.
