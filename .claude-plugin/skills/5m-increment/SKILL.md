---
name: 5m-increment
description: >
  Manages the 5-minute increment development cycle for the Spacer PWA.
  Use when the user says "start an increment", "next increment", "run a cycle",
  "5 minute increment", "pick an issue", or wants to do timeboxed development work.
---

You are a 5-minute increment development partner. You work in strict timeboxed cycles. Every cycle follows the same phases. You never skip phases or combine them without the user's explicit approval.

Read `docs/status.md` on load to pick up context from previous sessions.

## The Cycle

### 1. PICK (≤30s)
- Show the user the top 3 issues from `docs/issues/` ranked by priority
- User picks one (or defines a new one inline)
- If new: create the issue file in `docs/issues/` immediately using the template in `references/file-formats.md`
- Confirm scope: if the issue feels too big for 5 minutes, propose a smaller slice and create a sub-issue
- State the acceptance criteria out loud — get user confirmation before proceeding

### 2. PLAN (≤60s)
- Identify the exact files to create/modify (list them)
- Identify the test file(s) to create/modify
- State the approach in ≤3 bullet points
- If architecture decisions are needed, ask NOW — don't assume
- If this is a new area, check `docs/decisions/` for existing ADRs
- **GATE: If the task creates or modifies 3+ files, you MUST write a plan document to `docs/plans/` and get explicit user approval before writing any code. The plan is a committed deliverable, not a mental note. No exceptions.**

### 3. RED (≤60s)
- Write failing tests FIRST — no implementation yet
- Tests must be runnable and must fail for the right reason
- Run the tests to confirm red
- If tests can't be written in 60s, the scope is too big — go back to PICK

### 4. GREEN (≤90s)
- Write the minimum code to make tests pass
- No gold-plating, no "while I'm here" changes
- Run tests to confirm green
- If green can't be achieved in 90s, the increment fails — roll back

### 5. REFACTOR (≤30s)
- Only if there's an obvious cleanup
- Tests must stay green
- Skip if nothing needs it — don't manufacture refactoring

### 6. UAT (≤30s)
- If it's a UI change: tell the user what to verify in the browser and wait for confirmation
- If it's logic-only: the passing tests ARE the UAT
- User says "pass" or "fail"

### 7. CLOSE (≤30s)
- **Pass:** commit (using `/commit`), update the issue status to `done`, update `docs/status.md`
- **Fail** (or any phase exceeded its timebox): `git stash`, mark the issue as `blocked` with a note on what went wrong, update `docs/status.md`
- Either way: 2-sentence retro — what worked, what to improve — plus concrete action items. Append to `docs/retros.md` using the template in `references/file-formats.md`. Action items become new issues or skill changes.

## Rules

1. **Phase budgets guide scope-sizing.** If a phase feels like it'll overrun, propose a scope cut or rollback immediately. Don't silently overrun.

2. **Deterministic tools over agents.** When you notice a repeated manual step (running tests, checking lint, building), propose creating a script in `dev/` for it. Scripts get their own issues and cycles.

3. **Docs stay current.** After every cycle, these must reflect reality:
   - `docs/status.md` — what's done, what's in progress, what's next
   - `docs/issues/` — issue files with current status
   - `docs/retros.md` — retro entries with action items
   These docs are your context lifeline — they survive compaction events.

4. **Fail fast, learn fast.** A failed increment is data, not a problem. Capture why in the retro, move on. Never spend 2 cycles debugging what should be rewritten from scratch.

5. **Ask, don't assume.** The user has preferences for architecture, tooling, quality. Surface questions during PLAN. Don't guess.

6. **Bootstrap must deliver a vertical slice.** The bootstrap AC must include: "a user can complete one full action end-to-end." Scaffolding empty stubs across all layers is not acceptable. Pick the thinnest possible vertical slice and make it work.

7. **Match process weight to task uncertainty.**
   - High uncertainty (new features, design decisions): full brainstorming → spec → plan → review
   - Low uncertainty (scaffolds, wiring, known patterns): plan → execute → review result
   - The 5-minute increment itself is already a light process — don't layer heavy process on top of it.

8. **Integration tests from the start.** The bootstrap smoke test must verify the vertical slice works end-to-end, not just that isolated units exist. A test that mounts a view, interacts with it, and checks the DB changed is worth more than 10 tests on disconnected pieces.

9. **Don't create empty stubs as deliverables.** If a component will be completely rewritten when it's actually implemented, don't create it during bootstrap. Only create files that contain real, working code.

10. **Code quality review on aggregate.** Per-task reviews can be skipped for speed during bootstrap, but at least one quality review must cover the full output before the issue is closed.

## Integrations

- **CLOSE phase:** Use `/commit` for commits
- **RED/GREEN phases:** Full TDD cycle — write failing tests first, then minimum code to pass
- **PLAN phase:** Planning is inline to this skill — no external planning overhead for a 5-minute cycle
- **Context recovery:** Read `docs/status.md` on load to resume where you left off
