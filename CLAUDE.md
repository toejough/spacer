# Spacer PWA

Spaced repetition flashcard app. Vue 3 + TypeScript + Vite + Dexie (IndexedDB) + Tailwind CSS v4.

## Build Commands

```bash
dev/dev          # Start Vite dev server
dev/test         # Run Vitest
dev/build        # Type-check + production build
dev/check        # Type-check + tests
```

## Conventions

- **TDD always** — write failing tests first, then minimum code to pass
- **Vertical slices** — every increment delivers a working user-facing feature
- **No empty stubs** — only create files with real, working code
- **Integration tests** — at least one test per feature that exercises the full data path
- **Commits** — use `/commit`, conventional commits format
- **Issues** — `docs/issues/{number}-{slug}.md`
- **Status** — `docs/status.md` updated every cycle

## Tech Notes

- DB: Dexie with EntityTable for type-safe IndexedDB
- SM-2: pure function in `src/sm2.ts`
- Test env: happy-dom + fake-indexeddb
- Each test gets its own DB instance (no shared mutable state)
