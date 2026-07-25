## MODIFIED Requirements

### Requirement: Identify due items
The system SHALL identify items that are due for review.

#### Scenario: Review badge shows count
- **WHEN** items are due for review
- **THEN** the Review tab badge displays the number of due items

#### Scenario: No due items
- **WHEN** no items are due for review
- **THEN** the Review tab displays an empty state message

#### Scenario: Finished todos are excluded
- **WHEN** a todo item is marked as done
- **THEN** it is not counted as due for review
- **AND** it does not appear in the Review tab
