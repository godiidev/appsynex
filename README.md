# AppSynex API

AppSynex API là một hệ thống quản lý doanh nghiệp được xây dựng bằng Go, sử dụng Gin framework và MySQL database với hệ thống phân quyền chi tiết.

## Tính năng chính

- 🔐 **Hệ thống phân quyền chi tiết**: Module-based permissions với role và user permissions
- 👥 **Quản lý người dùng**: CRUD users, assign roles, manage permissions
- 🏷️ **Quản lý danh mục sản phẩm**: Hierarchical categories
- 📦 **Quản lý sản phẩm mẫu**: Sample products với filtering và search
- 🔑 **JWT Authentication**: Secure token-based authentication
- 🐳 **Docker Support**: Containerized development environment
- 📊 **Database Migrations**: Automated schema management

## Yêu cầu hệ thống

- Go 1.24.1+
- Docker & Docker Compose
- MySQL 8.0+
- Make (optional, nhưng khuyến nghị)

## Cài đặt và setup

### 1. Clone repository

```bash
git clone <repository-url>
cd appsynex
```

### 2. Setup lần đầu (Khuyến nghị)

```bash
make setup
```

Lệnh này sẽ:

- Copy file `.env.example` thành `.env`
- Cài đặt dependencies
- Khởi động Docker containers
- Setup database với migrations và seed data
- Tạo user admin mặc định

### 3. Setup thủ công (nếu cần)

```bash
# Copy environment file
cp .env.example .env

# Install dependencies
make install-deps

# Start Docker services
make docker-up

# Wait for MySQL to be ready (khoảng 15-20 giây)
sleep 15

# Setup database
make setup-db
```

# Khởi động tất cả services

docker-compose up -d

# Xem logs

docker-compose logs -f api

# Dừng services

docker-compose down

# Rebuild API sau khi thay đổi code

docker-compose up -d --build api

# Xóa hoàn toàn (bao gồm volumes)

docker-compose down -v

## Chạy ứng dụng

### Development mode

```bash
make run
# hoặc
go run cmd/api/main.go
```

### Quick start (cho setup đã có)

```bash
make quick-start
```

### Production build

```bash
make build
./bin/appsynex-api
```

## Thông tin đăng nhập mặc định

Sau khi setup thành công, bạn có thể đăng nhập với:

- **Username**: `admin`
- **Password**: `admin123`
- **API Base URL**: `http://localhost:8080/api/v1`

## API Endpoints

### Authentication

```bash
# Login
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "admin123"
}
```

### Users Management

```bash
# Get all users (requires USER_VIEW permission)
GET /api/v1/users

# Create user (requires USER_CREATE permission)
POST /api/v1/users

# Get user by ID (requires USER_VIEW permission)
GET /api/v1/users/{id}

# Update user (requires USER_UPDATE permission)
PUT /api/v1/users/{id}

# Delete user (requires USER_DELETE permission)
DELETE /api/v1/users/{id}
```

### Permission Management

```bash
# Get all permissions
GET /api/v1/permissions

# Get permissions by module
GET /api/v1/permissions/module/{module}

# Check user permission
POST /api/v1/permissions/check

# Assign permissions to role
POST /api/v1/roles/{roleId}/permissions

# Get user effective permissions
GET /api/v1/users/{userId}/effective-permissions
```

### Product Categories

```bash
# Get all categories
GET /api/v1/categories

# Create category
POST /api/v1/categories

# Get category by ID
GET /api/v1/categories/{id}

# Update category
PUT /api/v1/categories/{id}

# Delete category
DELETE /api/v1/categories/{id}
```

### Sample Products

```bash
# Get all samples with filtering
GET /api/v1/samples?page=1&limit=10&search=cotton&category=1

# Create sample
POST /api/v1/samples

# Get sample by ID
GET /api/v1/samples/{id}

# Update sample
PUT /api/v1/samples/{id}

# Delete sample
DELETE /api/v1/samples/{id}
```

## Hệ thống phân quyền

### Roles mặc định

- **SUPER_ADMIN**: Toàn quyền hệ thống
- **ADMIN**: Quản trị viên
- **MANAGER**: Quản lý với quyền hạn giới hạn
- **STAFF**: Nhân viên với quyền cơ bản

### Permission Modules

