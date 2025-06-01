-- +goose Up
-- Create extension first
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create account_type enum
CREATE TYPE account_type AS ENUM ('bank', 'cash', 'credit_card');

-- Create accounts table
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type account_type NOT NULL,
    balance DECIMAL(15,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for accounts
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_user_active ON accounts(user_id, is_active);

-- Create category_type enum
CREATE TYPE category_type AS ENUM ('income', 'expense');

-- Create categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    type category_type NOT NULL,
    color VARCHAR(7) DEFAULT '#6B7280',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    amount DECIMAL(15,2) NOT NULL,
    description VARCHAR(500) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for transactions
CREATE INDEX idx_transactions_account_date ON transactions(account_id, date);
CREATE INDEX idx_transactions_category_date ON transactions(category_id, date);
CREATE INDEX idx_transactions_date ON transactions(date);

-- Insert default categories
INSERT INTO categories (name, type, color) VALUES
    ('Salary', 'income', '#10B981'),
    ('Freelance', 'income', '#34D399'),
    ('Investment', 'income', '#6EE7B7'),
    ('Food & Dining', 'expense', '#EF4444'),
    ('Transportation', 'expense', '#F97316'),
    ('Shopping', 'expense', '#8B5CF6'),
    ('Entertainment', 'expense', '#EC4899'),
    ('Bills & Utilities', 'expense', '#6B7280'),
    ('Healthcare', 'expense', '#14B8A6'),
    ('Education', 'expense', '#F59E0B');

-- +goose Down
-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;

-- Drop enums
DROP TYPE IF EXISTS category_type;
DROP TYPE IF EXISTS account_type;
