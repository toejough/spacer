# Board view specification (delta)

This spec describes the board UI with columns and time-based swimlanes and the behaviors required.

## ADDED Requirements

### Requirement: Board view
The system SHALL present a Board view accessible from the Todos tab or navigation.

#### Scenario: Open the Board view
- **WHEN** the user opens the Todos tab and selects the Board view
- **THEN** the Board view displays three columns labelled Todo, Done, and Abandoned
- **AND** each column shows cards for the appropriate todo items

### Requirement: Swimlanes
The system SHALL partition columns into swimlanes based on closed/archived timestamps.

#### Scenario: Swimlane bucketing
- **WHEN** the Board view is rendered
- **THEN** each todo is placed into one of the swimlanes: open/closed-in-last-week, closed-in-last-month, closed-in-last-quarter, closed-in-last-year, everything-else
- **AND** swimlane calculation uses the latest of completed_at and archived_at where applicable

### Requirement: Drag-and-drop
The system SHALL support dragging cards between columns and swimlanes.

#### Scenario: Drag todo to Done
- **WHEN** the user drags a card from Todo and drops it into a Done swimlane
- **THEN** the system sets completed_at to the current timestamp
- **AND** the card appears in the correct Done swimlane

#### Scenario: Drag todo to Abandoned
- **WHEN** the user drags a card from Todo and drops it into an Abandoned swimlane
- **THEN** the system sets archived_at to the current timestamp
- **AND** the card is removed from active Todo lists

#### Scenario: Reorder within swimlane
- **WHEN** the user reorders cards inside a swimlane
- **THEN** the new order is persisted and reflected on subsequent loads

### Requirement: Accessibility and keyboard support
The system SHALL provide keyboard-accessible alternatives to drag-and-drop.

#### Scenario: Move via keyboard
- **WHEN** the user focuses a card and triggers the move-to-Done keyboard shortcut or selects it from the item's menu
- **THEN** the system marks the item done and announces the change via ARIA live region

### Requirement: Mobile and responsive layout
The Board view SHALL render appropriately on narrow viewports.

#### Scenario: Mobile fallback
- **WHEN** viewport width is small
- **THEN** the Board view collapses into a stacked list grouped by swimlane and column
- **AND** status changes are available via controls rather than drag-and-drop

### Requirement: API
The system SHALL provide APIs for status updates and reordering.

#### Scenario: Update status via API
- **WHEN** the client sends PATCH /api/todos/:id with {status: "done"} or {status: "archived"}
- **THEN** the server updates completed_at or archived_at timestamps and returns the updated todo

#### Scenario: Persist reorder
- **WHEN** the client POSTs to /api/todos/reorder with a swimlane id and ordered list of todo ids
- **THEN** the server stores ordering metadata for the swimlane and responds with success
