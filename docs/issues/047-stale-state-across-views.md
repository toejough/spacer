# 047 — Stale state after mutations across views

**Status:** closed
**Type:** quality / prevention
**Source:** #017 implementation premortem

## Problem

Views load data in `onMounted` and re-query after their own mutations, but there's no cross-view reactivity. Navigation between views shows whatever was last fetched on mount. As views multiply and display overlapping data (e.g., card counts on the home screen, stats dashboard, search results), mutations in one view leave others stale until the user manually navigates away and back.

## Principle

When multiple views display the same underlying data, mutations should be visible across views without requiring manual refresh. The mechanism should be proportional — this doesn't necessarily require a global store, but it does require some form of reactive data subscription.

## Guidance

Before implementing, research Dexie's `liveQuery` and its Vue integration — it provides reactive queries that auto-update when underlying data changes, which would solve this without a separate store layer. Consider how this interacts with #030 (data-access layer) — reactive queries could be the foundation of that layer. Read the current data flow and view structure to understand what's changed since this issue was filed.
