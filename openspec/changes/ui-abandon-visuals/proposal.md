# UI polish for Abandon action

## Why

We updated the app behavior to use an "Abandon" action instead of "Delete" and limited it to active todos. The next step is a focused UI polish so the control communicates intent clearly across platforms and assistive tech:
- Replace the trash icon with a neutral/marker icon (SVG) instead of an emoji to avoid 'delete' connotations.
- Ensure consistent appearance across platforms (SVG + CSS instead of emoji glyphs).
- Add accessible labels (aria-label) and keyboard focus styles.
- Update modal and list button placement and visual affordance so Abandon is distinct from permanent deletion.

## What Changes

This change is a UI-only refinement: provide an SVG icon, aria-labels, styling, and small layout tweaks. No behavior changes.

