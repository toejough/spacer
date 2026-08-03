## MODIFIED Requirements

### Requirement: Todo list ordering
The system SHALL order the todo list by relevance: state first (open, then done, then abandoned), then, for open items only, review urgency (soonest due review first, with items that have review disabled or no scheduled review ranked after any item with a due date), then recency (most recently updated first) as the final tiebreaker. For done and abandoned items, review urgency does not apply — ordering within those tiers is by recency alone.

#### Scenario: Open todos precede done and abandoned
- **WHEN** the todo list contains open, done, and abandoned items
- **THEN** all open items appear before any done item, and all done items appear before any abandoned item

#### Scenario: Soonest review due sorts first within the open tier
- **WHEN** two open todos have different upcoming review due dates
- **THEN** the todo due for review sooner appears first

#### Scenario: Items without a scheduled review sort after items with one, within the open tier
- **WHEN** an open todo has reviews disabled or no scheduled review, and another open todo has a due date
- **THEN** the item with a due date appears first

#### Scenario: Multi-cloze items use their earliest due date, within the open tier
- **WHEN** an open todo has multiple cloze deletions with different `next_review` dates
- **THEN** the todo's position is determined by the earliest of those dates

#### Scenario: Recency alone orders done and abandoned items
- **WHEN** two todos are both done, or both abandoned
- **THEN** the more recently updated todo appears first, regardless of any stale review scheduling data

#### Scenario: Recency breaks ties within the open tier's review-urgency ranking
- **WHEN** two open todos have the same review-urgency ranking
- **THEN** the more recently updated todo appears first

### Requirement: Done and Abandoned are mutually exclusive end-of-life states
The system SHALL treat Done and Abandoned as mutually exclusive end-of-life states: an item in one closed state cannot move directly to the other — it MUST be Reopened first. Both Done and Abandon SHALL be equally accessible actions, and Reopen SHALL be available from either closed state. Reopening an item SHALL reset its review scheduling to fresh defaults, since any prior schedule is stale after the item was closed.

#### Scenario: Direct transition between closed states is blocked
- **WHEN** a todo is Done
- **THEN** the system offers only Reopen, not Abandon, on that item
- **WHEN** a todo is Abandoned
- **THEN** the system offers only Reopen, not Done, on that item
- **AND** reaching the other closed state from either requires Reopen first

#### Scenario: Done and Abandoned are visually distinct
- **WHEN** a todo is Done or Abandoned
- **THEN** the two states use clearly different visual treatment (not label text alone), so a user can distinguish them at a glance

#### Scenario: Controls are accessible
- **WHEN** a user interacts with Done, Abandon, or Reopen controls
- **THEN** the controls are keyboard-focusable
- **AND** state changes are announced via ARIA live regions

#### Scenario: Reopen resets review scheduling
- **WHEN** the user reopens a Done or Abandoned todo
- **THEN** the todo's review scheduling (ease factor, interval, repetitions, and next review date, or the equivalent for each cloze deletion) is reset to fresh defaults
- **AND** the todo becomes due for review immediately

#### Scenario: Closed items show relative time since the state change, not review info
- **WHEN** a todo is Done or Abandoned
- **THEN** its list entry shows relative time since it was marked Done or Abandoned (e.g. "Completed 3d ago", "Abandoned yesterday") instead of review-due information
