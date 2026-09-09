#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OPENAPI_FILE="$ROOT_DIR/contracts/openapi/openapi.json"
SCHEMA_DIR="$ROOT_DIR/contracts/schemas"
GENERATED_DIR="$ROOT_DIR/admin/src/generated/api"

mkdir -p "$ROOT_DIR/contracts/openapi" "$SCHEMA_DIR"
(cd "$ROOT_DIR/backend" && go run ./cmd/openapi -output "$OPENAPI_FILE" -schema-dir "$SCHEMA_DIR")
node "$ROOT_DIR/scripts/generate-ts-client.mjs" "$OPENAPI_FILE" "$GENERATED_DIR"
