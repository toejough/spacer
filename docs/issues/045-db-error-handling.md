# 045 — DB operations assume success

**Status:** open
**Type:** quality / prevention
**Source:** #017 implementation premortem

## Problem

No DB operation in the codebase has error handling — `db.decks.add()`, `db.cards.update()`, `db.cards.where()` are all unguarded `await` calls. IndexedDB can fail: quota exceeded, transaction aborted, concurrent tab conflicts, or corruption. At current scale with tiny datasets this never triggers, but import/export, sync, or media storage would introduce real failure modes. When a DB call fails today, the view silently does nothing and the user has no idea what happened.

## Principle

System boundaries (DB, network, filesystem) can fail. The question isn't whether to handle errors but where — at the call site, at a data-access layer, or at a global error boundary. The strategy should be consistent so new features don't each invent their own approach.

## Guidance

Before implementing, research error handling patterns for Dexie/IndexedDB — global error handlers, transaction-level catches, and Vue error boundaries. Consider how this interacts with #032 (user action feedback) — error handling and error display are two sides of the same coin. The goal is a consistent strategy, not try/catch on every line. Read the current data flow to understand what's changed since this issue was filed.
