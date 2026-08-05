## Why

`todo.toejough.dev` renders a white screen. `adopt-standard-app-layout` repointed the server at `/app/public` and published four of five frontend assets there by symlink, but `/app/public/index.html` was left as a hand-written 13-line placeholder instead of a link to the real 197-line `srv/templates/index.html`. The server served the stub faithfully — an empty `<div id="app">`, broken relative asset paths, no PWA manifest link, no service worker registration — and every automated signal (systemd active, port open, `GET /` 200) looked healthy anyway. Filed as GitHub issue #26.

A second, unrelated regression compounded this: an in-progress, uncommitted edit to `srv/server.go` (from an abandoned change that misdiagnosed this same symptom as a routing bug) left the file in a non-compiling state and dropped the `Cache-Control: no-store` header from the root route. Since `/app/server/bin/serve` runs `make build && exec ./todo-srv` on every restart, a broken build meant `app.service` was crash-looping and falling back to a stale binary on every cycle.

## What Changes

- `/app/public/index.html` becomes a symlink to `srv/templates/index.html`, matching the publish-by-symlink pattern already used for `manifest.json`, `script.js`, `style.css`, and `sw.js`. This is a filesystem action on `/app/public` (a directory outside the git tree, populated by hand, per `adopt-standard-app-layout`'s design) — not a code change.
- `srv/server.go` is restored to a compiling state: unused imports (`io/fs`, `os`, `path/filepath`) removed, and the root route (`GET /{$}`) goes back to an explicit `handleIndex` that serves `index.html` from `TemplatesDir` with `Cache-Control: no-store`. The `/app/public` repointing (`StaticDir`/`TemplatesDir` consts) from `adopt-standard-app-layout` is kept as-is.
- The `openspec/changes/fix-root-handler-404` change (untracked, 0/7 tasks done, based on the false premise that the root route was 404ing `/script.js`-style paths) is superseded by this change and was deleted.

## Capabilities

### Modified Capabilities
- `web-delivery`: the entry page must be published the same way as every other asset — live-linked to its source, not a hand-authored duplicate — and the no-store cache header applies to the entry page, not only to files under `/static/`.

## Impact

- `srv/server.go` (root route handler, imports)
- `/app/public/index.html` (filesystem, outside git tree)
- Deletes `openspec/changes/fix-root-handler-404/`
- No database or API changes; no user-visible behavior changes beyond the page actually rendering
