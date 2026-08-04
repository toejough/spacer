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

### Requirement: Abandon a todo
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

### Requirement: Permanently delete a todo
The system SHALL allow the user to permanently delete a Done or Abandoned todo. Deletion SHALL NOT be available for an Open todo. A deleted todo is removed entirely and is not recoverable except via a prior manual export.

#### Scenario: Delete is not available for open todos
- **WHEN** a todo is Open
- **THEN** no delete action is available for it, via swipe or the edit modal

#### Scenario: Swipe threshold commits delete
- **WHEN** the user swipes a Done or Abandoned todo card left past the delete threshold and releases
- **THEN** the todo is removed from the list immediately
- **AND** an Undo toast appears

#### Scenario: Swipe reveals a delete affordance behind the card
- **WHEN** the user drags a Done or Abandoned todo card left
- **THEN** the space the card vacates shows a red background with a trash icon
- **AND** the trash icon is centered within that revealed space once the drag reaches the delete threshold

#### Scenario: An interrupted or vertical drag never deletes
- **WHEN** the user drags a Done or Abandoned todo card vertically, or a horizontal drag is interrupted (for example by the browser taking over for scrolling) before being released
- **THEN** the todo is not deleted, regardless of how far any horizontal movement had already gone
- **AND** the card returns to its original position

#### Scenario: Undo restores the deleted todo
- **WHEN** the user activates Undo from the toast after a swipe-delete
- **THEN** the todo is restored to its prior state, including its Done/Abandoned status
- **AND** the toast is dismissed

#### Scenario: Undo toast expires
- **WHEN** the Undo toast's time window elapses without the user activating Undo
- **THEN** the todo remains permanently deleted
- **AND** the toast is dismissed

#### Scenario: Delete Forever is available in the edit modal for closed todos
- **WHEN** the user opens the edit modal for a Done or Abandoned todo
- **THEN** the modal shows a "Delete Forever" control
- **AND** the control is keyboard-focusable

#### Scenario: Delete Forever is not available in the edit modal for open todos
- **WHEN** the user opens the edit modal for an Open todo
- **THEN** the modal does not show a "Delete Forever" control

#### Scenario: Confirming Delete Forever in the edit modal deletes immediately
- **WHEN** the user clicks "Delete Forever" in the edit modal
- **THEN** the todo is permanently deleted
- **AND** the modal closes
- **AND** the todo list is refreshed

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

