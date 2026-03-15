# 014 — Architecture Premortem

**Status:** open
**Type:** tech-debt / prevention

## Context

Spacer is small today: 3 views, 2 domain models, 1 pure algorithm, flat file structure. But the roadmap implies significant growth (import/export, tags, search, stats, sync, settings, media, etc.).

## The Exercise

Perform a premortem on the codebase architecture. Assume we've shipped 10+ features and the codebase is now a rats nest of difficult-to-follow call flows, duplication, and inconsistency — making it hard to add features and easy to introduce bugs and performance problems.

**Your job:** Figure out what led to that state.

### How to run the premortem

1. **Read the current codebase** — understand the actual file structure, how views call the DB, how SM-2 is wired into the review flow, how types are defined, how routing works. The codebase is small; read all of it.
2. **Imagine 10+ features added** — stats dashboard, tags/search, import/export, sync, settings, media cards, cram mode, streak tracking, multiple review modes, sharing, etc.
3. **Identify 3-5 specific architectural weaknesses** in the current codebase that would compound badly under that growth. Focus on things that are fine at current scale but would rot. Be concrete — reference actual files, patterns, and call flows, not generic advice.
4. **For each weakness**, describe: what specifically goes wrong, why the current structure enables it, and a concrete mitigation (with a recommendation on whether to adopt now or defer with a specific trigger).

### What makes a good premortem item

- Rooted in what the code actually does today, not hypothetical patterns
- Specific enough that you could point to the exact file/line where the rot starts
- The mitigation is proportional — lightweight for "adopt now", clearly scoped for "defer"

## Deliverable

- 3-5 premortem items with analysis and mitigations
- For each: a decision recommendation (adopt now / defer with trigger / reject)
- Any "adopt now" mitigations implemented
- Conventions documented in CLAUDE.md or a conventions doc

## Acceptance Criteria

- [ ] Current codebase fully read and understood
- [ ] 3-5 risks identified with concrete references to current code
- [ ] Each risk has a mitigation with adopt/defer/reject recommendation
- [ ] Decisions recorded and any immediate mitigations implemented
