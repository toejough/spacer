# search Specification

## Purpose

Find todos and notes by text.

## ADDED Requirements

### Requirement: Search items
The system SHALL allow the user to search across todos and notes.

#### Scenario: Search by title
- **WHEN** the user types text into the search field and submits
- **THEN** items whose title contains the search text are displayed

#### Scenario: Search by content
- **WHEN** the user types text into the search field and submits
- **THEN** items whose content contains the search text are displayed

#### Scenario: Empty search
- **WHEN** the search field is empty and the user submits
- **THEN** no results are shown or all items are shown
