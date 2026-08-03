# todo-list Specification

## Purpose

Manage todo items in the Remember Everything app.
## Requirements
### Requirement: Add a todo
The system SHALL allow the user to add a todo from the Todos tab.

#### Scenario: Add a todo by pressing Enter
- **WHEN** the user types text into the todo input and presses Enter
- **THEN** a new todo item is created with the given title
- **AND** the input is cleared
- **AND** the todo list is refreshed

#### Scenario: Add a todo by clicking the Add button
- **WHEN** the user types text into the todo input and clicks the Add button
- **THEN** a new todo item is created with the given title
- **AND** the input is cleared
- **AND** the todo list is refreshed

#### Scenario: Empty input is rejected
- **WHEN** the user attempts to add a todo with no text
- **THEN** no item is created

### Requirement: Complete a todo
The system SHALL allow the user to mark a todo as complete.

#### Scenario: Toggle completion
- **WHEN** the user clicks the Complete button on a todo
- **THEN** the todo's done state toggles between complete and incomplete
- **AND** the todo list is refreshed
- **AND** finished todos are removed from the review list

### Requirement: Edit a todo
The system SHALL allow the user to edit a todo's title.

#### Scenario: Open edit modal
- **WHEN** the user clicks a todo in the list
- **THEN** an edit modal opens populated with the todo's current fields

#### Scenario: Save changes
- **WHEN** the user edits the title and saves the modal
- **THEN** the todo is updated
- **AND** the modal closes
- **AND** the todo list is refreshed

### Requirement: Delete a todo
The system SHALL allow the user to abandon a todo (archive it from view).

#### Scenario: Delete from edit modal
- **WHEN** the user opens a todo's edit modal and clicks Abandon
- **THEN** the todo is marked archived (removed from active lists)
- **AND** the modal closes
- **AND** the todo list is refreshed

#### Scenario: Abandon is available in edit modal UI
- **WHEN** the user opens a todo's edit modal
- **THEN** the edit modal shows an Abandon control (button)
- **AND** the Abandon control is hidden or disabled when the todo is marked as Done
- **AND** confirming Abandon marks the todo archived, closes the modal, and refreshes lists

### Requirement: Done and Abandoned are mutually exclusive end-of-life states
The system SHALL treat Done and Abandoned as mutually exclusive end-of-life states: an item in one closed state cannot move directly to the other — it MUST be Reopened first. Both Done and Abandon SHALL be equally accessible actions, and Reopen SHALL be available from either closed state.

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

