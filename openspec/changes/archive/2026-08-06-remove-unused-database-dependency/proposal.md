## Why

Issue #28 started narrower: `db/dbgen`'s Item query layer (`CreateItem`, `GetItem`, `ListDueForReview`, etc.) has zero callers anywhere in `srv/*.go` or `cmd/srv/*.go` — `srv/server.go` opens the database only to run migrations, and every todo/note lives in browser `localStorage` instead. That issue asked whether this was a deliberate, permanent design or a stopgap that never got cleaned up.

The repo owner's answer: "I am not purposely keeping anything unused in the db. I don't even know why we have a db at this point." That resolves the open question and widens the scope. Investigating further:

- The `items` table itself has **0 rows** in the live production database (verified directly against `db.sqlite3`) — there's no data to migrate or preserve.
- Nothing anywhere in the module issues raw SQL against `items` outside the migration machinery itself (`grep -rn "\.Exec(\|\.Query(\|\.QueryRow(" srv/*.go cmd/srv/*.go db/db.go` shows only pragma-setting and migration-tracking queries in `db/db.go`).
- `README.md`'s "Database" section reads, verbatim, "This template uses sqlite (\`db.sqlite3\`). SQL queries are managed with sqlc." — scaffold-template boilerplate language, the same origin story as `welcome.html` and the `visitors` table (both already removed as dead exe.dev-scaffold remnants, issue #27).

Put together: nothing in this app has used the SQL database for its stated purpose since the todo/notes app was built on `localStorage` instead. `srv/server.go`'s own comment even says so directly: *"Still open DB so migrations table exists (template requirement), but we don't use it for items anymore."* That "template requirement" is inherited scaffold assumption, not a verified current one — no evidence anywhere (platform docs, `AGENTS.md`, `README.md`) that phone-llm or anything else actually requires this app to maintain a SQLite database.

## What Changes

- **BREAKING** (internal only, no user-facing effect): `srv.New()`'s signature drops its `dbPath` parameter — it becomes `New(hostname string) (*Server, error)`. `Server.DBPath` is removed.
- Delete the entire `db/` directory: `db.go`, `sqlc.yaml`, `dbgen/` (all 3 generated files), `migrations/` (all 4 files), `queries/items.sql`.
- Remove the `sqlc` build-time tool dependency (`tool github.com/sqlc-dev/sqlc/cmd/sqlc` in `go.mod`) and run `go mod tidy` to drop `modernc.org/sqlite` and the many indirect dependencies that existed purely to support sqlc's own multi-database code generation (MySQL/Postgres drivers, `cel-go`, `grpc`, `cobra`, the TiDB SQL parser, `wazero`, etc. — none of these are the app's own dependencies, they're sqlc's).
- Update `README.md`: remove the "Database" section, remove the `db` bullet from "Code layout".
- Update `openspec/config.yaml`'s `context` block: remove the SQLite/sqlc tech-stack lines.
- `cmd/srv/main.go` no longer passes `"db.sqlite3"` to `srv.New()`.
- `.gitignore`'s `db.sqlite3`/`*.db` entries become moot (nothing produces them anymore) — remove them.

## Capabilities

No capabilities, new or modified — the database was never wired into any observable behavior (items are and remain `localStorage`-only). `skip_specs: true` is set in `.openspec.yaml`.

## Impact

- `db/` (deleted entirely)
- `srv/server.go` (`New()` signature change, `DBPath` field removed, `db` import removed)
- `cmd/srv/main.go` (call site updated)
- `go.mod`/`go.sum` (sqlc tool dependency and its transitive tree removed via `go mod tidy`)
- `README.md`, `openspec/config.yaml`, `.gitignore` (stale DB references removed)
- No change to any user-facing behavior — items were already 100% `localStorage`-based.
