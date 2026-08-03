## Why

Every closing action in the app today — Done, Abandon — is fully reversible via Reopen. There's no way to actually get rid of something: old completed/abandoned clutter for housekeeping, or a genuine mistake (a fat-fingered note, a todo superseded by another). Users need a real, permanent Delete, distinct from Abandon's reversible archive.

## What Changes

- **BREAKING** (new capability, not a behavior removal): add permanent Delete for todos and notes. This is the first destructive, non-reversible operation in the app.
- Delete is reachable **only from Done or Abandoned** items — the two end-of-life states. Open items cannot be deleted directly; they must close first.
- **Touch**: swipe left past a distance threshold commits the delete immediately (no separate confirm tap) — the standard "reveal on the right" convention. As the card slides, a red panel with a trash icon is revealed behind it in the space it vacates, alongside the card itself reddening progressively the further it's dragged.
- **Undo**: after a swipe commits, a brief toast appears with an Undo action, since there's no backend or automatic backup — the only other recovery path is the existing manual JSON export/import.
- **Keyboard/screen-reader parity**: a "Delete Forever" control in the edit modal, alongside the existing Abandon control, visible only when the item is Done or Abandoned — giving non-touch users the same capability the swipe gives touch users.
- Applies identically to todos and notes, mirroring how Done/Abandon/Reopen already work the same way across both item types.

## Capabilities

### Modified Capabilities
- `todo-list`: new Delete requirement, coexisting with the existing "Done and Abandoned are mutually exclusive end-of-life states" requirement.
- `note-taking`: new Delete requirement, coexisting with the existing "Delete a note" (Abandon) requirement — naming note: the existing "Delete a note" requirement is actually about Abandon, not permanent deletion; this proposal's new requirement is the first true delete.

## Impact

- `srv/static/script.js`: `deleteItemPending`/`commitDelete`/`undoDelete` (pending-delete buffer, nothing removed from `localStorage` until the undo window elapses), the delegated swipe-gesture handler, undo-toast rendering, edit-modal "Delete Forever" wiring
- `srv/static/style.css`: card drag-tint (`.swiping`/`.swipe-armed`), the reveal-behind panel (`.item-card-wrapper`/`.swipe-reveal`), a static `touch-action: pan-y` on swipeable cards (so the browser's own scroll-vs-gesture arbitration handles vertical drags correctly, rather than racing custom JS logic), undo toast styling
- `srv/templates/index.html`: "Delete Forever" control in the edit modal
- No API or server changes — purely client-side, consistent with how Done/Abandon/Reopen already work.
