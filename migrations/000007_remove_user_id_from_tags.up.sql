-- Remove user_id column from tags table (many-to-many relationship is handled by user_tags table)
ALTER TABLE tags DROP COLUMN IF EXISTS user_id;
