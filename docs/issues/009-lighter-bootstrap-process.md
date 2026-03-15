# Streamline process for bootstrap-scale tasks

**Status:** wont-fix
**Priority:** p1
**Labels:** tooling, issue-1-retro
**Created:** 2026-03-14
**Closed:** 2026-03-14

## Description
The bootstrap went through the full brainstorming → spec → plan → spec review → plan review → re-review → execute pipeline. This was massively heavy for a scaffold — more time on process than on code. For tasks where the output is well-understood (scaffolds, migrations, dependency upgrades), the process should be lighter.

## Acceptance Criteria
- [ ] 5m-increment skill documents when to use full process vs. light process
- [ ] Light process defined: plan → execute → review result (skip spec doc, skip review loops)
- [ ] Full process reserved for: new features with design decisions, architectural changes, ambiguous requirements
- [ ] Decision captured as ADR in docs/decisions/

## Notes
The full process is valuable — the plan review caught 12 real issues. But it should be proportional to the uncertainty of the task. A scaffold has near-zero uncertainty.

Closing — superseded by rebuild. Lessons incorporated into updated process.
