# Exploration notes and decision matrix (updated)

## Updated decision factors
- Parity: how equally accessible Done and Abandon are
- Speed and discoverability for both actions
- Safety: protections against accidental Abandon
- Accessibility and implementation complexity

## Variant comparison (focused on parity)
- Variant A (Two visible buttons)
  - Parity: excellent (both visible)
  - Speed: high
  - Safety: medium — Abandon visible, so require confirm or undo
  - Accessibility: high
  - Complexity: low
- Variant B (Split-swipe)
  - Parity: high (left/right symmetric)
  - Speed: high
  - Safety: medium — swipes can be accidental; require undo and optional confirm for Abandon
  - Accessibility: medium — must provide keyboard alternatives
  - Complexity: moderate
- Variant D (Reveal on drag)
  - Parity: high, visually symmetric
  - Speed: medium
  - Safety: high if reveal requires an additional confirm tap
  - Accessibility: medium
  - Complexity: moderate

## Recommended approach for exploration
- Prototype Variant A (two visible buttons) and Variant B (split-swipe) for comparison in usability tests.
- For Abandon safety: use confirm modal for high-sensitivity items or immediate undo snackbar with an extended undo window (6–8s). Track undo rates to tune.
- Ensure keyboard/ARIA alternatives: each card has focusable Done and Abandon buttons; announce actions via ARIA live region.

## Metrics & events (suggested)
- event: todo.action_initiated {action: 'done'|'abandon', method: 'button'|'swipe'|'menu', todo_id}
- event: todo.action_confirmed {action, todo_id}
- event: todo.action_undone {action, todo_id}
- metric: undo_rate = undone_actions / actions
- metric: accidental_abandon_rate ≈ undo_rate for abandon actions

