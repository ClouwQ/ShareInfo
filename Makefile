.PHONY: migrate-create migrate-up migrate-down db-seed db-reset analyze-indexes test-health test-spaces test-all

DB_USER ?= postgres
DB_PASSWORD ?=
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= spaces_db

dev:
	go run cmd/api/main.go -env=dev

prod:
	go run cmd/api/main.go -env=prod

swagger:
	swag init -g cmd/api/main.go --parseInternal --parseDependency

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/repository/postgres/migrations -seq $$name

migrate-up:
	migrate -path internal/repository/postgres/migrations \
		-database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path internal/repository/postgres/migrations \
		-database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down 1

db-reset:
	make migrate-down
	make migrate-up

# Test health endpoint
test-health:
	@echo "Testing health endpoint..."
	@curl -s http://localhost:8001/health | jq

# Test spaces endpoint
test-spaces:
	@echo "Testing spaces endpoint..."
	@curl -s -H "X-User-Id: 1" -H "X-User-Role: user" \
		http://localhost:8001/api/v1/spaces | jq

# Create test space
test-create-space:
	@echo "Creating test space..."
	@curl -s -X POST http://localhost:8001/api/v1/spaces \
		-H "Content-Type: application/json" \
		-H "X-User-Id: 1" \
		-H "X-User-Role: user" \
		-d '{"title":"Test Space 2","description":"Test","space_name":"test-124","type":"solo","tags":["test"]}'

# Run all tests
test-all: test-health test-spaces