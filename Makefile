# File: Makefile
# Tạo tại: Makefile (root của project)  
# Mục đích: Các lệnh tiện ích để build, run, test project với enhanced permission system

.PHONY: help build run test clean setup-db seed docker-up docker-down

# Variables
BINARY_NAME=appsynex-api
MAIN_PATH=./cmd/api
BUILD_DIR=./bin

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

run: ## Run the application
	@echo "Running application..."
	@go run $(MAIN_PATH)

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build files
	@echo "Cleaning build files..."
	@rm -rf $(BUILD_DIR)
	@go clean

setup-db: ## Setup database with migrations and seed data
	@echo "Setting up database..."
	@go run scripts/setup.go

seed: ## Run database seeding only
	@echo "Running database seeding..."
	@go run scripts/seed_enhanced_permissions.go

docker-up: ## Start docker containers
	@echo "Starting docker containers..."
	@docker-compose up -d

docker-down: ## Stop docker containers
	@echo "Stopping docker containers..."
	@docker-compose down

docker-build: ## Build docker image
	@echo "Building docker image..."
	@docker-compose build

docker-logs: ## Show docker logs
	@docker-compose logs -f

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

dev: ## Start development environment
	@echo "Starting development environment..."
	@make docker-up
	@sleep 10
	@make setup-db
	@make run

setup: ## Setup project for first time
	@echo "Setting up project for first time..."
	@echo "1. Copying environment file..."
	@cp .env.example .env
	@echo "2. Installing dependencies..."
	@make install-deps
	@echo "3. Starting Docker services..."
	@make docker-up
	@echo "4. Waiting for database to be ready..."
	@sleep 15
	@echo "5. Setting up database and seeding data..."
	@make setup-db
	@echo ""
	@echo "Setup completed! You can now:"
	@echo "  - Run 'make run' to start the server"
	@echo "  - Login with username: admin, password: admin123"
	@echo "  - Access the API at http://localhost:8080"

reset-db: ## Reset database (WARNING: This will delete all data)
	@echo "WARNING: This will delete all data in the database!"
	@echo "Press Ctrl+C to cancel, or Enter to continue..."
	@read
	@echo "Resetting database..."
	@make docker-down
	@docker volume rm appsynex_mysql_data || true
	@make docker-up
	@sleep 15
	@make setup-db
	@echo "Database reset completed!"

logs: ## Show application logs
	@docker-compose logs -f api

mysql-shell: ## Access MySQL shell
	@docker-compose exec mysql mysql -u root -prootpassword appsynex

status: ## Show status of services
	@echo "Docker services status:"
	@docker-compose ps
	@echo ""
	@echo "Database connection test:"
	@docker-compose exec mysql mysqladmin -u root -prootpassword ping

backup-db: ## Backup database
	@echo "Creating database backup..."
	@mkdir -p backups
	@docker-compose exec mysql mysqldump -u root -prootpassword appsynex > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Backup completed in backups/ directory"

format: ## Format Go code
	@echo "Formatting Go code..."
	@go fmt ./...
	@echo "Code formatted"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

security-check: ## Run security checks
	@echo "Running security checks..."
	@gosec ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "Modules tidied"

generate: ## Generate code (if needed)
	@echo "Generating code..."
	@go generate ./...

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated"

release: ## Build release version
	@echo "Building release version..."
	@make clean
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	@CGO_ENABLED=0 GOOS=windows go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_PATH)
	@CGO_ENABLED=0 GOOS=darwin go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-macos $(MAIN_PATH)
	@echo "Release builds completed in $(BUILD_DIR)/"

check-env: ## Check if .env file exists
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Run 'make setup' first."; \
		exit 1; \
	fi

serve: check-env ## Start the server (alias for run)
	@make run

quick-start: ## Quick start for existing setup
	@echo "Quick starting AppSynex API..."
	@make docker-up
	@sleep 5
	@make run