## Why

The repo has accumulated code that's never run in production and spec text that no longer matches what the live app (`srv/static/script.js`) actually does. Both make the codebase harder to trust: dead code invites confusion about what's real, and drifted specs give a false picture of behavior to anyone (human or agent) reading them before making a change.

## What Changes

- **BREAKING** (repo cleanup only, no user-facing effect): delete `web-app/`, the Quasar/Vue PWA prototype. It isn't referenced by any build tooling (`Makefile`, `magefile.go`), isn't deployed, and its mobile equal-access UX was independently and separately built directly into the production app already (see `openspec/specs/todo-list/spec.md`).
- Remove the unused `/api/todos` REST handlers (`handleTodos`, `handleTodo`) from `srv/server.go`. Confirmed unused: `srv/static/script.js` has zero `fetch()` calls anywhere — everything runs on `localStorage`. The DB is still opened and migrated (kept for the migrations table), just not queried by these handlers.
- Delete the stale release artifact `releases/remember-everything-v24.tar.gz`; the project is at v27 per recent commit history and this tarball isn't referenced by any build or deploy tooling.
- Correct spec drift in `openspec/specs/todo-list/spec.md`:
  - "Complete a todo" → "Toggle completion" scenario says the user "clicks the checkbox," but production uses an icon button, not a checkbox.
  - "Mobile parity for Done and Abandon (IMPLEMENTED)" is written as if this were a separate mobile-only surface (references swipe gestures as an alternative, carries an "(IMPLEMENTED)" title suffix). In production this is just how the one todo list behaves, not a mobile-specific mode — reword to drop the mobile-only framing.
  - Consolidate the now-duplicated mutual-exclusion scenario that exists nearly identically under both "Delete a todo" and "Mobile parity for Done and Abandon" into one place.

## Capabilities

### Modified Capabilities
- `todo-list`: requirement text corrected to match actual production behavior (button not checkbox; mutual-exclusion behavior described once, without mobile-specific framing). No SHALL-level behavior changes — production already does this.

## Impact

- `web-app/` (deleted)
- `srv/server.go` (unused handlers removed)
- `releases/remember-everything-v24.tar.gz` (deleted)
- `openspec/specs/todo-list/spec.md` (wording corrections)
- No change to `srv/static/script.js` behavior, the Go server's actual serving path, or any user-facing functionality.
