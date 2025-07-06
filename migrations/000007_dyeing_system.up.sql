-- File: migrations/000007_dyeing_system.up.sql  
-- Tạo tại: migrations/000007_dyeing_system.up.sql
-- Mục đích: Tạo hệ thống quản lý nhuộm (dyeing_subcontractors, dyeing_lots)

CREATE TABLE IF NOT EXISTS dyeing_subcontractors (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    assigned_order_id INT UNSIGNED NULL,
    fabric_roll_id INT UNSIGNED NULL,
    dye_color VARCHAR(255) NULL,
    color_fastness VARCHAR(50) NULL,
    shrinkage VARCHAR(50) NULL,
    dyeing_cost_per_unit DECIMAL(15,2) NULL,
    dyeing_start_time TIMESTAMP NULL,
    dyeing_end_time TIMESTAMP NULL,
    dyeing_quantity INT NULL,
    special_requirements TEXT NULL,
    contact_info TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    INDEX idx_name (name),
    INDEX idx_assigned_order_id (assigned_order_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_dyeing_subcontractors_order FOREIGN KEY (assigned_order_id) REFERENCES orders (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS dyeing_lots (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    dyeing_subcontractor_id INT UNSIGNED NOT NULL,
    lot_code VARCHAR(100) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    weight DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    quality_issues TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    notes TEXT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_lot_code (lot_code),
    INDEX idx_dyeing_subcontractor_id (dyeing_subcontractor_id),
    INDEX idx_status (status),
    CONSTRAINT fk_dyeing_lots_subcontractor FOREIGN KEY (dyeing_subcontractor_id) REFERENCES dyeing_subcontractors (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;