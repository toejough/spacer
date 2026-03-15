# 030 — Introduce data-access query functions

**Status:** open
**Type:** tech-debt
**Source:** #014 architecture premortem

## Problem

All 3 views import `db` directly and write inline Dexie queries (`db.cards.where("deckId").equals(...)`, `db.decks.add(...)`, `db.cards.update(...)`). As features grow (stats, import/export, sync), query logic scatters across many files with no single place to understand or change data access patterns.

## Recommendation

Defer until a 4th consumer of card/deck queries appears. When triggered, extract query functions into `db.ts` alongside the schema — not a separate file. #028 is a first step in this direction.

**Trigger:** a 4th file needs to query cards or decks (e.g., stats dashboard, import/export module).

## Research

- Dexie.js patterns for [encapsulating database operations](https://dexie.org/docs/Tutorial/Design-Patterns) — Dexie's own docs recommend a "repository" pattern for larger apps.
- Vue composables as a query layer — look at how `useLiveQuery` from `dexie` integrates with Vue reactivity. Could replace manual `onMounted` + `ref` patterns with reactive queries that auto-update, which would also solve the "stale data after navigation" class of bugs.
