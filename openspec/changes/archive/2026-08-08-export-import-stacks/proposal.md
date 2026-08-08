## Why

A user exported their data, reinstalled the app (wiping local storage), and imported the backup back in — and lost every stack they'd created. `buildExportPayload()` only serializes `items` via `loadItems()`; it never calls `loadStacks()`, so stack names and membership were never in the backup file to begin with. The member cards survive (their `stack_id` field rides along inside each item), but with no stack record to resolve that id against, they silently render back as ordinary unstacked cards. Card-stacks shipped after data-portability's export/import feature, and the two were never reconciled — the export requirement says "all local todo/note data," but stacks are local data it doesn't cover.

Separately, the same investigation surfaced that import gives no feedback at all: no count of items added, none skipped as duplicates, nothing about stacks. That silence is what made this bug hard to even diagnose — the user had no way to tell, from the app itself, whether a restore was complete.

## What Changes

- Export SHALL include stacks (id, name, created/updated timestamps) alongside items, so a backup captures the full local dataset.
- Import SHALL restore stacks and re-link imported items to them. Since imported items already get fresh, locally-unique ids to avoid colliding with existing items, imported stacks need the same treatment: assign fresh stack ids on import and rewrite each imported item's `stack_id` to match, rather than trusting the exported ids to still be free.
- Import SHALL merge into any existing local stack that has the same name, instead of creating a duplicate stack with the same name (mirrors the existing item-level "existing wins on exact match" dedup rule, applied to stacks by name).
- Import SHALL report a summary after it completes: items added, items skipped as duplicates, stacks added, stacks merged. No more silent no-op.
- `data-portability`'s "export all data" and "import a file" requirements are updated to explicitly cover stacks, closing the gap left when `card-stacks` shipped without a matching update here.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
- `data-portability`: export payload gains a `stacks` array; import restores stacks (merging by name into existing stacks) and re-links imported items to the resulting stack ids; import surfaces a result summary instead of completing silently.

## Impact

- `srv/static/script.js`: `buildExportPayload`, `applyImportedItems`, `upgradeExportData`, `handleImportFileChange` (or equivalent) — export/import logic and its user-facing feedback.
- `srv/static/script.test.js`: new/updated tests for stack export, stack import (fresh install and merge-by-name-into-existing cases), item `stack_id` re-linking, and the import summary.
- `openspec/specs/data-portability/spec.md`: requirements updated via this change's delta and eventual archive-sync.
- No server-side (`cmd/srv`, `srv/server.go`) or storage-schema changes — this is entirely client-side localStorage logic; `EXPORT_SCHEMA_VERSION` likely needs a bump since the payload shape changes (existing exports without a `stacks` field must still import cleanly).
