DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_products_id;
DROP INDEX IF EXISTS idx_products_user_id;
DROP INDEX IF EXISTS idx_customers_id;
DROP INDEX IF EXISTS idx_customers_staff_id;
DROP INDEX IF EXISTS idx_transactions_id;
DROP INDEX IF EXISTS idx_transactions_customer_id;
DROP INDEX IF EXISTS idx_transaction_details_id;
DROP INDEX IF EXISTS idx_transaction_details_transaction_id;
DROP INDEX IF EXISTS idx_transaction_details_product_id;

DROP TABLE IF EXISTS transaction_details;

DROP TABLE IF EXISTS transactions;

DROP TABLE IF EXISTS customers;

DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS USERS;