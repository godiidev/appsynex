-- File: migrations/000013_enhanced_permissions.up.sql
-- Tạo tại: migrations/000013_enhanced_permissions.up.sql
-- Mục đích: Enhanced permission system tables

-- Drop and recreate permissions table with new structure
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;

CREATE TABLE IF NOT EXISTS permissions (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    module VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(100) NULL,
    permission_name VARCHAR(200) NOT NULL,
    description TEXT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_permission_name (permission_name),
    INDEX idx_module (module),
    INDEX idx_action (action),
    INDEX idx_module_action (module, action),
    INDEX idx_is_active (is_active),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Permission groups table
CREATE TABLE IF NOT EXISTS permission_groups (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    group_name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    description TEXT NULL,
    module VARCHAR(100) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_group_name (group_name),
    INDEX idx_module (module),
    INDEX idx_sort_order (sort_order),
    INDEX idx_is_active (is_active),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Enhanced role_permissions table
CREATE TABLE IF NOT EXISTS role_permissions (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    role_id INT UNSIGNED NOT NULL,
    permission_id INT UNSIGNED NOT NULL,
    granted_by INT UNSIGNED NULL,
    granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_role_permission (role_id, permission_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id),
    INDEX idx_granted_by (granted_by),
    INDEX idx_is_active (is_active),
    INDEX idx_expires_at (expires_at),
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_granted_by FOREIGN KEY (granted_by) REFERENCES users (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- User permissions table (for direct user permissions)
CREATE TABLE IF NOT EXISTS user_permissions (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id INT UNSIGNED NOT NULL,
    permission_id INT UNSIGNED NOT NULL,
    grant_type VARCHAR(20) NOT NULL DEFAULT 'GRANT', -- 'GRANT' or 'DENY'
    granted_by INT UNSIGNED NOT NULL,
    granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    reason TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_user_permission (user_id, permission_id),
    INDEX idx_user_id (user_id),
    INDEX idx_permission_id (permission_id),
    INDEX idx_grant_type (grant_type),
    INDEX idx_granted_by (granted_by),
    INDEX idx_is_active (is_active),
    INDEX idx_expires_at (expires_at),
    CONSTRAINT fk_user_permissions_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_permissions_granted_by FOREIGN KEY (granted_by) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT chk_grant_type CHECK (grant_type IN ('GRANT', 'DENY'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert predefined permissions
INSERT INTO permissions (module, action, permission_name, description) VALUES
-- User Management
('USER', 'VIEW', 'USER_VIEW', 'View users'),
('USER', 'CREATE', 'USER_CREATE', 'Create new users'),
('USER', 'UPDATE', 'USER_UPDATE', 'Update user information'),
('USER', 'DELETE', 'USER_DELETE', 'Delete users'),
('USER', 'ASSIGN_ROLES', 'USER_ASSIGN_ROLES', 'Assign roles to users'),
('USER', 'RESET_PASSWORD', 'USER_RESET_PASSWORD', 'Reset user passwords'),
('USER', 'VIEW_OWN', 'USER_VIEW_OWN', 'View own user profile'),
('USER', 'ASSIGN_PERMISSIONS', 'USER_ASSIGN_PERMISSIONS', 'Assign direct permissions to users'),

-- Role Management
('ROLE', 'VIEW', 'ROLE_VIEW', 'View roles'),
('ROLE', 'CREATE', 'ROLE_CREATE', 'Create new roles'),
('ROLE', 'UPDATE', 'ROLE_UPDATE', 'Update role information'),
('ROLE', 'DELETE', 'ROLE_DELETE', 'Delete roles'),
('ROLE', 'ASSIGN_PERMISSIONS', 'ROLE_ASSIGN_PERMISSIONS', 'Assign permissions to roles'),

-- Product Management
('PRODUCT', 'VIEW', 'PRODUCT_VIEW', 'View products'),
('PRODUCT', 'CREATE', 'PRODUCT_CREATE', 'Create new products'),
('PRODUCT', 'UPDATE', 'PRODUCT_UPDATE', 'Update product information'),
('PRODUCT', 'DELETE', 'PRODUCT_DELETE', 'Delete products'),
('PRODUCT', 'EXPORT', 'PRODUCT_EXPORT', 'Export product data'),
('PRODUCT', 'IMPORT', 'PRODUCT_IMPORT', 'Import product data'),

-- Product Category Management
('PRODUCT_CATEGORY', 'VIEW', 'PRODUCT_CATEGORY_VIEW', 'View product categories'),
('PRODUCT_CATEGORY', 'CREATE', 'PRODUCT_CATEGORY_CREATE', 'Create product categories'),
('PRODUCT_CATEGORY', 'UPDATE', 'PRODUCT_CATEGORY_UPDATE', 'Update product categories'),
('PRODUCT_CATEGORY', 'DELETE', 'PRODUCT_CATEGORY_DELETE', 'Delete product categories'),

-- Sample Management
('SAMPLE', 'VIEW', 'SAMPLE_VIEW', 'View samples'),
('SAMPLE', 'CREATE', 'SAMPLE_CREATE', 'Create new samples'),
('SAMPLE', 'UPDATE', 'SAMPLE_UPDATE', 'Update sample information'),
('SAMPLE', 'DELETE', 'SAMPLE_DELETE', 'Delete samples'),
('SAMPLE', 'DISPATCH', 'SAMPLE_DISPATCH', 'Dispatch samples to customers'),
('SAMPLE', 'TRACK', 'SAMPLE_TRACK', 'Track sample status'),

-- Customer Management
('CUSTOMER', 'VIEW', 'CUSTOMER_VIEW', 'View customers'),
('CUSTOMER', 'CREATE', 'CUSTOMER_CREATE', 'Create new customers'),
('CUSTOMER', 'UPDATE', 'CUSTOMER_UPDATE', 'Update customer information'),
('CUSTOMER', 'DELETE', 'CUSTOMER_DELETE', 'Delete customers'),
('CUSTOMER', 'VIEW_ACTIVITY', 'CUSTOMER_VIEW_ACTIVITY', 'View customer activity logs'),

-- Order Management
('ORDER', 'VIEW', 'ORDER_VIEW', 'View orders'),
('ORDER', 'CREATE', 'ORDER_CREATE', 'Create new orders'),
('ORDER', 'UPDATE', 'ORDER_UPDATE', 'Update order information'),
('ORDER', 'DELETE', 'ORDER_DELETE', 'Delete orders'),
('ORDER', 'APPROVE', 'ORDER_APPROVE', 'Approve orders'),
('ORDER', 'CANCEL', 'ORDER_CANCEL', 'Cancel orders'),
('ORDER', 'SHIP', 'ORDER_SHIP', 'Ship orders'),

-- Warehouse Management
('WAREHOUSE', 'VIEW', 'WAREHOUSE_VIEW', 'View warehouse data'),
('WAREHOUSE', 'CREATE', 'WAREHOUSE_CREATE', 'Create warehouse entries'),
('WAREHOUSE', 'UPDATE', 'WAREHOUSE_UPDATE', 'Update warehouse data'),
('WAREHOUSE', 'DELETE', 'WAREHOUSE_DELETE', 'Delete warehouse entries'),
('WAREHOUSE', 'TRANSFER', 'WAREHOUSE_TRANSFER', 'Transfer inventory'),

-- Financial Management
('FINANCE', 'VIEW', 'FINANCE_VIEW', 'View financial data'),
('FINANCE', 'CREATE', 'FINANCE_CREATE', 'Create financial records'),
('FINANCE', 'UPDATE', 'FINANCE_UPDATE', 'Update financial data'),
('FINANCE', 'DELETE', 'FINANCE_DELETE', 'Delete financial records'),
('FINANCE', 'APPROVE', 'FINANCE_APPROVE', 'Approve financial transactions'),

-- Reporting
('REPORT', 'VIEW', 'REPORT_VIEW', 'View reports'),
('REPORT', 'CREATE', 'REPORT_CREATE', 'Create custom reports'),
('REPORT', 'EXPORT', 'REPORT_EXPORT', 'Export reports'),

-- System Administration
('SYSTEM', 'VIEW', 'SYSTEM_VIEW', 'View system information'),
('SYSTEM', 'VIEW_LOGS', 'SYSTEM_VIEW_LOGS', 'View system logs'),
('SYSTEM', 'MANAGE_SETTINGS', 'SYSTEM_MANAGE_SETTINGS', 'Manage system settings'),
('SYSTEM', 'BACKUP', 'SYSTEM_BACKUP', 'Perform system backup'),
('SYSTEM', 'RESTORE', 'SYSTEM_RESTORE', 'Restore system from backup');

-- Insert permission groups
INSERT INTO permission_groups (group_name, display_name, description, module, sort_order) VALUES
('USER_MANAGEMENT', 'User Management', 'Manage users and their accounts', 'USER', 1),
('ROLE_MANAGEMENT', 'Role & Permission Management', 'Manage roles and permissions', 'ROLE', 2),
('PRODUCT_MANAGEMENT', 'Product Management', 'Manage products and inventory', 'PRODUCT', 3),
('CATEGORY_MANAGEMENT', 'Category Management', 'Manage product categories', 'PRODUCT_CATEGORY', 4),
('SAMPLE_MANAGEMENT', 'Sample Management', 'Manage sample products', 'SAMPLE', 5),
('CUSTOMER_MANAGEMENT', 'Customer Management', 'Manage customers and relationships', 'CUSTOMER', 6),
('ORDER_MANAGEMENT', 'Order Management', 'Manage orders and workflow', 'ORDER', 7),
('WAREHOUSE_MANAGEMENT', 'Warehouse Management', 'Manage warehouse and inventory', 'WAREHOUSE', 8),
('FINANCIAL_MANAGEMENT', 'Financial Management', 'Manage financial data and transactions', 'FINANCE', 9),
('REPORTING', 'Reports & Analytics', 'Access reports and analytics', 'REPORT', 10),
('SYSTEM_ADMINISTRATION', 'System Administration', 'System administration and maintenance', 'SYSTEM', 11);

-- Update ADMIN role with all permissions
INSERT INTO role_permissions (role_id, permission_id, granted_by, granted_at)
SELECT 
    r.id as role_id,
    p.id as permission_id,
    1 as granted_by,
    NOW() as granted_at
FROM roles r
CROSS JOIN permissions p
WHERE r.role_name = 'ADMIN'
AND p.is_active = true;

---

-- File: migrations/000013_enhanced_permissions.down.sql
-- Tạo tại: migrations/000013_enhanced_permissions.down.sql

DROP TABLE IF EXISTS user_permissions;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permission_groups;
DROP TABLE IF EXISTS permissions;

-- Recreate simple permissions table for rollback
CREATE TABLE IF NOT EXISTS permissions (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    permission_name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_permission_name (permission_name),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INT UNSIGNED NOT NULL,
    permission_id INT UNSIGNED NOT NULL,
    module VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id),
    INDEX idx_role_permissions_permission_id (permission_id),
    CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;