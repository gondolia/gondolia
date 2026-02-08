.PHONY: build run test migrate migrate-down migrate-version seed clean help

# Default database URL for local development
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/identity?sslmode=disable

# Build the service
build:
	go build -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate

# Run the service locally
run: build
	JWT_ACCESS_SECRET=dev-access-secret-32-chars-minimum \
	JWT_REFRESH_SECRET=dev-refresh-secret-32-chars-minimum \
	DATABASE_HOST=localhost \
	./bin/server

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run migrations up
migrate:
	DATABASE_HOST=localhost go run ./cmd/migrate -command=up

# Run migrations down (1 step)
migrate-down:
	DATABASE_HOST=localhost go run ./cmd/migrate -command=down -steps=1

# Show migration version
migrate-version:
	DATABASE_HOST=localhost go run ./cmd/migrate -command=version

# Seed database with test data
seed:
	DATABASE_HOST=localhost go run ./cmd/migrate -command=seed

# Full setup: migrate and seed
setup: migrate seed

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Lint code
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Start local PostgreSQL with Docker
db-start:
	docker run -d --name identity-postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=identity \
		-p 5432:5432 \
		postgres:15-alpine

# Stop local PostgreSQL
db-stop:
	docker stop identity-postgres || true
	docker rm identity-postgres || true

# Connect to local PostgreSQL
db-connect:
	docker exec -it identity-postgres psql -U postgres -d identity

# Show help
help:
	@echo "Identity Service Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          - Build the service and migration tool"
	@echo "  make run            - Run the service locally"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make migrate        - Run database migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make migrate-version- Show current migration version"
	@echo "  make seed           - Seed database with test data"
	@echo "  make setup          - Run migrations and seed"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make db-start       - Start local PostgreSQL"
	@echo "  make db-stop        - Stop local PostgreSQL"
	@echo "  make db-connect     - Connect to local PostgreSQL"
	@echo ""
	@echo "Environment variables:"
	@echo "  DATABASE_URL        - Database connection URL"
