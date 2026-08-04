## ADDED Requirements

### Requirement: Member cards within an expanded stack are relevance-sorted
The system SHALL order a stack's member cards, when expanded, using the same relevance rules applied to top-level cards (state, then review urgency, then recency).

#### Scenario: Expanded stack members follow relevance order
- **WHEN** the user expands a stack containing members in differing states, review urgencies, or update recency
- **THEN** the member cards are displayed in the same relative order they would appear in if they were unstacked top-level cards
