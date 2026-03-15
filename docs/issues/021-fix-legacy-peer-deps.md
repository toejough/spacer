# 021 — Fix --legacy-peer-deps requirement

**Status:** closed
**Type:** tech-debt
**Source:** retro #18 (also noted in retro #1)

## Context

Every `npm install` requires `--legacy-peer-deps` because `@tailwindcss/vite@4.2.1` hasn't declared support for Vite 8. This has been a known issue since the bootstrap and is still unresolved.

Options:
- Pin to a Tailwind version that supports Vite 8
- Downgrade Vite to 7.x
- Wait for `@tailwindcss/vite` to update their peer dep
- Document the workaround if we accept the status quo

## Acceptance Criteria

- [x] `npm install` works without `--legacy-peer-deps` — downgraded Vite from 8.x to 7.x
