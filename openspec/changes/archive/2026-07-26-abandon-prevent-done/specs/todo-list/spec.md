## MODIFIED Requirements

### Requirement: Delete a todo
The system SHALL allow the user to abandon a todo (archive it from view).

#### Scenario: Delete from edit modal
- **WHEN** the user opens a todo's edit modal and clicks Abandon
- **THEN** the todo is marked archived (removed from active lists)
- **AND** the modal closes
- **AND** the todo list is refreshed

#### Scenario: Abandon is not allowed for completed todos
- **WHEN** the user attempts to Abandon a todo that is marked as Done
- **THEN** the system prevents the abandon action
- **AND** the system displays an informative message to the user
- **AND** the todo remains unchanged

