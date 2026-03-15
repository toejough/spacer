# Add cards to a deck

**Status:** wont-fix
**Priority:** p1
**Labels:** feature, issue-1-retro
**Created:** 2026-03-14
**Closed:** 2026-03-14

## Description
Users need to be able to create flashcards within a deck. The CardEditor component should let users enter a front (prompt) and back (answer) and save the card to IndexedDB via Dexie.

## Acceptance Criteria
- [ ] CardEditor form with front and back text inputs
- [ ] Save button persists card to Dexie with correct deckId
- [ ] New card appears in the deck's card list
- [ ] Card gets default SM-2 state (easeFactor 2.5, interval 0, repetitions 0, nextReview now)
- [ ] Validation: both front and back must be non-empty

## Notes
Depends on DeckView being wired up to show cards. May need to create a deck first — consider whether deck creation is a prerequisite issue.

Closing — superseded by rebuild. Lessons incorporated into updated process.
