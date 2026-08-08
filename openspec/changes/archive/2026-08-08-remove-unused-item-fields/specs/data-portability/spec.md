## MODIFIED Requirements

### Requirement: Import data from a file
The system SHALL allow the user to import a previously exported JSON file, without requiring login. Import SHALL always append data — it SHALL NOT remove or overwrite any existing local item or stack. An item's content is defined as its type (todo or note) plus its title (including any cloze markup); everything else (done/archived state, and all spaced-repetition metadata) is metadata, not content. A stack's identity for import purposes is its name (case-sensitive, exact match after trimming).

#### Scenario: Import skips items that duplicate existing content
- **WHEN** the user imports a file containing an item (todo or note) whose type and title exactly match an existing local item's type and title
- **THEN** the imported item is not added
- **AND** the existing local item is left completely unchanged, including all of its metadata (done/archived state, next_review, ease_factor, interval_days, repetitions, cloze_data, last_reviewed, review_enabled, stack_id, timestamps) — its stack membership is not altered by anything in the imported file

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
