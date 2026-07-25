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

### Requirement: Delete a note
The system SHALL allow the user to delete a note.

#### Scenario: Delete from edit modal
- **WHEN** the user opens a note's edit modal and clicks Delete
- **THEN** the note is removed
- **AND** the modal closes
- **AND** the note list is refreshed

