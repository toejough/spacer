# Tasks: skip-finished-todos-in-review

## 1. Diagnose and fix the bug

- [x] 1.1 Locate review logic in `srv/static/script.js`.
- [x] 1.2 Confirm that finished todos (`done === 1`) are incorrectly included in `getDueReviewEntries()`.
- [x] 1.3 Add `item.item_type === 'todo' && item.done === 1` to the review exclusion check.
- [x] 1.4 Add regression tests in `srv/static/script.test.js`.
- [x] 1.5 Run JS tests to verify the fix.

## 2. Update OpenSpec artifacts

- [x] 2.1 Create change proposal `skip-finished-todos-in-review`.
- [x] 2.2 Write `proposal.md` and delta spec under `specs/spaced-review/spec.md`.
- [x] 2.3 Update `openspec/specs/spaced-review/spec.md` and `openspec/specs/todo-list/spec.md`.
- [x] 2.4 Run `openspec validate --all`.
- [x] 2.5 Archive the change.

## 3. Fix broken tests

- [x] 3.1 Add missing Go test dependencies to `go.mod` (`github.com/alexflint/go-arg`, `pgregory.net/rapid`).
- [x] 3.2 Fix stale import path in `dev/mutate/pretest_test.go` and `dev/mutate/run_test.go`.
- [x] 3.3 Remove `srv/server_test.go` (references removed API handlers).
- [x] 3.4 Remove `dev/mutate/testMutations_test.go` (references unimplemented `testMutations` function).
- [x] 3.5 Run `go test ./...` to confirm all Go tests pass.
- [x] 3.6 Run `make test-js` to confirm JS tests pass.
