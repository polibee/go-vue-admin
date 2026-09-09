#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

fail() {
  echo "module-check: $*" >&2
  exit 1
}

run_repository_checks() {
  pnpm run lint
  pnpm run typecheck
  pnpm run test
  pnpm run build
  pnpm run openapi:generate

  if find backend -name '*.go' -print -quit | grep -q .; then
    (cd backend && go test ./...)
  else
    echo "module-check: backend tests skipped (no Go packages yet)"
  fi

  if git rev-parse --verify HEAD >/dev/null 2>&1; then
    if ! git diff --quiet HEAD -- contracts/openapi contracts/generated admin/src/generated; then
      fail "generated OpenAPI or SDK files are out of date; run pnpm openapi:generate and commit the result"
    fi
  else
    echo "module-check: generated drift check skipped (no Git HEAD yet)"
  fi
}

check_ui_boundaries() {
  local layer file spec resolved
  for layer in admin/src/components/ui admin/src/components/ai-elements; do
    [[ -d "$layer" ]] || continue
    if rg -n --glob '*.{ts,tsx,js,jsx,vue}' \
      '(@|~)/(app|core|modules|pages|resource-engine|stores|generated|services|router|navigation|permissions?|auth|components/(admin|ui-extensions))(/|$)' \
      "$layer"; then
      fail "$layer imports application code; upstream UI layers must remain dependency-free"
    fi

    while IFS= read -r file; do
      while IFS= read -r spec; do
        resolved="$(realpath -m "$(dirname "$file")/$spec")"
        case "$resolved" in
          "$ROOT_DIR/admin/src/components/ui"/*|"$ROOT_DIR/admin/src/components/ai-elements"/*)
            ;;
          *)
            fail "$file imports outside the upstream UI layers through relative path: $spec"
            ;;
        esac
      done < <(rg -o --no-filename '((\.\.?/)[^"'"'"'(),;[:space:]]+)' "$file" | sed -E 's/["'"'"'`].*$//')
    done < <(rg --files "$layer" -g '*.{ts,tsx,js,jsx,vue}')
  done
}

canonical_module_path() {
  local input="$1"
  local absolute
  if [[ "$input" = /* ]]; then
    absolute="$input"
  else
    absolute="$ROOT_DIR/$input"
  fi

  [[ -d "$absolute" ]] || fail "module directory does not exist: $input"
  absolute="$(cd "$absolute" && pwd)"
  case "$absolute" in
    "$ROOT_DIR"/modules/*) printf '%s\n' "${absolute#"$ROOT_DIR"/}" ;;
    *) fail "module path must be under modules/: $input" ;;
  esac
}

check_module_contract() {
  local module_path="$1"
  local module_dir="$ROOT_DIR/$module_path"
  local required
  for required in module.yaml README.md acceptance.md; do
    [[ -f "$module_dir/$required" ]] || fail "$module_path is missing required contract file: $required"
  done
}

run_module_checks() {
  local module_path="$1"
  local module_dir="$ROOT_DIR/$module_path"

  check_module_contract "$module_path"

  if [[ -f "$module_dir/backend/go.mod" ]]; then
    (cd "$module_dir/backend" && go test ./...)
  elif find "$module_dir/backend" -name '*.go' -print -quit 2>/dev/null | grep -q .; then
    fail "$module_path/backend contains Go files but no go.mod"
  else
    echo "module-check: $module_path backend tests skipped (no Go package yet)"
  fi

  if [[ -f "$module_dir/admin/package.json" ]]; then
    pnpm --dir "$module_dir/admin" run --if-present lint
    pnpm --dir "$module_dir/admin" run --if-present typecheck
    pnpm --dir "$module_dir/admin" run --if-present test
    pnpm --dir "$module_dir/admin" run --if-present build
  elif [[ -d "$module_dir/admin" ]]; then
    fail "$module_path/admin exists but has no package.json"
  else
    echo "module-check: $module_path admin checks skipped (no package yet)"
  fi
}

changed_modules() {
  local diff_args=()
  if [[ -n "${MODULE_CHECK_BASE:-}" ]]; then
    diff_args=("${MODULE_CHECK_BASE}" HEAD)
  elif git rev-parse --verify HEAD^ >/dev/null 2>&1; then
    diff_args=(HEAD^ HEAD)
  else
    git ls-files modules
    git diff --name-only
    git diff --cached --name-only
    return
  fi

  {
    git diff --name-only "${diff_args[@]}"
    git diff --name-only
    git diff --cached --name-only
  } | awk -F/ '$1 == "modules" && NF >= 2 { print $1 "/" $2 }' | sort -u
}

if [[ $# -gt 1 ]]; then
  fail "usage: $0 [modules/<module>]"
fi

if [[ -n "${MODULE_CHECK_BASE:-}" ]]; then
  if [[ "$MODULE_CHECK_BASE" =~ ^0+$ ]]; then
    unset MODULE_CHECK_BASE
  else
    git rev-parse --verify "${MODULE_CHECK_BASE}^{commit}" >/dev/null 2>&1 \
      || fail "MODULE_CHECK_BASE is not a commit: ${MODULE_CHECK_BASE}"
  fi
fi

run_repository_checks
check_ui_boundaries

if [[ $# -eq 1 ]]; then
  module_path="$(canonical_module_path "$1")" || exit $?
  run_module_checks "$module_path"
else
  mapfile -t modules < <(changed_modules)
  if [[ ${#modules[@]} -eq 0 ]]; then
    echo "module-check: no changed modules detected"
  else
    for module_path in "${modules[@]}"; do
      run_module_checks "$module_path"
    done
  fi
fi

echo "module-check: passed"
