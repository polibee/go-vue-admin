#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"
RUN_DIR="$(mktemp -d /tmp/go-vue-admin-e2e.XXXXXX)"
backend_pid=''
admin_pid=''
cleanup() {
  [[ -z "$admin_pid" ]] || kill "$admin_pid" 2>/dev/null || true
  [[ -z "$backend_pid" ]] || kill "$backend_pid" 2>/dev/null || true
}
trap cleanup EXIT INT TERM
echo "E2E logs: $RUN_DIR"
(cd backend && go build -o "$RUN_DIR/backend" ./)
export APP_ENV=local APP_DEBUG=true APP_KEY=01234567890123456789012345678901
export APP_HOST=127.0.0.1 APP_PORT=3103 APP_URL=http://127.0.0.1:3103
export AUTH_BOOTSTRAP_EMAIL=admin@example.com AUTH_BOOTSTRAP_PASSWORD=test-only-password
export RESOURCE_PROVIDER=memory SESSION_DRIVER=file SESSION_HTTP_ONLY=true SESSION_SAME_SITE=lax
if curl -fsS --max-time 1 "$APP_URL/api/health" >/dev/null 2>&1; then
  echo 'E2E port 3103 is occupied; stop its test service first.' >&2
  exit 1
fi
(cd backend && exec "$RUN_DIR/backend") >"$RUN_DIR/backend.log" 2>&1 & backend_pid=$!
(cd admin && VITE_API_BASE_URL=http://127.0.0.1:5197 VITE_RESOURCE_PROVIDER=http VITE_BACKEND_URL="$APP_URL" exec ./node_modules/.bin/vite --host 127.0.0.1 --port 5197 --strictPort) >"$RUN_DIR/admin.log" 2>&1 & admin_pid=$!
ready=0
for attempt in $(seq 1 60); do
  kill -0 "$backend_pid" "$admin_pid"
  if curl -fsS --max-time 1 "$APP_URL/api/auth/bootstrap" >/dev/null 2>&1 && curl -fsS --max-time 1 http://127.0.0.1:5197/login >/dev/null 2>&1; then
    ready=1
    break
  fi
  sleep 1
done
[[ "$ready" == 1 ]] || { echo "E2E services failed to start; see $RUN_DIR" >&2; exit 1; }
ADMIN_BASE_URL=http://127.0.0.1:5197 node_modules/.bin/playwright test "$@"
