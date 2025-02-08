BEGIN;

ALTER TABLE users
DROP COLUMN verification_url,
DROP COLUMN verification_token,
DROP COLUMN verified_at;

COMMIT;
