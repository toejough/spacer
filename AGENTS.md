# Agent Instructions

This is a Go web application, running on the phone-llm platform.
See README.md for details on the structure and components.

## Queued work

Run `openspec list` to see what's queued, and `openspec status --change <name>` for details on a
specific one.

## OpenSpec workflow

This project uses OpenSpec (https://openspec.dev) for spec-driven development.
Specs are stored in `openspec/specs/` and change proposals in `openspec/changes/`.

When the user asks for a new feature or change:

1. Read the relevant specs in `openspec/specs/` first to understand current behavior.
2. Use `/opsx propose <name>` or `openspec new change <name>` to create a change proposal.
3. Write `proposal.md`, `design.md` (if needed), and `specs/<capability>/spec.md` deltas.
4. Mark tasks in `tasks.md` as complete as work progresses.
5. Run `/opsx validate` or `openspec validate --all` and then `openspec archive <name>` when done.
6. Commit the OpenSpec artifacts alongside the code.

Use the Makefile targets `openspec-status`, `openspec-validate`, and `openspec-list` to inspect the OpenSpec state.

## Bump the app version marker on every static-asset change

Whenever `srv/static/script.js`, `srv/static/style.css`, or `srv/templates/index.html` changes,
bump the shared version number in lockstep across every place it appears. Don't judge whether the
change "counts" as a behavior change — bump on any edit to those files, full stop. That judgment
call is exactly what got missed once already (see below).

Locations (all must carry the same number):
- `srv/templates/index.html`: the `?v=NN` query string on the stylesheet link, the manifest link,
  the `script.js` tag, and the service-worker registration URL; plus the visible
  `<footer class="app-version">vNN</footer>` text.
- `srv/static/sw.js`: `CACHE_NAME = 'remember-everything-vNN'` and the two `?v=NN` entries in
  `PRECACHE`.

Why this matters: the version number is the user's only signal for whether they're looking at a
stale cached copy or the current code. If it doesn't move, an absent behavior reads as "not
implemented yet" instead of "reload the page" — indistinguishable without checking git history.
This was missed once already (the `export-import-stacks` change, 2026-08-08) and had to be
fixed after the fact.

## Cross-capability integration check

Before scoping a change to a single capability, grep `openspec/specs/*/spec.md` for other
capabilities that touch the same data or make a broad claim over it (e.g. "export all data" in
`data-portability`, "search all items" in `search`). A capability that adds new state without
checking whether those broad-coverage capabilities need updating too creates a silent gap.

This already happened once: `card-stacks` added stacks as new local state, but `data-portability`
export/import was never updated to include them — backups silently dropped every stack until
`export-import-stacks` fixed it. When a change adds or reshapes state, treat "does anything else
claim to cover all of this?" as part of scoping the change, not an afterthought.
