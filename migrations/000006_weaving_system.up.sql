-- File: migrations/000006_weaving_system.up.sql
-- Tạo tại: migrations/000006_weaving_system.up.sql
-- Mục đích: Tạo hệ thống quản lý dệt (weaving_facilities, weaving_orders, weaving_shifts, weaving_staff, weaving_operations, weaving_quality_control, weaving_financials)

CREATE TABLE IF NOT EXISTS weaving_facilities (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    facility_name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    contact_info TEXT NULL,
    facility_type VARCHAR(50) NOT NULL DEFAULT 'internal', -- 'internal', 'external'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    INDEX idx_facility_name (facility_name),
    INDEX idx_facility_type (facility_type),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weaving_shifts (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    weaving_facility_id INT UNSIGNED NOT NULL,
    shift_name VARCHAR(100) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    shift_leader VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_weaving_facility_id (weaving_facility_id),
    INDEX idx_shift_name (shift_name),
    CONSTRAINT fk_weaving_shifts_facility FOREIGN KEY (weaving_facility_id) REFERENCES weaving_facilities (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weaving_staff (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    weaving_facility_id INT UNSIGNED NOT NULL,
    staff_name VARCHAR(255) NOT NULL,
    position VARCHAR(100) NOT NULL,
    contact_info TEXT NULL,
    shift_id INT UNSIGNED NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    INDEX idx_weaving_facility_id (weaving_facility_id),
    INDEX idx_shift_id (shift_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_weaving_staff_facility FOREIGN KEY (weaving_facility_id) REFERENCES weaving_facilities (id) ON DELETE CASCADE,
    CONSTRAINT fk_weaving_staff_shift FOREIGN KEY (shift_id) REFERENCES weaving_shifts (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weaving_orders (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    order_code VARCHAR(100) NOT NULL,
    weaving_facility_id INT UNSIGNED NOT NULL,
    linked_order_id INT UNSIGNED NOT NULL,
    yarn_box_id INT UNSIGNED NOT NULL,
    product_id INT UNSIGNED NOT NULL,
    start_date TIMESTAMP NULL,
    end_date TIMESTAMP NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_order_code (order_code),
    INDEX idx_weaving_facility_id (weaving_facility_id),
    INDEX idx_linked_order_id (linked_order_id),
    INDEX idx_yarn_box_id (yarn_box_id),
    INDEX idx_product_id (product_id),
    CONSTRAINT fk_weaving_orders_facility FOREIGN KEY (weaving_facility_id) REFERENCES weaving_facilities (id),
    CONSTRAINT fk_weaving_orders_order FOREIGN KEY (linked_order_id) REFERENCES orders (id),
    CONSTRAINT fk_weaving_orders_yarn FOREIGN KEY (yarn_box_id) REFERENCES yarn_boxes (id),
    CONSTRAINT fk_weaving_orders_product FOREIGN KEY (product_id) REFERENCES products (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weaving_operations (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    weaving_order_id INT UNSIGNED NOT NULL,
    weaving_facility_id INT UNSIGNED NOT NULL,
    weaving_staff_id INT UNSIGNED NOT NULL,
    yarn_box_id INT UNSIGNED NOT NULL,
    weaving_machine_id INT UNSIGNED NULL,
    cam_layout TEXT NULL,
    start_time TIMESTAMP NULL,
    end_time TIMESTAMP NULL,
    quantity INT NOT NULL DEFAULT 0,
    quality_check_status VARCHAR(50) DEFAULT 'pending',
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_weaving_order_id (weaving_order_id),
    INDEX idx_weaving_facility_id (weaving_facility_id),
    INDEX idx_weaving_staff_id (weaving_staff_id),
    INDEX idx_yarn_box_id (yarn_box_id),
    CONSTRAINT fk_weaving_operations_order FOREIGN KEY (weaving_order_id) REFERENCES weaving_orders (id) ON DELETE CASCADE,
    CONSTRAINT fk_weaving_operations_facility FOREIGN KEY (weaving_facility_id) REFERENCES weaving_facilities (id),
    CONSTRAINT fk_weaving_operations_staff FOREIGN KEY (weaving_staff_id) REFERENCES weaving_staff (id),
    CONSTRAINT fk_weaving_operations_yarn FOREIGN KEY (yarn_box_id) REFERENCES yarn_boxes (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weaving_quality_control (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    weaving_order_id INT UNSIGNED NOT NULL,
    quality_inspector VARCHAR(255) NOT NULL,
    quality_check_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    defect_found BOOLEAN NOT NULL DEFAULT FALSE,
    defect_type VARCHAR(255) NULL,
    corrective_action TEXT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'passed',
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_weaving_order_id (weaving_order_id),
    INDEX idx_quality_check_time (quality_check_time),
    CONSTRAINT fk_weaving_quality_order FOREIGN KEY (weaving_order_id) REFERENCES weaving_orders (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS weaving_financials (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    weaving_order_id INT UNSIGNED NOT NULL,
    customer_id INT UNSIGNED NOT NULL,
    total_cost DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    payment_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    due_date TIMESTAMP NULL,
    invoice_id INT UNSIGNED NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_weaving_order_id (weaving_order_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_payment_status (payment_status),
    CONSTRAINT fk_weaving_financials_order FOREIGN KEY (weaving_order_id) REFERENCES weaving_orders (id) ON DELETE CASCADE,
    CONSTRAINT fk_weaving_financials_customer FOREIGN KEY (customer_id) REFERENCES customers (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;