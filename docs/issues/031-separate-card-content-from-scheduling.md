# 031 — Separate card content from scheduling state

**Status:** open
**Type:** tech-debt
**Source:** #014 architecture premortem

## Problem

`Card` in `db.ts` joins two semantic concerns in one interface and one table: content (`front`, `back`, `deckId`) and SM-2 scheduling state (`easeFactor`, `interval`, `repetitions`, `nextReview`). Progressive disclosure principle says these should be held separate.

Currently harmless because every consumer needs both halves together, but the coupling will block clean feature work when the concerns diverge.

## Recommendation

Defer until triggered. Split scheduling state into a separate `ReviewState` table with a `cardId` FK when it happens.

**Triggers (any one):**
- A feature needs card content without scheduling (e.g., export, search index)
- A feature needs scheduling without card content (e.g., stats dashboard, algorithm swap)
- A third semantic concern wants to attach to Card (e.g., review history log)

## Research

- Dexie.js docs on [relations and references](https://dexie.org/docs/Tutorial/Design-Patterns) — patterns for 1:1 related tables and efficient joined reads.
- How Anki separates [notes vs. cards vs. review log](https://github.com/ankitects/anki) — Anki's three-table model (notes hold content, cards hold scheduling, revlog holds history) is the mature version of this split.
