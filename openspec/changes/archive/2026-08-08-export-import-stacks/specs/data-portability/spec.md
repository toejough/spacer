## MODIFIED Requirements

### Requirement: Export all data as a file
The system SHALL allow the user to export all their local todo/note data as a downloadable JSON file, without requiring login. Exported data SHALL include stacks (id, name, and timestamps) alongside items, so a backup captures the complete local dataset, not items alone.

#### Scenario: Export from Help tab
- **WHEN** the user clicks the "Export data" button in the Help tab
- **THEN** a JSON file is generated containing all items currently in local storage, all stacks currently in local storage, a schema version, and an export timestamp
- **AND** the browser downloads the file with a name that includes the current date

#### Scenario: Export with no items
- **WHEN** the user exports data while no items exist
- **THEN** a valid JSON file is still downloaded, containing an empty items list

#### Scenario: Export with no stacks
- **WHEN** the user exports data while no stacks exist
- **THEN** a valid JSON file is still downloaded, containing an empty stacks list

#### Scenario: Export includes stack membership
- **WHEN** the user exports data while at least one item is a member of a stack
- **THEN** the exported item retains its stack membership
- **AND** the exported stacks list includes that stack with its name

### Requirement: Import data from a file
The system SHALL allow the user to import a previously exported JSON file, without requiring login. Import SHALL always append data — it SHALL NOT remove or overwrite any existing local item or stack. An item's content is defined as its type (todo or note) plus its title (including any cloze markup); everything else (done/archived state, priority, due date, tags, and all spaced-repetition metadata) is metadata, not content. A stack's identity for import purposes is its name (case-sensitive, exact match after trimming).

#### Scenario: Import skips items that duplicate existing content
- **WHEN** the user imports a file containing an item (todo or note) whose type and title exactly match an existing local item's type and title
- **THEN** the imported item is not added
- **AND** the existing local item is left completely unchanged, including all of its metadata (done/archived state, priority, tags, due date, next_review, ease_factor, interval_days, repetitions, cloze_data, last_reviewed, review_enabled, stack_id, timestamps) — its stack membership is not altered by anything in the imported file

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

#### Scenario: Import restores a stack that has no local counterpart
- **WHEN** the user imports a file containing a stack whose name does not match any existing local stack's name
- **THEN** a new local stack is created with that name and a freshly assigned id
- **AND** every imported item that belonged to that stack in the file is linked to the new local stack's id

#### Scenario: Import merges into an existing stack with the same name
- **WHEN** the user imports a file containing a stack whose name exactly matches an existing local stack's name
- **THEN** no new stack is created
- **AND** every imported item that belonged to that stack in the file is linked to the existing local stack's id
- **AND** the existing local stack's name and timestamps are left unchanged

#### Scenario: Imported item's stack membership is preserved without id collisions
- **WHEN** an imported item belonged to a stack in the export file
- **THEN** after import, the item's stack reference resolves to a real local stack (either newly created or merged), never to a stale id from the export file
- **AND** items on the local device that already used the same numeric stack id the file happened to use are unaffected — their stack membership does not change or merge as a side effect of the import

### Requirement: Exported file is portable and versioned
The exported file SHALL be self-describing so it can be safely re-imported on a different host, browser, or app version.

#### Scenario: Schema version included
- **WHEN** data is exported
- **THEN** the file includes a `schema_version` field describing the export format

#### Scenario: Forward-compatible import
- **WHEN** an older-format export file is imported into a newer app version
- **THEN** the system upgrades the data to the current in-memory shape before loading it, without data loss

#### Scenario: Import of a pre-stacks export file
- **WHEN** the user imports a file exported before stacks were included in the export format (no `stacks` field present)
- **THEN** the import succeeds, treating the file as containing zero stacks
- **AND** all items in the file are imported per the normal item-import rules

## ADDED Requirements

### Requirement: Import reports a result summary
The system SHALL report a summary of what happened immediately after an import completes, so the user can tell whether the restore was complete without inspecting the file or local storage directly.

#### Scenario: Summary after a successful import
- **WHEN** an import completes without error
- **THEN** the user is shown how many items were added, how many items were skipped as duplicates, how many stacks were added, and how many stacks were merged into existing ones

#### Scenario: Summary reflects an import that added nothing new
- **WHEN** every item and stack in the imported file already matches existing local content
- **THEN** the summary reports zero items added and zero stacks added, rather than showing no feedback at all
