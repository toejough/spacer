## REMOVED Requirements

### Requirement: Import data from a file
**Reason**: Replaced below — todos are now deduplicated by content (type + title) on import, the same as notes, instead of being appended unconditionally.
**Migration**: No user action needed; this only changes de-duplication behavior for future imports.

## ADDED Requirements

### Requirement: Import data from a file
The system SHALL allow the user to import a previously exported JSON file, without requiring login. Import SHALL always append data — it SHALL NOT remove or overwrite any existing local item. An item's content is defined as its type (todo or note) plus its title (including any cloze markup); everything else (done/archived state, priority, due date, tags, and all spaced-repetition metadata) is metadata, not content.

#### Scenario: Import skips items that duplicate existing content
- **WHEN** the user imports a file containing an item (todo or note) whose type and title exactly match an existing local item's type and title
- **THEN** the imported item is not added
- **AND** the existing local item is left completely unchanged, including all of its metadata (done/archived state, priority, tags, due date, next_review, ease_factor, interval_days, repetitions, cloze_data, last_reviewed, review_enabled, timestamps)

#### Scenario: Import adds items with new content
- **WHEN** the user imports a file containing an item (todo or note) whose type+title does not match any existing local item
- **THEN** the item is added to local storage as a new item with a freshly assigned id

#### Scenario: Import an invalid file
- **WHEN** the user selects a file that is not valid JSON or does not match the expected export shape
- **THEN** the system shows an error message
- **AND** no local data is modified

#### Scenario: Single import action
- **WHEN** the user opens the Help tab
- **THEN** there is a single "Import" action (no separate Merge/Replace choice)
