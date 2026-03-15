# 050 — Document issue-writing conventions

**Status:** closed
**Type:** process

## Problem

We arrived at a good issue-writing style during the premortem cycle through trial and correction: enough context and rationale for the resolver to understand the problem, guidance to do external research on best practices, and an explicit instruction to review the codebase for current relevance rather than assuming the state when the issue was filed. This style isn't documented anywhere — it lives only in the examples set by the 20 issues we just filed.

The key properties of a well-written issue in this project:
- **Problem**: concrete description of what's wrong, with references to actual code/files
- **Principle**: the underlying design principle being violated
- **Guidance**: direction for the resolver, including external research recommendations, but explicitly leaving room to re-evaluate against the current codebase state

This style prevents two failure modes: (1) issues so vague that the resolver has to re-discover the problem, and (2) issues so prescriptive that the resolver implements a stale solution that doesn't fit the codebase as it's evolved.

## Principle

Process knowledge should be captured where it's needed. If every issue should follow a pattern, that pattern should be documented where issue creators will see it.

## Guidance

Before implementing, review the issues filed in this cycle (git history for #028–#047) as examples of the style. Decide where the convention belongs — CLAUDE.md (always loaded), a section in status.md, or as part of the `targ issue-new` template (#048). Consider whether it should be prescriptive (required sections) or descriptive (guidelines with examples). Read current CLAUDE.md conventions to understand the existing documentation style.
