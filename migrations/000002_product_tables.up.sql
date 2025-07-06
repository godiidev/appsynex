-- File: migrations/000002_product_tables.up.sql
-- Tạo tại: migrations/000002_product_tables.up.sql
-- Mục đích: Tạo bảng product_categories, product_names, products, sample_products

CREATE TABLE IF NOT EXISTS product_categories (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    category_name VARCHAR(255) NOT NULL,
    parent_category_id INT UNSIGNED NULL,
    description TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    INDEX idx_category_name (category_name),
    INDEX idx_parent_category_id (parent_category_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_product_categories_parent FOREIGN KEY (parent_category_id) REFERENCES product_categories (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS product_names (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    product_name_vi VARCHAR(255) NOT NULL,
    product_name_en VARCHAR(255) NOT NULL,
    sku_parent VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_sku_parent (sku_parent),
    INDEX idx_product_name_vi (product_name_vi),
    INDEX idx_product_name_en (product_name_en),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS products (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    product_name_id INT UNSIGNED NOT NULL,
    category_id INT UNSIGNED NOT NULL,
    sku VARCHAR(100) NOT NULL,
    sku_variant VARCHAR(100) NULL,
    description TEXT NULL,
    fabric_type VARCHAR(255) NULL,
    weight DECIMAL(10,2) NULL,
    width DECIMAL(10,2) NULL,
    color VARCHAR(150) NULL,
    quality VARCHAR(50) NULL,
    fiber_content VARCHAR(255) NULL,
    additional_info JSON NULL,
    price DECIMAL(15,2) NULL,
    sales_price DECIMAL(15,2) NULL,
    stock_quantity DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_sku (sku),
    INDEX idx_product_name_id (product_name_id),
    INDEX idx_category_id (category_id),
    INDEX idx_fabric_type (fabric_type),
    INDEX idx_color (color),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_products_product_name FOREIGN KEY (product_name_id) REFERENCES product_names (id),
    CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES product_categories (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS sample_products (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    sku VARCHAR(100) NOT NULL,
    product_name_id INT UNSIGNED NOT NULL,
    category_id INT UNSIGNED NOT NULL,
    description TEXT NULL,
    sample_type VARCHAR(255) NULL,
    weight DECIMAL(10,2) NULL,
    width DECIMAL(10,2) NULL,
    color VARCHAR(255) NULL,
    color_code VARCHAR(50) NULL,
    quality TEXT NULL,
    remaining_quantity INT NOT NULL DEFAULT 0,
    fiber_content VARCHAR(255) NULL,
    source VARCHAR(255) NULL,
    sample_location VARCHAR(255) NULL,
    barcode VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_sku (sku),
    INDEX idx_product_name_id (product_name_id),
    INDEX idx_category_id (category_id),
    INDEX idx_sample_type (sample_type),
    INDEX idx_color (color),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_sample_products_product_name FOREIGN KEY (product_name_id) REFERENCES product_names (id),
    CONSTRAINT fk_sample_products_category FOREIGN KEY (category_id) REFERENCES product_categories (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;