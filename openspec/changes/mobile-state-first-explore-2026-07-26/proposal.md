# Explore: Mobile-first state-focused Todo UX (equal accessibility for Done & Abandon)

## Goal
Explore mobile interactions that treat "Done" (completed) and "Abandon" (decided not to do) as equally accessible states — both should be reachable quickly while remaining safe and discoverable. Also ensure items can be reopened from either completed state.

## Why
Some users need to mark items as "decided not to do" as quickly as they mark them Done; hiding Abandon reduces discoverability and slows workflows. We must balance parity of access with safeguards against accidental destructive actions. Additionally, users must be able to quickly reopen items marked Done or Abandoned so that state changes are reversible and low-friction.

## Variants to explore (updated)
- Variant A: Two visible buttons per card — Done and Abandon — with distinct affordances and undo. When a card is in Done or Abandoned state it shows a primary Reopen action.
- Variant B: Split-swipe: swipe-right → Done, swipe-left → Abandon (both show undo; Abandon may require confirmation depending on sensitivity). Swiping an already-closed card left/right could reveal a Reopen affordance.
- Variant C: Two-finger swipe / hold for Abandon + single-finger swipe for Done (less discoverable; optional).
- Variant D: Reveal both actions on drag (drag left reveals Abandon, drag right reveals Done) — visually symmetric; closed cards show Reopen affordance.

## Success criteria
- Both Done and Abandon reachable in ≤2 taps for 90% of cases.
- Reopen reachable in ≤2 taps from Done or Abandon for 95% of cases.
- Accidental Abandon reduced via confirmation/undo to ≤5% of Abandon events.
- Accessibility: both actions operable via keyboard and screen reader.
- Metrics: parity in time-to-action between Done and Abandon, acceptable accidental rates.

