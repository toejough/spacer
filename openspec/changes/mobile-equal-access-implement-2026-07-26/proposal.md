# Implement: Mobile equal-access Todo UX (Done / Abandon parity + Reopen)

## Goal
Implement the mobile-first UX where Done and Abandon are equally accessible, with Reopen available from closed states. Deliver this change using TDD and include visual/accessibility analysis as part of CI.

## Scope
- Frontend: MobileTodoList + MobileTodoCard components; visual treatments for Open/Done/Abandoned; Reopen action; keyboard/ARIA support.
- Tests: Playwright E2E tests (mobile-equal-access.spec.ts) added first (TDD). Unit tests for components where applicable.
- Visual analysis: automated contrast checks and accessibility audit (axe or computed contrast) as part of test suite.
- Feature flag: gate rollout behind feature flag 'mobile_equal_access' (config/ENV) — default enabled for demo.

## Rollout
- Merge behind feature flag; run A/B experiment if desired.

