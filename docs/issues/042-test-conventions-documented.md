# 042 — No test conventions documented

**Status:** open
**Type:** quality / prevention
**Source:** #016 testing premortem

## Problem

There's no documented guidance on when to write a unit test vs. integration test vs. E2E test, where test files should live, how to name them, or what patterns to follow. The current test suite is too small to reveal this gap, but as features grow, contributors will make inconsistent choices.

## Principle

Test conventions should be documented in CLAUDE.md so that every new feature gets tests written the same way. The conventions should cover: test types and when to use each, file naming and location, assertion style, and DI expectations.

## Guidance

This issue should be resolved after #037 (test standards) and #040 (file organization) are settled, since the conventions document should reflect those decisions. The deliverable is a section in CLAUDE.md that captures the agreed-upon testing approach.
