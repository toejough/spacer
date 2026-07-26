# Column-based Todo Board with Swimlanes and Drag-and-Drop

## Why

Users have requested a more visual, status-oriented view of their todos: a board with columns for Todo / Done / Abandoned and time-based swimlanes that help them see recently closed items at a glance. This improves discoverability of recent work, makes triage easier, and gives a compact overview for review and cleanup.

## What changes

- Add a new Board view (tab) that presents todos in three columns: Todo, Done, Abandoned.
- Each column is divided into swimlanes (rows) based on when an item was last closed:
  - Open / closed in last week
  - Closed in last month
  - Closed in last quarter
  - Closed in last year
  - Everything else
- Allow drag-and-drop from the Todo column into Done or Abandoned swimlanes. Dropping into Done marks item as complete with a closed timestamp; dropping into Abandoned archives the item (Abandon action) and records an archived timestamp.
- Support reordering inside a swimlane and moving between swimlanes when timestamps change.
- Provide keyboard and screen-reader accessible alternatives to drag-and-drop (move actions from context menu / keyboard shortcuts).
- Provide a responsive layout that falls back to a stacked list on narrow screens.

## Impact

- UI: a new Board view and associated components, DnD interactions, responsive CSS, and accessibility support.
- API: endpoints to update status (done/archived) with timestamps and to reorder items within a swimlane/column.
- Data: rely on existing timestamps (completed_at, archived_at); ensure queries can compute swimlane buckets efficiently.
- Tests: unit & e2e tests for DnD, keyboard flows, and swimlane bucketing.

## Rollout

- Feature flag the Board view (disabled by default) for gradual rollout.
- Start with client-side implementation that consumes existing APIs; add API optimizations if necessary.

