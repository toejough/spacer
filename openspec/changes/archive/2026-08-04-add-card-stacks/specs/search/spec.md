## ADDED Requirements

### Requirement: Stacks appear in search results
The system SHALL include a stack in search results, as a collapsed tile, when at least one of its members (of any item type) matches the search criteria.

#### Scenario: Mixed-type stack matches via one member
- **WHEN** a stack contains a todo and a note, and only the note matches the search text
- **THEN** the stack appears in search results as a collapsed tile

#### Scenario: Expanding a search result stack shows all matching-view members
- **WHEN** the user expands a stack shown in search results
- **THEN** all of the stack's members (regardless of item type) are shown as individual cards
