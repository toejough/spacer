## 1. Entrypoint

- [x] 1.1 Rewrite `/app/run` to `make build` then `exec ./todo-srv -listen :8080`
- [x] 1.2 Confirm `db.sqlite3`'s relative path still resolves correctly under `run`'s working directory

## 2. Cache headers

- [x] 2.1 Change `Cache-Control` from `no-cache` to `no-store` in `handleIndex` (srv/server.go)
- [x] 2.2 Change `Cache-Control` from `no-cache` to `no-store` in the `sw.js` handler (srv/server.go)
- [x] 2.3 Add `Cache-Control: no-store` to the `/static/` file server (srv/server.go)

## 3. Cleanup

- [x] 3.1 Delete `srv.service`
- [x] 3.2 Rewrite README.md's deployment section for phone-llm conventions; drop exe.dev systemd install steps and `X-ExeDev-*` auth mentions
- [x] 3.3 Update AGENTS.md/README.md tech-stack references from "exe.dev" to phone-llm where they describe deployment

## 4. Verify

- [x] 4.1 `make build && make test` pass locally
- [x] 4.2 Restart via `systemctl restart app` (or next `app.service` cycle) and confirm `https://todo.toejough.dev` loads
- [x] 4.3 Confirm `/api/todos` still responds correctly
- [x] 4.4 Confirm response headers show `Cache-Control: no-store` on `/`, `/sw.js`, and a `/static/*` asset
