# Don't skip code quality review on bootstrap output

**Status:** wont-fix
**Priority:** p2
**Labels:** tooling, issue-1-retro
**Created:** 2026-03-14
**Closed:** 2026-03-14

## Description
During the bootstrap, code quality reviews were skipped for most tasks ("it's generated code", "it's just a scaffold"). This set a bad precedent. At minimum, a single quality review pass should cover the full bootstrap output before closing the issue.

## Acceptance Criteria
- [ ] Process documented: bootstrap gets at least one full code quality review at the end
- [ ] Spec review can be skipped per-task for scaffolds, but quality review happens at least once on the aggregate

## Notes
Skipping per-task reviews during a scaffold is reasonable for speed. Skipping ALL reviews is not.

Closing — superseded by rebuild. Lessons incorporated into updated process.
