# Add structural plan gate to 5m-increment skill

**Status:** done
**Priority:** p1
**Labels:** tooling, issue-2-retro
**Created:** 2026-03-15
**Closed:** 2026-03-15

## Description
During the rebuild (#12), the LLM jumped straight to execution despite having the retro, updated skill, and prompt all saying "plan first." Information alone doesn't change behavior — a structural gate is needed.

## Acceptance Criteria
- [x] 5m-increment skill PLAN phase has a hard gate: 3+ files → written plan in `docs/plans/` + user approval
- [x] Plans are committed deliverables in `docs/plans/`, not mental notes
- [x] Existing plan moved from `docs/superpowers/plans/` to `docs/plans/`
- [x] Retro updated with root cause and action item

## Notes
Fixed in commit 5ce78c0. Root cause: "I know what to build" ≠ "skip the plan." The plan exists for the next session, not just the current one.
