## Why

The app is moving off exe.dev to the self-hosted phone-llm platform. The platform-level move is already done — this checkout runs on the `todo` app machine and `todo.toejough.dev` is already routed and Access-gated. What's left is adapting the app itself: it still carries exe.dev-era deployment assumptions (a stale Vite-based `run` starter, `no-cache` instead of the platform's required `no-store`, and a systemd unit/docs pointing at exe.dev paths) that don't match how phone-llm actually serves and caches apps.

## What Changes

- Rewrite `/app/run` to build and exec the Go server bound to `:8080` (`make build` then `./todo-srv -listen :8080`), replacing the platform-default `npm run dev --host 0.0.0.0 --port 8080` starter left over from `newapp` (this repo has no npm/Vite).
- Change `Cache-Control` from `no-cache` to `no-store` on `handleIndex` and the `sw.js` handler in `srv/server.go`, and add an explicit `Cache-Control: no-store` to the `/static/` file server, which currently sets no cache header at all. Cloudflare edge-caches `.js`/`.css`/`.svg`/`.png` for hours otherwise, which can hide deploys or poison cached asset URLs.
- **BREAKING** (deployment only, no user-facing effect): retire `srv.service`, the exe.dev-specific systemd unit — phone-llm's own `app.service` runs `/app/run` and restarts it, making this unit dead and misleading.
- Update `README.md`/`AGENTS.md` to drop exe.dev-specific deployment instructions (systemd install steps, `X-ExeDev-*` auth header mentions) that no longer apply. The auth-header text was already unused in code (never read anywhere), so removing it has no functional effect — it's stale documentation only.
- Verify `db.sqlite3`'s relative path still resolves under the new `run` script's working directory.

## Capabilities

No spec-level behavior changes — this is a deployment/infra adaptation of an existing template to a new host. `skip_specs: true` is set in `.openspec.yaml`.

## Impact

- `run` (platform entrypoint script)
- `srv/server.go` (Cache-Control headers)
- `srv.service` (removed)
- `README.md`, `AGENTS.md` (deployment docs)
- No API, schema, or user-facing behavior changes.
