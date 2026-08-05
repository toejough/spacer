## Why

This repo's GitHub remote is still named `spacer`, from its original life as a standalone Go mutation-testing CLI (2022, TDD "RED"/"GREEN" commits). That tool — `dev/mutate`, `dev/protest`, `dev/fuzz.fish`, `dev/mutate.fish` — was scaffolded with real dependency-injected code but abandoned mid-implementation in December 2023 (`testMutations()` is a literal `panic("...is unimplemented")`). Every commit since has been the "Remember Everything" todo/notes app; the module was even renamed `spacer` → `srv.exe.dev`. Six GitHub issues (#9, #10, #11, #12, #14, #24) tracked ideas for finishing that tool; the repo owner decided not to resume it, closed those issues, and asked for the tooling itself to be removed so dead code stops inviting confusion about what's actually running.

## What Changes

- Delete the entire `dev/` directory: `dev/mutate/` (`mutate.go`, `pretest.go`, `pretest_test.go`, `run.go`, `run_test.go`), `dev/protest/protest.go`, `dev/fuzz.fish`, `dev/mutate.fish`, `dev/golangci.toml`, `dev/dev-install.fish`, `dev/dev-install.sh`.
- Delete `magefile.go` entirely.
- **Revised from the original proposal**: earlier drafts of this change kept `magefile.go` and `dev/golangci.toml`/`dev/dev-install.fish` on the theory that their general targets (`Lint`, `Test`, `Tidy`, `Monitor`, `Check`, `InstallTools`, `Clean`) were current tooling for the todo app, by analogy with the precedent set by `cleanup-unused-code-and-spec-drift` (which used "referenced by `Makefile`/`magefile.go`" as its bar for dead code). Tracing actual callers found that reasoning was wrong: `magefile.go` itself has no caller. The documented, actually-used build path (`Makefile`: `build`/`clean`/`test`/`test-js`/`openspec-*`, referenced from `README.md`) never invokes `mage` or anything under `dev/`. No CI exists. The only way any `mage` target ever runs is a human manually typing `mage <target>`, and that isn't documented anywhere current. So the whole mage-based toolchain is exactly as orphaned as the mutation-testing tool it was built to gate — none of it is current production code or tooling for the todo app.
- Run `go mod tidy` to drop three now-unused dependencies: `github.com/alexflint/go-arg` (only imported by `dev/mutate/mutate.go`), `pgregory.net/rapid` (only imported by `dev/mutate/pretest_test.go`), and `github.com/magefile/mage` (only imported by `magefile.go`, under its `//go:build mage` tag).

## Capabilities

No capabilities, new or modified — this is pure dev-tooling and test deletion with zero effect on the app's behavior, API, or specs. `skip_specs: true` is set in `.openspec.yaml`.

## Impact

- `dev/` (deleted entirely — 11 files)
- `magefile.go` (deleted)
- `go.mod` / `go.sum` (three dependencies dropped via `go mod tidy`)
- No change to `srv/`, `db/`, `cmd/`, `Makefile`, or any user-facing behavior.
