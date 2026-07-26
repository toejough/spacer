# Explore: Mobile-first state-focused Todo UX

## Goal
Explore alternatives to the 3-column board for mobile: prioritize safe, discoverable, and efficient ways to change todo state (Todo → Done / Abandoned) while preserving the time-based swimlane grouping for context.

## Why
Columns are great on desktop but fragment the mobile experience. We want a mobile-first interaction that:
- Prevents accidental destructive actions (Abandon)
- Makes common actions (Done) fast and discoverable
- Preserves temporal grouping for review (week/month/quarter/year)
- Supports bulk operations and accessibility

## Variants to explore
- Variant A (recommended): Stacked swimlane list + visible Done button + long-press/overflow menu for Abandon + undo snackbar
- Variant B (swipe-first): Swipe-right to Done (with undo), overflow menu for Abandon (no swipe for Abandon)
- Variant C (action-sheet): Long-press or tap to open action-sheet with Done/Abandon and metadata (compact, discoverable for power users)

## Success criteria
- Users can mark a todo Done in ≤2 taps for 90% of cases
- Abandon requires explicit confirm or provides undo for 95% of accidental cases
- Accessibility: all actions reachable by keyboard and screen reader
- Metrics: reduction in accidental Abandon events vs current baseline

