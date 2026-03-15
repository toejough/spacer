# 043 — Type assertions bypass safety at the DB boundary

**Status:** open
**Type:** quality / prevention
**Source:** #017 implementation premortem

## Problem

`DeckView.vue` uses `as Card` and `HomeView.vue` uses `as Deck` when calling `db.*.add()`. The flow test uses `as any`. These casts tell TypeScript "trust me" and suppress errors for missing or mistyped fields. If `Card` or `Deck` gains a required field via schema change, these call sites will silently write incomplete records with no compile-time error.

## Principle

The DB boundary is a system boundary — it should be validated, not cast through. Types at insertion points should be structurally verified by the compiler, not asserted by the developer.

## Guidance

Before implementing, investigate why the `as` casts are needed — likely because Dexie's `add()` expects a type with `id` but the caller doesn't provide one (auto-incremented). Research Dexie's `EntityTable` and whether there's a clean way to type inserts without the `id` field (e.g., `Omit<Card, "id">` or Dexie's built-in creation types). The fix should make incomplete inserts a compile error. Read the current DB types and insertion sites to understand what's changed since this issue was filed.
