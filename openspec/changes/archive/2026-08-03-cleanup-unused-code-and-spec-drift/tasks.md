## 1. Delete unused code

- [x] 1.1 Delete `web-app/` (entire directory)
- [x] 1.2 Remove `handleTodos` and `handleTodo` from `srv/server.go`, and their mux registrations (`/api/todos`, `/api/todos/`)
- [x] 1.3 Delete `releases/remember-everything-v24.tar.gz`

## 2. Verify nothing else references what was removed

- [x] 2.1 `make build && make test` pass after removing the `/api/todos` handlers
- [x] 2.2 Confirm no remaining references to `web-app/` in `Makefile`, `magefile.go`, `README.md`, `AGENTS.md`, or CI/build config (also fixed a stale mention in `openspec/config.yaml`)
- [x] 2.3 Confirm no remaining references to the deleted handlers or `releases/remember-everything-v24.tar.gz` anywhere in the repo

## 3. Spec accuracy

- [x] 3.1 Confirm main spec sync applies the "checkbox" → "Complete button" wording fix
- [x] 3.2 Confirm main spec sync replaces "Mobile parity for Done and Abandon (IMPLEMENTED)" with "Done and Abandoned are mutually exclusive end-of-life states," with the duplicated mutual-exclusion scenario consolidated into the new requirement only
