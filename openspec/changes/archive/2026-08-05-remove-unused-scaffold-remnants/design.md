## Context

See proposal.md for motivation and the connection between `welcome.html`, `visitors.sql`, and the `visitors` table — one abandoned scaffold feature, not three unrelated ones. The one piece here that isn't pure file deletion is dropping a live database table, which is why this gets a design doc despite being otherwise mechanical.

## Goals / Non-Goals

**Goals:**
- Remove all three pieces of the dead feature completely, not just the file `srv/templates/welcome.html` that issue #27 originally named.
- Drop the `visitors` table without risking data loss or breaking already-deployed databases.

**Non-Goals:**
- Auditing the rest of the schema for other unused tables/columns — this change is scoped to the one connected feature traced from issue #27, not a general schema audit.

## Decisions

### 1. Drop the table via a new migration, not by editing `001-base.sql`

This repo's migrations are sequential and already-applied (the running app's `db.sqlite3` has migrations `001-base` through `003-add-timestamps` recorded in the `migrations` tracking table — `003-add-timestamps.sql` already existed at implementation time, one number past what this proposal assumed). Editing `001-base.sql` in place would make a fresh database's schema diverge from what already-deployed databases have — a fresh DB would never create `visitors` at all, while the deployed one would still have it until manually intervened on. A new migration (`004-drop-visitors.sql`) is the only way that's consistent for both fresh and already-deployed databases: both end up without the table, via the same recorded migration step.

### 2. Confirming the drop is safe before writing the migration

`UpsertVisitor` is the only code path that ever writes to `visitors`, and it has zero callers anywhere in `srv/*.go` or `cmd/srv/*.go` (confirmed by grep across the whole module). Since the table can only ever have been populated by that function, and that function has never been called, the table is guaranteed empty in every deployment — there is no scenario where dropping it discards real data. This reasoning is worth stating explicitly in the migration file itself (a comment), since "drop this table" reads as more dangerous than it is without it.

### 3. Regenerate `db/dbgen/` with sqlc rather than hand-editing the generated file

`db/dbgen/visitors.sql.go` is sqlc-generated output (`db/sqlc.yaml` globs `queries/*.sql`). Deleting `db/queries/visitors.sql` and running `sqlc generate` is the correct way to remove it — hand-deleting the generated file would work today but leave the repo one `sqlc generate` away from it silently reappearing out of sync with the query source.

## Risks / Trade-offs

- **Migrations, unlike file deletions, aren't reversible by `git revert` alone once applied to a live database** → mitigated by the safety argument in decision 2: the table is provably empty, so there's nothing to lose either direction.
- **The running production database still needs this migration applied** → `db.RunMigrations` already runs on every server start (per `srv/server.go`'s `New()`), so the next deploy picks it up automatically; no manual step needed.
