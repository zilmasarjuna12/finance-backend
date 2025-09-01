-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    created_at bigint,
    updated_at bigint,
    deleted_at bigint
);

CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(36) NOT NULL,
    session_token TEXT UNIQUE NOT NULL,
    expires_at bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    deleted_at bigint,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS wallets (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    balance DECIMAL(15,2) NOT NULL CHECK (balance >= 0),

    created_at bigint,
    updated_at bigint,
    deleted_at bigint
);

CREATE TABLE IF NOT EXISTS has_wallets (
    user_id VARCHAR(36) NOT NULL,
    wallet_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (user_id, wallet_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS budgets (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    type VARCHAR(100) NOT NULL,
    category VARCHAR(100) NOT NULL,

    created_at bigint,
    updated_at bigint,
    deleted_at bigint
);

CREATE TABLE IF NOT EXISTS has_budgets (
    user_id VARCHAR(36) NOT NULL,
    budget_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (user_id, budget_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (budget_id) REFERENCES budgets(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    amount DECIMAL(15,2) NOT NULL CHECK (amount >= 0),
    type VARCHAR(100) NOT NULL,
    note TEXT,
    transaction_date bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    deleted_at bigint,

    budget_id VARCHAR(36),
    wallet_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE RESTRICT,
    FOREIGN KEY (budget_id) REFERENCES budgets(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS has_transactions (
    user_id VARCHAR(36) NOT NULL,
    transaction_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (user_id, transaction_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_sessions_session_token ON sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_sessions_deleted_at ON sessions(deleted_at);

CREATE INDEX IF NOT EXISTS idx_transactions_note ON transactions(note);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS has_wallets;
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS has_budgets;
DROP TABLE IF EXISTS budgets;
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
