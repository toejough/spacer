## Context

See proposal.md - Why. Relevant current shape (`srv/static/script.js`):

- `loadItems`/`saveItems` and `loadStacks`/`saveStacks` are independent localStorage-backed stores (`remember_everything_items`, `remember_everything_stacks`). Items reference stacks one-directionally via `item.stack_id`; stacks don't list their members.
- `buildExportPayload()` returns `{ schema_version, exported_at, items }` — `EXPORT_SCHEMA_VERSION` is currently `1`.
- `applyImportedItems(items)` dedupes by `itemContentKey` (`type::trimmed title`) against locally-existing items, appends the rest with fresh sequential ids from `getNextId`, and gives no return value or user feedback.
- `getNextId(items)` / `getNextStackId(stacks)` both compute `max(existing ids) + 1`, so imported items and imported stacks each need their own id-remapping pass to land on ids that don't collide with what's already local.
- `handleImportFileChange` wraps the whole import in try/catch and `alert()`s on failure; there's no success path message today.

## Goals / Non-Goals

**Goals:**
- Export/import round-trips stacks and stack membership losslessly.
- Old export files (no `stacks` field) keep importing cleanly.
- The user gets a visible summary after import, using the app's existing feedback idioms rather than a new UI subsystem.

**Non-Goals:**
- Not changing how stacks are created/renamed/merged during normal use (drag-and-drop, edit modal) — only import-time stack creation/merging.
- Not deduplicating items against each other *within* one import file (existing behavior; out of scope, not something this change touches).
- Not building a preview/confirm-before-import step. Import stays a single action, per the existing "Single import action" requirement.

## Decisions

**Stack id remapping mirrors the existing item id remapping.** Build an `oldStackId -> localStackId` map before touching items: for each stack in the imported file, look up an existing local stack by trimmed exact-name match. If found, map to its id (merge). If not, create a new local stack (fresh id via `getNextStackId`) and map to that. This is the same shape as `applyImportedItems`'s existing `nextIdCounter` handling for items — one counter/lookup pass up front, then a straightforward rewrite pass — so it stays consistent with code already in the file rather than introducing a different pattern.

Alternative considered: trust the imported stack ids as-is (only remap on collision). Rejected — the whole bug we're fixing is "an id from a stale export silently pointing at the wrong local thing"; always remapping through the name-keyed map removes that failure mode entirely instead of narrowing it.

**Merge key is stack name, exact match after trim (case-sensitive).** Matches the item-level content-key rule (`type::trimmed title`, case-sensitive) already in the spec, so the two dedup rules read the same way to anyone auditing the code. Case-insensitive matching was considered and rejected: it would silently merge two stacks the user deliberately named differently by case, and nothing else in the app treats titles/names case-insensitively.

**Rewrite `item.stack_id` for every imported item that had one, using the map from step 1 — including items that get skipped as content-duplicates.** Skipped items don't touch local storage at all (existing local item wins, per spec), so there's nothing to rewrite for them; the map is only consulted while building `toAppend` for items that are actually added. This falls out of the existing per-item loop structure, not a new mechanism.

**Schema version bumps to 2, with `upgradeExportData` defaulting a missing `stacks` field to `[]`.** `EXPORT_SCHEMA_VERSION` becomes `2`; `buildExportPayload` adds `stacks: loadStacks()`. `upgradeExportData` currently just validates `data.items` is an array and returns it; it changes to return `{ items, stacks: Array.isArray(data.stacks) ? data.stacks : [] }` so a v1 file (no `stacks` key) imports as "zero stacks" per the spec's forward-compatibility requirement, without a version-number branch — the field is simply optional in the input shape rather than gated on the version number. This keeps the upgrade path a plain shape-normalization instead of growing an if/else per schema version for a single-field addition.

**Import summary uses `alert()`, matching the existing failure-path idiom.** `handleImportFileChange` already `alert()`s on failure; success now `alert()`s a one-line summary (e.g. "Imported: 12 items added, 3 skipped as duplicates, 2 stacks added, 1 stack merged."). Alternative considered: a toast via the existing `undoToastContainer`. Rejected for this change — toasts in this codebase are purpose-built for undo (they carry a countdown/cancel action), and building a second toast variant just to avoid a blocking `alert()` is a UI investment beyond what "make failures/results visible" requires. Nothing here forecloses upgrading to a toast later if `alert()` proves annoying in practice.

## Risks / Trade-offs

- [`alert()` is blocking/modal, which is mildly heavier than a toast] → Acceptable: import is a rare, deliberate action (Help tab, explicit file picker), not a hot-path interaction where a blocking dialog would be disruptive.
- [Always remapping stack ids (rather than trusting them) means two devices' exports can never be merged by id, only by name] → Intentional; matches how item merging already works (by content, not id) and avoids the collision bug this change exists to fix.
- [Merging by exact name means a stack renamed on one device and not the other creates a second stack on next import, rather than merging] → Same shape as existing item-content dedup (a retitled item also stops matching); consistent, not a new class of surprise.

## Migration Plan

No data migration needed — existing localStorage on any given device is untouched by this change; only the export/import code path changes. Rollout is a normal code deploy (`make build` + restart via the existing deploy path). No rollback concerns beyond reverting the commit: old export files remain importable (that's the explicit "pre-stacks export file" scenario in the spec), and this change doesn't alter the existing `items`/`stacks` localStorage keys or shapes, only the export payload and the import routine that reads it.
