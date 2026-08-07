## Why

GitHub issue #34 (github.com/toejough/spacer): the deployed app binds `0.0.0.0:8080`, so it is
reachable from every other app sharing the OrbStack machine bridge and from the local LAN, not just
the tailnet, which is the only path meant to reach it. The live bind is controlled by
`/app/server/bin/serve`, a launch script that passes `-listen 0.0.0.0:8080` explicitly. That script
is outside git (not part of this repo, same as `/app/public`), but it is not platform-locked — it's
a plain file this agent can edit directly, so closing the exposure is in scope here, not separate
ops work for someone else to do.

Separately, but caused by the same root habit, the repo's own default is unsafe: `cmd/srv/main.go`
defines `-listen` with a default of `:8000`, which Go binds on all interfaces. Anyone who runs the
binary without passing `-listen` explicitly — a fresh clone, a local dev run, a future deployment —
gets the same exposure by accident. The default should be safe without requiring the caller to know
to override it.

## What Changes

- Change the `-listen` flag's default value in `cmd/srv/main.go` from `:8000` to `127.0.0.1:8000`, so
  the server only listens on loopback unless a caller explicitly opts into a wider bind.
- Establish a `tailscale serve --bg --tcp 8080 tcp://127.0.0.1:8080` mapping so the tailnet path
  keeps working once the app stops listening on all interfaces.
- Edit `/app/server/bin/serve` to pass `-listen 127.0.0.1:8080` instead of `-listen 0.0.0.0:8080`,
  and restart `app.service` to pick it up.
- Ordering matters: the tailscale serve mapping must be established and verified *before* the bind
  flips, or the app is unreachable in between (per issue #34).

## Capabilities

### New Capabilities
- `network-binding`: the app's listen-address defaults and how they're overridden.

### Modified Capabilities
(none)

## Impact

- Affected code: `cmd/srv/main.go` (flag default only).
- Affects anyone who runs `todo-srv` without `-listen`: they now get loopback-only instead of
  all-interfaces.
- Affected deployment: `/app/server/bin/serve` (outside git) and live `app.service` state, plus
  tailscaled's persistent serve-mapping state. Closes the machine-bridge and LAN exposure described
  in issue #34.
