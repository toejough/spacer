## MODIFIED Requirements

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

