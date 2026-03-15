# 017 — Implementation Details Premortem

**Status:** closed
**Type:** quality / prevention

## Context

Spacer's implementation is intentionally minimal today:

- **Data layer:** `db.ts` — Dexie DB class with 2 tables (decks, cards), schema + types in one file, singleton `db` export. Views call Dexie directly with inline queries.
- **Domain logic:** `sm2.ts` — pure function, no side effects, clean interface/result types.
- **Views:** 3 SFCs (`HomeView`, `DeckView`, `ReviewView`) with all logic inline — data fetching, mutations, computed state, and templates in one file each. No composables, no shared components.
- **Wiring:** `main.ts` — router defined inline (3 lazy routes), app created and mounted. No stores, no plugins, no middleware.
- **Styling:** Tailwind v4 via Vite plugin, single `@import "tailwindcss"` entry point, all styles as inline utility classes.
- **Build:** Vite 8, vue-tsc for type checking, no linter configured.
- **Dependencies:** Vue 3, vue-router 5, Dexie 4, Tailwind 4. Dev: Vitest, happy-dom, fake-indexeddb, @vue/test-utils (installed but unused).

## The Exercise

Perform a premortem on the implementation details. Assume we've shipped 10+ features and the codebase has accumulated subtle, hard-to-debug problems — data corruption, performance cliffs, type-safety holes, dependency tangles, and code that works by coincidence rather than by design. Bugs are hard to reproduce and fixes in one area break another.

**Your job:** Figure out what led to that state.

### How to run the premortem

1. **Read all production code** — every `.ts` and `.vue` file. Understand the actual data flow: how a card goes from creation to DB to review screen to SM-2 to DB update. Note where types are asserted (`as any`, `as Card`), where errors are silently ignored, where state could get inconsistent.
2. **Examine the dependency and config choices** — `package.json`, `vite.config.ts`, `tsconfig.json`. Note what's configured, what's missing, and what's installed but unused.
3. **Imagine 10+ features added** — each adding data models, queries, UI state, and async operations. Consider: DB migrations, concurrent writes, offline behavior, large datasets, complex computed state, error boundaries.
4. **Identify 3-5 specific implementation weaknesses** that would compound badly under growth. Focus on things that work by coincidence at current scale but would cause real bugs, data loss, or performance problems. Be concrete — reference actual code patterns, type assertions, missing error handling, or implicit assumptions.
5. **For each weakness**, describe: what specifically goes wrong, why the current code enables it, and a concrete mitigation (with a recommendation on whether to adopt now or defer with a specific trigger).

### What makes a good premortem item

- Rooted in actual code, not hypothetical patterns — cite the file and the line
- About correctness, data integrity, or reliability — not just code style
- The failure mode is concrete: "this causes X bug" not "this is messy"
- The mitigation addresses root cause, not symptoms

### Things to consider

- How are types enforced at the DB boundary? What happens when the schema evolves?
- Where does the code assume success (no error handling on DB ops, no validation on inputs)?
- How does the singleton `db` export affect testability and isolation?
- What implicit coupling exists between views and the data layer?
- Are there race conditions or stale-state risks in the async flows?
- What's the DB migration story when models change?

## Deliverable

- 3-5 premortem items with analysis and mitigations
- For each: a decision recommendation (adopt now / defer with trigger / reject)
- Any "adopt now" mitigations implemented
- Implementation conventions documented for future feature work

## Acceptance Criteria

- [ ] All production code fully read and understood
- [ ] 3-5 risks identified with concrete references to current code
- [ ] Each risk has a mitigation with adopt/defer/reject recommendation
- [ ] Decisions recorded and any immediate mitigations implemented
