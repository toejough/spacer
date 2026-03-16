# 052 — LLM skips interactive/discovery steps

**Status:** open
**Type:** process

## Problem

The #014 premortem was documented as an "interactive walk-through" in status.md, and the LLM's own CLAUDE.md says "don't skip interview/discovery steps." Despite both instructions, the LLM read the codebase, generated findings, and went straight to implementing code changes without any discussion. The user had to intervene to course-correct. This is a recurring pattern — LLMs default to action over dialogue, especially when the task description includes a deliverable.

## Principle

Analysis tasks and implementation tasks require different modes. When a task says "walk through this together" or "discuss before acting," the LLM should present findings and wait for feedback before making changes. The bias toward action is strong and needs structural counterweight.

## Guidance

Before addressing, consider where the enforcement should live. Options: (1) issue templates could have a "mode: interactive" flag that the LLM checks, (2) CLAUDE.md could have stronger language about analysis-vs-implementation tasks, (3) the premortem issues themselves could have been written more explicitly to prevent autonomous execution. Research how other AI-assisted workflows distinguish "think together" tasks from "go do it" tasks. Read current CLAUDE.md warnings about skipping discovery to assess whether stronger language would help or whether this needs structural enforcement.
