---
name: openspec
description: Use when the user wants to plan, scope, or implement changes using OpenSpec spec-driven development in a project that has an openspec/ directory.
---

# OpenSpec skill

OpenSpec (https://openspec.dev) is a lightweight, spec-driven development framework.
When the project has an `openspec/` directory, use OpenSpec for any non-trivial change.

## When to use

- The user asks for a new feature, change, or refactor.
- The project already contains `openspec/specs/` or `openspec/config.yaml`.
- The user explicitly mentions OpenSpec, specs, or proposals.

## Workflow

1. **Read existing specs**: Check `openspec/specs/` for relevant capabilities. Use `openspec list --specs` if needed.
2. **Create a change proposal**: Run `openspec new change <kebab-name>` (or use the `/opsx propose <name>` slash hook if available) with a short description.
3. **Write the artifacts** in `openspec/changes/<name>/`:
   - `proposal.md` — why, what changes, capabilities impacted.
   - `specs/<capability>/spec.md` — delta specs using `## ADDED`, `## MODIFIED`, `## REMOVED`, or `## RENAMED Requirements`.
   - `design.md` — architectural decisions (optional for small changes).
   - `tasks.md` — trackable checkboxes (`- [ ] N.M description`).
4. **Implement** the code changes, marking tasks as complete as you go.
5. **Validate and archive**: Run `openspec validate --all` then `openspec archive -y <name>` to merge the deltas into `openspec/specs/`.
6. **Commit** both the code and the OpenSpec artifacts together.

## Conventions

- Use SHALL/MUST in requirements. Use `#### Scenario: <name>` with WHEN/THEN/AND.
- One capability per `openspec/specs/<capability>/spec.md` folder.
- Keep the first change small; archive frequently so specs stay current.
- If the project has no `openspec/` directory yet, run `openspec init .` first and add project context to `openspec/config.yaml`.
