-- ============================================
-- RESET DATABASE - Delete All Data
-- ============================================

-- Delete all data (in correct order due to foreign keys)
DELETE FROM stock_logs;
DELETE FROM products;
DELETE FROM users;

-- Reset sequences
DROP SEQUENCE user_seq;
DROP SEQUENCE product_seq;
DROP SEQUENCE stock_log_seq;

CREATE SEQUENCE user_seq START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE product_seq START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE stock_log_seq START WITH 1 INCREMENT BY 1 NOCACHE;

-- Insert default users with PLAIN TEXT passwords (for testing)
INSERT INTO users (username, password_hash, full_name, role, status)
VALUES ('admin', 'admin123', 'System Administrator', 'ADMIN', 'ACTIVE');

INSERT INTO users (username, password_hash, full_name, role, status)
VALUES ('staff', 'staff123', 'Staff User', 'STAFF', 'ACTIVE');

INSERT INTO users (username, password_hash, full_name, role, status)
VALUES ('test', 'test1', 'Test User', 'ADMIN', 'ACTIVE');

-- Insert sample products
INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by)
VALUES ('SKU001', 'Coca Cola 330ml', 'Carbonated soft drink', 15.00, 10.00, 100, 'ACTIVE', 1, 1);

INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by)
VALUES ('SKU002', 'Pepsi 330ml', 'Carbonated soft drink', 15.00, 10.00, 80, 'ACTIVE', 1, 1);

INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by)
VALUES ('SKU003', 'Lays Chips Original 50g', 'Potato chips', 20.00, 12.00, 150, 'ACTIVE', 1, 1);

INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by)
VALUES ('SKU004', 'Snickers Bar 50g', 'Chocolate bar', 25.00, 15.00, 200, 'ACTIVE', 1, 1);

INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by)
VALUES ('SKU005', 'Mineral Water 600ml', 'Drinking water', 10.00, 5.00, 300, 'ACTIVE', 1, 1);

COMMIT;

-- Verify reset
SELECT 'Users:' as info, COUNT(*) as count FROM users
UNION ALL
SELECT 'Products:', COUNT(*) FROM products
UNION ALL
SELECT 'Stock Logs:', COUNT(*) FROM stock_logs;

SELECT 'Database reset complete!' as status FROM dual;
