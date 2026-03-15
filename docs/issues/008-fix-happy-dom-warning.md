# Fix happy-dom --localstorage-file warning in test output

**Status:** wont-fix
**Priority:** p2
**Labels:** tooling, issue-1-retro
**Created:** 2026-03-14
**Closed:** 2026-03-14

## Description
Every test run prints `Warning: --localstorage-file was provided without a valid path` from happy-dom. This is noise that obscures real test output.

## Acceptance Criteria
- [ ] Test output is clean — no happy-dom warnings
- [ ] All existing tests still pass

## Notes
May need to configure happy-dom settings in vitest.config.ts, or switch to jsdom if happy-dom's localStorage handling is broken with current versions.

Closing — superseded by rebuild. Lessons incorporated into updated process.
