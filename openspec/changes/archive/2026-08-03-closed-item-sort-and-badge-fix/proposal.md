## Why

The relevance-sort ordering shipped in `relevance-sort-todos-notes` ranks items by state, then review urgency, then recency. Within the done/abandoned tiers, review urgency is a stale, meaningless field — done items are excluded from spaced repetition, so `next_review` is frozen at whatever it was before completion — but because it outranks recency, a just-completed item can sort below older completed items instead of near the top of its subgroup. The review-info badge on closed cards also still shows "Next review: …", which is misleading once review no longer applies.

## What Changes

- `getReviewUrgencyRank` returns a constant for done/archived items, collapsing the review-urgency tier to a tie so recency decides ordering directly within the done and abandoned tiers. Open items are unaffected.
- The review-info badge on done/abandoned item cards (both todos and notes) shows relative time since completion/abandonment (e.g. "Completed 3d ago", "Abandoned yesterday") instead of next-review info, reusing the existing `formatDate` relative-time formatting on `updated_at`.
- `reopenItem` resets review scheduling to fresh defaults (`ease_factor: 2.5, interval_days: 0, repetitions: 0, next_review: now`), applied per-cloze for multi-cloze items, instead of leaving the stale pre-close schedule in place. Today a reopened item silently carries whatever review state it had before closing, which can leave it arbitrarily overdue or not due for months.
- Bump asset version (currently v32) to reflect the shipped change, following the existing pattern.

## Capabilities

### Modified Capabilities
- `todo-list`: ordering requirement refined so recency, not review urgency, governs relative order within the done and abandoned tiers; Reopen now resets review scheduling to fresh defaults.
- `note-taking`: same ordering refinement; review-info display requirement updated for the abandoned state; Reopen now resets review scheduling to fresh defaults.

## Impact

- `srv/static/script.js`: `getReviewUrgencyRank`, `renderTodoCard`, `renderNoteCard`, `reopenItem`
- `srv/templates/index.html`: asset version bump
- No API or data model changes — `updated_at` is already stamped by the existing `toggleTodo`/`archiveItem`/`reopenItem` functions.
