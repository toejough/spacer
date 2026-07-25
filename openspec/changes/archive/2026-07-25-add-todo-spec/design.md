# Design: add-todo-spec

## Context

The Remember Everything app is a client-side single-page application. Items are stored in `localStorage`, and review logic is implemented in `srv/static/script.js`. This change introduces a spec library to document the existing behavior without altering the implementation.

## Goals / Non-Goals

**Goals:**
- Create a durable, version-controlled specification library under `openspec/specs/`.
- Capture the four primary user-facing capabilities: todos, notes, spaced review, and search.
- Provide a template for future spec-driven changes.

**Non-Goals:**
- Changing any code behavior.
- Migrating data storage from `localStorage` to SQLite.
- Updating the Quasar PWA prototype under `web-app/`.

## Decisions

- **Keep specs in the repo.** OpenSpec stores specs alongside code, so agents can read them in context without external tools.
- **One capability per spec folder.** Each capability gets its own `openspec/specs/<capability>/spec.md` file, matching OpenSpec's default convention.
- **No delta specs needed.** Because this is the first spec library, all requirements are added rather than modified.
- **Focus on user-facing behavior.** Implementation details such as the exact SM-2 formula constants are described in the spec but not exhaustively formalized; the existing code is the reference.

## Risks / Trade-offs

- [Risk] Specs drift from implementation as the code evolves. → Mitigation: Make future changes through OpenSpec change proposals, which produce spec deltas.
- [Risk] Scope creep into documenting every UI detail. → Mitigation: Keep the four specs at a high level; refine incrementally.
