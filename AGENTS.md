# Agent Instructions

This is a Go web application, running on the phone-llm platform.
See README.md for details on the structure and components.

## OpenSpec workflow

This project uses OpenSpec (https://openspec.dev) for spec-driven development.
Specs are stored in `openspec/specs/` and change proposals in `openspec/changes/`.

A Shelley slash hook is installed at `~/.config/shelley/hooks/slash/opsx` (source in `shelley-hooks/slash/opsx`).
The user can run commands like `/opsx specs`, `/opsx validate`, and `/opsx propose <name>` directly in the chat.
If the hook is missing, install it with `./shelley-hooks/install.sh`.

When the user asks for a new feature or change:

1. Read the relevant specs in `openspec/specs/` first to understand current behavior.
2. Use `/opsx propose <name>` or `openspec new change <name>` to create a change proposal.
3. Write `proposal.md`, `design.md` (if needed), and `specs/<capability>/spec.md` deltas.
4. Mark tasks in `tasks.md` as complete as work progresses.
5. Run `/opsx validate` or `openspec validate --all` and then `openspec archive <name>` when done.
6. Commit the OpenSpec artifacts alongside the code.

Use the Makefile targets `openspec-status`, `openspec-validate`, and `openspec-list` to inspect the OpenSpec state.
