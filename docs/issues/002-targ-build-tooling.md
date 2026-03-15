# Adopt targ as build tooling

**Status:** wont-fix
**Priority:** p1
**Labels:** tooling, issue-1-retro
**Created:** 2026-03-14
**Closed:** 2026-03-14

## Description
Replace direct npm/npx commands with targ as the build system. All build interactions (test, dev server, build, lint, etc.) should go through targ commands. Update `dev/` scripts to use targ, and ensure CLAUDE.md references targ as the standard build tool.

## Acceptance Criteria
- [ ] targ is installed and configured for the project
- [ ] `targ test` runs Vitest
- [ ] `targ dev` starts the Vite dev server
- [ ] `targ build` runs the production build
- [ ] `dev/` scripts updated to use targ
- [ ] CLAUDE.md or project docs note targ as the standard build interface

## Notes
Joe's standard: use `targ` for all build/test/check operations, not npm/npx directly.

Closing — superseded by rebuild. Lessons incorporated into updated process.
