-- First, check the current state
SELECT * FROM schema_migrations;

-- Force the version to clean state
UPDATE schema_migrations SET dirty = false WHERE version = 3;

-- Optional: If you want to completely roll back to before version 3
DELETE FROM schema_migrations WHERE version = 3;
