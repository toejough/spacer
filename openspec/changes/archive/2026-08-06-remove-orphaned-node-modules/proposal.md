## Why

`node_modules/` (372 files) and `.env.example` sit at the repo root, tracked in git, with no `package.json` anywhere in the repo to explain them (confirmed: `find . -maxdepth 2 -iname package.json` finds none outside `node_modules` itself and the unrelated `run-remember-everything` skill's own separate Playwright-driver `node_modules`). Nothing invokes any of it — no references in `Makefile`, `README.md`, `AGENTS.md`, or any script. `.env.example`'s one line (`VITE_MOBILE_EQUAL_ACCESS=true`) has no Go-side consumer; no Go code in this repo reads any environment variable at all.

This is leftover from the already-archived `cleanup-unused-code-and-spec-drift` change, which deleted `web-app/` (the Quasar/Vue PWA prototype) but didn't catch these two paths since they live at the repo root rather than inside `web-app/`. Confirmed via issue #29.

Kept as its own change (not batched with the other cleanup issues) purely for diff legibility — 372 unrelated file deletions would bury the much smaller, unrelated fixes in `remove-stale-ops-and-agent-files` and `remove-dead-css-rules`.

## What Changes

- Delete `node_modules/` (all 372 tracked files).
- Delete `.env.example`.
- Add `node_modules/` to `.gitignore` so this can't silently recur if something local ever runs `npm install` at the repo root again.

## Capabilities

No capabilities, new or modified — nothing currently reachable depends on any of this. `skip_specs: true` is set in `.openspec.yaml`.

## Impact

- `node_modules/` (deleted, 372 files)
- `.env.example` (deleted)
- `.gitignore` (one line added)
- No change to `srv/`, `db/`, `cmd/`, or any user-facing behavior.
