#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMP_ROOT="$(mktemp -d)"
trap 'rm -rf "$TEMP_ROOT"' EXIT

MODULE_ROOT="$TEMP_ROOT/modules" bash "$ROOT_DIR/scripts/make-module.sh" inventory
MODULE_DIR="$TEMP_ROOT/modules/inventory"

diff -u \
  <(printf '%s\n' \
    "README.md" \
    "acceptance.md" \
    "admin/locales/README.md" \
    "admin/module.ts" \
    "admin/package.json" \
    "admin/resources/README.md" \
    "admin/routes/README.md" \
    "backend/README.md" \
    "backend/go.mod" \
    "backend/module.go" \
    "backend/module_test.go" \
    "module.yaml") \
  <(cd "$MODULE_DIR" && find . -type f -printf '%P\n' | sort)

grep -Fx 'id: inventory' "$MODULE_DIR/module.yaml" >/dev/null
grep -Fx 'name: Inventory Module' "$MODULE_DIR/module.yaml" >/dev/null
grep -Fx 'runtime: compiled' "$MODULE_DIR/module.yaml" >/dev/null
grep -Fx 'dependencies: []' "$MODULE_DIR/module.yaml" >/dev/null

if MODULE_ROOT="$TEMP_ROOT/modules" bash "$ROOT_DIR/scripts/make-module.sh" inventory >/dev/null 2>&1; then
  echo "expected duplicate module generation to fail" >&2
  exit 1
fi

echo "make-module: snapshot passed"
