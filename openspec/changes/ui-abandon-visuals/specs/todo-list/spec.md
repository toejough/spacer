## MODIFIED Requirements

### Requirement: Delete a todo
The system SHALL allow the user to abandon a todo (archive it from view).

#### Scenario: Abandon UI uses neutral marker
- **WHEN** the Abandon control is rendered for a todo
- **THEN** it uses a stable, platform-independent icon (SVG) that does not evoke irreversible deletion
- **AND** the control includes an accessible label (aria-label) for screen readers
- **AND** keyboard focus styles are visible and meet contrast requirements

