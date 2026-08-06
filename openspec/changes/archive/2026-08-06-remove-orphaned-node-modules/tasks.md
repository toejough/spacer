## 1. Delete

- [x] 1.1 `git rm -r node_modules` (repo root only — do not touch `.claude/skills/run-remember-everything/node_modules`, a separate and legitimate dependency tree)
- [x] 1.2 `git rm .env.example`
- [x] 1.3 Add `node_modules/` to `.gitignore`

## 2. Verify

- [x] 2.1 Confirm no `package.json` referenced `node_modules` at the root (`find . -maxdepth 2 -iname package.json` still shows none outside skills/)
- [x] 2.2 Confirm `.claude/skills/run-remember-everything/node_modules` is untouched and its own tests/driver still work
- [x] 2.3 `make build && make test` pass
- [x] 2.4 `git status` shows a clean, expected diff (372+ deletions, one `.gitignore` addition, no unexpected files touched)
