# 039 — No view/component testing layer

**Status:** open
**Type:** quality / prevention
**Source:** #016 testing premortem

## Problem

The views contain real logic — due-card filtering, input validation (empty string guards), the review flow state machine (flip → rate → advance → done), and route param parsing. None of this is tested. The flow test validates DB+SM-2 but not Vue wiring. Bugs in template conditionals, reactive state transitions, or route param handling would only be caught by E2E tests (which are slow and only cover offline).

## Principle

View logic should be testable without a browser. The right approach depends on how much logic lives in views vs. extracted composables/functions. DI should be used to inject a test DB rather than hitting real IndexedDB where possible.

## Guidance

Before implementing, assess how much logic currently lives in the views and whether it should be tested in-place (component tests with injected dependencies) or extracted and tested as pure functions. Research Vue 3 component testing patterns with Vitest — particularly how to inject mock Dexie instances and route params. Consider which presentation properties should hold (e.g., "rating buttons only visible after flip", "done state shown after last card").

Should follow standards established in #037 if that's been resolved.
