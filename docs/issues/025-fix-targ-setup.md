# 025 — Fix targ setup: shell script shadows system targ

**Status:** closed
**Type:** bug

## Context

Issue #024 migrated `dev/` scripts to a `./targ` shell script dispatcher at the project root. This created two problems:

1. **`./targ` shadows the system `targ` binary** — Running `targ` (without `./`) invokes the system-wide Go-based targ build tool at `/Users/joe/go/bin/targ`, not the local shell script. This is confusing and defeats the purpose of using targ.
2. **Broken `targs.go` at project root** — A `targs.go` file was generated with package name `5mincrements` (derived from directory `5m-increments`), which is an invalid Go identifier. This causes `targ` to error on any invocation.
3. **`dev/` directory deleted** — The original `dev/` scripts were removed, but the replacement should have been `dev/targets.go` (proper targ target file), not a root-level shell script.

The whole point of targ is that you write Go target files in `dev/` and the system `targ` binary discovers and runs them — no local binary needed.

## Acceptance Criteria

- [x] `./targ` shell script removed from project root
- [x] `targs.go` (broken) removed from project root (was never committed)
- [x] `dev/targets.go` created with proper targ string targets
- [x] `targ --help` shows all targets (dev, test, build, check, clean)
- [x] `targ test` runs successfully
- [x] CLAUDE.md updated: `./targ` → `targ`
