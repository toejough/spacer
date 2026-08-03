## 1. Shared comparator

- [x] 1.1 Add `compareByRelevance(a, b)` in `srv/static/script.js`: state rank (open=0, done=1, archived=2), then review-urgency rank (soonest `next_review`/earliest `cloze_data` due date, or `Infinity` if `review_enabled === false` or no date), then `updated_at` descending
- [x] 1.2 Wire `compareByRelevance` into `loadTodos`, replacing the existing `archived/done/priority/created_at` sort

## 2. Notes: surface archived, add Reopen

- [x] 2.1 Remove the `!i.archived` filter from `loadNotes`; wire `compareByRelevance` into its sort
- [x] 2.2 Add a Reopen button to `renderNoteCard`, shown only when the note is archived (mirroring the todo card's `btn-reopen`)
- [x] 2.3 Add an `.abandoned` class to archived note cards in `renderNoteCard`
- [x] 2.4 Add CSS for archived note cards in `srv/static/style.css` (reuse or mirror the existing todo `.item-card.abandoned` treatment) — no change needed; `.item-card.abandoned` is already a shared, non-todo-scoped selector and note cards use the same `.item-card` base class

## 3. Verify

- [x] 3.1 Manual check: mixed open/done/abandoned todos sort correctly (state, then review urgency, then recency) — verified via VM-harness script exercising `compareByRelevance` directly
- [x] 3.2 Manual check: active and abandoned notes sort correctly, abandoned notes visible and visually distinct — verified via `loadNotes()` + rendered HTML inspection
- [x] 3.3 Manual check: an item with reviews disabled or no scheduled review sorts after a same-tier item with a due date — verified
- [x] 3.4 Manual check: reopening an abandoned note moves it back into the active tier and clears the archived visual state — verified via `reopenItem()`
- [x] 3.5 `make test-js` (or equivalent) passes if any JS tests cover list rendering — all 24 existing tests pass, none broken by this change
- [x] 3.6 Bump asset version (footer + `?v=` cache-busters in `srv/templates/index.html`) from v31 to v32
