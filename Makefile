.PHONY: help test build clean lint fmt docker-build docker-push

# Default target
help:
	@echo "Gondolia - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make test          - Run all tests"
	@echo "  make test-cover    - Run tests with coverage"
	@echo "  make lint          - Run linters"
	@echo "  make fmt           - Format code"
	@echo ""
	@echo "Build:"
	@echo "  make build         - Build all services"
	@echo "  make clean         - Remove build artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  - Build Docker images"
	@echo "  make docker-push   - Push Docker images to registry"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make k3d-up        - Start local K3d cluster"
	@echo "  make k3d-down      - Stop local K3d cluster"
	@echo "  make deploy-dev    - Deploy to local K3d"

# Registry configuration
REGISTRY := ghcr.io/gondolia/gondolia
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Service list
SERVICES := identity

# --- Testing ---

test:
	@echo "🧪 Running tests..."
	go test -v ./...

test-cover:
	@echo "📊 Running tests with coverage..."
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

test-race:
	@echo "🏁 Running tests with race detector..."
	go test -race ./...

# --- Linting & Formatting ---

lint:
	@echo "🔍 Running linters..."
	golangci-lint run ./...

fmt:
	@echo "✨ Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# --- Build ---

build:
	@echo "🔨 Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd services/$$service && go build -o ../../bin/$$service ./cmd/server && cd ../..; \
	done
	@echo "✅ Build complete"

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf coverage.out coverage.html
	@echo "✅ Clean complete"

# --- Docker ---

docker-build:
	@echo "🐳 Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building $$service:$(VERSION)..."; \
		docker build -t $(REGISTRY)-$$service:$(VERSION) -f services/$$service/Dockerfile .; \
	done
	@echo "✅ Docker images built"

docker-push:
	@echo "📤 Pushing Docker images..."
	@for service in $(SERVICES); do \
		echo "Pushing $$service:$(VERSION)..."; \
		docker push $(REGISTRY)-$$service:$(VERSION); \
	done
	@echo "✅ Docker images pushed"

# --- Infrastructure ---

k3d-up:
	@echo "🚀 Starting K3d cluster..."
	k3d cluster create gondolia --config infrastructure/kubernetes/k3d-config.yaml
	@echo "✅ K3d cluster started"

k3d-down:
	@echo "🛑 Stopping K3d cluster..."
	k3d cluster delete gondolia
	@echo "✅ K3d cluster stopped"

deploy-dev:
	@echo "🚢 Deploying to local K3d..."
	kubectl apply -f infrastructure/kubernetes/base/namespaces.yaml
	kubectl apply -k infrastructure/kubernetes/dev/
	@echo "✅ Deployment complete"

# --- Dependencies ---

deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	@echo "✅ Dependencies downloaded"

deps-tidy:
	@echo "🧹 Tidying dependencies..."
	go mod tidy
	@echo "✅ Dependencies tidied"

# --- Database Migrations ---

migrate-up:
	@echo "🔼 Running migrations..."
	cd services/identity && goose -dir migrations postgres "$(DB_URL)" up
	@echo "✅ Migrations applied"

migrate-down:
	@echo "🔽 Rolling back migrations..."
	cd services/identity && goose -dir migrations postgres "$(DB_URL)" down
	@echo "✅ Migrations rolled back"

# --- Frontend ---

frontend-install:
	@echo "📦 Installing frontend dependencies..."
	cd frontend && npm install
	@echo "✅ Frontend dependencies installed"

frontend-dev:
	@echo "🎨 Starting frontend dev server..."
	cd frontend && npm run dev

frontend-build:
	@echo "🔨 Building frontend..."
	cd frontend && npm run build
	@echo "✅ Frontend built"
