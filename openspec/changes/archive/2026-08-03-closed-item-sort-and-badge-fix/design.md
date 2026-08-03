## Context

`compareByRelevance`/`getReviewUrgencyRank` and the reviewInfo badge logic live once in `srv/static/script.js`, shared by `renderTodoCard` and `renderNoteCard` (per `relevance-sort-todos-notes`'s design decision to keep one comparator for both). `reopenItem` currently only clears `archived`/`done` and re-stamps `updated_at`; it never touches `next_review`, `ease_factor`, `repetitions`, or per-cloze `cloze_data`. See proposal.md - Why.

## Goals / Non-Goals

**Goals:**
- Recency governs ordering within the done/abandoned tiers directly, not as a rarely-reached tiebreaker behind a stale field.
- Reopen produces a predictable, fresh review state rather than silently resurrecting whatever was there before closing.

**Non-Goals:**
- No change to the SM-2 algorithm itself (`calculateSM2`) — only to what state an item starts from after Reopen.
- No change to open-item ordering — review urgency still governs there exactly as `relevance-sort-todos-notes` specified.

## Decisions

- **Neutralize review-urgency for closed items inside `getReviewUrgencyRank`** (proposal's option 2, chosen over a state-conditional comparator): `if (item.done || item.archived) return 0;` before the existing logic. One-line change, keeps `compareByRelevance` itself untouched, and keeps the "what ranks first" logic in the one function already responsible for review-urgency ranking.
- **Reopen resets to the same fresh-item defaults used elsewhere in the codebase**: `ease_factor: 2.5, interval_days: 0, repetitions: 0, next_review: now` — these are the exact literals `ensureClozeData` already uses when seeding a new cloze, so this isn't a new convention, just applying the existing one on Reopen too. For multi-cloze items, reset every entry in `cloze_data` the same way rather than only the item-level fields.
- **Badge text branches on `item.done || item.archived`**, reusing `formatDate(item.updated_at)` rather than adding a new relative-time formatter or a new timestamp field. `updated_at` is already the single field both `toggleTodo` and `archiveItem` stamp at the moment of the state change, so it's already exactly "time since completed/abandoned" with no additional bookkeeping.

## Risks / Trade-offs

- [Resetting review scheduling on Reopen discards the item's prior ease factor / interval history] → intentional: that history describes performance *before* the item went dormant, which is stale context once it re-enters rotation. Starting fresh is more honest than carrying forward a number that no longer reflects anything real.
- [Badge text change means done/abandoned items no longer show any review-due info, even for someone curious about it] → acceptable; that information isn't actionable once an item is closed (it can't be reviewed in that state), and `showHistory` (already wired to the badge's click handler) remains available for anyone who wants the full record.

## Migration Plan

1. Edit `getReviewUrgencyRank`, `renderTodoCard`, `renderNoteCard`, `reopenItem` in `srv/static/script.js`.
2. Bump asset version.
3. `make test-js` to confirm no existing tests broke; add/extend verification for the done/abandoned-tier recency ordering, badge text, and Reopen's scheduling reset.
4. Rollback: revert the commit — no data migration, since this only changes what `reopenItem` writes going forward and how existing fields are read for sorting/display.
