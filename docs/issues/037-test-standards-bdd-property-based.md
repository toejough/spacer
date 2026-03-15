# 037 — Establish test standards: BDD, property-based, DI, expressive matchers

**Status:** closed
**Type:** quality / prevention
**Source:** #016 testing premortem

## Problem

The current test suite has no documented standards for how tests should be written. As features grow, each contributor will write tests in whatever style they're comfortable with — some imperative, some declarative, some testing implementation details, some testing behavior. The suite becomes inconsistent and hard to read.

## Desired Standards

**BDD style:** Tests should read as behavioral specifications. Describe what the system does, not how the code works internally. Structure tests as given/when/then scenarios.

**Hamcrest-style matchers:** Assertions should be concise and expressive — closer to natural language than method chains. Prefer matchers that describe the expected shape/property rather than exact equality where possible.

**Dependency injection over real I/O:** Unit and behavioral tests should inject dependencies (DB, clock, etc.) rather than hitting real I/O. Reserve actual I/O for dedicated integration and E2E tests. This keeps the fast tests fast and the slow tests intentional.

**Property-based testing:** Identify the behavioral, structural, and presentation properties that must hold true, express them as BDD-style requirements, and verify them with property-based tests. Examples for Spacer:
- Behavioral: SM-2 ease factor never drops below 1.3 regardless of input sequence
- Structural: every card belongs to exactly one deck; due-card queries return a subset of deck cards
- Presentation: review always shows exactly one card; rating buttons only appear after flip

## Guidance

Before implementing, research the TypeScript/Vitest ecosystem for:
- BDD-style test structure (Vitest's `describe`/`it` is a start, but consider whether a more expressive layer helps)
- Hamcrest-style matcher libraries for JS/TS (e.g., `jest-extended`, `chai` with plugins, or Vitest's `expect.extend`)
- Property-based testing libraries (e.g., `fast-check`) and how they integrate with Vitest
- DI patterns for Vue + Dexie that don't require a full DI framework

Identify the specific properties that must hold for the current codebase, express them in BDD language, then pick the lightest libraries that make expressing those properties trivial in code. Don't adopt a library just because it exists — adopt it because it makes expressing a real property easier than the alternative.

Document the chosen standards in CLAUDE.md so future test work follows them.
