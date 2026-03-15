# 040 — No convention for test file organization

**Status:** open
**Type:** quality / prevention
**Source:** #016 testing premortem

## Problem

All integration tests are in a single `src/__tests__/flow.test.ts`. As features grow, either this file becomes a monolith, or new test files appear with no convention for naming, location, or scope. Some tests will end up next to source files, some in `__tests__/`, and the relationship between tests and production code becomes unclear.

## Principle

Test organization should make it obvious where to find tests for a given module, where to add tests for a new feature, and what kind of test (unit/integration/E2E) a file contains.

## Guidance

Before implementing, review common conventions for test organization in Vitest/Vue projects — co-located vs. `__tests__/` directory, naming conventions (`.test.ts` vs. `.spec.ts` for different test types), and how test file structure maps to source structure. Pick a convention that scales to 10+ features and document it. Read the current test setup and file structure to understand what's changed since this issue was filed.

Should follow standards established in #037 if that's been resolved.
