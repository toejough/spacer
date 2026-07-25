# spaced-review Specification

## Purpose

Review todos and notes using the SM-2 spaced repetition algorithm.
## Requirements
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

### Requirement: Cloze review
The system SHALL review each cloze deletion independently.

#### Scenario: One cloze per day per note
- **WHEN** a note has multiple clozes due
- **THEN** only the earliest-due cloze is shown in the Review tab
- **AND** the user is informed which cloze is being reviewed

#### Scenario: Reveal cloze answer
- **WHEN** the user is reviewing a cloze
- **THEN** the answer is hidden initially
- **AND** the user can click Reveal to see the answer before rating

### Requirement: SM-2 ratings
The system SHALL use the SM-2 algorithm to schedule reviews.

#### Scenario: Failed review (rating < 3)
- **WHEN** the user rates an item 0, 1, or 2
- **THEN** repetitions reset to 0
- **AND** the interval resets to 0
- **AND** the ease factor is unchanged

#### Scenario: Successful review (rating >= 3)
- **WHEN** the user rates an item 3, 4, or 5
- **THEN** repetitions increase by 1
- **AND** the interval is set according to the SM-2 formula
- **AND** the ease factor is updated with the SM-2 formula, clamped to a minimum of 1.3

