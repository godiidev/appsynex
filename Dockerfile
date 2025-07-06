# File: Dockerfile
# Tạo tại: Dockerfile (root của project)
# Mục đích: Containerize ứng dụng Go API

# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Copy config files
COPY --from=builder /app/.env.example ./.env

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]

---

# File: docker-compose.yml
# Tạo tại: docker-compose.yml (root của project)
# Mục đích: Setup môi trường development với MySQL và Go API

version: '3.8'

services:
mysql:
image: mysql:8.0
container_name: appsynex_mysql
restart: always
environment:
MYSQL_ROOT_PASSWORD: rootpassword
MYSQL_DATABASE: appsynex
MYSQL_USER: appsynex_user
MYSQL_PASSWORD: appsynex_password
ports:
- "3306:3306"
volumes:
- mysql_data:/var/lib/mysql
- ./init.sql:/docker-entrypoint-initdb.d/init.sql
networks:
- appsynex_network

api:
build: .
container_name: appsynex_api
restart: always
ports:
- "8080:8080"
environment:
- DB_HOST=mysql
- DB_PORT=3306
- DB_USER=appsynex_user
- DB_PASS=appsynex_password
- DB_NAME=appsynex
- JWT_SECRET=your_super_secret_jwt_key_here
- PORT=8080
- ENV=development
depends_on:
- mysql
networks:
- appsynex_network
volumes:
- .:/app
command: ["./main"]

phpmyadmin:
image: phpmyadmin/phpmyadmin
container_name: appsynex_phpmyadmin
restart: always
ports:
- "8081:80"
environment:
PMA_HOST: mysql
PMA_PORT: 3306
PMA_USER: root
PMA_PASSWORD: rootpassword
depends_on:
- mysql
networks:
- appsynex_network

volumes:
mysql_data:

networks:
appsynex_network:
driver: bridge

---

# File: init.sql
# Tạo tại: init.sql (root của project)
# Mục đích: Initialize database với charset phù hợp

CREATE DATABASE IF NOT EXISTS appsynex
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

USE appsynex;

-- Grant permissions
GRANT ALL PRIVILEGES ON appsynex.* TO 'appsynex_user'@'%';
FLUSH PRIVILEGES;