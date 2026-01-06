-- Rollback: Remove references to waiting_for_feedback and completed statuses
-- Update any matches with these statuses back to accepted
UPDATE matches 
SET status = 'accepted' 
WHERE status IN ('waiting_for_feedback', 'completed');
