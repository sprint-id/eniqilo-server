BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS USERS (
    id UUID PRIMARY KEY,
    phone_number VARCHAR UNIQUE,
    name VARCHAR,
    password VARCHAR,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES USERS(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    sku VARCHAR NOT NULL,
    category VARCHAR NOT NULL,
    image_url VARCHAR NOT NULL,
    notes VARCHAR NOT NULL,
    price DECIMAL NOT NULL,
    stock INTEGER NOT NULL,
    location VARCHAR NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    staff_id UUID REFERENCES USERS(id) ON DELETE CASCADE,
    phone_number VARCHAR UNIQUE,
    name VARCHAR NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id) ON DELETE CASCADE,
    paid DECIMAL NOT NULL,
    change DECIMAL NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE INDEX idx_users_phone ON users (phone_number);
CREATE INDEX idx_users_id ON users (id);
CREATE INDEX idx_products_id ON products (id);
CREATE INDEX idx_products_user_id ON products (user_id);
CREATE INDEX idx_customers_id ON customers (id);
CREATE INDEX idx_customers_staff_id ON customers (staff_id);
CREATE INDEX idx_transactions_id ON transactions (id);
CREATE INDEX idx_transactions_customer_id ON transactions (customer_id);
CREATE INDEX idx_transaction_details_id ON transaction_details (id);
CREATE INDEX idx_transaction_details_transaction_id ON transaction_details (transaction_id);
CREATE INDEX idx_transaction_details_product_id ON transaction_details (product_id);

COMMIT TRANSACTION;