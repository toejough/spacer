## Why

Every todo and note carries four fields — `content`, `priority`, `due_date`, `tags` — that are set once at creation (`quickAdd`) and never read or written again anywhere in the app: no UI exposes them, no sort or filter logic uses them, and no spec (other than the two touched by this change) describes them as a feature. They've been present since the very first commit and appear to be scaffold from whatever the app started as. `content` is the one exception worth calling out: it's referenced once, in `doSearch()`, as `i.content.toLowerCase().includes(q)` — but since `content` is always `''`, that clause is permanently `false` and never contributes a match. It's a dead branch that still executes on every search, not just unused data.

That dead branch exists because `openspec/specs/search/spec.md` has a "Search by content" requirement that was seemingly never actually implemented — no UI has ever given the user a way to write to `item.content` (todos and notes are title-only). The scenario has been unfulfillable since it was written. This change treats that as dead/aspirational rather than a feature to build now: it removes the field and the requirement together, rather than leaving a spec promise the code can't keep.

## What Changes

- Remove `content`, `priority`, `due_date`, and `tags` from the item shape everywhere they're set (`quickAdd` in `srv/static/script.js`, and the equivalent test fixture `baseItem` in `srv/static/script.test.js`).
- Remove the dead `i.content.toLowerCase().includes(q)` clause from `doSearch()`; search continues to match on title only, which is the only thing it has ever actually matched on in practice.
- **BREAKING**: Remove the "Search by content" requirement from the `search` capability spec — it documented behavior that was never reachable, since nothing has ever populated `content`. Search remains title-only, now honestly documented as such.
- Update `data-portability`'s "Import data from a file" requirement text and its "existing wins on duplicate" scenario, which currently name `priority`, `tags`, and `due date` in their metadata lists — drop those three from the list, since the fields they refer to no longer exist. The rest of that requirement (append-only import, content-based dedup, existing-item-wins) is unaffected.
- No change to `todo-list` or `note-taking` specs — their "due date" language refers to spaced-repetition `next_review`, a distinct, still-used field, not the removed `due_date` field.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
- `search`: removes the "Search by content" requirement (dead/unreachable — nothing has ever written to `content`). Search remains title-only.
- `data-portability`: the "Import data from a file" requirement's metadata-field lists drop `priority`, `tags`, and `due date`, since those fields are removed from the item shape. Append-only import behavior and content-based (type+title) dedup are unchanged.

## Impact

- `srv/static/script.js`: `quickAdd` (item shape), `doSearch` (drop dead clause).
- `srv/static/script.test.js`: `baseItem` fixture, and any assertions referencing `priority`/`due_date`/`tags`/`content` on items.
- `openspec/specs/search/spec.md`, `openspec/specs/data-portability/spec.md`: updated via this change's deltas and eventual archive-sync.
- Existing localStorage data and old export files may still contain these fields for items created before this change — that's harmless; they're simply never read, same as today, and import continues to work (it spreads whatever fields an imported item happens to have). No migration needed.
- No server-side (`cmd/srv`, `srv/server.go`) changes — out of scope here, tracked as a separate change for the unrelated `Server.Hostname` dead field.
