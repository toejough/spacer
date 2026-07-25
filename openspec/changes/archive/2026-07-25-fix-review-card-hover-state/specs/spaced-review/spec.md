## MODIFIED Requirements

### Requirement: Review an item
The system SHALL present a due item and accept a quality rating.

#### Scenario: Show due item
- **WHEN** the user opens the Review tab
- **THEN** the earliest-due item is shown
- **AND** a rating scale from 0 to 5 is presented

#### Scenario: Submit rating
- **WHEN** the user submits a rating for an item
- **THEN** the SM-2 algorithm is applied to update the item's ease factor, interval, and repetitions
- **AND** the next review date is set based on the interval
- **AND** the next due item is shown
- **AND** the next due item does not display any rating option as preselected or visually selected from the previous card
