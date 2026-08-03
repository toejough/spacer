## RENAMED Requirements

- FROM: `### Requirement: Delete a note`
- TO: `### Requirement: Abandon a note`

## ADDED Requirements

### Requirement: Permanently delete a note
The system SHALL allow the user to permanently delete a note that has been Abandoned. Deletion SHALL NOT be available for an active note. A deleted note is removed entirely and is not recoverable except via a prior manual export.

#### Scenario: Delete is not available for active notes
- **WHEN** a note is active (not archived)
- **THEN** no delete action is available for it, via swipe or the edit modal

#### Scenario: Swipe threshold commits delete
- **WHEN** the user swipes an Abandoned note card left past the delete threshold and releases
- **THEN** the note is removed from the list immediately
- **AND** an Undo toast appears

#### Scenario: Swipe reveals a delete affordance behind the card
- **WHEN** the user drags an Abandoned note card left
- **THEN** the space the card vacates shows a red background with a trash icon
- **AND** the trash icon is centered within that revealed space once the drag reaches the delete threshold

#### Scenario: An interrupted or vertical drag never deletes
- **WHEN** the user drags an Abandoned note card vertically, or a horizontal drag is interrupted (for example by the browser taking over for scrolling) before being released
- **THEN** the note is not deleted, regardless of how far any horizontal movement had already gone
- **AND** the card returns to its original position

#### Scenario: Undo restores the deleted note
- **WHEN** the user activates Undo from the toast after a swipe-delete
- **THEN** the note is restored to its prior state, including its Abandoned status
- **AND** the toast is dismissed

#### Scenario: Undo toast expires
- **WHEN** the Undo toast's time window elapses without the user activating Undo
- **THEN** the note remains permanently deleted
- **AND** the toast is dismissed

#### Scenario: Delete Forever is available in the edit modal for abandoned notes
- **WHEN** the user opens the edit modal for an Abandoned note
- **THEN** the modal shows a "Delete Forever" control
- **AND** the control is keyboard-focusable

#### Scenario: Delete Forever is not available in the edit modal for active notes
- **WHEN** the user opens the edit modal for an active note
- **THEN** the modal does not show a "Delete Forever" control

#### Scenario: Confirming Delete Forever in the edit modal deletes immediately
- **WHEN** the user clicks "Delete Forever" in the edit modal
- **THEN** the note is permanently deleted
- **AND** the modal closes
- **AND** the note list is refreshed
