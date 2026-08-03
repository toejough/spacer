## MODIFIED Requirements

### Requirement: Complete a todo
The system SHALL allow the user to mark a todo as complete.

#### Scenario: Toggle completion
- **WHEN** the user clicks the Complete button on a todo
- **THEN** the todo's done state toggles between complete and incomplete
- **AND** the todo list is refreshed
- **AND** finished todos are removed from the review list

### Requirement: Delete a todo
The system SHALL allow the user to abandon a todo (archive it from view).

#### Scenario: Delete from edit modal
- **WHEN** the user opens a todo's edit modal and clicks Abandon
- **THEN** the todo is marked archived (removed from active lists)
- **AND** the modal closes
- **AND** the todo list is refreshed

#### Scenario: Abandon is available in edit modal UI
- **WHEN** the user opens a todo's edit modal
- **THEN** the edit modal shows an Abandon control (button)
- **AND** the Abandon control is hidden or disabled when the todo is marked as Done
- **AND** confirming Abandon marks the todo archived, closes the modal, and refreshes lists

## REMOVED Requirements

### Requirement: Mobile parity for Done and Abandon (IMPLEMENTED)
**Reason**: Not a mobile-specific surface — this is core behavior of the one production todo list. Replaced by "Done and Abandoned are mutually exclusive end-of-life states," which drops the mobile/swipe-gesture framing and consolidates the mutual-exclusion scenario that was duplicated with "Delete a todo."
**Migration**: No behavior change; see the replacement requirement below.

## ADDED Requirements

### Requirement: Done and Abandoned are mutually exclusive end-of-life states
The system SHALL treat Done and Abandoned as mutually exclusive end-of-life states: an item in one closed state cannot move directly to the other — it MUST be Reopened first. Both Done and Abandon SHALL be equally accessible actions, and Reopen SHALL be available from either closed state.

#### Scenario: Direct transition between closed states is blocked
- **WHEN** a todo is Done
- **THEN** the system offers only Reopen, not Abandon, on that item
- **WHEN** a todo is Abandoned
- **THEN** the system offers only Reopen, not Done, on that item
- **AND** reaching the other closed state from either requires Reopen first

#### Scenario: Done and Abandoned are visually distinct
- **WHEN** a todo is Done or Abandoned
- **THEN** the two states use clearly different visual treatment (not label text alone), so a user can distinguish them at a glance

#### Scenario: Controls are accessible
- **WHEN** a user interacts with Done, Abandon, or Reopen controls
- **THEN** the controls are keyboard-focusable
- **AND** state changes are announced via ARIA live regions
