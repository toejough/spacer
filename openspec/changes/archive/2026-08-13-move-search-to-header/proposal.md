## Why

On narrow viewports the primary tab row (Review, Todos, Notes, Search) overflows its available width, and Search — last in the row — is pushed fully off-screen with no visible hint that a fourth tab exists. This directly violates the existing `navigation` spec requirement that content tabs stay reachable without hidden horizontal overflow. Search is also conceptually not a peer of Review/Todos/Notes: it's an invoked action/mode, not a content destination, matching how mobile apps generally treat search (an icon in the app bar) versus content tabs (Gmail, Twitter/X, Instagram, Google Play, iOS Mail). This app already implements exactly that pattern for Help, which lives as a header icon button rather than a tab row entry.

## What Changes

- Remove Search from the primary tab row (`nav.tabs`). Review, Todos, and Notes remain as content tabs and fit within mobile-width tab bars without overflow.
- Add a Search icon button to `.header-actions`, next to the existing Help button, reusing the same open/active mechanism Help already uses (view swap, active state, aria-pressed).
- Clicking the Search button opens the same Search view that exists today (search input, submit, results list) — the search behavior itself (`search` capability) is unchanged.
- The Search button, Help button, and content tabs become mutually exclusive active states: opening Search clears Help's and the content tabs' active state, and vice versa.
- Drop the ambient Search tab badge. It only ever reflected a nonzero count while already viewing Search results, so as an ambient "glance from elsewhere" indicator it was always empty and hidden — consistent with Help, which has no badge.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
- `navigation`: the "content tabs" that must fit on mobile without overflow become Review/Todos/Notes only (Search is removed from that set); a new reachability scenario for the Search button (mirroring Help's) is added; the active-state requirement is extended to cover mutual exclusivity between Search, Help, and the content tabs.

## Impact

- `srv/templates/index.html`: move the Search button markup out of `nav.tabs` into `.header-actions`; the `#tab-search` section and its contents are unchanged.
- `srv/static/style.css`: style the Search header button like `.help-btn`; drop now-unused `#tabBadgeSearch` badge styling (or leave the shared `.tab-badge` rule intact since other tabs still use it).
- `srv/static/script.js`: extend `switchTab` to toggle a `searchBtn` active state the same way it toggles `helpBtn`; stop writing/showing the `tabBadgeSearch` badge in `updateTabCounts`.
- `srv/static/script.test.js`: update any tests asserting Search is part of the tab row or has a badge.
- `openspec/specs/navigation/spec.md`: delta spec updates the two affected requirements.
