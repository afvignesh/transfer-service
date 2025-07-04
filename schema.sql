-- schema.sql
CREATE TABLE IF NOT EXISTS accounts (
    id INT PRIMARY KEY,
    balance NUMERIC(15,5) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transaction table (To log transactions)
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    source_account_id INT NOT NULL,
    destination_account_id INT NOT NULL,
    amount NUMERIC(15,5) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
