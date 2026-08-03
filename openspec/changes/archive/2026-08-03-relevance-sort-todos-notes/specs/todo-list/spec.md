## ADDED Requirements

### Requirement: Todo list ordering
The system SHALL order the todo list by relevance: state first (open, then done, then abandoned), then review urgency (soonest due review first, with items that have review disabled or no scheduled review ranked after any item with a due date), then recency (most recently updated first) as the final tiebreaker.

#### Scenario: Open todos precede done and abandoned
- **WHEN** the todo list contains open, done, and abandoned items
- **THEN** all open items appear before any done item, and all done items appear before any abandoned item

#### Scenario: Soonest review due sorts first within a state tier
- **WHEN** two open todos have different upcoming review due dates
- **THEN** the todo due for review sooner appears first

#### Scenario: Items without a scheduled review sort after items with one
- **WHEN** a todo has reviews disabled or no scheduled review, and another todo in the same state tier has a due date
- **THEN** the item with a due date appears first

#### Scenario: Multi-cloze items use their earliest due date
- **WHEN** a todo has multiple cloze deletions with different `next_review` dates
- **THEN** the todo's position is determined by the earliest of those dates

#### Scenario: Recency breaks ties within a review-urgency tier
- **WHEN** two todos have the same review-urgency ranking
- **THEN** the more recently updated todo appears first
