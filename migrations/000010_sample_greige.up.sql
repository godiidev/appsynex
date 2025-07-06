-- File: migrations/000010_sample_greige.up.sql
-- Tạo tại: migrations/000010_sample_greige.up.sql
-- Mục đích: Tạo hệ thống quản lý mẫu và vải mộc (sample_images, sample_transactions, sample_dispatches, sample_labels, greige_fabrics, greige_transactions, greige_labels)

-- Sample images table (already exists in 000002 but need to add here for completeness)
CREATE TABLE IF NOT EXISTS sample_images (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    sample_product_id INT UNSIGNED NOT NULL,
    image_url TEXT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_sample_product_id (sample_product_id),
    CONSTRAINT fk_sample_images_sample FOREIGN KEY (sample_product_id) REFERENCES sample_products (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS sample_transactions (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    sample_product_id INT UNSIGNED NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, -- 'IN', 'OUT'
    quantity DECIMAL(10,2) NOT NULL,
    transaction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    remaining_quantity DECIMAL(10,2) NOT NULL,
    handled_by INT UNSIGNED NOT NULL,
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_sample_product_id (sample_product_id),
    INDEX idx_transaction_date (transaction_date),
    INDEX idx_handled_by (handled_by),
    CONSTRAINT fk_sample_transactions_sample FOREIGN KEY (sample_product_id) REFERENCES sample_products (id) ON DELETE CASCADE,
    CONSTRAINT fk_sample_transactions_user FOREIGN KEY (handled_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS sample_dispatches (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    sample_product_id INT UNSIGNED NOT NULL,
    customer_id INT UNSIGNED NOT NULL,
    dispatch_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dispatch_quantity DECIMAL(10,2) NOT NULL,
    dispatch_weight DECIMAL(10,2) NOT NULL,
    dispatch_color VARCHAR(255) NULL,
    lot_number VARCHAR(100) NULL,
    tracking_number VARCHAR(100) NULL,
    dispatch_notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_sample_product_id (sample_product_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_dispatch_date (dispatch_date),
    CONSTRAINT fk_sample_dispatches_sample FOREIGN KEY (sample_product_id) REFERENCES sample_products (id) ON DELETE CASCADE,
    CONSTRAINT fk_sample_dispatches_customer FOREIGN KEY (customer_id) REFERENCES customers (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS sample_labels (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    sample_product_id INT UNSIGNED NOT NULL,
    label_type VARCHAR(50) NOT NULL, -- 'barcode', 'info'
    label_content JSON NOT NULL,
    printed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    printed_by INT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_sample_product_id (sample_product_id),
    INDEX idx_printed_at (printed_at),
    INDEX idx_printed_by (printed_by),
    CONSTRAINT fk_sample_labels_sample FOREIGN KEY (sample_product_id) REFERENCES sample_products (id) ON DELETE CASCADE,
    CONSTRAINT fk_sample_labels_user FOREIGN KEY (printed_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Greige fabrics system
CREATE TABLE IF NOT EXISTS greige_fabrics (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    greige_code VARCHAR(100) NOT NULL,
    fabric_type VARCHAR(255) NOT NULL,
    weight DECIMAL(10,2) NOT NULL,
    width DECIMAL(10,2) NOT NULL,
    color VARCHAR(255) NULL,
    quality VARCHAR(50) NULL,
    fiber_content VARCHAR(255) NULL,
    yarn_type VARCHAR(255) NULL,
    weave_pattern VARCHAR(255) NULL,
    cam_layout TEXT NULL,
    weaving_machine VARCHAR(255) NULL,
    weaving_facility_id INT UNSIGNED NULL,
    source VARCHAR(255) NULL,
    location VARCHAR(255) NULL,
    remaining_quantity DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    barcode VARCHAR(100) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_greige_code (greige_code),
    INDEX idx_fabric_type (fabric_type),
    INDEX idx_weaving_facility_id (weaving_facility_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_greige_fabrics_facility FOREIGN KEY (weaving_facility_id) REFERENCES weaving_facilities (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS greige_transactions (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    greige_fabric_id INT UNSIGNED NOT NULL,
    transaction_type VARCHAR(50) NOT NULL, -- 'IN', 'OUT', 'MOVE'
    quantity DECIMAL(10,2) NOT NULL,
    transaction_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    remaining_quantity DECIMAL(10,2) NOT NULL,
    handled_by INT UNSIGNED NOT NULL,
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_greige_fabric_id (greige_fabric_id),
    INDEX idx_transaction_date (transaction_date),
    INDEX idx_handled_by (handled_by),
    CONSTRAINT fk_greige_transactions_greige FOREIGN KEY (greige_fabric_id) REFERENCES greige_fabrics (id) ON DELETE CASCADE,
    CONSTRAINT fk_greige_transactions_user FOREIGN KEY (handled_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS greige_labels (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    greige_fabric_id INT UNSIGNED NOT NULL,
    label_type VARCHAR(50) NOT NULL, -- 'barcode', 'info'
    label_content JSON NOT NULL,
    printed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    printed_by INT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_greige_fabric_id (greige_fabric_id),
    INDEX idx_printed_at (printed_at),
    INDEX idx_printed_by (printed_by),
    CONSTRAINT fk_greige_labels_greige FOREIGN KEY (greige_fabric_id) REFERENCES greige_fabrics (id) ON DELETE CASCADE,
    CONSTRAINT fk_greige_labels_user FOREIGN KEY (printed_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;