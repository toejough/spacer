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
- **WHEN** the user clicks the checkbox on a todo
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
The system SHALL allow the user to delete a todo.

#### Scenario: Delete from edit modal
- **WHEN** the user opens a todo's edit modal and clicks Delete
- **THEN** the todo is removed
- **AND** the modal closes
- **AND** the todo list is refreshed

