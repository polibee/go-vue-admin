.DEFAULT_GOAL := check

.PHONY: lint typecheck test build openapi-generate check module-check

lint:
	pnpm run lint

typecheck:
	pnpm run typecheck

test:
	pnpm run test
	@if find backend -name '*.go' -print -quit | grep -q .; then \
		(cd backend && go test ./...); \
	else \
		echo "backend: test skipped (no Go packages yet)"; \
	fi

build:
	pnpm run build

openapi-generate:
	pnpm run openapi:generate

module-check:
	bash scripts/module-check.sh

check: lint typecheck test build openapi-generate module-check
