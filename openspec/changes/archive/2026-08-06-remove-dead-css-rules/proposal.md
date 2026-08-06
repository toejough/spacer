## Why

Four CSS rules in `srv/static/style.css` have zero matching references anywhere in `srv/templates/index.html` or any template string in `srv/static/script.js` (confirmed by grep across both, including JS-built HTML, not just static markup):

- `.item-content` (line 108)
- `.item-tag` (line 111)
- `.item-date` (line 113)
- `.item-priority` and its `.p1`/`.p2`/`.p3` modifiers (lines 114–117)

`.item-priority`/`.item-tag` line up with the `priority` and `tags` columns in the `items` SQL table schema — leftover styling from before the card UI was redesigned to its current form, and part of the same "SQL schema has fields the live app doesn't use" story being resolved in `remove-unused-database-dependency`. Confirmed via issue #31.

Kept as its own change rather than folded into `remove-stale-ops-and-agent-files` because this is the only cleanup issue touching something the browser actually loads and renders — isolating it keeps a quick visual re-check of the app scoped to one small, easy-to-reason-about change.

## What Changes

- Delete the four dead rules from `srv/static/style.css`.
- **Found during implementation, not in the original audit**: `.review-card .item-content` (a scoped compound-selector variant reusing the same dead `item-content` class) was missed because the original audit's grep only checked `index.html`/`script.js` for the bare class name, not `style.css` itself for other selectors reusing it. Deleted alongside the four originally-scoped rules for the same reason — nothing anywhere renders an element with class `item-content`.

## Capabilities

No capabilities, new or modified — the rules have no matching elements anywhere, so removing them changes no rendered output. `skip_specs: true` is set in `.openspec.yaml`.

## Impact

- `srv/static/style.css` (5 rules removed: the original 4 plus the scoped `.review-card .item-content` variant found during implementation)
- No change to `srv/templates/index.html`, `srv/static/script.js`, or any other file.
