# 044 — No DB migration pattern

**Status:** open
**Type:** quality / prevention
**Source:** #017 implementation premortem

## Problem

`db.ts` defines `version(1)` with no established pattern for schema evolution. The first time `Card` or `Deck` changes shape (adding tags, media, settings), someone needs to bump the version and write an upgrade function — or forget and corrupt existing user data in IndexedDB. Since this is a PWA with persistent local data, bad migrations can't be fixed by "just redeploy."

## Principle

Schema changes in a client-side DB are irreversible for users — there's no DBA to run a fix. The migration pattern should be established before it's needed under pressure.

## Guidance

Before implementing, research Dexie's [versioning and migration API](https://dexie.org/docs/Dexie/Dexie.version()) — how to add fields with defaults, rename properties, and backfill data. Look at how other Dexie-based PWAs handle migrations. The deliverable should be a documented pattern (not just a convention — an actual example migration) so the first real migration has a template to follow. Read the current schema to understand what's changed since this issue was filed.
