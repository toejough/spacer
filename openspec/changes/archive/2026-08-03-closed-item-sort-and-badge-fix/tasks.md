## 1. Sort fix

- [x] 1.1 In `getReviewUrgencyRank`, return `0` immediately when `item.done || item.archived`, before the existing review-urgency logic, so recency governs ordering directly within the done/abandoned tiers

## 2. Badge text

- [x] 2.1 In `renderTodoCard`, when the todo is done or archived, show relative time since `updated_at` ("Completed …ago" / "Abandoned …ago") instead of the review-info badge
- [x] 2.2 In `renderNoteCard`, when the note is archived, show relative time since `updated_at` ("Abandoned …ago") instead of the review-info badge

## 3. Reopen resets review scheduling

- [x] 3.1 In `reopenItem`, reset `ease_factor: 2.5, interval_days: 0, repetitions: 0, next_review: <now>` on the item
- [x] 3.2 In `reopenItem`, if the item has cloze deletions, reset every entry in `cloze_data` to the same fresh defaults

## 4. Verify

- [x] 4.1 `make test-js` passes — all 24 existing tests pass
- [x] 4.2 Verify: two done items with different stale `next_review` dates now sort by recency, not review date — verified via VM-harness script exercising `compareByRelevance` directly
- [x] 4.3 Verify: two abandoned items sort by recency, not review date — verified
- [x] 4.4 Verify: done/abandoned item badges show relative completed/abandoned time, not next-review info — verified via `renderTodoCard`/`renderNoteCard` output inspection
- [x] 4.5 Verify: reopening a done or abandoned item resets its review scheduling to fresh defaults and it becomes immediately due — verified via `reopenItem()`, including the multi-cloze case
- [x] 4.6 Bump asset version (footer + `?v=` cache-busters in `srv/templates/index.html`) from v32 to v33
