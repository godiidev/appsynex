-- File: init.sql  
-- Tạo tại: init.sql (root của project)
-- Mục đích: Initialize database với charset phù hợp

CREATE DATABASE IF NOT EXISTS appsynex
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

USE appsynex;

-- Grant permissions
GRANT ALL PRIVILEGES ON appsynex.* TO 'appsynex_user'@'%';
FLUSH PRIVILEGES;