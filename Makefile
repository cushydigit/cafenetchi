include configs/dev.env
export

# API Development
build:
	@echo "building cafenetchi-api server..."
	@cd ./api && go build -o ./bin/cafenetchi-api ./cmd

run: build
	@echo "running cafenetchi-api server..."
	@cd ./api && ./bin/cafenetchi-api

dev: 
	@echo "starting development environment..."
	@docker compose --env-file ./configs/dev.env up -d
	@cd api && air -c .air.toml

prod:
	@echo "starting production environment..."
	@docker compose --env-file ./configs/prod.env up

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
	docker compose --env-file ./configs/dev.env up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-reset:
	docker compose down -v
	docker compose --env-file ./configs/dev.env up -d --build

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

