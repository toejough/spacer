# data-portability Specification

## Purpose
TBD - created by archiving change data-export-import. Update Purpose after archive.
## Requirements
### Requirement: Export all data as a file
The system SHALL allow the user to export all their local todo/note data as a downloadable JSON file, without requiring login.

#### Scenario: Export from Help tab
- **WHEN** the user clicks the "Export data" button in the Help tab
- **THEN** a JSON file is generated containing all items currently in local storage, a schema version, and an export timestamp
- **AND** the browser downloads the file with a name that includes the current date

#### Scenario: Export with no items
- **WHEN** the user exports data while no items exist
- **THEN** a valid JSON file is still downloaded, containing an empty items list

### Requirement: Import data from a file
The system SHALL allow the user to import a previously exported JSON file, without requiring login.

#### Scenario: Import in replace mode
- **WHEN** the user selects a valid export file and chooses "Replace"
- **THEN** the user is asked to confirm the destructive action
- **AND** upon confirmation, all current local items are discarded and replaced with the imported items
- **AND** the UI refreshes to reflect the imported data

#### Scenario: Import in merge mode
- **WHEN** the user selects a valid export file and chooses "Merge"
- **THEN** imported items whose id does not match an existing item are added to local storage
- **AND** imported items whose id matches an existing item are left unchanged (existing item wins)
- **AND** the UI refreshes to reflect the merged data

#### Scenario: Import an invalid file
- **WHEN** the user selects a file that is not valid JSON or does not match the expected export shape
- **THEN** the system shows an error message
- **AND** no local data is modified

### Requirement: Exported file is portable and versioned
The exported file SHALL be self-describing so it can be safely re-imported on a different host, browser, or app version.

#### Scenario: Schema version included
- **WHEN** data is exported
- **THEN** the file includes a `schema_version` field describing the export format

#### Scenario: Forward-compatible import
- **WHEN** an older-format export file is imported into a newer app version
- **THEN** the system upgrades the data to the current in-memory shape before loading it, without data loss

