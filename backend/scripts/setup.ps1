# TaskMan Backend Setup Script

Write-Host "🚀 Setting up TaskMan Backend..." -ForegroundColor Green

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "✅ Go is installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go is not installed. Please install Go 1.21 or higher." -ForegroundColor Red
    exit 1
}

# Create .env file if it doesn't exist
if (-not (Test-Path ".env")) {
    Write-Host "📝 Creating .env file from template..." -ForegroundColor Yellow
    Copy-Item "env.example" ".env"
    Write-Host "⚠️  Please update .env file with your configuration" -ForegroundColor Yellow
} else {
    Write-Host "✅ .env file already exists" -ForegroundColor Green
}

# Download dependencies
Write-Host "📦 Downloading dependencies..." -ForegroundColor Yellow
go mod download
go mod tidy

# Create bin directory
if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin"
}

Write-Host "✅ Setup complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Update .env file with your database configuration" -ForegroundColor White
Write-Host "2. Run database migrations from internal/database/migrations.sql" -ForegroundColor White
Write-Host "3. Start the server with: go run main.go" -ForegroundColor White
Write-Host ""
Write-Host "Available commands:" -ForegroundColor Cyan
Write-Host "  go run main.go     - Run the application" -ForegroundColor White
Write-Host "  go build -o bin/taskman-backend main.go - Build the application" -ForegroundColor White
Write-Host "  go test ./...      - Run tests" -ForegroundColor White
