Title: UI polish for Abandon action (icon/ARIA/focus)

Created: 2026-07-26
Assignee: TBD

Description:
We changed the app behavior to use an "Abandon" action and prevented abandoning completed todos. The next step is a UI polish so the control communicates intent clearly and is accessible.

Goals:
- Replace emoji trash/flag with a neutral SVG icon (abandon.svg) in list and modal.
- Add aria-label="Abandon" to buttons and ensure screen-readers announce the action.
- Ensure keyboard focus styles and contrast meet accessibility guidelines.
- Keep behavior unchanged: Abandon should only be available for active todos.

Attachments:
- Proposal: openspec/changes/ui-abandon-visuals/proposal.md
- Tasks: openspec/changes/ui-abandon-visuals/tasks.md
- Screenshot(s): (to be added)

Priority: medium

