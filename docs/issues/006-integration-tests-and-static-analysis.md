# Add integration tests and static analysis

**Status:** wont-fix
**Priority:** p0
**Labels:** tooling, issue-1-retro
**Created:** 2026-03-14
**Closed:** 2026-03-14

## Description
The bootstrap has unit tests for isolated pieces (SM-2, DB schema, router config) but nothing that verifies the app actually works as a whole. A user-facing smoke test and static analysis would have caught that nothing was wired up.

## Acceptance Criteria
- [ ] Integration test: mount Home view with real (fake-indexeddb) DB, verify deck list renders
- [ ] Integration test: create a deck, add a card, verify it appears
- [ ] Integration test: start a review, rate a card, verify SM-2 state updates
- [ ] TypeScript strict mode enabled (`strict: true` in tsconfig) — verify no errors
- [ ] ESLint configured with Vue + TypeScript rules
- [ ] `targ check` (or equivalent) runs type-check + lint + tests in one command
- [ ] CI-ready: all checks pass in a single command

## Notes
The unit tests gave false confidence — everything "passed" but the app did nothing. Integration tests that exercise the real data path (view → store → Dexie) prevent this. Static analysis catches unused imports, untyped boundaries, and dead code.

Closing — superseded by rebuild. Lessons incorporated into updated process.
