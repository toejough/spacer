# Wire up core flow: create deck, add cards, review

**Status:** wont-fix
**Priority:** p0
**Labels:** feature, issue-1-retro
**Created:** 2026-03-14
**Closed:** 2026-03-14

## Description
The bootstrap created all the pieces (DB, stores, SM-2, views, components) but nothing is connected. A user visiting the app can't do anything — no deck creation, no card adding, no reviewing. Wire the pieces together into a usable end-to-end flow.

## Acceptance Criteria
- [ ] Home view lists decks from Dexie, shows due card count per deck
- [ ] User can create a new deck from Home view
- [ ] DeckView shows cards in the selected deck
- [ ] User can add cards (front/back) from DeckView via CardEditor
- [ ] User can start a review session (cards where nextReview ≤ now)
- [ ] ReviewView shows cards one at a time, flip to reveal, rate 1-5
- [ ] Rating updates SM-2 state and persists to Dexie
- [ ] Review session shows progress and summary at end
- [ ] Navigation between views works (home → deck → review → home)

## Notes
This is the critical path — without this, the app is a static page. Consider breaking into sub-issues per screen if any single piece exceeds 5 minutes.

Closing — superseded by rebuild. Lessons incorporated into updated process.
