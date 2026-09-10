#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
output="$(RELEASE_CHECK_LIST=1 bash "$ROOT_DIR/scripts/release-check.sh")"

grep -F "backend tests" <<<"$output" >/dev/null
grep -F "admin checks" <<<"$output" >/dev/null
grep -F "OpenAPI drift" <<<"$output" >/dev/null
grep -F "module contract" <<<"$output" >/dev/null
grep -F "E2E: opt-in" <<<"$output" >/dev/null

if grep -Eq 'github|gh run|workflow_dispatch' <<<"$output"; then
  echo "release-check must not invoke GitHub CI" >&2
  exit 1
fi

echo "release-check: dry-run passed"
