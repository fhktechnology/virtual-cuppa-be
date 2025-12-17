CREATE TABLE IF NOT EXISTS match_availabilities (
    id SERIAL PRIMARY KEY,
    match_id INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    availability JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(match_id, user_id)
);

CREATE INDEX idx_match_availabilities_match_id ON match_availabilities(match_id);
CREATE INDEX idx_match_availabilities_user_id ON match_availabilities(user_id);
