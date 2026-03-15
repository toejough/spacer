# 032 — No feedback for user actions

**Status:** open
**Type:** ux / prevention
**Source:** #015 UX/design premortem

## Problem

User actions (create deck, add card, rate card) succeed silently — the item appears in a list, but there's no explicit confirmation, error message, or loading indicator. At current scale (3 actions) this is tolerable because the UI updates immediately. As actions multiply (delete, import, export, sync, share), some will be slow or fallible, and users won't know if something worked, failed, or is still processing.

There are also no error states anywhere — if a Dexie operation fails, the UI silently does nothing.

## Principle

Users should always know what happened in response to their action. The feedback pattern should be consistent across the app so that new features inherit it rather than inventing their own.

## Guidance

Before implementing, review current UX best practices for feedback in mobile-first PWAs — toast patterns, optimistic UI, inline validation, and error boundaries. Consider what's proportional for a small app vs. what creates a foundation for growth. Read the current codebase to understand what's changed since this issue was filed.
