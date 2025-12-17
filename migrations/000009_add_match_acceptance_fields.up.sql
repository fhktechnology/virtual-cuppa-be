-- Add fields to track individual user acceptances
ALTER TABLE matches
ADD COLUMN user1_accepted BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN user2_accepted BOOLEAN NOT NULL DEFAULT false;

-- Update existing 'accepted' matches to have both users accepted
UPDATE matches
SET user1_accepted = true, user2_accepted = true
WHERE status = 'accepted';

-- Update existing 'rejected' matches - mark user1 as rejected (we don't know which one rejected)
UPDATE matches
SET user1_accepted = false, user2_accepted = false
WHERE status = 'rejected';
