-- File: migrations/000005_yarn_system.up.sql
-- Tạo tại: migrations/000005_yarn_system.up.sql
-- Mục đích: Tạo hệ thống quản lý sợi (yarn_boxes, yarn_inventory_transactions, yarn_orders)

CREATE TABLE IF NOT EXISTS yarn_boxes (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    box_code VARCHAR(100) NOT NULL,
    yarn_type VARCHAR(255) NOT NULL,
    cone_quantity INT NOT NULL DEFAULT 0,
    total_weight DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    warehouse_location VARCHAR(255) NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'available',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_box_code (box_code),
    INDEX idx_yarn_type (yarn_type),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS yarn_inventory_transactions (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    yarn_box_id INT UNSIGNED NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, -- 'IN', 'OUT'
    transaction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    quantity INT NOT NULL,
    weight DECIMAL(15,2) NOT NULL,
    remaining_quantity INT NOT NULL,
    remaining_weight DECIMAL(15,2) NOT NULL,
    handled_by INT UNSIGNED NOT NULL,
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_yarn_box_id (yarn_box_id),
    INDEX idx_transaction_date (transaction_date),
    INDEX idx_handled_by (handled_by),
    CONSTRAINT fk_yarn_transactions_box FOREIGN KEY (yarn_box_id) REFERENCES yarn_boxes (id) ON DELETE CASCADE,
    CONSTRAINT fk_yarn_transactions_user FOREIGN KEY (handled_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS yarn_orders (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_id INT UNSIGNED NOT NULL,
    yarn_box_id INT UNSIGNED NOT NULL,
    yarn_type VARCHAR(255) NOT NULL,
    cone_quantity INT NOT NULL,
    total_weight DECIMAL(15,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'allocated',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_order_id (order_id),
    INDEX idx_yarn_box_id (yarn_box_id),
    CONSTRAINT fk_yarn_orders_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    CONSTRAINT fk_yarn_orders_box FOREIGN KEY (yarn_box_id) REFERENCES yarn_boxes (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;