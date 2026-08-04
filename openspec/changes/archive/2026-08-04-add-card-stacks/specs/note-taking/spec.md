## ADDED Requirements

### Requirement: Stacks appear in the Notes tab
The system SHALL render stacks containing at least one note member as collapsed stack tiles in the Notes tab, interleaved with unstacked note cards according to the tab's relevance ordering (see Note list ordering), and SHALL show only note members when a stack is expanded on this tab.

#### Scenario: Stack with note members appears on the Notes tab
- **WHEN** a stack has at least one note member
- **THEN** the stack's collapsed tile appears in the Notes tab, positioned per the tab's relevance ordering

#### Scenario: Stack with no note members is omitted from the Notes tab
- **WHEN** a stack has no note members (all members are todos)
- **THEN** the stack does not appear in the Notes tab

#### Scenario: Expanding a stack on the Notes tab shows only its note members
- **WHEN** the user expands a stack on the Notes tab
- **THEN** only the stack's note members are shown as individual cards
