
# Force bash globally
SHELL := /bin/bash

# Environment
DEV_ENV=./configs/dev.env
PRO_ENV=./configs/pro.env
TEST_ENV=./configs/test.env

DEV_COMPOSE=docker compose --env-file $(DEV_ENV)
PRO_COMPOSE=docker compose -f docker-compose.prod.yml --env-file $(PRO_ENV)

# Macro
define load_env
	set -a && source $(1) && set +a
endef


# API Development
build:
	@echo "Building cafenetchi-api..."
	@$(call load_env, $(DEV_ENV)) && \
	cd api && go build -o ./bin/cafenetchi-api ./cmd

run: build
	@echo "Running cafenetchi-api server..."
	@$(call load_env, $(PRO_ENV)) && \
	cd ./api && ./bin/cafenetchi-api

dev: 
	@echo "Starting development environment..."
	@$(DEV_COMPOSE) up -d
	@$(call load_env, $(DEV_ENV)) && \
	cd api && air -c .air.toml

prod:
	@echo "Starting production environment..."
	@$(PRO_COMPOSE) up -d --build

# Testing
test:
	@echo "Exporting test environment variables..."
	@$(call load_env, $(TEST_ENV)) && \
	cd api && go test ./... -v

test-coverage:
	@echo "Generating coverage report..."
	@$(call load_env, $(TEST_ENV)) && \
	cd ./api && go test ./... -coverprofile=coverage.out
	@echo "Opening coverage report..."
	@$(call load_env, $(TEST_ENV)) && \
	cd ./api && go tool cover -html=coverage.out

test-db-shell:
	@echo "Connecting to database..."
	@$(call load_env, $(TEST_ENV)) && \
	psql "$$DB_URL"
	
# Database
sqlc:
	@echo "Generating sqlc files..."
	@cd api && sqlc generate

migrate-up:
	@echo "Migrating development database..."
	@$(call load_env, $(DEV_ENV)) && \
	cd api && goose -dir internal/db/migrations postgres "$$DB_URL" up

migrate-status:
	@echo "Development migration status..."
	@$(call load_env, $(DEV_ENV)) && \
	cd api && goose -dir internal/db/migrations postgres "$$DB_URL" status

migrate-down:
	@echo "Rolling back development database..."
	@$(call load_env, $(DEV_ENV)) && \
	cd api && goose -dir internal/db/migrations postgres "$$DB_URL" down

migrate-create:
	@echo "Creating new migration..."
	@cd api && goose -dir internal/db/migrations create $(name) sql

migrate-test-up:
	@echo "Migrating test database..."
	@$(call load_env, $(TEST_ENV)) && \
	cd api && goose -dir internal/db/migrations postgres "$$DB_URL" up

migrate-test-down:
	@echo "Rolling back test database..."
	@$(call load_env, $(TEST_ENV)) && \
	cd api && goose -dir internal/db/migrations postgres "$$DB_URL" down

migrate-test-status:
	@echo "Test migration status"
	@$(call load_env, $(TEST_ENV)) && \
	cd api && goose -dir internal/db/migrations postgres "$$DB_URL" status

# Docker
docker-up:
	$(DEV_COMPOSE) up -d

docker-down:
	$(DEV_COMPOSE) down

docker-logs:
	$(DEV_COMPOSE) logs -f

docker-reset:
	$(DEV_COMPOSE) down -v
	$(DEV_COMPOSE) up -d --build

docker-prod-up:
	$(PRO_COMPOSE) up -d --build

docker-build-prod:
	$(PRO_COMPOSE) down

docker-prod-logs:
	$(PRO_COMPOSE) logs -f

# ===============================================
# Utilities
# ===============================================
clean:
	@echo "Cleaning up..."
	@rm -rf ./api/bin
	@rm -rf ./api/tmp

lint:
	@echo "Linting..."
	@cd api && golangci-lint run

fmt:
	@echo "Formatting..."
	@cd api && go fmt ./...

tidy:
	@echo "Tidying modules..."
	@cd api && go mod tidy

