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
- **#024** Migrate dev/ to targ → created `./targ` dispatcher with dev/test/build/check/clean subcommands, removed `dev/` scripts, updated CLAUDE.md, added `test-results/` to `.gitignore`

### Cycle 8 — Premortems (#014–#017) ← next

Interactive walk-through of four premortem documents to identify weaknesses before the next feature cycle.

| Order | # | Title | Phase | Type |
|-------|---|-------|-------|------|
| 1 | 014 | Architecture premortem | queued | Prevention |
| 2 | 015 | UX/design premortem | queued | Prevention |
| 3 | 016 | Testing premortem | queued | Prevention |
| 4 | 017 | Implementation premortem | queued | Prevention |
