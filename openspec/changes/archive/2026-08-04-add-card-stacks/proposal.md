## Why

Related todos and notes currently have no way to be grouped together — they compete individually for position in the relevance-sorted list, making it hard to see that several cards belong to the same errand, project, or topic. A lightweight, user-defined grouping ("stack") lets related cards travel together in the UI without changing how each card's own state, review schedule, or type is tracked.

## What Changes

- New concept: a **stack** — a named, user-created group of cards (todos and/or notes, mixed types allowed).
- Cards can be assigned to a stack via the existing edit modal (accessible/complete path).
- Drag-and-drop (new to this app): each card gets a drag handle.
  - Dragging card A onto card B prompts for a stack name and creates a new stack containing both; canceling the prompt aborts the merge (a stack is never created without a name).
  - Dragging a card onto an existing stack adds it to that stack, no prompt.
  - Dragging a card out of a stack (drag handle, dropped outside the stack) removes it from the stack.
  - Dragging a stack onto another stack is not supported.
- Stacks render as a single collapsed tile per tab (name + count of members visible on that tab), expandable/collapsible to reveal member cards.
- A stack's position in a tab's relevance-sorted list is governed by its most urgent member that is visible on that tab (matches the tab's type filter). A stack with no members of the current tab's type does not appear on that tab.
- Search results can show mixed-type stacks (a stack with both todo and note members), also collapsed.
- Stack renaming is available inline (e.g. via the collapsed tile) and via the edit modal.

## Capabilities

### New Capabilities
- `card-stacks`: Defines stacks as a first-class grouping of cards — creation (via drag-merge or edit modal), naming/renaming, membership (add/remove via drag or edit modal), collapse/expand, and the rule that a stack sorts by its most urgent visible member.

### Modified Capabilities
- `todo-list`: Todos tab must render stack membership (collapsed stack tiles, per-tab visible-member counts) instead of always rendering every todo as an independent flat card, and must incorporate stacks into the tab's relevance sort.
- `note-taking`: Notes tab gains the equivalent stack rendering and sort behavior for note cards.
- `search`: Search results must be able to show mixed-type stacks (todo + note members together), collapsed like tab views.

## Impact

- Data model: new `stack` entity (id, name, ordered member ids) persisted in `localStorage` alongside the existing items array; items gain a `stack_id` field.
- Frontend (`srv/static/script.js`, `srv/templates/index.html`, `srv/static/style.css`): new drag-and-drop infrastructure (drag handles, pointer-based drag/drop, drop-target detection) coexisting with the existing swipe-to-delete gesture; new stack rendering (collapsed tile, expand/collapse, inline rename); changes to the relevance-sort logic to account for stack grouping; edit modal gains a stack field.
- No backend/API or database changes — this app's item and stack data lives entirely client-side.
