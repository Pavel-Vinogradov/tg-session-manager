.PHONY: build run test clean deps docker-build docker-run

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build Docker image
docker-build:
	docker build -t tg-session-manager .

# Run with Docker
docker-run:
	docker run --rm -p 50051:50051 -v $(PWD)/sessions:/app/sessions tg-session-manager

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Check for vulnerabilities
sec:
	gosec ./...

# Development setup
dev-setup: deps install-tools
	cp .env.example .env
	@echo "Development environment setup complete!"
	@echo "Don't forget to edit .env file with your Telegram API credentials"
