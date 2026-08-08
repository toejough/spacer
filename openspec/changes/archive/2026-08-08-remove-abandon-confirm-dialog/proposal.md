## Why

Abandoning a todo or note pops a native `confirm('Abandon this item?')` dialog before it takes effect. That's the only `confirm()` call anywhere in the app — even permanent delete, the one genuinely hard-to-recover action, doesn't use a blocking dialog; it deletes immediately and offers an Undo toast instead. Abandon is strictly less risky than that: an abandoned item stays visible in its list (sorted last, visually marked "Abandoned"), and Reopen is always available from the card or the edit modal. The confirmation dialog protects against a mistake that was never actually costly, and it's inconsistent with the immediate-action-plus-recovery pattern the app already uses everywhere else.

## What Changes

- Remove the `confirm('Abandon this item?')` call from `archiveItem` (`srv/static/script.js`); abandoning becomes an immediate, single-click action from both the card's Abandon button and the edit modal's Abandon button (both call `archiveItem`).
- Update the two tests that currently assert a confirm dialog fires on abandon (`archiveItem allows abandoning completed todo`, `abandonFromModal allows archiving done todo`) to reflect that no dialog is shown.
- Tighten the `todo-list` spec's "confirming Abandon" wording (in the "Abandon is available in edit modal UI" scenario) to "clicking Abandon," removing language that could be misread as implying a confirmation step — the adjacent "Delete from edit modal" scenario in both `todo-list` and `note-taking` already describes Abandon as a direct action with no dialog, so this aligns the ambiguous sentence with what was already specified elsewhere.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
- `todo-list`: "Abandon a todo" requirement's "Abandon is available in edit modal UI" scenario no longer implies a confirmation step.

## Impact

- `srv/static/script.js`: `archiveItem`.
- `srv/static/script.test.js`: the two tests asserting `confirms.length === 1` on abandon.
- `openspec/specs/todo-list/spec.md`: updated via this change's delta and eventual archive-sync. (`note-taking`'s parallel "Abandon a note" requirement already has no confirm-implying language, so it needs no delta.)
