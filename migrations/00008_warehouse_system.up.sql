-- File: migrations/000008_warehouse_system.up.sql
-- Tạo tại: migrations/000008_warehouse_system.up.sql  
-- Mục đích: Tạo hệ thống quản lý kho (warehouses, lots, fabric_rolls, warehouse_transactions, lot_histories, lot_quality_issues, lot_splits)

CREATE TABLE IF NOT EXISTS warehouses (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    capacity INT NOT NULL DEFAULT 0,
    created_by INT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    INDEX idx_name (name),
    INDEX idx_created_by (created_by),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_warehouses_user FOREIGN KEY (created_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lots (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    lot_code VARCHAR(100) NOT NULL,
    origin VARCHAR(255) NOT NULL, -- 'Dyeing', 'Purchase'
    product_id INT UNSIGNED NOT NULL,
    color_code VARCHAR(100) NULL,
    fabric_type VARCHAR(255) NULL,
    quantity INT NOT NULL DEFAULT 0,
    total_weight DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    location VARCHAR(255) NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'available',
    quality_status VARCHAR(50) NOT NULL DEFAULT 'good',
    origin_order_id INT UNSIGNED NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    notes TEXT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_lot_code (lot_code),
    INDEX idx_product_id (product_id),
    INDEX idx_origin_order_id (origin_order_id),
    INDEX idx_status (status),
    CONSTRAINT fk_lots_product FOREIGN KEY (product_id) REFERENCES products (id),
    CONSTRAINT fk_lots_origin_order FOREIGN KEY (origin_order_id) REFERENCES orders (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS fabric_rolls (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    lot_id INT UNSIGNED NOT NULL,
    warehouse_id INT UNSIGNED NOT NULL,
    barcode VARCHAR(100) NOT NULL,
    weight DECIMAL(10,2) NOT NULL,
    length DECIMAL(10,2) NOT NULL,
    location VARCHAR(255) NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'available',
    created_by INT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    notes TEXT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_barcode (barcode),
    INDEX idx_lot_id (lot_id),
    INDEX idx_warehouse_id (warehouse_id),
    INDEX idx_created_by (created_by),
    CONSTRAINT fk_fabric_rolls_lot FOREIGN KEY (lot_id) REFERENCES lots (id) ON DELETE CASCADE,
    CONSTRAINT fk_fabric_rolls_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses (id),
    CONSTRAINT fk_fabric_rolls_user FOREIGN KEY (created_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS warehouse_transactions (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    warehouse_id INT UNSIGNED NOT NULL,
    fabric_roll_id INT UNSIGNED NULL,
    lot_id INT UNSIGNED NULL,
    transaction_type VARCHAR(50) NOT NULL, -- 'IN', 'OUT', 'MOVE', 'ADJUST'
    transaction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    quantity INT NOT NULL,
    weight DECIMAL(15,2) NOT NULL,
    remaining_quantity INT NOT NULL,
    remaining_weight DECIMAL(15,2) NOT NULL,
    handled_by INT UNSIGNED NOT NULL,
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_warehouse_id (warehouse_id),
    INDEX idx_fabric_roll_id (fabric_roll_id),
    INDEX idx_lot_id (lot_id),
    INDEX idx_transaction_date (transaction_date),
    INDEX idx_handled_by (handled_by),
    CONSTRAINT fk_warehouse_transactions_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses (id),
    CONSTRAINT fk_warehouse_transactions_fabric FOREIGN KEY (fabric_roll_id) REFERENCES fabric_rolls (id) ON DELETE SET NULL,
    CONSTRAINT fk_warehouse_transactions_lot FOREIGN KEY (lot_id) REFERENCES lots (id) ON DELETE SET NULL,
    CONSTRAINT fk_warehouse_transactions_user FOREIGN KEY (handled_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lot_histories (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    lot_id INT UNSIGNED NOT NULL,
    action_type VARCHAR(255) NOT NULL,
    previous_lot_code VARCHAR(100) NULL,
    new_lot_code VARCHAR(100) NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    changed_by INT UNSIGNED NOT NULL,
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_lot_id (lot_id),
    INDEX idx_changed_at (changed_at),
    INDEX idx_changed_by (changed_by),
    CONSTRAINT fk_lot_histories_lot FOREIGN KEY (lot_id) REFERENCES lots (id) ON DELETE CASCADE,
    CONSTRAINT fk_lot_histories_user FOREIGN KEY (changed_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lot_quality_issues (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    lot_id INT UNSIGNED NOT NULL,
    issue_type VARCHAR(255) NOT NULL,
    detected_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    resolved_by INT UNSIGNED NULL,
    resolution_action TEXT NULL,
    new_lot_code VARCHAR(100) NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_lot_id (lot_id),
    INDEX idx_detected_at (detected_at),
    INDEX idx_resolved_by (resolved_by),
    CONSTRAINT fk_lot_quality_lot FOREIGN KEY (lot_id) REFERENCES lots (id) ON DELETE CASCADE,
    CONSTRAINT fk_lot_quality_user FOREIGN KEY (resolved_by) REFERENCES users (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lot_splits (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    original_lot_id INT UNSIGNED NOT NULL,
    new_lot_id INT UNSIGNED NOT NULL,
    quantity INT NOT NULL,
    weight DECIMAL(15,2) NOT NULL,
    fabric_roll_info JSON NULL,
    split_reason TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_original_lot_id (original_lot_id),
    INDEX idx_new_lot_id (new_lot_id),
    CONSTRAINT fk_lot_splits_original FOREIGN KEY (original_lot_id) REFERENCES lots (id) ON DELETE CASCADE,
    CONSTRAINT fk_lot_splits_new FOREIGN KEY (new_lot_id) REFERENCES lots (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;