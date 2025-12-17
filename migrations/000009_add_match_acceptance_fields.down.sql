-- Remove acceptance tracking fields
ALTER TABLE matches
DROP COLUMN IF EXISTS user1_accepted,
DROP COLUMN IF EXISTS user2_accepted;
