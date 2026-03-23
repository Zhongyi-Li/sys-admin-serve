.PHONY: dev dev-down dev-infra-up dev-infra-down docker-dev docker-prod migrate-up migrate-down migrate-version migrate-force migrate-new seed

APP_CONFIG ?= configs/config.local.yaml
MIGRATION_STEPS ?= 1
MIGRATION_VERSION ?= 1
NAME ?=

dev:
	sh ./scripts/dev-up.sh

dev-down:
	sh ./scripts/dev-down.sh

dev-infra-up:
	docker compose -f docker-compose.infra.dev.yml up -d

dev-infra-down:
	docker compose -f docker-compose.infra.dev.yml down

docker-dev:
	docker compose -f docker-compose.dev.yml up --build

docker-prod:
	docker compose -f docker-compose.prod.yml up --build -d

migrate-up:
	APP_CONFIG=$(APP_CONFIG) go run ./cmd/migrate up

migrate-down:
	APP_CONFIG=$(APP_CONFIG) go run ./cmd/migrate down $(MIGRATION_STEPS)

migrate-version:
	APP_CONFIG=$(APP_CONFIG) go run ./cmd/migrate version

migrate-force:
	APP_CONFIG=$(APP_CONFIG) go run ./cmd/migrate force $(MIGRATION_VERSION)

migrate-new:
	sh ./scripts/new-migration.sh $(NAME)

seed:
	APP_CONFIG=$(APP_CONFIG) go run ./cmd/seed