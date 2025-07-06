-- File: migrations/000011_lapdip_system.up.sql
-- Tạo tại: migrations/000011_lapdip_system.up.sql  
-- Mục đích: Tạo hệ thống kiểm tra lap dip (lap_dip_tests, lap_dip_dispatches, lap_dip_final_selection)

CREATE TABLE IF NOT EXISTS lap_dip_tests (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    ld_code VARCHAR(100) NOT NULL,
    customer_id INT UNSIGNED NOT NULL,
    color_name VARCHAR(255) NOT NULL,
    color_code VARCHAR(100) NOT NULL,
    greige_fabric_id INT UNSIGNED NOT NULL,
    dyeing_subcontractor_id INT UNSIGNED NOT NULL,
    test_start_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    test_end_date TIMESTAMP NULL,
    ld_results JSON NULL, -- Array of lap dip codes from dyeing facility
    selected_ld_code VARCHAR(100) NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'testing',
    initial_sample_image TEXT NULL,
    result_image TEXT NULL,
    customer_approval_image TEXT NULL,
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_ld_code (ld_code),
    INDEX idx_customer_id (customer_id),
    INDEX idx_greige_fabric_id (greige_fabric_id),
    INDEX idx_dyeing_subcontractor_id (dyeing_subcontractor_id),
    INDEX idx_test_start_date (test_start_date),
    CONSTRAINT fk_lap_dip_tests_customer FOREIGN KEY (customer_id) REFERENCES customers (id),
    CONSTRAINT fk_lap_dip_tests_greige FOREIGN KEY (greige_fabric_id) REFERENCES greige_fabrics (id),
    CONSTRAINT fk_lap_dip_tests_dyeing FOREIGN KEY (dyeing_subcontractor_id) REFERENCES dyeing_subcontractors (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lap_dip_dispatches (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    ld_test_id INT UNSIGNED NOT NULL,
    dispatch_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    dispatch_code VARCHAR(100) NOT NULL, -- Combined code from dyeing facility and company
    customer_response TEXT NULL,
    response_date TIMESTAMP NULL,
    final_ld_code VARCHAR(100) NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'dispatched',
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_ld_test_id (ld_test_id),
    INDEX idx_dispatch_date (dispatch_date),
    INDEX idx_response_date (response_date),
    CONSTRAINT fk_lap_dip_dispatches_test FOREIGN KEY (ld_test_id) REFERENCES lap_dip_tests (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS lap_dip_final_selection (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    ld_test_id INT UNSIGNED NOT NULL,
    final_ld_code VARCHAR(100) NOT NULL,
    selected_by_customer_id INT UNSIGNED NOT NULL,
    order_id INT UNSIGNED NULL,
    final_selection_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT NULL,
    PRIMARY KEY (id),
    INDEX idx_ld_test_id (ld_test_id),
    INDEX idx_selected_by_customer_id (selected_by_customer_id),
    INDEX idx_order_id (order_id),
    CONSTRAINT fk_lap_dip_final_test FOREIGN KEY (ld_test_id) REFERENCES lap_dip_tests (id) ON DELETE CASCADE,
    CONSTRAINT fk_lap_dip_final_customer FOREIGN KEY (selected_by_customer_id) REFERENCES customers (id),
    CONSTRAINT fk_lap_dip_final_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;