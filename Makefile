.PHONY: dev build test check generate migrate up down clean build-cli run-cli install-cli build-prod

# --- Build ---
build: build-api build-cli

build-api:
	@mkdir -p tmp
	cd apps/core && go build -o ../../tmp/server ./cmd/server

build-cli:
	@mkdir -p tmp
	cd apps/core && go build -o ../../tmp/smokery ./cmd/smokery

# Production build: embeds the frontend static site into the server binary.
# Requires: cd apps/web && bun run build (outputs to apps/core/internal/frontend/dist/)
build-prod:
	@mkdir -p tmp
	cd apps/web && bun run build
	cp -r apps/web/build/* apps/core/internal/frontend/dist/
	cd apps/core && go build -tags embed_frontend -o ../../tmp/server ./cmd/server

run-cli: build-cli
	./tmp/smokery $(ARGS)

install-cli: build-cli
	cp tmp/smokery $${GOBIN:-$(HOME)/go/bin}/smokery

# --- Dev ---
dev:
	@mkdir -p tmp
	cd configs && air -c .air.toml

# --- Generate ---
generate: generate-openapi generate-types

generate-openapi:
	cd apps/core && go run ./cmd/openapi/ ../web/openapi.json

generate-types:
	cd apps/web && bun run generate

generate-sqlc:
	cd apps/core && sqlc generate

# --- Test & Check ---
test:
	cd apps/core && go test ./...

test-integration:
	cd apps/core && go test -tags integration ./internal/adapter/postgres/ -v -count=1

check:
	cd apps/web && bun run check

lint: test check

lint-go:
	cd apps/core && golangci-lint run ./...

# --- Database ---
up:
	docker compose -f configs/docker-compose.yml up -d

down:
	docker compose -f configs/docker-compose.yml down

migrate:
	cd apps/core && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path db/migrations -database "$$DATABASE_URL" up

# --- Install ---
install:
	cd apps/core && go mod tidy
	bun install

# --- Clean ---
clean:
	rm -rf tmp/
	rm -rf apps/web/.svelte-kit
