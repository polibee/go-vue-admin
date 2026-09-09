#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

fail() {
  echo "module-finish: $*" >&2
  exit 1
}

[[ $# -eq 2 ]] || fail "usage: $0 <modules/<module>> <commit-message>"
MODULE_PATH="$1"
COMMIT_MESSAGE="$2"

if ! gh auth status; then
  fail "GitHub authentication is unavailable. Repair gh auth login before this command can check, create, or push a repository."
fi

bash scripts/module-check.sh "$MODULE_PATH"

git diff --quiet --cached || fail "staged changes already exist; review or commit them before module-finish"
git add -- "$MODULE_PATH"
git diff --cached --quiet && fail "no changes staged for $MODULE_PATH"
git commit -m "$COMMIT_MESSAGE"

branch="$(git branch --show-current)"
[[ -n "$branch" ]] || fail "current branch is detached; switch to a branch before pushing"

if gh repo view polibee/go-vue-admin --json nameWithOwner >/dev/null 2>&1; then
  if git remote get-url origin >/dev/null 2>&1; then
    :
  elif [[ -n "$(git remote)" ]]; then
    fail "a local remote already exists but origin is absent; refusing to change remote configuration"
  else
    git remote add origin https://github.com/polibee/go-vue-admin.git
  fi

  if [[ -n "$(git ls-remote --heads origin)" ]]; then
    fail "origin contains branches; refusing to overwrite a non-empty remote"
  fi
  git push -u origin "$branch"
else
  [[ -z "$(git remote)" ]] || fail "a local remote already exists; refusing to create or replace origin"
  gh repo create polibee/go-vue-admin --private --source . --remote origin --push
fi
