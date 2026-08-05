## Context

See proposal.md for motivation. Two independent facts shape the approach:

- `/app/public` is a directory outside the git tree (a sibling of `/app/repo`), populated only by hand-run `ln -s` commands during `adopt-standard-app-layout` — no script in this repo creates or maintains those links. Fixing `index.html` there is a filesystem action, not a commit.
- `srv/server.go`'s working tree was independently broken by an abandoned, uncommitted attempt (the deleted `fix-root-handler-404` change) that misdiagnosed this same symptom.

## Goals / Non-Goals

**Goals:**
- Restore the entry page to the same publish-by-symlink guarantee already made for every other asset.
- Get `srv/server.go` back to the state `adopt-standard-app-layout` intended: `/app/public` repointing kept, nothing else changed.

**Non-Goals:**
- Redesigning how assets are published (symlink-by-hand is `adopt-standard-app-layout`'s decision, not revisited here).
- Automating symlink creation with a script. Worth considering separately, but out of scope for a one-file fix — see Open Questions.

## Decisions

### Revert the routing change instead of fixing it forward

`fix-root-handler-404`'s premise — that the root route 404s on `/script.js`, `/style.css`, etc. — is false. Go 1.22+'s `ServeMux` pattern `/{$}` matches only the exact path `/`; it was never intercepting other paths. Those assets are correctly served under `/static/`, which the real `index.html` template already references with absolute paths. The replacement handler (`http.FileServer(http.Dir(TemplatesDir))` on `/{$}`) is therefore a no-op against the described bug, and it silently dropped `Cache-Control: no-store` from the root route — a regression, since the app's `?v=NN` cache-busting scheme depends on the HTML document itself never being cached by the platform's CDN. The correct move is reverting to the explicit `handleIndex` handler, not patching the FileServer approach.

### Symlink the entry page rather than embedding it

`adopt-standard-app-layout`'s design already chose symlink-by-hand as the publishing mechanism for every other asset (editing the source is live on the next request, no rebuild). Embedding `index.html` into the binary (`embed.FS`, as that change's task list mentions in passing for templates generally) would break that live-edit property and diverge from the pattern the other four assets use. Consistency wins for a one-file fix.

## Risks / Trade-offs

- **The symlink fix lives outside the repo** → nothing in version control prevents this exact regression from recurring if `/app/public` is ever rebuilt from scratch. Mitigated only by this change's spec delta making the expectation explicit; a durable, scripted fix is future work (see Open Questions).
- **Status codes and process health don't catch content regressions** → this bug shipped with systemd active, the port open, and `GET /` returning 200. Verification for this change is done by loading the page in a browser and checking rendered content (manifest link, service worker state, console errors), not by checking process/HTTP status alone.

## Open Questions

- Should publishing to `/app/public` be scripted (e.g. as part of `make build` or `bin/serve`) instead of relying on someone running `ln -s` by hand during a future layout change? Doesn't change this fix or its specs — worth a separate proposal if it recurs.
