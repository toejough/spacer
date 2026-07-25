## Why

The todo/review app has no written specification. Without specs, it is hard for agents (or new developers) to understand intended behavior, review correctness, or plan changes without re-reading the codebase each time. This change establishes a lightweight OpenSpec library covering the current behavior so future work can be spec-driven.

## What Changes

- Introduce an `openspec/specs/` library describing current capabilities.
- Add a proposal, design, and task breakdown for the spec creation work.
- No code behavior changes; this is documentation-only.

## Capabilities

### New Capabilities
- `todo-list`: Creating, completing, editing, and deleting todo items.
- `note-taking`: Creating notes with optional cloze deletions.
- `spaced-review`: Reviewing items with the SM-2 spaced repetition algorithm.
- `search`: Searching todos and notes by title and content.

### Modified Capabilities
- None.

## Impact

- New files under `openspec/specs/` and `openspec/changes/add-todo-spec/`.
- No impact on running code, build, or deployment.
