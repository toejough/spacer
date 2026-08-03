## Context

`srv/static/script.js` has two separate, ad hoc sort comparators today: `loadTodos` sorts by `archived, done, priority desc, created_at desc`; `loadNotes` filters out archived items entirely and sorts only by `updated_at desc`. Both lists render from the same `localStorage`-backed `loadItems()` array with the same per-item shape (`archived`, `done`, `updated_at`, `next_review`, `cloze_data`, `review_enabled`), just filtered by `item_type`. See proposal.md - Why.

## Goals / Non-Goals

**Goals:**
- One shared comparator function drives both `loadTodos` and `loadNotes`, so the relevance rule is defined once and can't drift between the two lists.
- Archived notes become visible (sorted last) with a Reopen affordance, bringing notes to parity with how todos already handle Abandon/Reopen.

**Non-Goals:**
- No change to the SM-2 review algorithm or how `next_review`/`cloze_data` get computed — this only reads those fields to rank items.
- No change to the Search tab's result ordering (it filters `!i.archived` independently and isn't part of this proposal's scope).

## Decisions

- **Single shared comparator, parameterized by nothing** — `compareByRelevance(a, b)`. Todos and notes share the same item shape (`done` is simply `undefined`/falsy on notes), so one function handles both without a type-specific branch. Simpler than two near-duplicate comparators that could quietly diverge.
- **State rank as a small int**: `archived ? 2 : (done ? 1 : 0)`. Works unmodified for notes since `done` is never set there.
- **Review-urgency rank as a timestamp**: soonest `next_review` (or earliest across `cloze_data` for multi-cloze items) as milliseconds; `Infinity` when `review_enabled === false` or no date exists, so those items naturally sort after anything with a real due date via plain numeric comparison — no special-casing needed in the comparator itself.
- **Notes gain a Reopen button and `.abandoned` class**, mirroring the todo card's existing `btn-reopen`/`.abandoned` treatment, rather than inventing a different pattern for notes.

## Risks / Trade-offs

- [Showing archived notes for the first time changes what a returning user sees in the Notes tab] → acceptable and intentional per this change's proposal; the state tier still pushes them to the bottom, so active notes remain the visually dominant content.
- [Computing review-urgency rank re-derives cloze due dates on every sort] → the item counts here are small (personal single-user app), so recomputing per sort is not a performance concern.

## Migration Plan

1. Add `compareByRelevance` in `srv/static/script.js`; wire it into `loadTodos` and `loadNotes`.
2. Remove the `!i.archived` filter from `loadNotes`; add Reopen button + `.abandoned` class to `renderNoteCard`.
3. Add `.item-card.abandoned` styling for notes in `srv/static/style.css` (can likely reuse the existing todo `.abandoned` rule if selectors are shared, or add an equivalent note-scoped rule).
4. Manual check: verify ordering across mixed open/done/abandoned todos and active/abandoned notes with varying review due dates.
5. Rollback: revert the commit — no data migration involved, since sorting/visibility are purely derived from existing fields.
