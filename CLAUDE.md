# Spacer PWA

Spaced repetition flashcard app. Vue 3 + TypeScript + Vite + Dexie (IndexedDB) + Tailwind CSS v4.

## Build Commands

```bash
targ dev             # Start Vite dev server
targ test            # Run Vitest
targ build           # Type-check + production build
targ check           # Type-check + tests
targ clean           # Remove build/test artifacts
targ issues          # List open issues
targ issue-new SLUG  # Scaffold a new issue file with next number
targ issue-close N   # Close and archive issue (requires status.md entry)
targ history         # List deleted docs from git history
targ history-show P  # Show a deleted file from git history
```

## Conventions

- **TDD always** — write failing tests first, then minimum code to pass
- **Vertical slices** — every increment delivers a working user-facing feature
- **No empty stubs** — only create files with real, working code
- **Integration tests** — at least one test per feature that exercises the full data path
- **Commits** — use `/commit`, conventional commits format
- **Issues** — `docs/issues/{NNN}-{slug}.md`, open issues only; use `targ issue-new` to scaffold, `targ issue-close` to close+archive. Every issue has three sections: **Problem** (concrete description with code/file references), **Principle** (the underlying design principle), **Guidance** (direction for the resolver — including external research recommendations — but leave room to re-evaluate against the codebase as it exists when the issue is picked up, not when it was filed)
- **Status** — `docs/status.md` is a narrative project log (timeline + rationale), updated every cycle
- **Plans/specs** — write to `docs/plans/`, not `docs/superpowers/`
- **Docs lifecycle** — HEAD contains only current/open docs. `targ issue-close` closes and archives in one step (two commits). `docs/status.md` indexes all historical work. To find deleted docs: `targ history` then `targ history-show <path>`
- **Retro format** — `docs/retros/{date}-{number}-{slug}.md`. Sections: **Plus** (reinforce/repeat), **Delta** (change next time), **Other observations**, **Action items**. Every action item must either be filed as an issue or absorbed into CLAUDE.md before the retro can be archived. Retros are archived once all action items are handled.
- **Test runner isolation** — when adding a new test runner, update all existing test configs to exclude the new runner's directory in the same commit
- **Claims need code** — if CLAUDE.md, package.json, or project descriptions claim a capability, working code must back it. Don't document features that don't exist yet.
- **Fix errors, don't report them** — when a command, build, or test fails: (1) read the error, (2) identify the cause, (3) fix it, (4) re-run to verify. Only stop and ask if you've attempted a fix and it failed twice, or if the error is genuinely ambiguous with multiple plausible causes.

## Architecture

- **No stores** — views query Dexie directly. Add a store layer only when cross-component reactivity is needed.
- **No extracted components** — all UI lives in 3 view files (Home, Deck, Review). Extract components when there's duplication or complexity, not preemptively.
- **Routes** — `/` (Home), `/deck/:id` (Deck), `/review/:deckId` (Review), all lazy-loaded
- **PWA** — Workbox via `vite-plugin-pwa`, `generateSW` strategy, precaches all build assets including lazy chunks

## Tech Notes

- DB: Dexie with EntityTable for type-safe IndexedDB
- SM-2: pure function in `src/sm2.ts`
- Unit tests: Vitest + happy-dom + fake-indexeddb; each test gets its own DB instance
- E2E tests: Playwright in `e2e/`, runs against production build (SW requires it)
