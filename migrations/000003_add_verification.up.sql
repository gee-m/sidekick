BEGIN;

ALTER TABLE users
ADD COLUMN verification_token UUID,
ADD COLUMN verified_at TIMESTAMPTZ;

CREATE INDEX idx_users_verification_token ON users(verification_token)
WHERE verification_token IS NOT NULL;

COMMIT;
