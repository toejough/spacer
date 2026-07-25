# Tasks: fix-review-card-hover-state

## 1. Diagnose and fix the bug

- [x] 1.1 Locate review rendering and submission logic in `srv/static/script.js`.
- [x] 1.2 Identify that hover/focus state carries over to the next card because the DOM is replaced under the pointer.
- [x] 1.3 Add `lockReviewPointerEvents()` and call it from `submitReview()` to suppress pointer events during re-render.
- [x] 1.4 Blur the active review button in `submitReview()` to prevent focus from jumping to the next card.
- [x] 1.5 Add the corresponding CSS rule to `srv/static/style.css`.
- [x] 1.6 Add a regression test in `srv/static/script.test.js`.
- [x] 1.7 Run `node ./srv/static/script.test.js` to verify the fix.
- [x] 1.8 Run `go test ./...` to verify Go tests still pass.

## 2. Update OpenSpec artifacts

- [x] 2.1 Create change proposal `fix-review-card-hover-state`.
- [x] 2.2 Write `proposal.md`, `design.md`, and delta spec under `specs/spaced-review/spec.md`.
- [x] 2.3 Run `openspec validate --all`.
- [x] 2.4 Archive the change.

## 3. Bump version

- [x] 3.1 Update `srv/templates/index.html` from `v18` to `v19`.
- [x] 3.2 Update `srv/static/sw.js` cache name from `remember-everything-v18` to `remember-everything-v19`.
