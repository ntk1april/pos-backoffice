-- ============================================
-- POS BACKOFFICE - COMPLETE SETUP SCRIPT
-- Run this as SYSTEM user to set up everything
-- ============================================

-- 1. CLEANUP & USER CREATION
-- ==========================
PROMPT Creating pos_user...

BEGIN
    EXECUTE IMMEDIATE 'DROP USER pos_user CASCADE';
EXCEPTION
    WHEN OTHERS THEN NULL;
END;
/

CREATE USER pos_user IDENTIFIED BY pos_password
    DEFAULT TABLESPACE USERS
    TEMPORARY TABLESPACE TEMP
    QUOTA UNLIMITED ON USERS;

GRANT CONNECT, RESOURCE, DBA TO pos_user;
GRANT CREATE SESSION TO pos_user;
GRANT CREATE TABLE TO pos_user;
GRANT CREATE VIEW TO pos_user;
GRANT CREATE SEQUENCE TO pos_user;
GRANT UNLIMITED TABLESPACE TO pos_user;

PROMPT User created successfully.

-- 2. SWITCH TO NEW USER
-- =====================
PROMPT Connecting as pos_user...
CONNECT pos_user/pos_password@XEPDB1;

-- 3. SCHEMA SETUP
-- ===============
PROMPT Creating tables and data...

-- Drop existing (just in case)
BEGIN EXECUTE IMMEDIATE 'DROP TABLE transactions CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE products CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE stores CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE users CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE user_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE product_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE store_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE transaction_seq'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- Create Sequences
CREATE SEQUENCE user_seq START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE product_seq START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE store_seq START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE transaction_seq START WITH 1 INCREMENT BY 1 NOCACHE;

-- Create Tables
CREATE TABLE users (
    id NUMBER DEFAULT user_seq.NEXTVAL PRIMARY KEY,
    username VARCHAR2(50) UNIQUE NOT NULL,
    password_hash VARCHAR2(255) NOT NULL,
    full_name VARCHAR2(100) NOT NULL,
    role VARCHAR2(20) NOT NULL CHECK (role IN ('ADMIN', 'STAFF')),
    status VARCHAR2(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stores (
    id NUMBER DEFAULT store_seq.NEXTVAL PRIMARY KEY,
    code VARCHAR2(20) UNIQUE NOT NULL,
    name VARCHAR2(100) NOT NULL,
    address VARCHAR2(255),
    phone VARCHAR2(20),
    status VARCHAR2(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by NUMBER,
    updated_by NUMBER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

CREATE TABLE products (
    id NUMBER DEFAULT product_seq.NEXTVAL PRIMARY KEY,
    sku VARCHAR2(50) UNIQUE NOT NULL,
    name VARCHAR2(100) NOT NULL,
    description VARCHAR2(255),
    price NUMBER(10,2) NOT NULL,
    cost NUMBER(10,2) NOT NULL,
    stock NUMBER DEFAULT 0,
    status VARCHAR2(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by NUMBER,
    updated_by NUMBER,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
);

CREATE TABLE transactions (
    id NUMBER DEFAULT transaction_seq.NEXTVAL PRIMARY KEY,
    transaction_type VARCHAR2(20) NOT NULL CHECK (transaction_type IN ('INCREASE', 'DECREASE')),
    product_id NUMBER NOT NULL,
    store_id NUMBER,
    quantity NUMBER NOT NULL,
    unit_price NUMBER(10,2) NOT NULL,
    total_amount NUMBER(10,2) NOT NULL,
    notes VARCHAR2(255),
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by NUMBER NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- 4. INSERT DATA
-- ==============

-- Users
INSERT INTO users (username, password_hash, full_name, role, status) VALUES ('admin', 'admin123', 'System Administrator', 'ADMIN', 'ACTIVE');
INSERT INTO users (username, password_hash, full_name, role, status) VALUES ('staff', 'staff123', 'Staff User', 'STAFF', 'ACTIVE');

-- Stores
INSERT INTO stores (code, name, address, phone, status, created_by, updated_by) VALUES ('MB001', 'Main Branch', '123 Main Street, Bangkok', '02-123-4567', 'ACTIVE', 1, 1);
INSERT INTO stores (code, name, address, phone, status, created_by, updated_by) VALUES ('CP002', 'Central Plaza', '456 Central Plaza, Bangkok', '02-234-5678', 'ACTIVE', 1, 1);
INSERT INTO stores (code, name, address, phone, status, created_by, updated_by) VALUES ('MM003', 'Mega Mall', '789 Mega Mall, Bangkok', '02-345-6789', 'ACTIVE', 1, 1);

-- Products
INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by) VALUES ('SKU001', 'Coca Cola 330ml', 'Carbonated soft drink', 15.00, 10.00, 500, 'ACTIVE', 1, 1);
INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by) VALUES ('SKU002', 'Pepsi 330ml', 'Carbonated soft drink', 15.00, 10.00, 450, 'ACTIVE', 1, 1);
INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by) VALUES ('SKU003', 'Lays Chips 50g', 'Potato chips', 20.00, 12.00, 300, 'ACTIVE', 1, 1);
INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by) VALUES ('SKU004', 'Snickers Bar 50g', 'Chocolate bar', 25.00, 15.00, 400, 'ACTIVE', 1, 1);
INSERT INTO products (sku, name, description, price, cost, stock, status, created_by, updated_by) VALUES ('SKU005', 'Mineral Water 600ml', 'Drinking water', 10.00, 5.00, 600, 'ACTIVE', 1, 1);

-- Transactions (INCREASE)
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('INCREASE', 1, NULL, 500, 10.00, 5000.00, 'Initial stock', 1);
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('INCREASE', 2, NULL, 450, 10.00, 4500.00, 'Initial stock', 1);
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('INCREASE', 3, NULL, 300, 12.00, 3600.00, 'Initial stock', 1);
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('INCREASE', 4, NULL, 400, 15.00, 6000.00, 'Initial stock', 1);
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('INCREASE', 5, NULL, 600, 5.00, 3000.00, 'Initial stock', 1);

-- Transactions (DECREASE)
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('DECREASE', 1, 1, 50, 15.00, 750.00, 'Sold to Main Branch', 1);
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('DECREASE', 1, 2, 30, 15.00, 450.00, 'Sold to Central Plaza', 1);
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('DECREASE', 2, 1, 40, 15.00, 600.00, 'Sold to Main Branch', 1);
INSERT INTO transactions (transaction_type, product_id, store_id, quantity, unit_price, total_amount, notes, created_by) VALUES ('DECREASE', 3, 3, 25, 20.00, 500.00, 'Sold to Mega Mall', 1);

COMMIT;

PROMPT
PROMPT ====================================
PROMPT SETUP COMPLETE SUCCESSFULLY!
PROMPT ====================================
PROMPT
