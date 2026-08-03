## Why

The Notes tab's count badge shares the same orange (`--warning`) styling as the Review and Todos badges. For Review and Todos, orange correctly signals "action needed" (items due, open todos). For Notes, the count is just informational — how many notes exist — and orange makes it read as if something needs attention when nothing does.

## What Changes

- Give `#tabBadgeNotes` a neutral color distinct from `.tab-badge`'s default orange, while leaving Review, Todos, and Search badges unchanged.

## Capabilities

No spec-level behavior change — this is a visual styling adjustment only. `skip_specs: true` is set in `.openspec.yaml`.

## Impact

- `srv/static/style.css` only. No JS, HTML structure, or behavior changes.
