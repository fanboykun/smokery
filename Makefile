.PHONY: dev build test check generate migrate up down clean

# --- Build ---
build: build-api

build-api:
	@mkdir -p tmp
	cd apps/api && go build -o ../../tmp/server ./cmd/server

# --- Dev ---
dev:
	@mkdir -p tmp
	cd configs && air -c .air.toml

# --- Generate ---
generate: generate-openapi generate-types

generate-openapi:
	cd apps/api && go run ./cmd/openapi/ ../web/openapi.json

generate-types:
	cd apps/web && bun run generate

generate-sqlc:
	cd apps/api && sqlc generate

# --- Test & Check ---
test:
	cd apps/api && go test ./...

check:
	cd apps/web && bun run check

lint: test check

# --- Database ---
up:
	docker compose -f configs/docker-compose.yml up -d

down:
	docker compose -f configs/docker-compose.yml down

migrate:
	cd apps/api && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path db/migrations -database "$$DATABASE_URL" up

# --- Install ---
install:
	cd apps/api && go mod tidy
	bun install

# --- Clean ---
clean:
	rm -rf tmp/
	rm -rf apps/web/.svelte-kit
