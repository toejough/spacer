# Spacer — Project Log

High-level timeline of project decisions and actions. Updated every cycle.

## Cycle Phases

Each issue progresses through: **pick → plan → implement (red/green/refactor) → retro → file action items → close**

## Timeline

### Cycle 1 — Bootstrap (#1)

Scaffolded the project with horizontal layers: Vite, Vue Router, Pinia, Dexie, SM-2, PWA manifest, Tailwind. Created stubs for all components. No working end-to-end flow.

**Retro:** Filed 10 issues (#2–#11) covering missing vertical slice, code quality, tooling gaps. Realized the horizontal-layer approach left nothing usable. Reverted all application code and closed #1–#11.

### Cycle 2 — Rebuild (#12)

Started fresh with a vertical-slice approach: one working flow (create deck → add card → review with SM-2) delivered incrementally. Each commit added a working layer.

**Retro:** Process worked much better. Filed #13 (plan gate) to prevent future multi-file changes without a written plan.

### Cycle 3 — Plan Gate (#13)

Added structural plan gate to the 5-minute increment skill: any change touching 3+ files requires a written plan in `docs/plans/` with user approval. Closed immediately — process change, no code.

### Cycle 4 — Premortems (#14–#17)

Filed four premortem issues to identify weaknesses before they compound: architecture (#14), UX/design (#15), testing (#16), implementation details (#17). All still open — intended as living documents.

### Cycle 5 — PWA Offline Fix (#18)

**Picked** because the app crashed without a dev server — the PWA promise was broken. Added service worker with Workbox precaching, fixed manifest, added E2E tests with Playwright.

**Retro:** Filed 5 action items (#19–#23) covering status doc maintenance, spec location, legacy-peer-deps, unused test-utils dep, and test runner isolation convention.

### Cycle 6 — Status Doc (#19)

**Picked** this process issue. Decision: keep `docs/status.md` as a narrative project log (not a duplicated issue board). Redesigned format to track timeline and per-issue cycle phases.

### Cycle 7 — Housekeeping Sprint (#022, #023, #021, #020, #024)

Batched five quick wins as autonomous work. Rationale: clear friction and hygiene issues first so premortems can inform the next feature cycle cleanly.

- **#021** Fix --legacy-peer-deps → downgraded Vite 8.x to 7.x; `npm install` now works clean
- **#022** Remove unused vue/test-utils → `npm uninstall @vue/test-utils`
- **#023** Test runner isolation convention → added to CLAUDE.md
- **#020** Move specs to docs/plans/ → migrated `docs/superpowers/` content, removed directory, added convention to CLAUDE.md
- **#024** Migrate dev/ to targ → build targets in `dev/targets.go`, system `targ` discovers them; added `test-results/` to `.gitignore` (approach fixed by #025)
- **#025** Fix targ setup → replaced root `./targ` shell script with proper `dev/targets.go` string targets; added doc lifecycle targets (`issues`, `issue-close`, `issue-archive`, `history`, `history-show`); adopted "HEAD = current state" convention for docs

- **#026** Fix errors, don't report them → convention added to CLAUDE.md: read error, identify cause, fix, re-run; only stop after two failed attempts or genuine ambiguity
- **#027** Claims need code → convention added to CLAUDE.md: don't document capabilities that aren't backed by working code

### Cycle 8 — Premortems (#014–#017)

Interactive walk-through of four premortem documents to identify weaknesses before the next feature cycle. Also merged `issue-close` and `issue-archive` targ targets into a single `issue-close` command.

- **#014** Architecture premortem → filed #028–#031 (due-card query duplication, SM-2 default factory, data-access layer, card model separation)
- **#015** UX/design premortem → filed #032–#036 (action feedback, design tokens, navigation structure, review rating labels, view state patterns)
- **#016** Testing premortem → filed #037–#042 (test standards with BDD/property-based/DI/expressive matchers, SM-2 unit tests, view testing, test file organization, E2E coverage gaps, test conventions doc)
- **#017** Implementation premortem → filed #043–#047 (type safety at DB boundary, DB migration pattern, DB error handling, input validation, stale state across views)

**Retro:** Process-focused retro identified 5 process/tooling gaps. Filed #048–#052 (issue template scaffolding, streamlined issue-close flow, issue-writing conventions, commit skill adherence, interactive discovery enforcement). Prioritized remaining issues by structural impact and dependency overlap to minimize rework.

### Cycle 9 — Issue Conventions (#050, #048)

Documented issue-writing conventions (Problem/Principle/Guidance template) in CLAUDE.md. Added `targ issue-new` command to scaffold issue files with auto-incremented numbers, scanning both current files and git history for highest number.

### Cycle 10 — Test Standards (#037, absorbing #040, #042) ← next

Establish BDD, property-based, DI, and expressive matcher conventions before any code work. Research ecosystem (fast-check, chai/vitest matchers), pick libraries, document conventions. #040 (file organization) and #042 (conventions doc) are naturally part of this — the output is both the standards and where tests live.

### Cycle 11 — Data-Access Layer (#030 + #047 + #028 + #029)

Biggest structural change. Adopt Dexie liveQuery as the reactive data-access pattern, rewriting how all views consume data. Subsumes #028 (due-card query) and #029 (SM-2 factory) into the access layer. Solves #047 (stale state) via reactive queries. Creates the seam where #043 (type safety) and #045 (error handling) will later be enforced.

### Cycle 12 — Type Safety at DB Boundary (#043)

With the access layer in place from cycle 11, enforce types at the insertion/query boundary. Eliminate `as Card` / `as any` casts. Small scope after the access layer exists.

### Cycle 13 — UI Foundations (#033 + #036)

Establish design tokens and view state patterns. Views were just rewritten for data (cycle 11), now standardize their styling and loading/error/empty representation. All subsequent UX issues (#032, #034, #035) inherit these patterns.
