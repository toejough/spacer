# Fix mobile tab overflow hiding the Help tab

## Why

On mobile-width screens the app has five tabs (Review, Todos, Notes, Search, Help) laid out in a horizontal `overflow-x:auto` bar. Only about four tabs fit the viewport, so the Help tab is scrolled completely off the right edge with no scroll affordance (no fade, no arrow, no visible scrollbar on touch). A user on a phone cannot discover Help without knowing to swipe the tab bar — it looks broken and low quality.

## What Changes

- Move the Help entry out of the scrollable tab bar and into a persistent icon button in the header, so it is always visible regardless of screen width.
- Keep Review, Todos, Notes, and Search as the primary in-bar tabs (four tabs fit comfortably on mobile).
- Tapping the header help button opens the existing Help view via `switchTab('help')`, preserving the current help content and styling.
- Make the header help button reflect active state when the Help view is shown, and restore tab-bar active state when a content tab is selected.
- Bump the asset version (`v21` -> `v22`) in `index.html` asset query strings and footer, and the cache name in `sw.js`, so the service worker serves fresh assets.

## Capabilities

### New Capabilities

- `navigation`: Captures the cross-cutting tab/view navigation behavior that was previously implicit, including the requirement that the Help entry remain reachable on all screen sizes.

### Modified Capabilities

None of the existing capabilities (`note-taking`, `search`, `spaced-review`, `todo-list`) change behavior; only the navigation affordance around them changes.

## Impact

- Affected code: `srv/templates/index.html` (move Help tab to a header button, bump version), `srv/static/style.css` (header help-button styles, version bump), `srv/static/script.js` (help button opens Help view and syncs active state, version bump), `srv/static/sw.js` (cache version bump).
- No API or database changes.
- No impact on the inactive `web-app` Quasar prototype.
