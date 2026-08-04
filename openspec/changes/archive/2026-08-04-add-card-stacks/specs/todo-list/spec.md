## ADDED Requirements

### Requirement: Stacks appear in the Todos tab
The system SHALL render stacks containing at least one todo member as collapsed stack tiles in the Todos tab, interleaved with unstacked todo cards according to the tab's relevance ordering (see Todo list ordering), and SHALL show only todo members when a stack is expanded on this tab.

#### Scenario: Stack with todo members appears on the Todos tab
- **WHEN** a stack has at least one todo member
- **THEN** the stack's collapsed tile appears in the Todos tab, positioned per the tab's relevance ordering

#### Scenario: Stack with no todo members is omitted from the Todos tab
- **WHEN** a stack has no todo members (all members are notes)
- **THEN** the stack does not appear in the Todos tab

#### Scenario: Expanding a stack on the Todos tab shows only its todo members
- **WHEN** the user expands a stack on the Todos tab
- **THEN** only the stack's todo members are shown as individual cards
