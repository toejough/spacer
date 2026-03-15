# 023 — Add test runner isolation convention to CLAUDE.md

**Status:** closed
**Type:** process
**Source:** retro #18

## Context

When Playwright was added alongside Vitest, Vitest picked up the Playwright test files and crashed. The fix was adding `exclude: ["e2e/**"]` to vitest.config.ts, but this should have been done proactively.

Add a convention: when introducing a new test runner, update all existing test configs to exclude the new runner's directory in the same commit.

## Acceptance Criteria

- [x] Convention documented in CLAUDE.md
