-- Add completed_at and archived_at timestamps to items
ALTER TABLE items ADD COLUMN completed_at TEXT;
ALTER TABLE items ADD COLUMN archived_at TEXT;

INSERT OR IGNORE INTO migrations (migration_number, migration_name)
VALUES (003, '003-add-timestamps');
