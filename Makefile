# File: Makefile
# Tạo tại: Makefile (root của project)
# Mục đích: Các lệnh tiện ích để build, run, test project

.PHONY: help build run test clean migrate-up migrate-down seed docker-up docker-down

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

migrate-up: ## Run up migrations
	@echo "Running up migrations..."
	@go run scripts/migrate.go up

migrate-down: ## Run down migrations
	@echo "Running down migrations..."
	@go run scripts/migrate.go down

seed: ## Run database seeding
	@echo "Running database seeding..."
	@go run scripts/seed.go
	@go run scripts/seed_sample_data.go

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
	@make migrate-up
	@make seed
	@make run

setup: ## Setup project for first time
	@echo "Setting up project..."
	@cp .env.example .env
	@make install-deps
	@make docker-up
	@sleep 15
	@make migrate-up
	@make seed
	@echo "Setup completed! Run 'make run' to start the server"

---

# File: build.sh
# Tạo tại: build.sh (root của project)
# Mục đích: Script build cho production

#!/bin/bash

# Build script for AppSynex API
set -e

echo "Starting build process..."

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf bin/
mkdir -p bin/

# Build for Linux (production)
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/appsynex-api-linux ./cmd/api

# Build for Windows (development)
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o bin/appsynex-api-windows.exe ./cmd/api

# Build for macOS (development)
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o bin/appsynex-api-macos ./cmd/api

echo "Build completed successfully!"
echo "Binaries available in bin/ directory:"
ls -la bin/

---

# File: deploy.sh  
# Tạo tại: deploy.sh (root của project)
# Mục đích: Script deploy lên aaPanel

#!/bin/bash

# Deploy script for aaPanel
set -e

SERVER_HOST=${1:-"your-server-ip"}
SERVER_USER=${2:-"root"}
APP_DIR="/www/wwwroot/appsynex-api"

echo "Deploying to $SERVER_HOST..."

# Build application
echo "Building application..."
make build

# Create deployment package
echo "Creating deployment package..."
tar -czf appsynex-api.tar.gz bin/ migrations/ .env.example

# Upload to server
echo "Uploading to server..."
scp appsynex-api.tar.gz $SERVER_USER@$SERVER_HOST:/tmp/

# Deploy on server
echo "Deploying on server..."
ssh $SERVER_USER@$SERVER_HOST << 'EOF'
cd /tmp
tar -xzf appsynex-api.tar.gz

# Stop existing service
systemctl stop appsynex-api || true

# Backup current version
if [ -d "/www/wwwroot/appsynex-api" ]; then
    mv /www/wwwroot/appsynex-api /www/wwwroot/appsynex-api.backup.$(date +%Y%m%d_%H%M%S)
fi

# Create app directory
mkdir -p /www/wwwroot/appsynex-api
cd /www/wwwroot/appsynex-api

# Move files
mv /tmp/bin ./
mv /tmp/migrations ./
mv /tmp/.env.example ./

# Set permissions
chmod +x bin/appsynex-api-linux
chown -R www:www /www/wwwroot/appsynex-api

# Setup systemd service
cat > /etc/systemd/system/appsynex-api.service << 'SERVICE_EOF'
[Unit]
Description=AppSynex API Server
After=network.target

[Service]
Type=simple
User=www
Group=www
WorkingDirectory=/www/wwwroot/appsynex-api
ExecStart=/www/wwwroot/appsynex-api/bin/appsynex-api-linux
Restart=always
RestartSec=5

Environment=PORT=8080
Environment=ENV=production

[Install]
WantedBy=multi-user.target
SERVICE_EOF

# Reload systemd and start service
systemctl daemon-reload
systemctl enable appsynex-api
systemctl start appsynex-api

echo "Deployment completed successfully!"
EOF

# Cleanup
rm -f appsynex-api.tar.gz

echo "Deployment script completed!"