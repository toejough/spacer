# Testing Strategy

## Test Boundaries

| Type | What it tests | I/O | Speed | Runner |
|------|--------------|-----|-------|--------|
| **Unit** | Pure functions, no dependencies | None | Instant | Vitest |
| **Behavioral** | Component/module logic with injected deps | Mocked/injected | Fast | Vitest + happy-dom |
| **Integration** | Real dependencies, real data flow | Real (or faithful fakes) | Moderate | Vitest + fake-indexeddb |
| **E2E** | User journeys through the actual running system | Real everything | Slow | Playwright |

## File Organization

```
tests/
  unit/           # Pure functions, no I/O
  behavior/       # Component logic with injected deps
  integration/    # Real data flow with faithful fakes
  e2e/            # Playwright against production build
  test-setup.ts   # Shared Vitest setup (fake-indexeddb, etc.)
```

Naming: `.test.ts` for Vitest (unit, behavior, integration). `.spec.ts` for Playwright (E2E).

## BDD Style

Tests are planned and written as behavioral specifications:

1. **Identify properties** — behavioral, structural, and presentation invariants
2. **Write specs in given/when/then** — in a plan doc (e.g., `docs/plans/sm2-test-spec.md`)
3. **Commit the spec** — it lives in git history
4. **Copy specs into test files as comments** — comments carry the given/when/then structure
5. **Implement tests under the comments** — code mirrors the spec
6. **Delete the spec doc, commit** — consumed, retrievable via `targ history`

## Test Naming

State the invariant directly in `it()` strings. No "should."

```ts
it("never drops ease factor below 1.3")
it("resets repetitions on failed review")
it("shows rating buttons only after flip")
```

## Assertion Style

**Vitest built-in + jest-extended.** Use the most expressive matcher available. Add domain-specific custom matchers via `expect.extend()` when a real need emerges.

## Dependency Injection

Unit and behavioral tests inject dependencies as function/composable parameters. No module-level singletons in tests. Integration tests use faithful fakes (e.g., fake-indexeddb). E2E tests use the real system.

```ts
// Composable with explicit dependency
function useCards(db: SpacerDB) { ... }

// Test injects a test DB
const db = createTestDB()
const cards = useCards(db)
```

## Property-Based Testing

Use `@fast-check/vitest` to verify invariants across random inputs. Identify properties first, express as BDD specs, then implement with fast-check arbitraries.

Property categories:
- **Behavioral** — algorithm invariants (e.g., ease factor bounds, interval monotonicity)
- **Structural** — data relationship invariants (e.g., every card belongs to one deck)
- **Presentation** — UI state invariants (e.g., exactly one card shown during review)

## Coverage

On-demand via `targ coverage`. 80% per-function threshold as a smell signal. Will become a lint-time gate when a lint pipeline is added. No coverage on E2E tests.

## Mutation Testing

On-demand via Stryker Mutator (`targ mutate`). Verifies that tests actually catch implementation changes. Scope to specific files/modules to keep runtime reasonable. Use periodically, not on every commit.

## Fuzz Testing

Desired but the JS/TS ecosystem lacks a viable OSS coverage-guided fuzzer (Jazzer.js was discontinued). fast-check's property-based testing is our best approximation — structured random input with shrinking. If the project grows a backend in Go or Rust, adopt their built-in fuzzers. Revisit the JS fuzzing ecosystem when features process untrusted external input (import, sync).
