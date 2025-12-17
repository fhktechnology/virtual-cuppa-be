CREATE TABLE match_feedbacks (
    id SERIAL PRIMARY KEY,
    match_id INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(match_id, user_id)
);

CREATE INDEX idx_match_feedbacks_match_id ON match_feedbacks(match_id);
CREATE INDEX idx_match_feedbacks_user_id ON match_feedbacks(user_id);
CREATE INDEX idx_match_feedbacks_deleted_at ON match_feedbacks(deleted_at);
