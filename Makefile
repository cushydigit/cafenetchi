# ===============================================
# Development - api
# ===============================================
build:
	@echo "building cafenetchi-api server..."
	@cd ./api && go build -o ./bin/cafenetchi-api ./cmd

run: build
	@echo "running cafenetchi-api server..."
	@cd ./api && ./bin/cafenetchi-api

test:
	@echo "testing.."
	@cd ./api && go test ./... -v

test-coverage:
	@cd ./api && go test ./... -coverprofile=coverage.out
	@cd ./api && go tool cover -html=coverage.out

# ===============================================
# Docker
# ===============================================
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-reset: docker-down docker-up

docker-logs:
	docker compose logs -f

# ===============================================
# Utilities
# ===============================================
clean:
	@cd ./api && rm -rf bin/
	@cd ./api && go clean -cache

lint:
	@cd ./api && golangci-lint run

fmt:
	@cd ./api && go fmt ./...

tidy:
	@cd ./api && go mod tidy

# ===============================================
# Full reset (useful during early development)
# ===============================================
reset: docker-down
	rm -rf ./api/bin/ db_data/ minio_data/
	docker compose up -d --build
	@echo "Project reset complete!"
