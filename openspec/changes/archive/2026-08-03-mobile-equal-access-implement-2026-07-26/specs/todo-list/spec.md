## MODIFIED Requirements

### Requirement: Mobile parity for Done and Abandon (IMPLEMENTED)
The system SHALL make both Done and Abandon equally accessible on mobile while providing safeguards for Abandon, and SHALL provide a Reopen action from either closed state. Done and Abandoned are mutually exclusive end-of-life states: an item in one closed state cannot be moved directly to the other — it MUST be Reopened first.

#### Scenarios implemented
- Two visible actions with Reopen on closed cards.
- Split-swipe left/right is NOT part of the default implementation; swipe-based interactions are archived as experimental. The default implementation uses two visible buttons per card (Done and Abandon) and Reopen on closed cards.
- Keyboard and ARIA announcements implemented: controls are focusable and state changes use ARIA live regions.
- Visual treatments for each state provided and contrast verified.

#### Scenario: Done and Abandoned cannot transition directly to each other
- **WHEN** a todo is Done
- **THEN** the system offers only Reopen, not Abandon, on that card
- **WHEN** a todo is Abandoned
- **THEN** the system offers only Reopen, not Done, on that card
- **AND** reaching the other closed state from either requires Reopen first

#### Scenario: Done and Abandoned are visually distinct
- **WHEN** a todo is Done or Abandoned
- **THEN** the two states use clearly different visual treatment (not label text alone), so a user can distinguish them at a glance

### Requirement: Delete a todo
The system SHALL allow the user to abandon a todo (archive it from view).

#### Scenario: Delete from edit modal
- **WHEN** the user opens a todo's edit modal and clicks Abandon
- **THEN** the todo is marked archived (removed from active lists)
- **AND** the modal closes
- **AND** the todo list is refreshed

#### Scenario: Done and Abandoned block direct transition to each other
- **WHEN** the user attempts to Abandon a todo that is marked Done
- **THEN** the system prevents the abandon action
- **AND** the system displays an informative message to the user
- **AND** the todo remains unchanged
- **WHEN** the user attempts to mark a todo Done while it is Abandoned
- **THEN** the system prevents the action
- **AND** the system displays an informative message to the user
- **AND** the todo remains unchanged
- **AND** in both cases, the todo must be Reopened before the other closed state can be set

#### Scenario: Abandon is available in edit modal UI
- **WHEN** the user opens a todo's edit modal
- **THEN** the edit modal shows an Abandon control (button)
- **AND** the Abandon control is hidden or disabled when the todo is marked as Done
- **AND** confirming Abandon marks the todo archived, closes the modal, and refreshes lists
