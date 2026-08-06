## Context

See proposal.md for the evidence trail (zero rows, zero raw-SQL callers, scaffold-boilerplate README text) and the repo owner's direct answer resolving issue #28's open question. This design covers the two real decisions: how far removal goes, and how to sequence it safely.

## Goals / Non-Goals

**Goals:**
- Remove every trace of the SQL database dependency, not just the unused `Item` query layer — matching the repo owner's actual question ("why do we have a db at all") rather than only the narrower thing issue #28 originally named.
- Leave `srv.New()`/`cmd/srv/main.go` with an honest signature: no `dbPath` parameter for a database the app doesn't use.

**Non-Goals:**
- Reintroducing server-side storage in any form. If a future need for server-side persistence emerges (multi-device sync, backup, etc.), that's a fresh design decision with fresh requirements — not a reason to keep today's unused scaffold around "just in case."

## Decisions

### 1. Remove the whole `db/` subsystem, not just the `items` query layer

Issue #28 was scoped narrowly (delete `db/queries/items.sql` and the generated `db/dbgen/items.sql.go`, matching the `visitors` precedent from issue #27). The repo owner's answer widens that: the database itself has no justified reason to exist, not just the query layer built on top of it. Stopping at "delete the item queries but keep `db.Open`/`RunMigrations`/the `items` table" would leave exactly the same shape of dead weight this session has been removing all day — a database that exists because it always has, not because anything needs it.

**Alternative considered**: keep `db.Open`/`RunMigrations` and the `migrations` table, drop only the item-specific pieces (matching how `visitors` was handled — schema dropped via a new migration, machinery kept). Rejected: unlike `visitors` (one unused table sitting alongside real usage), there is no remaining real usage once the item layer is gone — keeping the migration machinery running for a schema with nothing in it is the same "template requirement" cargo-culting the repo owner just said isn't intentional.

### 2. Remove sqlc as a build-time tool, not just its generated output

`go.mod`'s `tool github.com/sqlc-dev/sqlc/cmd/sqlc` line, and the ~25 indirect dependencies that exist to support it (MySQL/Postgres drivers, `cel-go`, `grpc`, `cobra`, the TiDB SQL parser, `wazero`, etc.), only exist because something needs code-generated from SQL. Once `db/queries/` is empty, there's nothing left for sqlc to generate. Leaving the tool declared but unused is the same shape of problem as the abandoned `mage`/mutation-testing toolchain removed earlier this session (issues #9-#12, #14, #24) — a tool nothing invokes, quietly carrying a large dependency tree.

### 3. Sequencing: delete code first, then `go mod tidy`, then verify

Delete `db/` and update the two Go call sites (`srv/server.go`, `cmd/srv/main.go`) first, confirm `go build ./...` succeeds (proving nothing else references the `db` package), *then* run `go mod tidy` to let Go's own dependency analysis determine exactly what's now unused — rather than hand-guessing which indirect dependencies belong to sqlc versus something else the app still needs. This is the same approach used for the mutation-testing tooling and scaffold-remnants removals earlier this session, and it's mechanically verifiable rather than relying on manual judgment about a ~25-entry indirect dependency list.

## Risks / Trade-offs

- **This is the largest and most structurally significant change of the four cleanup proposals** (touches `go.mod`'s dependency tree, two Go source files' signatures, and an entire top-level directory) → mitigated by the evidence already gathered (zero rows, zero callers, explicit repo-owner sign-off) and by verifying `make build && make test` plus a live browser check before considering it done, the same bar applied to every other change this session.
- **`db.sqlite3` itself is untracked** (already gitignored) — nothing to delete from git, but worth confirming no deploy script or systemd unit expects the file to exist at a specific path before this lands (it won't be created at all once `db.Open` is gone).
- **If some future feature genuinely needs server-side storage**, this removal means starting over rather than resuming dormant infrastructure. Accepted: resuming code built for a different, never-realized use case is not meaningfully cheaper than writing fresh code for whatever the actual future requirement turns out to be, and carrying it forward unused has a continuous cost (confusion, this exact investigation) that a future rebuild wouldn't.
