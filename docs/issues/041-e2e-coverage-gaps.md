# 041 — E2E tests only cover offline happy path

**Status:** open
**Type:** quality / prevention
**Source:** #016 testing premortem

## Problem

The 2 Playwright tests in `e2e/offline.spec.ts` verify that the app loads and navigates offline. The core user journey (create deck → add card → review → see scheduling effect) has no E2E coverage. Neither does error behavior, edge cases, or any feature beyond offline. As features multiply, the E2E suite stays frozen on offline while regressions in actual user flows go undetected.

## Principle

E2E tests should cover the critical user journeys that, if broken, make the app useless. They're expensive to write and slow to run, so they should be selective — but the core happy path should always be covered.

## Guidance

Before implementing, identify the 2-3 critical user journeys that warrant E2E coverage (the core review loop is the obvious first candidate). Research Playwright best practices for PWA testing — particularly how to structure E2E tests so they're maintainable as features grow. Read the current E2E setup and user flows to understand what's changed since this issue was filed.
