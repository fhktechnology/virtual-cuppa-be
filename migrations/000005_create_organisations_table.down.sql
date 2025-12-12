ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_organisation;
DROP INDEX IF EXISTS idx_users_organisation_id;
ALTER TABLE users DROP COLUMN IF EXISTS organisation_id;
ALTER TABLE users ADD COLUMN IF NOT EXISTS organisation VARCHAR(255);

DROP TABLE IF EXISTS organisations;
