#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

fail() {
  echo "module-finish: $*" >&2
  exit 1
}

EXPECTED_ORIGIN="https://github.com/polibee/go-vue-admin.git"

is_expected_origin() {
  case "$1" in
    https://github.com/polibee/go-vue-admin|https://github.com/polibee/go-vue-admin.git|git@github.com:polibee/go-vue-admin.git|ssh://git@github.com/polibee/go-vue-admin.git)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

validate_existing_origin() {
  local remote_url push_url refs
  if ! git remote get-url origin >/dev/null 2>&1; then
    return 0
  fi

  remote_url="$(git remote get-url origin)" || fail "cannot read origin URL"
  is_expected_origin "$remote_url" || fail "origin points to an unexpected repository: $remote_url"

  if push_url="$(git remote get-url --push origin 2>/dev/null)"; then
    is_expected_origin "$push_url" || fail "origin push URL points to an unexpected repository: $push_url"
  fi

  refs="$(git ls-remote --refs origin 2>/dev/null)" || fail "cannot inspect origin; refusing to push"
  [[ -z "$refs" ]] || fail "origin contains refs; refusing to overwrite a non-empty remote"
}

[[ $# -eq 2 ]] || fail "usage: $0 <modules/<module>> <commit-message>"
MODULE_PATH="$1"
COMMIT_MESSAGE="$2"

if ! gh auth status; then
  fail "GitHub authentication is unavailable. Repair gh auth login before this command can check, create, or push a repository."
fi

# Reject unsafe local remotes before creating a local commit.
validate_existing_origin

bash scripts/module-check.sh "$MODULE_PATH"

git diff --quiet --cached || fail "staged changes already exist; review or commit them before module-finish"
git add -- "$MODULE_PATH"
git diff --cached --quiet && fail "no changes staged for $MODULE_PATH"
git commit -m "$COMMIT_MESSAGE"

branch="$(git branch --show-current)"
[[ -n "$branch" ]] || fail "current branch is detached; switch to a branch before pushing"

repo_probe_error="$(mktemp)"
trap 'rm -f "$repo_probe_error"' EXIT
if gh api repos/polibee/go-vue-admin --silent >/dev/null 2>"$repo_probe_error"; then
  if git remote get-url origin >/dev/null 2>&1; then
    :
  elif [[ -n "$(git remote)" ]]; then
    fail "a local remote already exists but origin is absent; refusing to change remote configuration"
  else
    git remote add origin "$EXPECTED_ORIGIN"
  fi

  validate_existing_origin
  git push -u origin "$branch"
else
  if ! rg -qi '404|not found|could not resolve to a repository' "$repo_probe_error"; then
    cat "$repo_probe_error" >&2
    fail "unable to verify whether polibee/go-vue-admin exists"
  fi
  [[ -z "$(git remote)" ]] || fail "a local remote already exists; refusing to create or replace origin"
  gh repo create polibee/go-vue-admin --private --source . --remote origin --push
fi
