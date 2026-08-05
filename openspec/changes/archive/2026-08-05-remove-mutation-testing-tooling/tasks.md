## 1. Delete the orphaned mage-based toolchain

- [x] 1.1 Delete `dev/mutate/` (`mutate.go`, `pretest.go`, `pretest_test.go`, `run.go`, `run_test.go`)
- [x] 1.2 Delete `dev/protest/` (`protest.go`)
- [x] 1.3 Delete `dev/fuzz.fish`
- [x] 1.4 Delete `dev/mutate.fish`
- [x] 1.5 Delete `dev/golangci.toml`
- [x] 1.6 Delete `dev/dev-install.fish` and `dev/dev-install.sh`
- [x] 1.7 Confirm `dev/` is now empty and remove the directory
- [x] 1.8 Delete `magefile.go`

## 2. Clean up dependencies

- [x] 2.1 Run `go mod tidy` and confirm `github.com/alexflint/go-arg`, `pgregory.net/rapid`, and `github.com/magefile/mage` are dropped from `go.mod`/`go.sum`

## 3. Verify nothing is left dangling

- [x] 3.1 `make build && make test` pass
- [x] 3.2 `grep -rniE "mutate|fuzz|protest|magefile|\bmage\b"` across the repo (excluding `node_modules/`, `.git/`, and archived openspec changes, which are historical record) returns nothing referencing the deleted tooling
- [x] 3.3 Confirm no remaining references to `dev/`, `magefile.go`, or `mage` in `Makefile`, `README.md`, `AGENTS.md`, or `openspec/config.yaml`
