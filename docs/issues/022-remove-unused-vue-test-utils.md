# 022 — Remove or use @vue/test-utils

**Status:** closed
**Type:** tech-debt
**Source:** retro #18

## Context

`@vue/test-utils` is listed in devDependencies but is never imported anywhere. Either write component tests that use it, or remove it to keep dependencies honest.

## Acceptance Criteria

- [x] `@vue/test-utils` removed from package.json
