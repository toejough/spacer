## Why

`cmd/srv/main.go` fetches the machine's hostname via `os.Hostname()` (with an "unknown" fallback on error) and threads it into `srv.New(hostname)`, which stores it on `Server.Hostname`. Nothing in `srv/server.go` ever reads that field — not in `Serve()`, not in `handleIndex()`, nowhere. It's data plumbed all the way in at construction time and never used for anything: no logging, no response header, no routing decision. Pure dead weight, carried since the app's first commit.

## What Changes

- Remove the `Hostname` field from `Server` (`srv/server.go`).
- Change `srv.New` to take no arguments (or drop it entirely if a zero-argument constructor turns out to add nothing over a bare struct literal — a call made during implementation, not a spec-level concern).
- Remove the `os.Hostname()` call and its error-fallback branch from `cmd/srv/main.go`.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
(none — this removes an internal implementation detail with no externally observable behavior. Nothing reads `Hostname` today, so nothing observable changes when it's gone. `skip_specs: true` is set in this change's `.openspec.yaml`.)

## Impact

- `srv/server.go`: `Server` struct, `New` constructor.
- `cmd/srv/main.go`: `run()`.
- No other callers of `srv.New` exist in the repo.
