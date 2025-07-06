-- File: migrations/000001_init_schema.up.sql
-- Tạo tại: migrations/000001_init_schema.up.sql
-- Mục đích: Tạo schema ban đầu cho users, roles, user_roles

CREATE TABLE IF NOT EXISTS roles (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    role_name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_role_name (role_name),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS users (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NULL,
    last_login TIMESTAMP NULL,
    account_status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_account_status (account_status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INT UNSIGNED NOT NULL,
    role_id INT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    INDEX idx_user_roles_role_id (role_id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default roles
INSERT INTO roles (role_name, description) VALUES
('ADMIN', 'Administrator with full system access'),
('MANAGER', 'Manager with limited administrative access'),
('STAFF', 'Staff with basic access');