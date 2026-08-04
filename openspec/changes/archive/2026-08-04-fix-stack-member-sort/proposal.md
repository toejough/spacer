## Why

Cards inside an expanded stack currently render in raw storage/insertion order instead of relevance order, unlike every other card list in the app. This is a bug, not a design choice: the outer list (`buildDisplayList` in `srv/static/script.js`) is correctly sorted by `compareByRelevance` (state → review urgency → recency), and a stack occupies the list position of its most urgent member — but the stack's own `members` array, built by a plain `.filter()` with no `.sort()`, skips that same relevance sort. The result is inconsistent behavior: a stack's position among other cards/stacks is relevance-driven, but its contents are not.

## What Changes

- Apply the same `compareByRelevance` sort used for top-level cards to a stack's member list before rendering, so cards within an expanded stack order identically (state, then review urgency, then recency) to plain cards.
- No manual/drag ordering is introduced for stack members — relevance sort governs, consistent with top-level cards.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
- `card-stacks`: adds a requirement that a stack's member cards are ordered by the same relevance rules as top-level cards when the stack is expanded.

## Impact

- `srv/static/script.js`: `buildDisplayList` (member-list construction, ~line 211) — sort `memberPool.filter(...)` results with the existing `compareByRelevance` comparator.
- No other call sites are affected: `getStackMembers` (script.js:63-65) is used only by `gcStacks` for an existence check (`.length > 0`), not for display order.
