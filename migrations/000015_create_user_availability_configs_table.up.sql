CREATE TABLE IF NOT EXISTS user_availability_configs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    monday_morning BOOLEAN DEFAULT false,
    monday_afternoon BOOLEAN DEFAULT false,
    tuesday_morning BOOLEAN DEFAULT false,
    tuesday_afternoon BOOLEAN DEFAULT false,
    wednesday_morning BOOLEAN DEFAULT false,
    wednesday_afternoon BOOLEAN DEFAULT false,
    thursday_morning BOOLEAN DEFAULT false,
    thursday_afternoon BOOLEAN DEFAULT false,
    friday_morning BOOLEAN DEFAULT false,
    friday_afternoon BOOLEAN DEFAULT false,
    saturday_morning BOOLEAN DEFAULT false,
    saturday_afternoon BOOLEAN DEFAULT false,
    sunday_morning BOOLEAN DEFAULT false,
    sunday_afternoon BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_user_availability_configs_user_id ON user_availability_configs(user_id);
