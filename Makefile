include configs/dev.env
export

# Environment
DEV_ENV=./configs/dev.env
PRO_ENV=./configs/pro.env

DEV_COMPOSE=docker compose --env-file $(DEV_ENV)
PRO_COMPOSE=docker compose -f docker-compose.prod.yml --env-file $(PRO_ENV)

# API Development
build:
	@echo "building cafenetchi-api server..."
	@cd ./api && go build -o ./bin/cafenetchi-api ./cmd

run: build
	@echo "running cafenetchi-api server..."
	@cd ./api && ./bin/cafenetchi-api

dev: 
	@echo "starting development environment..."
	$(DEV_COMPOSE) up -d
	@cd api && air -c .air.toml

prod:
	@echo "starting production environment..."
	@$(PRO_COMPOSE) up -d -build

# Testing
test:
	@echo "testing.."
	@cd ./api && go test ./... -v

test-coverage:
	@cd ./api && go test ./... -coverprofile=coverage.out
	@cd ./api && go tool cover -html=coverage.out
	
# Database
sqlc:
	@echo "generating sqlc files..."
	@cd api && sqlc generate

migrate-up:
	@echo "migrating database up..."
	@cd api && goose -dir internal/db/migrations postgres "$(DB_URL)" up

migrate-check-status:
	@echo "migrating check database status..."
	@cd api && goose -dir internal/db/migrations postgres "$(DB_URL)" status


migrate-down:
	@echo "migrating database down..."
	@cd api && goose -dir internal/db/migrations postgres "$(DB_URL)" down

migrate-create:
	@echo "creating new migration..."
	@cd api && goose -dir internal/db/migrations create $(name) sql

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
	@echo "cleaning up..."
	@rm -rf ./api/bin
	@rm -rf ./api/tmp

lint:
	@echo "linting..."
	@cd api && golangci-lint run

fmt:
	@echo "formatting..."
	@cd api && go fmt ./...

tidy:
	@echo "tidying..."
	@cd api && go mod tidy

