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

### Requirement: Accessibility and keyboard support
The system SHALL provide keyboard-accessible alternatives to drag-and-drop.

#### Scenario: Move via keyboard
- **WHEN** the user focuses a card and triggers the move-to-Done keyboard shortcut or selects it from the item's menu
- **THEN** the system marks the item done and announces the change via ARIA live region

