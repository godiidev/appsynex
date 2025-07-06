-- File: migrations/000009_shipping_sales.up.sql
-- Tạo tại: migrations/000009_shipping_sales.up.sql
-- Mục đích: Tạo hệ thống vận chuyển và bán hàng (packing_lists, shipping, returned_goods, sales)

CREATE TABLE IF NOT EXISTS packing_lists (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_id INT UNSIGNED NOT NULL,
    fabric_roll_id INT UNSIGNED NOT NULL,
    delivery_date TIMESTAMP NULL,
    total_weight DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    packing_list_file TEXT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_order_id (order_id),
    INDEX idx_fabric_roll_id (fabric_roll_id),
    INDEX idx_delivery_date (delivery_date),
    CONSTRAINT fk_packing_lists_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    CONSTRAINT fk_packing_lists_fabric FOREIGN KEY (fabric_roll_id) REFERENCES fabric_rolls (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS shipping (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_id INT UNSIGNED NOT NULL,
    packing_list_id INT UNSIGNED NOT NULL,
    fabric_roll_id INT UNSIGNED NOT NULL,
    delivery_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    tracking_info JSON NULL,
    received_by VARCHAR(255) NULL,
    returned_quantity INT NULL DEFAULT 0,
    returned_weight DECIMAL(15,2) NULL DEFAULT 0.00,
    return_reason TEXT NULL,
    return_action TEXT NULL,
    transportation_method VARCHAR(100) NULL,
    vehicle_number VARCHAR(100) NULL,
    shipping_date TIMESTAMP NULL,
    shipping_cost DECIMAL(15,2) NULL DEFAULT 0.00,
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_order_id (order_id),
    INDEX idx_packing_list_id (packing_list_id),
    INDEX idx_fabric_roll_id (fabric_roll_id),
    INDEX idx_shipping_date (shipping_date),
    CONSTRAINT fk_shipping_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    CONSTRAINT fk_shipping_packing FOREIGN KEY (packing_list_id) REFERENCES packing_lists (id),
    CONSTRAINT fk_shipping_fabric FOREIGN KEY (fabric_roll_id) REFERENCES fabric_rolls (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS returned_goods (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    fabric_roll_id INT UNSIGNED NOT NULL,
    lot_id INT UNSIGNED NOT NULL,
    return_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    quantity INT NOT NULL,
    weight DECIMAL(15,2) NOT NULL,
    reason TEXT NOT NULL,
    action_taken TEXT NULL,
    new_fabric_roll_id INT UNSIGNED NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'received',
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_fabric_roll_id (fabric_roll_id),
    INDEX idx_lot_id (lot_id),
    INDEX idx_return_date (return_date),
    INDEX idx_new_fabric_roll_id (new_fabric_roll_id),
    CONSTRAINT fk_returned_goods_fabric FOREIGN KEY (fabric_roll_id) REFERENCES fabric_rolls (id),
    CONSTRAINT fk_returned_goods_lot FOREIGN KEY (lot_id) REFERENCES lots (id),
    CONSTRAINT fk_returned_goods_new_fabric FOREIGN KEY (new_fabric_roll_id) REFERENCES fabric_rolls (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS sales (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    customer_id INT UNSIGNED NOT NULL,
    fabric_roll_id INT UNSIGNED NOT NULL,
    quantity INT NOT NULL,
    weight DECIMAL(15,2) NOT NULL,
    sales_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    price DECIMAL(15,2) NOT NULL,
    invoice_id INT UNSIGNED NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'completed',
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_fabric_roll_id (fabric_roll_id),
    INDEX idx_sales_date (sales_date),
    CONSTRAINT fk_sales_customer FOREIGN KEY (customer_id) REFERENCES customers (id),
    CONSTRAINT fk_sales_fabric FOREIGN KEY (fabric_roll_id) REFERENCES fabric_rolls (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;