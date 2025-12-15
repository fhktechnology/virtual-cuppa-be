-- Add user_id column back (for rollback)
ALTER TABLE tags ADD COLUMN user_id INTEGER;
ALTER TABLE tags ADD CONSTRAINT tags_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_tags_user_id ON tags(user_id);
