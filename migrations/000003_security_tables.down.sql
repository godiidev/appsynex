-- File: migrations/000003_security_tables.down.sql
-- Tạo tại: migrations/000003_security_tables.down.sql
DROP TABLE IF EXISTS password_resets;
DROP TABLE IF EXISTS security_settings;
DROP TABLE IF EXISTS access_logs;
DROP TABLE IF EXISTS sessions;

-- File: migrations/000004_order_system.down.sql
-- Tạo tại: migrations/000004_order_system.down.sql
DROP TABLE IF EXISTS customer_saved_items;
DROP TABLE IF EXISTS customer_search_history;
DROP TABLE IF EXISTS customer_activity_log;
DROP TABLE IF EXISTS customer_orders;
DROP TABLE IF EXISTS order_histories;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS order_statuses;
DROP TABLE IF EXISTS customers;

-- File: migrations/000005_yarn_system.down.sql
-- Tạo tại: migrations/000005_yarn_system.down.sql
DROP TABLE IF EXISTS yarn_orders;
DROP TABLE IF EXISTS yarn_inventory_transactions;
DROP TABLE IF EXISTS yarn_boxes;

-- File: migrations/000006_weaving_system.down.sql
-- Tạo tại: migrations/000006_weaving_system.down.sql
DROP TABLE IF EXISTS weaving_financials;
DROP TABLE IF EXISTS weaving_quality_control;
DROP TABLE IF EXISTS weaving_operations;
DROP TABLE IF EXISTS weaving_orders;
DROP TABLE IF EXISTS weaving_staff;
DROP TABLE IF EXISTS weaving_shifts;
DROP TABLE IF EXISTS weaving_facilities;

-- File: migrations/000007_dyeing_system.down.sql
-- Tạo tại: migrations/000007_dyeing_system.down.sql
DROP TABLE IF EXISTS dyeing_lots;
DROP TABLE IF EXISTS dyeing_subcontractors;

-- File: migrations/000008_warehouse_system.down.sql
-- Tạo tại: migrations/000008_warehouse_system.down.sql
DROP TABLE IF EXISTS lot_splits;
DROP TABLE IF EXISTS lot_quality_issues;
DROP TABLE IF EXISTS lot_histories;
DROP TABLE IF EXISTS warehouse_transactions;
DROP TABLE IF EXISTS fabric_rolls;
DROP TABLE IF EXISTS lots;
DROP TABLE IF EXISTS warehouses;

-- File: migrations/000009_shipping_sales.down.sql
-- Tạo tại: migrations/000009_shipping_sales.down.sql
DROP TABLE IF EXISTS sales;
DROP TABLE IF EXISTS returned_goods;
DROP TABLE IF EXISTS shipping;
DROP TABLE IF EXISTS packing_lists;

-- File: migrations/000010_sample_greige.down.sql
-- Tạo tại: migrations/000010_sample_greige.down.sql
DROP TABLE IF EXISTS greige_labels;
DROP TABLE IF EXISTS greige_transactions;
DROP TABLE IF EXISTS greige_fabrics;
DROP TABLE IF EXISTS sample_labels;
DROP TABLE IF EXISTS sample_dispatches;
DROP TABLE IF EXISTS sample_transactions;
DROP TABLE IF EXISTS sample_images;

-- File: migrations/000011_lapdip_system.down.sql
-- Tạo tại: migrations/000011_lapdip_system.down.sql
DROP TABLE IF EXISTS lap_dip_final_selection;
DROP TABLE IF EXISTS lap_dip_dispatches;
DROP TABLE IF EXISTS lap_dip_tests;