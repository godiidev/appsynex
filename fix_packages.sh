#!/bin/bash

# File: fix_packages.sh
# Mục đích: Fix các lỗi package và setup project

echo "🔧 Fixing AppSynex package issues..."

# 1. Clean go modules
echo "📦 Cleaning Go modules..."
go clean -modcache
go mod tidy

# 2. Download dependencies
echo "📥 Downloading dependencies..."
go mod download

# 3. Check for any obvious issues
echo "🔍 Checking for issues..."

# Check if all required files exist
if [ ! -f ".env.example" ]; then
    echo "❌ .env.example not found"
    exit 1
fi

if [ ! -f "cmd/api/main.go" ]; then
    echo "❌ main.go not found"
    exit 1
fi

# 4. Try to build
echo "🔨 Testing build..."
if go build -o /tmp/appsynex-test ./cmd/api; then
    echo "✅ Build successful"
    rm -f /tmp/appsynex-test
else
    echo "❌ Build failed"
    echo "Common issues to check:"
    echo "  1. Make sure all files are saved correctly"
    echo "  2. Check for package declaration mismatches"
    echo "  3. Verify import paths are correct"
    exit 1
fi

# 5. Check if Docker is available
if command -v docker &> /dev/null; then
    echo "🐳 Docker found"
    if command -v docker-compose &> /dev/null; then
        echo "🐳 Docker Compose found"
    else
        echo "⚠️  Docker Compose not found. Please install docker-compose."
    fi
else
    echo "⚠️  Docker not found. Please install Docker for full functionality."
fi

# 6. Setup env file if it doesn't exist
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created"
else
    echo "📝 .env file already exists"
fi

echo ""
echo "🎉 Package fix completed!"
echo ""
echo "Next steps:"
echo "  1. Run 'make setup' for first-time setup"
echo "  2. Or run 'make docker-up && make setup-db' if you have existing setup"
echo "  3. Then run 'make run' to start the server"
echo ""
echo "Default login: admin / admin123"