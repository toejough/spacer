## MODIFIED Requirements

### Requirement: Delete a note
The system SHALL allow the user to abandon a note (archive it from view).

#### Scenario: Abandon from edit modal
- **WHEN** the user opens a note's edit modal and clicks Abandon
- **THEN** the note is marked archived (removed from active lists)
- **AND** the modal closes
- **AND** the note list is refreshed