- **USER**: Quản lý người dùng
- **ROLE**: Quản lý roles và permissions
- **PRODUCT**: Quản lý sản phẩm
- **PRODUCT_CATEGORY**: Quản lý danh mục
- **SAMPLE**: Quản lý sản phẩm mẫu
- **CUSTOMER**: Quản lý khách hàng
- **ORDER**: Quản lý đơn hàng
- **WAREHOUSE**: Quản lý kho
- **FINANCE**: Quản lý tài chính
- **REPORT**: Báo cáo
- **SYSTEM**: Quản trị hệ thống

### Permission Actions

- **VIEW**: Xem
- **CREATE**: Tạo mới
- **UPDATE**: Cập nhật
- **DELETE**: Xóa
- **EXPORT**: Xuất dữ liệu
- **IMPORT**: Nhập dữ liệu
- **APPROVE**: Phê duyệt
- **ASSIGN_ROLES**: Gán roles
- **ASSIGN_PERMISSIONS**: Gán permissions

## Lệnh Make hữu ích

```bash
# Xem tất cả lệnh available
make help

# Development
make run                  # Chạy server
make test                 # Chạy tests
make build               # Build binary

# Database
make setup-db           # Setup database với migrations và seeding
make seed               # Chạy seeding data
make reset-db           # Reset database (XÓA TẤT CẢ DỮ LIỆU)
make backup-db          # Backup database

# Docker
make docker-up          # Start containers
make docker-down        # Stop containers
make docker-logs        # Xem logs
make mysql-shell        # Access MySQL shell

# Code quality
make format             # Format code
make lint               # Run linter
make security-check     # Security checks

# Utilities
make status             # Xem status services
make clean              # Clean build files
make deps-update        # Update dependencies
```

## Cấu trúc thư mục

```
appsynex/
├── cmd/api/            # Application entry point
├── config/             # Configuration management
├── internal/
│   ├── api/           # HTTP handlers và middleware
│   │   ├── handlers/  # API handlers
│   │   ├── middleware/ # Authentication, CORS, permissions
│   │   └── router/    # Route definitions
│   ├── domain/        # Business logic
│   │   ├── models/    # Data models
│   │   └── services/  # Business services
│   ├── dto/           # Data Transfer Objects
│   │   ├── request/   # Request DTOs
│   │   └── response/  # Response DTOs
│   └── repository/    # Data access layer
│       ├── interfaces/ # Repository interfaces
│       └── mysql/     # MySQL implementations
├── migrations/        # Database migrations
├── pkg/              # Shared packages
│   ├── auth/         # JWT utilities
│   ├── logger/       # Logging utilities
│   └── utils/        # Utility functions
├── scripts/          # Setup và utility scripts
├── docker-compose.yml # Docker configuration
├── Dockerfile        # Container build file
├── Makefile         # Build và utility commands
└── README.md        # Documentation
```

## Environment Variables

Chỉnh sửa file `.env` để cấu hình:

```env
# Server Configuration
PORT=8080
ENV=development
LOG_LEVEL=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=appsynex_user
DB_PASS=appsynex_password
DB_NAME=appsynex
DB_CHARSET=utf8mb4

# JWT Configuration
JWT_SECRET=your_secret_key_here
JWT_EXPIRES_IN=24h
```

## Troubleshooting

### Database connection issues

```bash
# Check database status
make status

# Restart database
make docker-down
make docker-up

# Check logs
make docker-logs
```

### Permission denied errors

```bash
# Reset database và permissions
make reset-db
```

### Build issues

```bash
# Clean và rebuild
make clean
make install-deps
make build
```

## Development

### Thêm migration mới

1. Tạo file migration trong thư mục `migrations/`
2. Đặt tên theo format: `000XXX_description.up.sql` và `000XXX_description.down.sql`
3. Chạy `make setup-db` để apply

### Thêm permission mới

1. Thêm vào `PreDefinedPermissions` trong `internal/domain/models/permission.go`
2. Chạy `make setup-db` để update

### Thêm endpoint mới

1. Tạo handler trong `internal/api/handlers/v1/`
2. Thêm route trong `internal/api/router/router.go`
3. Thêm permission check nếu cần

## Production Deployment

### Build production

```bash
make release
```

### Deploy với Docker

```bash
docker build -t appsynex-api .
docker run -p 8080:8080 --env-file .env appsynex-api
```

## Liên hệ

- Email: support@appsynex.vn
- Documentation: [Link to docs]
- Issues: [Link to issues]
