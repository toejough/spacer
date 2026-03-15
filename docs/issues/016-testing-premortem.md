# 016 — Testing Premortem

**Status:** open
**Type:** quality / prevention

## Context

Spacer has one test file today: `src/__tests__/flow.test.ts` — an integration test that exercises the create-deck → add-card → review → SM-2-update flow directly against Dexie (via fake-indexeddb). There are no unit tests, no component tests, and no E2E tests. The test infrastructure is minimal: Vitest + happy-dom + fake-indexeddb, with a `createTestDB()` helper for DB isolation.

## The Exercise

Perform a premortem on the test suite. Assume we've shipped 10+ features and the tests are now a liability — slow, flaky, hard to understand, full of gaps, and failing for reasons unrelated to actual bugs. Developers skip running them locally and CI is routinely red. Regressions ship regularly.

**Your job:** Figure out what led to that state.

### How to run the premortem

1. **Read the current test infrastructure** — the test file, vitest config, test-setup, and the production code being tested. Understand what's covered, what's not, how tests are structured, what patterns are established (or missing), and how DB isolation works.
2. **Imagine 10+ features added** — each presumably adding tests. Consider: stats, tags, search, import/export, sync, settings, media cards, cram mode, multiple review modes, offline behavior, etc.
3. **Identify 3-5 specific testing weaknesses** that would compound badly under that growth. Focus on things that are fine with 1 test file but would rot at scale. Be concrete — reference actual test patterns, setup code, missing coverage categories, and structural decisions.
4. **For each weakness**, describe: what specifically goes wrong as the suite grows, why the current approach enables it, and a concrete mitigation (with a recommendation on whether to adopt now or defer with a specific trigger).

### What makes a good premortem item

- Rooted in what the tests actually do (and don't do) today
- Specific enough to point to the exact pattern or gap where problems start
- Considers the developer experience (speed, debuggability, trust in the suite)
- The mitigation is proportional — lightweight for "adopt now", clearly scoped for "defer"

### Things to consider

- What layers of the app have test coverage? What layers don't?
- How do the tests relate to the production code structure (views, db, sm2)?
- What happens when tests need to render Vue components with routing, stores, or async data?
- How does the DB isolation pattern scale?
- What kinds of bugs would the current tests catch? What kinds would slip through?

## Deliverable

- 3-5 premortem items with analysis and mitigations
- For each: a decision recommendation (adopt now / defer with trigger / reject)
- Any "adopt now" mitigations implemented
- Testing conventions documented for future feature work

## Acceptance Criteria

- [ ] Current test infrastructure fully read and understood
- [ ] 3-5 risks identified with concrete references to current tests/patterns
- [ ] Each risk has a mitigation with adopt/defer/reject recommendation
- [ ] Decisions recorded and any immediate mitigations implemented
