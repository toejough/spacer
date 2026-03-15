# 048 — Add `targ issue-new` to scaffold issue files

**Status:** closed
**Type:** tooling

## Problem

Issue files follow a consistent template (title, status, type, source, problem, principle, guidance) but are created manually each time. This means the LLM has to remember the template, and deviations creep in — inconsistent section names, missing fields, wrong status markers. We filed 20 issues in one session and the template was consistent only because one LLM held context across all of them.

## Principle

Structural consistency should be enforced by tooling, not by memory. If every issue file has the same skeleton, a tool should produce that skeleton.

## Guidance

Before implementing, review the existing issue files in git history (`targ history`) and current issues to identify the actual common structure. Consider whether the template should be rigid (fixed sections) or flexible (required fields + optional sections). Look at how other small projects handle issue templates — GitHub issue templates are one reference point, but this is local-file-based so the analogy is lightweight scaffolding, not a form. Read the current `dev/targets.go` to understand the existing target patterns.
