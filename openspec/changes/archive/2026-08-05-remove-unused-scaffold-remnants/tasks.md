## 1. Delete the dead template

- [x] 1.1 Delete `srv/templates/welcome.html`
- [x] 1.2 Confirm nothing references it: `grep -rn "welcome" --include="*.go" .` and `grep -rn "welcome" srv/static/` both return nothing

## 2. Remove the visitor-counter query layer

- [x] 2.1 Delete `db/queries/visitors.sql`
- [x] 2.2 Regenerate `db/dbgen/` (`sqlc generate` from `db/`) and confirm `db/dbgen/visitors.sql.go` is gone
- [x] 2.3 Confirm `db/dbgen/` still builds and no other generated file references `Visitor`

## 3. Drop the `visitors` table

- [x] 3.1 Add `db/migrations/004-drop-visitors.sql` (renumbered from the planned `003` — `db/migrations/003-add-timestamps.sql` already existed and was already applied to the live database, which the proposal's investigation missed): `DROP TABLE IF EXISTS visitors`, plus the migration-tracking insert matching the pattern in `001-base.sql`/`002-todos-notes.sql`, plus a comment noting why it's safe (table has zero callers, therefore always empty)
- [x] 3.2 `make build && make test` pass, confirming `db.RunMigrations` applies the new migration without error against the existing `db.sqlite3`
- [x] 3.3 Confirm the `visitors` table is actually gone after a run: query `sqlite3 db.sqlite3 ".tables"` (or equivalent) and see it absent

## 4. Verify

- [x] 4.1 `grep -rn "welcome\|visitor" --include="*.go" --include="*.html" --include="*.sql" .` (excluding `node_modules/`, `.git/`, and this change's own planning docs) returns nothing outside the new migration file's `DROP TABLE` line
- [x] 4.2 `make build && make test` pass
- [x] 4.3 The running app still serves correctly (`GET /` renders, no console errors) — this change touches the database layer, so confirm the app still starts and functions, not just that it compiles
