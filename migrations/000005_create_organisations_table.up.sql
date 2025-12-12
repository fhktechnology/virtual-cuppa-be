CREATE TABLE IF NOT EXISTS organisations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    company_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_organisations_deleted_at ON organisations(deleted_at);

ALTER TABLE users DROP COLUMN IF EXISTS organisation;
ALTER TABLE users ADD COLUMN IF NOT EXISTS organisation_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_users_organisation_id ON users(organisation_id);
ALTER TABLE users ADD CONSTRAINT fk_users_organisation FOREIGN KEY (organisation_id) REFERENCES organisations(id) ON DELETE SET NULL;
