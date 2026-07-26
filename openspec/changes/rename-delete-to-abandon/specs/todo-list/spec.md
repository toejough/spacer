## CHANGED Requirements

### Requirement: Abandon a todo
The system SHALL allow the user to abandon a todo (archive it from view).

#### Scenario: Abandon from edit modal
- **WHEN** the user opens a todo's edit modal and clicks Abandon
- **THEN** the todo is marked archived (removed from active lists)
- **AND** the modal closes
- **AND** the todo list is refreshed

