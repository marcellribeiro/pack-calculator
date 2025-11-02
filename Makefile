.PHONY: build test run clean docker-build docker-run

# Build the application
build:
	@echo "🔨 Building Pack Calculator..."
	@mkdir -p bin
	@go build -o bin/pack-calculator ./cmd/api
	@echo "✅ Build complete! Binary at: bin/pack-calculator"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Run the application
run:
	@echo "🚀 Starting Pack Calculator..."
	@go run cmd/api/main.go

# Clean build artifacts
clean:
	@rm -rf bin coverage.out coverage.html
	@rm -f pack-calculator coverage.out coverage.html
	@go clean
	@echo "✅ Clean complete!"

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t pack-calculator:latest .
	@echo "✅ Docker image built!"

# Run Docker container
docker-run:
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 pack-calculator:latest

# Run with docker-compose
docker-compose-up:
	@echo "🐳 Starting with docker-compose..."
	@docker-compose up --build

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed!"

# Format code
fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

# Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

# Help
help:
	@echo "Pack Calculator - Available commands:"
	@echo "  make build             - Build the application"
	@echo "  make test              - Run all tests"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make run               - Run the application"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make docker-build      - Build Docker image"
	@echo "  make docker-run        - Run Docker container"
	@echo "  make docker-compose-up - Run with docker-compose"
	@echo "  make deps              - Install dependencies"
	@echo "  make fmt               - Format code"
	@echo "  make help              - Show this help message"

