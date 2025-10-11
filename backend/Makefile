# TaskMan Backend Makefile

.PHONY: help build run test clean deps docker-build docker-run

# Default target
help:
	@echo "Available targets:"
	@echo "  deps        - Download dependencies"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build: deps
	go build -o bin/taskman-backend main.go

# Run the application
run: deps
	go run main.go

# Run tests
test: deps
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t taskman-backend .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env taskman-backend

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
