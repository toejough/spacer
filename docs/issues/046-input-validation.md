# 046 — No input validation beyond empty-string checks

**Status:** open
**Type:** quality / prevention
**Source:** #017 implementation premortem

## Problem

Card and deck creation only check for empty strings. There are no length limits, no content sanitization, no duplicate detection. As features grow (markdown in cards, media URLs, import from external formats), unvalidated input risks rendering bugs, storage bloat from accidental large payloads, and potential XSS if raw HTML is ever rendered.

## Principle

Validate at system boundaries — user input is a system boundary. Validation rules should live close to the domain model (not scattered across views) so that every input path (UI, import, sync) applies the same constraints.

## Guidance

Before implementing, assess what validation is actually needed for the current feature set vs. what's hypothetical. Don't add validation for scenarios that don't exist yet. Research lightweight validation approaches for Vue forms — built-in HTML validation attributes, simple composables, or libraries like `valibot`/`zod` for schema validation. Consider where validation logic should live so it's reusable across input paths. Read the current input handling to understand what's changed since this issue was filed.
