# Board view specification (delta)

This spec describes the board UI with columns and time-based swimlanes and the behaviors required.

## Requirement: Board view

The system SHALL present a Board view accessible from the Todos tab or navigation.

### Columns
- The Board view SHALL show three columns labelled: Todo, Done, Abandoned.
- Each column SHALL display cards representing todo items.

### Swimlanes
- Each column SHALL be partitioned into swimlanes by the item's last closed/archived timestamp (or "Open" for active items in Todo):
  - Open / closed in last week (0–7 days)
  - Closed in last month (8–30 days)
  - Closed in last quarter (31–90 days)
  - Closed in last year (91–365 days)
  - Everything else (>365 days or null)
- The system SHALL compute swimlane for a todo using the latest of completed_at and archived_at timestamps where applicable.

### Drag-and-drop
- The system SHALL allow dragging a card from the Todo column to a target swimlane inside Done or Abandoned.
- Dropping into Done SHALL mark the todo as done (set completed_at to now) and refresh the board so the card appears in the correct swimlane.
- Dropping into Abandoned SHALL mark the todo archived (set archived_at to now) and remove it from active Todo lists.
- The system SHALL allow reordering cards within a swimlane; reorder order SHALL be persisted.

### Keyboard & accessibility
- All drag-and-drop actions SHALL have keyboard-accessible alternatives (move to Done/Abandon via item menu and keyboard shortcuts) and ARIA live announcements for status changes.

### Mobile and responsiveness
- On narrow viewports the Board view SHALL collapse into a stacked list grouped by swimlane and column, preserving the ability to change status via controls rather than drag-and-drop.

### API notes
- Add/extend endpoints:
  - PATCH /api/todos/:id — accept fields { status: "done" | "archived", completed_at?: timestamp, archived_at?: timestamp }
  - POST /api/todos/reorder — accept payload for ordering within a swimlane
- Existing endpoints that return lists SHALL accept a query param view=board that returns items grouped/ordered suitable for rendering the board (or the client may compute swimlanes from a flat list).

