# 051 — LLM drifts from `/commit` skill over long sessions

**Status:** closed
**Type:** process

## Problem

CLAUDE.md says "use `/commit`" for all commits. In a 20-commit session, the LLM used `/commit` once and then switched to manual `git commit` for the remaining commits. The manual commits were fine — correct format, good messages — but the drift means any logic in the commit skill (message formatting, staging rules, validation) gets bypassed silently. This is a general problem: LLM adherence to "always do X" degrades over long sessions as the instruction competes with efficiency pressure.

## Principle

If a process step is important enough to document as mandatory, it should either be enforced by tooling (so skipping it is impossible) or the instruction should be strong enough that it survives context pressure. If neither is practical, maybe it shouldn't be mandatory.

## Guidance

Before addressing, consider the three options: (1) make the commit skill lightweight enough that using it is zero-friction, removing the temptation to skip it, (2) enforce it via a hook that rejects commits not made through the skill, or (3) downgrade "use `/commit`" from mandatory to recommended. Research whether Claude Code hooks can intercept git commits. Read the current commit skill to understand what value it adds beyond `git commit`. The right answer depends on whether the skill provides real guardrails or is just ceremony.
