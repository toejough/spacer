## Why

Three unrelated pieces of stale ops/agent-integration cruft, all confirmed dead in a repo-wide audit (issues #30, #32, #33), plus one adjacent finding caught while reviewing the same file: `AGENTS.md`'s "Queued work" section still names `adopt-standard-app-layout` as "the current one," but that change has been archived since this session's very first fix. Batched together because each is a trivial, zero-risk file/text deletion — none of them are loaded, served, built, or tested by anything, so the only verification needed is a grep sweep, not app behavior testing.

- **Issue #30**: the tracked `run` script at repo root is the pre-`adopt-standard-app-layout` version (`exec ./todo-srv -listen 0.0.0.0:8080`) and no longer matches the real deploy entrypoint `/app/run` (`exec /app/server/bin/serve`). Nothing references `./run` by path anywhere.
- **Issue #32**: `issues/0001-ui-abandon-visuals.md` references an openspec change deleted long ago (`1900fd9 chore: delete superseded ui-abandon-visuals change`). The repo now tracks issues via `gh issue` (GitHub Issues), evidenced by issues #9–#33 filed this session; this file is a one-off leftover from before that convention existed.
- **Issue #33**: Shelley (a separate coding agent) isn't used in this environment. `shelley-hooks/` (`install.sh`, `slash/opsx`) and its two `AGENTS.md` references should go. (Note: `.pi/` is *not* part of this — that belongs to a different client the repo owner does actively use, unrelated to Shelley.)
- **Adjacent finding**: `AGENTS.md`'s "Queued work" section describes `adopt-standard-app-layout` as current/queued and repeats two findings from that change's own docs (compile-time template path resolution, silent `make build` failure under systemd) as if they're still open concerns to read up on before touching anything. Both were resolved when that change was implemented. The section is now actively misleading, not just stale.

## What Changes

- Delete `run` (repo root).
- Delete `issues/0001-ui-abandon-visuals.md` and the now-empty `issues/` directory.
- Delete `shelley-hooks/` (`install.sh`, `slash/opsx`).
- Edit `AGENTS.md`: remove the two Shelley hook lines, and remove/rewrite the stale "Queued work" section (no change is currently queued — `openspec list` shows none active before this session's new proposals).

## Capabilities

No capabilities, new or modified — pure ops/doc/agent-integration cleanup with zero effect on the running app. `skip_specs: true` is set in `.openspec.yaml`.

## Impact

- `run` (deleted)
- `issues/` (deleted entirely)
- `shelley-hooks/` (deleted entirely)
- `AGENTS.md` (Shelley lines removed, "Queued work" section removed/rewritten)
- No change to `srv/`, `db/`, `cmd/`, or any user-facing behavior.
