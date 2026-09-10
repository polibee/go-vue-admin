#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

fail() {
  echo "release-check: $*" >&2
  exit 1
}

run_step() {
  local label="$1"
  shift
  echo "release-check: $label"
  "$@"
}

list_steps() {
  cat <<'EOF'
backend tests
admin checks
OpenAPI drift
module contract
UI compatibility
E2E: opt-in
EOF
}

if [[ "${RELEASE_CHECK_LIST:-0}" == "1" ]]; then
  list_steps
  exit 0
fi

run_step "backend tests" bash -c 'cd backend && go test ./...'

run_step "admin checks" bash -c 'pnpm run lint && pnpm run typecheck && pnpm run test && pnpm run build'

run_step "OpenAPI generation" go -C backend run ./cmd/admin-gen api --root ..
if ! git diff --quiet HEAD -- contracts/openapi contracts/schemas admin/src/generated; then
  fail "generated OpenAPI, schema, or SDK files changed; commit regenerated output before release"
fi

run_step "API audit" go -C backend run ./cmd/admin-gen audit --root ..

run_step "module contract" bash scripts/module-check.sh modules/example

run_step "generator contracts" go -C backend test ./app/core/generator ./cmd/admin-gen -count=1

run_step "UI compatibility" bash -c '
  if rg -n --glob "*.{ts,tsx,js,jsx,vue,css}" "element-plus|ant-design-vue|naive-ui|primevue|vuetify" admin/src plugins modules; then
    echo "release-check: a second UI framework was found" >&2
    exit 1
  fi
  if rg -n "^\\s+(pull_request|push):" .github/workflows/ci.yml; then
    echo "release-check: GitHub Actions must remain manual-only" >&2
    exit 1
  fi
'

if [[ "${RELEASE_CHECK_E2E:-0}" == "1" ]]; then
  run_step "E2E" pnpm run test:e2e
else
  echo "release-check: E2E: opt-in (set RELEASE_CHECK_E2E=1 with local services)"
fi

echo "release-check: passed"
