CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    account_type VARCHAR(10) NOT NULL DEFAULT 'User',
    organisation VARCHAR(255),
    is_confirmed BOOLEAN DEFAULT false,
    refresh_token TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
