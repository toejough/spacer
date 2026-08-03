## Context

`/app/run` is `app.service`'s entrypoint on phone-llm — it must bind `:8080`, the port the platform's Cloudflare tunnel talks to. The current `run` is `newapp`'s default JS starter (`npm run dev --host 0.0.0.0 --port 8080`), which doesn't apply since this repo builds a Go binary via `make build`. `X-ExeDev-*` auth headers are referenced only in docs, never read in code, so no code depends on them. See proposal.md - Why.

## Goals / Non-Goals

**Goals:**
- `run` reliably builds and starts the current code on `:8080` under `app.service` (restart-on-crash, per the platform).
- Static assets and HTML responses are never edge-cached past a deploy.

**Non-Goals:**
- No change to app behavior, API, schema, or the SQLite data model.
- No re-architecting of the Go server beyond the two header fixes.

## Decisions

- **Build-then-exec in `run`, not a long-lived dev-server wrapper.** Go doesn't need HMR; `run` should `make build` and `exec ./todo-srv -listen :8080` (`exec` so `app.service`'s process supervision applies directly to the binary, not a wrapping shell). Rebuilding on every `app.service` (re)start keeps `run` simple and matches how `make build` is already the project's build step — no separate CI/build artifact pipeline needed for a single-owner app.
- **`no-store`, not `no-cache`, on dynamic/static responses.** `no-cache` still permits Cloudflare to cache-and-revalidate; `no-store` is what the platform's README calls for to guarantee deploys are visible immediately. Applied to `handleIndex`, `sw.js`, and the `/static/` `http.FileServer` (which currently sets no header, so it's implicitly cacheable).
- **Delete `srv.service` rather than leave it.** It references `/home/exedev` paths that don't exist on this machine and duplicates what `app.service` already does; keeping it invites someone to `systemctl enable` a unit that will never work here.

## Risks / Trade-offs

- [Rebuilding on every `run` invocation adds a few seconds to restart time] → acceptable for a single-owner app; `make build` is already fast (small Go binary, no heavy codegen at build time — sqlc runs separately, not part of `make build`).
- [`no-store` disables all caching, including for genuinely static assets that rarely change] → acceptable trade-off per the platform's own guidance; this app's static footprint is small enough that re-fetching every request isn't a real cost.
- [Removing `srv.service` loses the reference for anyone spinning this template up on a plain exe.dev/systemd host in the future] → the file is preserved in git history if ever needed again; the README's deployment section is being rewritten for phone-llm anyway, so keeping a now-wrong instruction set alongside a correct one is a worse trade.

## Migration Plan

1. Edit `run`, `srv/server.go`, README/AGENTS docs; delete `srv.service`.
2. `make build && make test` locally to confirm the binary still builds and existing tests pass.
3. Restart via the platform's normal path: `systemctl restart app` (or let the next `app.service` cycle pick up the new `run`).
4. Verify `https://todo.toejough.dev` loads, `/api/todos` responds, and response headers show `Cache-Control: no-store` on `/`, `/sw.js`, and a `/static/*` asset.
5. Rollback: `git revert` the commit and restart `app.service` again — no data migration involved, so rollback is a plain code revert.
