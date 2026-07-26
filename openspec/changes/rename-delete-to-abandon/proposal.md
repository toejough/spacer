# Rename 'Delete' action to 'Abandon'

## Why

The UI currently labels the removal action as "Delete", which is ambiguous in this app because items are archived rather than permanently destroyed. Using the term "Abandon" better communicates that the item is being set aside (archived) and can be recovered later if needed. It also aligns wording across UI and specs.

## What Changes

- Update the UI label/tooltips and confirmation text from "Delete" to "Abandon".
- Update OpenSpec requirements to replace "Delete" with "Abandon" in the `todo-list` and `note-taking` specs.
- Bump the web asset/service worker version markers from `v22` to `v23` so clients pick up the updated wording.

## Impact

- Code: `srv/static/script.js`, `srv/templates/index.html`, `srv/static/sw.js`.
- Specs: `openspec/specs/todo-list/spec.md`, `openspec/specs/note-taking/spec.md` will be updated via this change and archived to the main specs.
- No DB or API changes.

