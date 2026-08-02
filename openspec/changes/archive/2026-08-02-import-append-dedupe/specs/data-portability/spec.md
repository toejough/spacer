## REMOVED Requirements

### Requirement: Import data from a file
**Reason**: Replaced by a simpler, safer import model below — import is always append-only (no destructive Replace mode, no id-based Merge), with note-content-based dedup that keeps existing local review metadata.
**Migration**: Users who previously relied on Replace to fully restore onto a clean browser can clear site data/local storage first (or use a new browser profile), then use Import, which will append into an empty store — producing the same end result without a destructive code path.

## ADDED Requirements

### Requirement: Import data from a file
The system SHALL allow the user to import a previously exported JSON file, without requiring login. Import SHALL always append data — it SHALL NOT remove or overwrite any existing local item.

#### Scenario: Import appends todos unconditionally
- **WHEN** the user imports a file containing todo items
- **THEN** every imported todo is added to local storage as a new item with a freshly assigned id
- **AND** existing todos are left completely unchanged, even if an imported todo has the same title as an existing todo

#### Scenario: Import skips notes that duplicate existing note content
- **WHEN** the user imports a file containing a note whose content exactly matches the content of an existing local note
- **THEN** the imported note is not added
- **AND** the existing local note is left completely unchanged, including its review metadata (next_review, ease_factor, interval_days, repetitions, cloze_data, last_reviewed, review_enabled)

#### Scenario: Import adds notes with new content
- **WHEN** the user imports a file containing a note whose content does not match any existing local note
- **THEN** the note is added to local storage as a new item with a freshly assigned id

#### Scenario: Import an invalid file
- **WHEN** the user selects a file that is not valid JSON or does not match the expected export shape
- **THEN** the system shows an error message
- **AND** no local data is modified

#### Scenario: Single import action
- **WHEN** the user opens the Help tab
- **THEN** there is a single "Import" action (no separate Merge/Replace choice)
