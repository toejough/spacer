-- Drop the unused visitors table.
--
-- Leftover from the exe.dev scaffold's default landing page (welcome.html,
-- removed alongside this migration). The only write path, UpsertVisitor, has
-- zero callers anywhere in this repo's history, so the table has never had a
-- row written to it in any deployment -- dropping it discards no data.
DROP TABLE IF EXISTS visitors;

INSERT OR IGNORE INTO migrations (migration_number, migration_name)
VALUES (004, '004-drop-visitors');
