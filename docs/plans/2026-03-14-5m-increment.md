# Plan: 5-Minute Increment Skill

**Date:** 2026-03-14
**Source prompt:** `docs/prompts/2026-03-14-5m-increment.md`
**Goal:** Create a Claude Code skill that manages the 5-minute increment development cycle.

## Scope

One skill: cycle management (PICK → PLAN → RED → GREEN → REFACTOR → UAT → CLOSE). No bootstrap/discovery — that's a separate concern for later.

## Deliverables

### 1. Skill structure

```
skills/5m-increment/
├── SKILL.md              # Core workflow (~1500 words)
└── references/
    └── file-formats.md   # Issue, status, and retro file templates
```

### 2. SKILL.md contents

**Frontmatter:**
- name: `5m-increment`
- description: Trigger phrases — "start an increment", "next increment", "run a cycle", "5 minute increment", "pick an issue"

**Body — the cycle phases:**

| Phase | Budget | What happens |
|-------|--------|-------------|
| PICK | ≤30s | Show top 3 issues from `docs/issues/`, user picks one, confirm AC |
| PLAN | ≤60s | List files to touch, state approach in ≤3 bullets, surface ADRs |
| RED | ≤60s | Write failing tests, run them, confirm red |
| GREEN | ≤90s | Minimum code to pass, run tests, confirm green |
| REFACTOR | ≤30s | Optional cleanup, tests stay green |
| UAT | ≤30s | UI → user verifies; logic → tests are UAT |
| RETRO | ≤30s | 2-sentence retro (what worked, what to improve) + concrete action items. Actions become issues or skill changes |
| CLOSE | ≤30s | Pass → commit + update issue/status. Fail → stash + mark blocked |

**Rules (trimmed to essentials):**
- Phase budgets guide scope-sizing — if a phase feels like it'll overrun, propose a scope cut or rollback immediately
- Deterministic tools over agents — propose scripts for repeated manual steps
- Docs stay current — `docs/status.md`, `docs/issues/`, `docs/retros.md` updated every cycle
- Fail fast — failed increment = data, not a problem. Capture in retro, move on
- Ask, don't assume — surface architecture/preference questions during PLAN

**Integrations:**
- Use `/commit` during CLOSE phase
- RED/GREEN phases: lean on TDD skill if it's pulling its weight; if it adds friction for the tight timeboxes, recommend disabling it and inlining a lighter TDD flow in this skill
- PLAN phase: skip `/project` — its overhead doesn't fit a 5-minute cycle. This skill owns planning inline
- Read `docs/status.md` on load to pick up context
- Evaluate existing skills during use — if one gets in the way more than it helps, recommend disabling it and propose a replacement scoped to this workflow

### 3. references/file-formats.md

Templates for:
- Issue files (`docs/issues/{number}-{slug}.md`)
- Status file (`docs/status.md`)
- Retro entries (appended to `docs/retros.md`) — includes action items that feed back into the backlog

### 4. Plugin manifest

Minimal `.claude-plugin/plugin.json` (already created).

## Out of scope

- Bootstrap/discovery session (tech stack interview, ADRs, scaffolding)
- Timer/countdown enforcement (phase budgets are guidance only)
- Backlog seeding

## Steps

1. Write `references/file-formats.md` with the three templates from the prompt
2. Write `SKILL.md` with frontmatter + cycle workflow + rules
3. Delete the pre-created directories we don't need (`assets/`, etc.)
4. Test: verify skill loads by asking Claude to "start an increment"
