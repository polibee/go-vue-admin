#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Repository already initialized: $ROOT_DIR"
  exit 0
fi

git init
echo "Initialized local repository: $ROOT_DIR"
echo "Run bash scripts/module-finish.sh <module-path> <commit-message> after GitHub authentication is valid."
