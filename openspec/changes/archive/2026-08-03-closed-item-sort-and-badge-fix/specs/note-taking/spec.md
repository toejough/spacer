## MODIFIED Requirements

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

### Requirement: Delete a note
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
