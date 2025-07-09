# Makefile for AppSynex API
# Optimized for Docker development workflow

.PHONY: help build run test clean docker-up docker-down setup dev

# Variables
BINARY_NAME=appsynex-api
MAIN_PATH=./cmd/api
BUILD_DIR=./bin

help: ## Show this help message
	@echo 'AppSynex API - Development Commands'
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Development:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# === DEVELOPMENT COMMANDS ===
dev: ## Start full development environment (recommended)
	@echo "🚀 Starting development environment..."
	@make docker-up
	@echo "✅ Development environment ready!"
	@echo "   API: http://localhost:8081"
	@echo "   phpMyAdmin: http://localhost:8082"
	@echo "   MySQL: localhost:3307"

setup: ## First time setup
	@echo "🔧 Setting up project for first time..."
	@echo "1. Copying environment file..."
	@cp .env.example .env 2>/dev/null || echo "   .env already exists"
	@echo "2. Installing dependencies..."
	@go mod download && go mod tidy
	@echo "3. Starting development environment..."
	@make dev
	@echo ""
	@echo "✅ Setup completed!"
	@echo "   Access API at: http://localhost:8081"
	@echo "   Access phpMyAdmin at: http://localhost:8082"

# === DOCKER COMMANDS ===
docker-up: ## Start all Docker services
	@echo "🐳 Starting Docker services..."
	@docker-compose up -d --build
	@echo "⏳ Waiting for services to be ready..."
	@sleep 10
	@make status

docker-down: ## Stop all Docker services
	@echo "🛑 Stopping Docker services..."
	@docker-compose down

docker-restart: ## Restart Docker services
	@echo "🔄 Restarting Docker services..."
	@make docker-down
	@make docker-up

docker-rebuild: ## Rebuild and restart Docker services
	@echo "🔨 Rebuilding Docker services..."
	@docker-compose down
	@docker-compose build --no-cache
	@docker-compose up -d
	@make status

# === LOCAL DEVELOPMENT ===
run-local: ## Run API locally (requires Docker services)
	@echo "🏃 Running API locally..."
	@echo "📋 Make sure Docker services are running: make docker-services"
	@DB_HOST=localhost DB_PORT=3307 DB_USER=appsynex_user DB_PASS=appsynex_password DB_NAME=appsynex PORT=8080 go run $(MAIN_PATH)

docker-services: ## Start only MySQL and phpMyAdmin
	@echo "🗄️ Starting database services..."
	@docker-compose up -d mysql phpmyadmin
	@echo "⏳ Waiting for MySQL to be ready..."
	@sleep 10
	@echo "✅ Database services ready!"

build: ## Build binary locally
	@echo "🔨 Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# === TESTING & QUALITY ===
test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

format: ## Format Go code
	@echo "🎨 Formatting Go code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint: ## Run linter (requires golangci-lint)
	@echo "🔍 Running linter..."
	@golangci-lint run ./... || echo "Install golangci-lint first: https://golangci-lint.run/usage/install/"

tidy: ## Tidy go modules
	@echo "🧹 Tidying go modules..."
	@go mod tidy
	@echo "✅ Modules tidied"

# === MONITORING & DEBUGGING ===
logs: ## Show API logs
	@echo "📋 Showing API logs..."
	@docker-compose logs -f api

logs-all: ## Show all service logs
	@echo "📋 Showing all service logs..."
	@docker-compose logs -f

status: ## Show service status
	@echo "📊 Service Status:"
	@docker-compose ps
	@echo ""
	@echo "🔍 Health Check:"
	@curl -s http://localhost:8081/ > /dev/null && echo "✅ API is responding" || echo "❌ API is not responding"
	@curl -s http://localhost:8082/ > /dev/null && echo "✅ phpMyAdmin is responding" || echo "❌ phpMyAdmin is not responding"

mysql-shell: ## Access MySQL shell
	@echo "🗄️ Opening MySQL shell..."
	@docker-compose exec mysql mysql -u root -prootpassword appsynex

# === CLEANUP ===
clean: ## Clean build files and Docker resources
	@echo "🧹 Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@docker-compose down --remove-orphans
	@docker system prune -f

clean-all: ## Clean everything including volumes (WARNING: deletes data)
	@echo "⚠️  WARNING: This will delete all data!"
	@echo "Press Ctrl+C to cancel, or Enter to continue..."
	@read
	@echo "🧹 Cleaning everything..."
	@docker-compose down -v --remove-orphans
	@docker system prune -af
	@docker volume prune -f
	@rm -rf $(BUILD_DIR)
	@echo "✅ Everything cleaned!"

# === UTILITIES ===
backup-db: ## Backup database
	@echo "💾 Creating database backup..."
	@mkdir -p backups
	@docker-compose exec mysql mysqladump -u root -prootpassword appsynex > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "✅ Backup completed in backups/ directory"

install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/cmd/gosec@latest
	@echo "✅ Tools installed"

health: ## Check system health
	@echo "🏥 System Health Check:"
	@echo "📋 Docker status:"
	@docker --version
	@echo "📋 Go version:"
	@go version
	@echo "📋 Services:"
	@make status

# === QUICK COMMANDS ===
start: dev ## Alias for dev
stop: docker-down ## Alias for docker-down
restart: docker-restart ## Alias for docker-restart
rebuild: docker-rebuild ## Alias for docker-rebuild