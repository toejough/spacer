# 020 — Move spec location from docs/superpowers/ to docs/plans/

**Status:** open
**Type:** process
**Source:** retro #18

## Context

The brainstorming skill defaults to writing specs to `docs/superpowers/specs/`, tying our repo's doc structure to a specific plugin. We're already using `docs/plans/` for this purpose. Add an override to CLAUDE.md so specs go to `docs/plans/` and remove/migrate any existing `docs/superpowers/` content.

## Acceptance Criteria

- [ ] CLAUDE.md specifies `docs/plans/` as the spec location
- [ ] Any existing content in `docs/superpowers/` migrated or removed
- [ ] Future brainstorming specs land in `docs/plans/`
