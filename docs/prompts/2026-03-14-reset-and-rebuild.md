# Reset and Rebuild: Spacer PWA

## Context

The first bootstrap attempt (issue #1) produced 11 commits, 16 passing tests, and an app that does nothing. It took ~15-20 minutes — 3-4x the 5-minute target. The retro surfaced 10 follow-up issues, indicating the initial approach had fundamental problems:

- **Horizontal layers, no vertical slice.** Every architectural piece was scaffolded but nothing was connected. A user visits the app and sees placeholder text.
- **Heavy process for low-uncertainty work.** Full brainstorming → spec → plan → review → re-review pipeline for a scaffold.
- **False confidence from unit tests.** 16 tests passed on isolated pieces. Zero integration coverage. No test verified the app actually works.
- **Empty stubs counted as deliverables.** Component and store stubs were created but will be completely rewritten.

## What to do

### Phase 1: Clean up

1. Close all open issues (#2-#11) as `wont-fix` with note: "Closing — superseded by rebuild. Lessons incorporated into updated process."
2. Update `docs/status.md` to reflect the reset.
3. Git: revert all code commits from issue #1 (keep docs/retros, docs/issues, docs/prompts, docs/plans, docs/superpowers, .claude-plugin). The goal is to remove all application code (src/, package.json, vite configs, dev/, public/, etc.) while preserving project memory.
4. Commit the cleanup.

### Phase 2: Apply lessons to process

Update `.claude-plugin/skills/5m-increment/SKILL.md` with these learnings:

1. **Bootstrap must deliver a vertical slice.** The bootstrap AC must include: "a user can complete one full action end-to-end." Scaffolding empty stubs across all layers is not acceptable. Pick the thinnest possible vertical slice (e.g., create a deck with one card and review it) and make that work.

2. **Match process weight to task uncertainty.**
   - High uncertainty (new features, design decisions): full brainstorming → spec → plan → review
   - Low uncertainty (scaffolds, wiring, known patterns): plan → execute → review result
   - The 5-minute increment itself is already a light process — don't layer heavy process on top of it.

3. **Integration tests from the start.** The bootstrap smoke test must verify the vertical slice works end-to-end, not just that isolated units exist. A test that mounts a view, interacts with it, and checks the DB changed is worth more than 10 tests on disconnected pieces.

4. **Don't create empty stubs as deliverables.** If a component will be completely rewritten when it's actually implemented, don't create it during bootstrap. Only create files that contain real, working code.

5. **Code quality review on aggregate.** Per-task reviews can be skipped for speed during bootstrap, but at least one quality review must cover the full output before the issue is closed.

Also update `references/file-formats.md` to add an `issue-N-retro` label convention for issues surfaced during retros.

### Phase 3: Rebuild

Start a new increment for the rebuild. This time:

- **One vertical slice as the bootstrap.** The AC is: a user can open the app, create a deck, add a card, and review it. All in one increment.
- **Use the light process.** Plan → execute → review result. No spec document, no review loops.
- **Integration test as the smoke test.** The primary test should exercise the full flow: mount app → create deck → add card → navigate to review → flip card → rate → verify SM-2 state updated in DB.
- **Only create files with real code.** No empty stubs, no placeholder components.
- **targ from the start.** Set up targ as the build tool during scaffold, not as a follow-up issue.
- **Project CLAUDE.md from the start.** Create it during scaffold so future sessions have context immediately.

### Success criteria for the rebuild

- [ ] A user can open the app, create a deck, add a card with front/back, and review it
- [ ] Review updates SM-2 state and persists to IndexedDB
- [ ] At least one integration test verifies the full flow
- [ ] targ is the build tool interface
- [ ] Project CLAUDE.md exists with build commands and conventions
- [ ] Total time ≤ 10 minutes (bootstrap is allowed to exceed 5, but not by 3x)
- [ ] ≤ 3 follow-up issues from the retro (not 10)

## Reference

- Retro: `docs/retros/2026-03-14-001-bootstrap.md`
- Original prompt: `docs/prompts/2026-03-14-5m-increment.md`
- Original plan: `docs/plans/2026-03-14-5m-increment.md`
