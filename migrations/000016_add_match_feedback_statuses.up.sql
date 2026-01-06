-- Add new match statuses: waiting_for_feedback and completed
-- No actual schema change needed, just ensuring the application can use these values
-- The status column in matches table already supports varchar values

-- This migration is a placeholder to document the addition of new status values:
-- - waiting_for_feedback: Match accepted and expired, waiting for feedback from users
-- - completed: Both users submitted feedback
