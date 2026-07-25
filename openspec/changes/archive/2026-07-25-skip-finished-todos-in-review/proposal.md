# Skip finished todos in review

## Why

Completed todos are not actionable for spaced repetition. Showing finished todos in the Review tab clutters the queue and makes the review badge count misleading.

## What Changes

- `srv/static/script.js`: `getDueReviewEntries()` now skips todo items with `done === 1`.
- `srv/static/script.test.js`: regression tests ensure finished todos are excluded from review entries and badge count, while unfinished todos and notes remain eligible.
- `openspec/specs/spaced-review/spec.md`: adds a scenario that finished todos are excluded from review.
- `openspec/specs/todo-list/spec.md`: updates the toggle-completion scenario to require removal from review.
- `Makefile`: adds a `test-js` target for the static-script test suite.
- `srv/templates/index.html` and `srv/static/sw.js`: bump version from v17 to v18.
- `go.mod`/`go.sum`: add missing test dependencies (`github.com/alexflint/go-arg`, `pgregory.net/rapid`).
- `dev/mutate/pretest_test.go` and `dev/mutate/run_test.go`: fix stale import path from `spacer/dev/protest` to `srv.exe.dev/dev/protest`.
- `srv/server_test.go` and `dev/mutate/testMutations_test.go`: remove stale tests for removed API handlers and unimplemented mutation functions.


Finished todos still appear in the Review tab and in the review badge count. A todo marked as done (`done === 1`) is not useful to review, so it should be excluded from the spaced-repetition queue.

## Scope

- Update `srv/static/script.js` so that `getDueReviewEntries()` skips items where `item_type === 'todo'` and `done === 1`.
- Add a regression test in `srv/static/script.test.js` verifying that finished todos are excluded from review entries.
- Update the OpenSpec `spaced-review` spec to require this behavior.
- Mark the `todo-list` spec's "Complete a todo" requirement as implying removal from review.

## Design

In `getDueReviewEntries()`, after the `item.archived` check, add a check that skips todo items with `done === 1`. This keeps the review list focused on actionable items and avoids showing completed todos for SM-2 rating.

