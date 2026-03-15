# Bootstrap should deliver a vertical slice, not horizontal layers

**Status:** wont-fix
**Priority:** p1
**Labels:** tooling, issue-1-retro
**Created:** 2026-03-14
**Closed:** 2026-03-14

## Description
The bootstrap built every architectural layer (DB, stores, router, SM-2, PWA, components) but nothing was connected. The result: 16 passing tests and an app that does nothing. A better bootstrap would end with one working user flow.

## Acceptance Criteria
- [ ] 5m-increment skill's bootstrap section updated to require a working vertical slice
- [ ] Bootstrap AC includes: "a user can complete one full action end-to-end"
- [ ] Horizontal-only scaffolds (empty stubs, disconnected layers) are explicitly discouraged
- [ ] Decision captured as ADR in docs/decisions/

## Notes
The counter-argument is that horizontal layers unblock parallel work. But for a solo developer using 5-minute increments, vertical slices give faster feedback and catch integration issues immediately.
