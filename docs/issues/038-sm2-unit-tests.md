# 038 — SM-2 is untested as a unit

**Status:** open
**Type:** quality / prevention
**Source:** #016 testing premortem

## Problem

`sm2()` in `src/sm2.ts` is the most important pure function in the app — it determines when every card is next shown. It's only exercised indirectly through `flow.test.ts` with quality=4 and quality=1. Edge cases are uncovered: quality=0, quality=5, the boundary at quality=3 (pass/fail threshold), ease factor floor at 1.3, large interval growth over many repetitions.

A subtle math bug for extreme inputs would ship undetected.

## Principle

Pure functions are the easiest and highest-value targets for thorough testing, especially property-based testing. SM-2 has clear invariants that should hold for any input sequence.

## Guidance

Before implementing, identify the behavioral properties of SM-2 that must hold (e.g., ease factor bounds, interval monotonicity for consistent good ratings, reset behavior on failure). Express these as BDD-style specifications. Consider property-based testing (e.g., `fast-check`) to verify invariants across random input sequences rather than just hand-picked examples. Read the current SM-2 implementation to understand what's changed since this issue was filed.

Should follow standards established in #037 if that's been resolved.
