-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('employee', 'manager')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create expenses table
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    amount_idr INTEGER NOT NULL CHECK (amount_idr >= 10000 AND amount_idr <= 50000000),
    description TEXT NOT NULL,
    receipt_url VARCHAR(512),
    status VARCHAR(50) NOT NULL CHECK (status IN ('awaiting_approval', 'approved', 'rejected', 'completed')),
    auto_approved BOOLEAN NOT NULL DEFAULT FALSE,
    submitted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    payment_id VARCHAR(255),
    payment_external_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create approvals table
CREATE TABLE IF NOT EXISTS approvals (
    id SERIAL PRIMARY KEY,
    expense_id INTEGER NOT NULL REFERENCES expenses(id),
    approver_id INTEGER NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL CHECK (status IN ('approved', 'rejected')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    expense_id INTEGER NOT NULL REFERENCES expenses(id),
    user_id INTEGER REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    old_status VARCHAR(50),
    new_status VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_expenses_user_id ON expenses(user_id);
CREATE INDEX idx_expenses_status ON expenses(status);
CREATE INDEX idx_expenses_payment_external_id ON expenses(payment_external_id);
CREATE INDEX idx_approvals_expense_id ON approvals(expense_id);
CREATE INDEX idx_audit_logs_expense_id ON audit_logs(expense_id);
CREATE INDEX idx_users_email ON users(email);
