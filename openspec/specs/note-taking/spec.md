# note-taking Specification

## Purpose

Manage notes in the Remember Everything app, including optional cloze deletions for review.
## Requirements
### Requirement: Add a note
The system SHALL allow the user to add a note from the Notes tab.

#### Scenario: Add a note by pressing Enter
- **WHEN** the user types text into the note input and presses Enter without Shift
- **THEN** a new note is created with the given title
- **AND** the input is cleared
- **AND** the note list is refreshed

#### Scenario: Add a note by clicking the Add button
- **WHEN** the user types text into the note input and clicks the Add button
- **THEN** a new note is created with the given title
- **AND** the input is cleared
- **AND** the note list is refreshed

#### Scenario: Empty input is rejected
- **WHEN** the user attempts to add a note with no text
- **THEN** no item is created

### Requirement: Edit a note
The system SHALL allow the user to edit a note's title and content.

#### Scenario: Open edit modal
- **WHEN** the user clicks a note in the list
- **THEN** an edit modal opens populated with the note's current fields

#### Scenario: Save changes
- **WHEN** the user edits the title or content and saves the modal
- **THEN** the note is updated
- **AND** the modal closes
- **AND** the note list is refreshed

### Requirement: Create cloze deletions
The system SHALL allow the user to mark portions of a note as cloze deletions.

#### Scenario: Create cloze from selection
- **WHEN** the user selects text in the title while editing a note and triggers the cloze command
- **THEN** the selected text is wrapped in `{{...}}` markers
- **AND** the note is treated as a cloze-based review item

#### Scenario: Multiple clozes in one note
- **WHEN** a note contains multiple `{{...}}` markers
- **THEN** each marker represents an independent cloze deletion
- **AND** each cloze is reviewed separately

### Requirement: Abandon a note
The system SHALL allow the user to abandon a note (archive it). An abandoned note remains visible in the note list, sorted last per the state tier of the note list ordering, and can be reopened. Reopening a note SHALL reset its review scheduling to fresh defaults, since any prior schedule is stale after the note was archived.

#### Scenario: Delete from edit modal
- **WHEN** the user opens a note's edit modal and clicks Abandon
- **THEN** the note is marked archived
- **AND** the modal closes
- **AND** the note list is refreshed

#### Scenario: Abandoned note remains visible
- **WHEN** a note is marked archived
- **THEN** the note still appears in the note list, sorted after all active notes
- **AND** the note is visually distinct from active notes
- **AND** the note's list entry shows relative time since it was abandoned (e.g. "Abandoned 3d ago") instead of review-due information

#### Scenario: Reopen an abandoned note
- **WHEN** the user clicks Reopen on an abandoned note
- **THEN** the note's archived state is cleared
- **AND** the note's review scheduling (ease factor, interval, repetitions, and next review date, or the equivalent for each cloze deletion) is reset to fresh defaults
- **AND** the note becomes due for review immediately
- **AND** the note list is refreshed

### Requirement: Permanently delete a note
The system SHALL allow the user to permanently delete a note that has been Abandoned. Deletion SHALL NOT be available for an active note. A deleted note is removed entirely and is not recoverable except via a prior manual export.

#### Scenario: Delete is not available for active notes
- **WHEN** a note is active (not archived)
- **THEN** no delete action is available for it, via swipe or the edit modal

#### Scenario: Swipe threshold commits delete
- **WHEN** the user swipes an Abandoned note card left past the delete threshold and releases
- **THEN** the note is removed from the list immediately
- **AND** an Undo toast appears

#### Scenario: Swipe reveals a delete affordance behind the card
- **WHEN** the user drags an Abandoned note card left
- **THEN** the space the card vacates shows a red background with a trash icon
- **AND** the trash icon is centered within that revealed space once the drag reaches the delete threshold

#### Scenario: An interrupted or vertical drag never deletes
- **WHEN** the user drags an Abandoned note card vertically, or a horizontal drag is interrupted (for example by the browser taking over for scrolling) before being released
- **THEN** the note is not deleted, regardless of how far any horizontal movement had already gone
- **AND** the card returns to its original position

#### Scenario: Undo restores the deleted note
- **WHEN** the user activates Undo from the toast after a swipe-delete
- **THEN** the note is restored to its prior state, including its Abandoned status
- **AND** the toast is dismissed

#### Scenario: Undo toast expires
- **WHEN** the Undo toast's time window elapses without the user activating Undo
- **THEN** the note remains permanently deleted
- **AND** the toast is dismissed

#### Scenario: Delete Forever is available in the edit modal for abandoned notes
- **WHEN** the user opens the edit modal for an Abandoned note
- **THEN** the modal shows a "Delete Forever" control
- **AND** the control is keyboard-focusable

#### Scenario: Delete Forever is not available in the edit modal for active notes
- **WHEN** the user opens the edit modal for an active note
- **THEN** the modal does not show a "Delete Forever" control

#### Scenario: Confirming Delete Forever in the edit modal deletes immediately
- **WHEN** the user clicks "Delete Forever" in the edit modal
- **THEN** the note is permanently deleted
- **AND** the modal closes
- **AND** the note list is refreshed

### Requirement: Note list ordering
The system SHALL order the note list by relevance: state first (active before abandoned), then, for active items only, review urgency (soonest due review first, with items that have review disabled or no scheduled review ranked after any item with a due date), then recency (most recently updated first) as the final tiebreaker. For abandoned items, review urgency does not apply — ordering within that tier is by recency alone.

#### Scenario: Active notes precede abandoned notes
- **WHEN** the note list contains active and abandoned notes
- **THEN** all active notes appear before any abandoned note

#### Scenario: Soonest review due sorts first within the active tier
- **WHEN** two active notes have different upcoming review due dates
- **THEN** the note due for review sooner appears first

#### Scenario: Items without a scheduled review sort after items with one, within the active tier
- **WHEN** an active note has reviews disabled or no scheduled review, and another active note has a due date
- **THEN** the item with a due date appears first

#### Scenario: Multi-cloze items use their earliest due date, within the active tier
- **WHEN** an active note has multiple cloze deletions with different `next_review` dates
- **THEN** the note's position is determined by the earliest of those dates

#### Scenario: Recency alone orders abandoned notes
- **WHEN** two notes are both abandoned
- **THEN** the more recently updated note appears first, regardless of any stale review scheduling data

#### Scenario: Recency breaks ties within the active tier's review-urgency ranking
- **WHEN** two active notes have the same review-urgency ranking
- **THEN** the more recently updated note appears first

### Requirement: Stacks appear in the Notes tab
The system SHALL render stacks containing at least one note member as collapsed stack tiles in the Notes tab, interleaved with unstacked note cards according to the tab's relevance ordering (see Note list ordering), and SHALL show only note members when a stack is expanded on this tab.

#### Scenario: Stack with note members appears on the Notes tab
- **WHEN** a stack has at least one note member
- **THEN** the stack's collapsed tile appears in the Notes tab, positioned per the tab's relevance ordering

#### Scenario: Stack with no note members is omitted from the Notes tab
- **WHEN** a stack has no note members (all members are todos)
- **THEN** the stack does not appear in the Notes tab

#### Scenario: Expanding a stack on the Notes tab shows only its note members
- **WHEN** the user expands a stack on the Notes tab
- **THEN** only the stack's note members are shown as individual cards

