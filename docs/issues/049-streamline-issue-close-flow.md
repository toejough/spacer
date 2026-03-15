# 049 — Streamline the issue-close workflow

**Status:** open
**Type:** tooling

## Problem

Closing an issue currently requires a manual multi-step dance: (1) commit related work, (2) manually update status.md with a summary, (3) commit the status update, (4) run `targ issue-close`. These steps always happen together in the same order, but they're not automated — the LLM has to remember the sequence, and mistakes (forgetting the status update, committing in the wrong order) are likely across sessions.

## Principle

Multi-step sequences that always happen together should be a single command or at least a guided flow. The tool should do the mechanical parts; the human/LLM should only provide the content (the status.md summary text).

## Guidance

Before implementing, consider what can be automated vs. what requires judgment. The status.md summary is content that requires thought — a tool can't generate it well. But the tool could prompt for it, insert it, and handle the commit sequence. Research how CLI tools handle guided multi-step workflows. Read the current `targ issue-close` implementation and the actual close sequences in git history to understand the real pattern. Be cautious about over-automating — the goal is fewer manual steps, not a fragile script that breaks when the workflow varies.
