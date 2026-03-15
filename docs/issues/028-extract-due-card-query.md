# 028 — Extract shared due-card query

**Status:** open
**Type:** tech-debt
**Source:** #014 architecture premortem

## Problem

`DeckView.vue` and `ReviewView.vue` both independently filter `c.nextReview <= now` to determine which cards are due. If due-card semantics change (e.g., cram mode, priority ordering, timezone handling), both views must be updated in lockstep or they drift.

## Recommendation

Extract a `getDueCards(db, deckId)` function co-located in `db.ts`. Both views call it instead of inlining the filter.

## Research

- Dexie.js docs on [WhereClause](https://dexie.org/docs/WhereClause/WhereClause) — check whether `nextReview` can be filtered via index (`where("nextReview").belowOrEqual(now)`) combined with a compound query on `deckId`, which would be more efficient than filtering in JS as card count grows.
