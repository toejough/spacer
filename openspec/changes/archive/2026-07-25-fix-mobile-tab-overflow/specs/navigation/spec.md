## ADDED Requirements

### Requirement: Reachable navigation on small screens
The system SHALL keep all primary navigation entries reachable without hidden horizontal overflow on mobile-width screens.

#### Scenario: Help is reachable on mobile
- **WHEN** the app is viewed on a narrow viewport (e.g. 390px)
- **THEN** the Help entry is visible without horizontally scrolling the tab bar
- **AND** tapping it opens the Help view

#### Scenario: Content tabs fit on mobile
- **WHEN** the app is viewed on a narrow viewport
- **THEN** the content tabs (Review, Todos, Notes, Search) are all visible within the tab bar width
- **AND** no primary tab is scrolled off-screen

### Requirement: Help entry active state
The system SHALL indicate when the Help view is active.

#### Scenario: Help button shows active state
- **WHEN** the user opens the Help view
- **THEN** the help button reflects an active state
- **AND** no content tab shows an active state

#### Scenario: Content tab clears help active state
- **WHEN** the user selects a content tab while Help is open
- **THEN** the help button active state is cleared
- **AND** the selected content tab shows an active state
