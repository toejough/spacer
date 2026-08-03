# Todo

A personal todo/notes app ("Remember Everything"), running on the
[phone-llm](https://github.com/toejough/phone-llm) platform at
`todo.toejough.dev`.

## Building and Running

Build with `make build`, then run `./todo-srv`. The server listens on port
8000 by default (`-listen` flag to override).

## Deployment

This app runs as the `todo` app machine on phone-llm. `app.service` on that
machine runs `/app/run`, which builds and execs the server bound to `:8080`
(what the platform's Cloudflare tunnel talks to), and restarts it if it dies.

To restart after code changes:

```bash
systemctl restart app
```

Access is gated by Cloudflare Access at the platform level (owner-only) —
the app itself does not implement authorization.

## Database

This template uses sqlite (`db.sqlite3`). SQL queries are managed with sqlc.

## Code layout

- `cmd/srv`: main package (binary entrypoint)
- `srv`: HTTP server logic (handlers)
- `srv/templates`: Go HTML templates
- `db`: SQLite open + migrations (001-base.sql)
