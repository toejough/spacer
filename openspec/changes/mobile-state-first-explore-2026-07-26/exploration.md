# Exploration notes and decision matrix

## Decision factors
- Speed: how quickly a user can mark Done
- Safety: protection against accidental Abandon
- Discoverability: are actions obvious to new users
- Accessibility: keyboard / screen reader support
- Implementation complexity: frontend and backend changes

## Quick comparison
- Variant A: Visible button + long-press
  - Speed: high (visible Done)
  - Safety: high (Abandon behind menu + confirm)
  - Discoverability: high
  - Complexity: moderate
- Variant B: Swipe-first
  - Speed: highest
  - Safety: medium (swipes can be accidental) — mitigated by undo
  - Discoverability: medium
  - Complexity: low–moderate
- Variant C: Action-sheet
  - Speed: medium
  - Safety: high
  - Discoverability: low–medium
  - Complexity: low

## Metrics to collect in experiment
- Time-to-first-action when viewing a swimlane
- Fraction of actions that are Abandon with/without immediate undo
- Rate of undone actions per cohort (undo indicates accidental)
- Conversion: items moved to Done vs Abandoned

## Recommended experiment setup
- Implement Variant A as default; enable Variant B for an experiment group (A/B)
- Track the metrics above for 2–4 weeks with minimum N users
- Safety: Log counts of Abandon confirmations and undo events

