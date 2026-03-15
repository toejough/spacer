# 024 — Migrate dev/ scripts to targ and add clean target

**Status:** closed
**Type:** tooling

## Context

Build commands are currently separate binaries in `dev/` (`dev/build`, `dev/check`, `dev/dev`, `dev/test`). These should be unified under `targ` so we can call `targ clean`, `targ dev`, `targ check`, `targ build`, etc.

Additionally, build and test tools leave artifacts that accumulate in the working tree: `test-results/` (Playwright), `dist/` (Vite build), and potentially others as tooling grows. These clutter `git status` and can confuse future sessions. A `targ clean` subcommand should handle removal.

## Acceptance Criteria

- [x] `targ` dispatches to subcommands (dev, build, test, check, clean, etc.)
- [x] `targ clean` removes build and test artifacts (dist/, test-results/, etc.)
- [x] Existing `dev/*` scripts migrated to targ subcommands
- [x] CLAUDE.md updated with new build commands
