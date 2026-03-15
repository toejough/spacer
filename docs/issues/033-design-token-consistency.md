# 033 — No shared design tokens or component patterns

**Status:** open
**Type:** ux / prevention
**Source:** #015 UX/design premortem

## Problem

Visual styling is applied via copy-pasted Tailwind classes across views. `bg-indigo-600` appears on every button and the header. Spacing, rounding, and shadow values are repeated verbatim. There's no single source of truth for "what a primary button looks like" or "what a card container looks like."

After 10+ features, slight deviations accumulate — `bg-indigo-500` here, `rounded-lg` there — and the app looks like a patchwork. A color or spacing change requires a find-and-replace across every view.

## Principle

Design tokens (colors, spacing, typography) and recurring UI patterns (buttons, cards, form inputs) should have a single source of truth. The mechanism doesn't need to be heavy — Tailwind's `@apply` in a utility class, CSS custom properties, or even a documented convention can work at this scale.

## Guidance

Before implementing, review Tailwind CSS v4's approach to theming and design tokens (CSS custom properties, `@theme`). Also look at how small apps handle this without a full component library — the goal is consistency, not abstraction for its own sake. Read the current codebase to understand what patterns have emerged since this issue was filed.
