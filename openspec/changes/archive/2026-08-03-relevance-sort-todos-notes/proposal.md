## Why

Todos currently sort by archived/done/priority/created_at, and notes sort only by updated_at with archived notes hidden entirely. Neither ordering reflects what's actually relevant right now: items still open, items due for review soon, and items touched recently should surface first, in that priority order, for both todos and notes.

## What Changes

- Replace both list sort comparators with a single three-tier ordering, applied to todos and notes alike:
  1. **State**: open before done before abandoned (notes have no "done" state, so this collapses to open before abandoned for notes).
  2. **Review urgency**: items due for review soonest come first; items with review disabled or nothing scheduled rank after any item with an actual due date. Multi-cloze items use the earliest `next_review` across their clozes.
  3. **Recency**: within a tier, most recently updated (`updated_at`) first.
- **Notes list now shows archived notes** (sorted last, per the state tier), a change from today's behavior where archived notes are hidden entirely from the Notes tab.
- Add a Reopen action to note cards (mirroring the todo card's existing Reopen button), since archived notes are now visible and need a way back to active.
- Give archived note cards a distinct visual treatment (mirroring the todo card's `.abandoned` styling), so an archived note doesn't look identical to an active one.

## Capabilities

### Modified Capabilities
- `todo-list`: sort order requirement changes from priority/created_at-based to the state/review/recency scheme.
- `note-taking`: "Delete a note" no longer removes the note from all active-list views — it remains visible (sorted last) with a Reopen action, mirroring how todos already handle Abandon.

## Impact

- `srv/static/script.js`: `loadTodos`, `loadNotes`, `renderNoteCard` (Reopen action added), a shared sort comparator
- `srv/static/style.css`: archived-note visual treatment
- No API or data model changes — sorting and visibility are purely client-side/localStorage-driven, as the rest of the app already is.
