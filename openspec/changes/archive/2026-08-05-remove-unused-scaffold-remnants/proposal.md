## Why

This repo's single "Add todo app" commit (2026-07-17) bulk-imported the exe.dev Go web template scaffold together with the real todo app. Two pieces of that scaffold never got wired up and have sat completely dead ever since — surfaced while investigating GitHub issue #26 (the white-screen bug), tracked as issue #27, and confirmed here:

- `srv/templates/welcome.html` — the scaffold's default "hello world" landing page. `git log` shows exactly one commit ever touched it (the import itself); no Go handler has ever served it in this repo's history. It references `{{.Hostname}}`, `{{.LoginURL}}`, `{{.Headers}}`, and the `exe.dev`/`Shelley` platform this app migrated away from (see `1767603 feat: migrate deployment from exe.dev to phone-llm platform`). It also has zero matching CSS classes in `style.css`, so it would render completely unstyled even if something served it.
- `db/queries/visitors.sql` / the generated `db/dbgen/visitors.sql.go` (`UpsertVisitor`, `VisitorWithID`) — the visit-counter backing `welcome.html`'s `{{.VisitCount}}`. Zero callers anywhere in `srv/*.go` or `cmd/srv/*.go`.
- The `visitors` table, created by `db/migrations/001-base.sql`. Since nothing has ever called `UpsertVisitor`, the table has never had a row written to it in any deployment — dropping it loses no data.

Same shape as the mutation-testing tooling removed earlier (issues #9-#12, #14, #24): code imported once, never wired to anything current, safe to delete outright rather than carry forward as confusion-inviting dead weight.

## What Changes

- Delete `srv/templates/welcome.html`.
- Delete `db/queries/visitors.sql`, regenerate `db/dbgen/` with `sqlc generate` so the now-orphaned `db/dbgen/visitors.sql.go` and any `Visitor`-related generated types disappear.
- Add a new migration (`db/migrations/004-drop-visitors.sql` — renumbered from an initially-assumed `003` once implementation found `003-add-timestamps.sql` already existed) dropping the `visitors` table, rather than editing `001-base.sql` in place — migrations in this repo are append-only/sequential, and rewriting an already-applied migration would desync already-deployed databases from what a fresh `db.sqlite3` would get.

## Capabilities

No capabilities, new or modified — nothing currently reachable depends on any of this, so removing it changes no observable behavior. `skip_specs: true` is set in `.openspec.yaml`.

## Impact

- `srv/templates/welcome.html` (deleted)
- `db/queries/visitors.sql` (deleted)
- `db/dbgen/visitors.sql.go` (deleted, via regeneration)
- `db/migrations/004-drop-visitors.sql` (new migration)
- No change to `srv/server.go`, routes, or any user-facing behavior.
