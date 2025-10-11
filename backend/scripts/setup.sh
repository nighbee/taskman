#!/bin/bash

# TaskMan Backend Setup Script

echo "🚀 Setting up TaskMan Backend..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | cut -d'o' -f2)
REQUIRED_VERSION="1.21"
if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go version $GO_VERSION is compatible"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "⚠️  Please update .env file with your configuration"
else
    echo "✅ .env file already exists"
fi

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

# Create bin directory
mkdir -p bin

echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Update .env file with your database configuration"
echo "2. Run database migrations from internal/database/migrations.sql"
echo "3. Start the server with: make run"
echo ""
echo "Available commands:"
echo "  make run        - Run the application"
echo "  make build      - Build the application"
echo "  make test       - Run tests"
echo "  make docker-build - Build Docker image"
echo "  make docker-run   - Run Docker container"
