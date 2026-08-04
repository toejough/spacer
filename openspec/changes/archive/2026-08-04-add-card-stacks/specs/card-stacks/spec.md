## Purpose

Lets the user group related todos and notes together into a named stack, so related cards travel together in the UI instead of competing individually for position in the relevance sort.

## ADDED Requirements

### Requirement: Create a stack by dragging a card onto another card
The system SHALL allow the user to create a new stack by dragging a card (via its drag handle) and dropping it onto another card that is not already in a stack.

#### Scenario: Drop prompts for a stack name
- **WHEN** the user drags card A by its drag handle and drops it onto card B, and neither card is currently in a stack
- **THEN** the system prompts the user for a stack name before completing the merge

#### Scenario: Confirming the name creates the stack
- **WHEN** the user enters a name in the stack-name prompt and confirms
- **THEN** a new stack is created with that name
- **AND** card A and card B both become members of the new stack

#### Scenario: Canceling the prompt aborts the merge
- **WHEN** the user cancels the stack-name prompt
- **THEN** no stack is created
- **AND** card A and card B remain independent, unstacked cards

#### Scenario: Dropping a card onto a card already in a stack adds to that stack
- **WHEN** the user drags an unstacked card and drops it onto a card that is already a member of a stack
- **THEN** the dragged card is added to that existing stack
- **AND** no stack-name prompt is shown

### Requirement: Add a card to an existing stack by dragging onto the stack
The system SHALL allow the user to add a card to an existing stack by dragging the card and dropping it onto that stack's collapsed tile.

#### Scenario: Drop onto a collapsed stack tile
- **WHEN** the user drags a card by its drag handle and drops it onto a stack's collapsed tile
- **THEN** the card becomes a member of that stack
- **AND** no stack-name prompt is shown

### Requirement: Stacks are not merged with each other via drag-and-drop
The system SHALL NOT support dragging one stack onto another stack to merge them.

#### Scenario: Dragging a stack onto another stack has no merge effect
- **WHEN** the user attempts to drag one stack's collapsed tile onto another stack's collapsed tile
- **THEN** the stacks are not merged
- **AND** each stack retains its own membership unchanged

### Requirement: Remove a card from a stack by dragging it out
The system SHALL allow the user to remove a card from a stack by dragging it (via its drag handle) out of the stack and dropping it outside the stack.

#### Scenario: Drag out un-stacks the card
- **WHEN** the user drags a card that is a member of a stack, by its drag handle, and drops it outside that stack
- **THEN** the card is removed from the stack
- **AND** the card becomes an independent, unstacked card
- **AND** the remaining stack members are unaffected

### Requirement: Assign or change a card's stack via the edit modal
The system SHALL allow the user to view and change a card's stack membership from the card's edit modal, as an accessible alternative to drag-and-drop.

#### Scenario: Edit modal shows current stack
- **WHEN** the user opens the edit modal for a card that is a member of a stack
- **THEN** the modal displays the card's current stack

#### Scenario: Change stack from the edit modal
- **WHEN** the user selects a different stack (or no stack) in the edit modal and saves
- **THEN** the card's stack membership is updated accordingly
- **AND** the modal closes
- **AND** the relevant list is refreshed

### Requirement: Stacks require a name
The system SHALL require every stack to have a non-empty name at all times; a stack SHALL NOT exist without a name.

#### Scenario: New stack cannot be created without a name
- **WHEN** the user is prompted for a stack name during a drag-merge and provides no name
- **THEN** the stack is not created

### Requirement: Rename a stack
The system SHALL allow the user to rename a stack inline from its collapsed tile, and from the edit modal of any of its member cards.

#### Scenario: Inline rename from the collapsed tile
- **WHEN** the user activates the rename control on a stack's collapsed tile, enters a new name, and confirms
- **THEN** the stack's name is updated
- **AND** the updated name is reflected everywhere the stack appears

#### Scenario: Canceling an inline rename keeps the prior name
- **WHEN** the user activates the rename control and cancels without confirming
- **THEN** the stack's name is unchanged

### Requirement: Stacks render as a collapsed tile that can expand and collapse
The system SHALL render a stack as a single collapsed tile showing its name and the count of members visible on the current view, with a control to expand it and reveal its member cards, and to collapse it again.

#### Scenario: Collapsed tile shows visible-member count
- **WHEN** a stack is rendered in a tab or in search results
- **THEN** the collapsed tile shows the stack's name and the count of the stack's members that match the current view's type filter

#### Scenario: Expand reveals member cards
- **WHEN** the user expands a stack's collapsed tile
- **THEN** the stack's member cards (matching the current view's type filter) are shown, each individually actionable as an ordinary card

#### Scenario: Collapse hides member cards again
- **WHEN** the user collapses an expanded stack
- **THEN** the member cards are hidden and only the collapsed tile is shown

### Requirement: A stack's list position follows its most urgent visible member
The system SHALL position a stack in a relevance-sorted list at the position its most urgent member would occupy, considering only members that match the current view's type filter.

#### Scenario: Stack sorts by its most urgent visible member
- **WHEN** a stack has multiple members visible in the current view with differing relevance ranks
- **THEN** the stack's position in the sorted list matches the position of its highest-ranked (most urgent) visible member

#### Scenario: Stack with no visible members is omitted from a view
- **WHEN** none of a stack's members match the current view's type filter
- **THEN** the stack does not appear in that view

### Requirement: Mixed-type stacks are supported and shown in search
The system SHALL allow a stack to contain members of different item types (todos and notes together), and SHALL show such stacks, collapsed, in search results.

#### Scenario: Mixed-type stack membership
- **WHEN** the user adds a note to a stack that already contains a todo (or vice versa)
- **THEN** the stack contains both members regardless of their differing types

#### Scenario: Mixed-type stack appears in search
- **WHEN** a mixed-type stack has at least one member matching the search criteria
- **THEN** the stack appears in search results as a collapsed tile
