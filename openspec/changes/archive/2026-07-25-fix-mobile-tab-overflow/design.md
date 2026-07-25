# Design: Fix mobile tab overflow hiding the Help tab

## Problem

The tab bar is `<nav class="tabs">` containing five `<button class="tab">` elements styled with `display:flex; gap:4px; overflow-x:auto; white-space:nowrap`. Each tab is ~90-110px wide with icon + label. On a 390px mobile viewport only ~4 tabs fit; the Help tab (the 5th) is scrolled off the right. Mobile browsers hide the horizontal scrollbar, and there is no fade or arrow affordance, so the overflow is invisible. The result: Help is effectively undiscoverable on phones.

## Solution

Remove the Help tab from the scrollable bar and promote it to a persistent icon button in the header `<header>`, next to the review badge. The header already uses `justify-content:space-between` and wraps, so a help button sits naturally on the right. The four remaining content tabs (Review, Todos, Notes, Search) fit comfortably on mobile without overflow.

Tapping the help button calls the existing `switchTab('help')`, which shows `#tab-help` exactly as today. We add a small amount of active-state syncing:

- `switchTab` already toggles the `.active` class on `.tab` buttons. We extend it to also toggle an `.active` class on the header help button (`#helpBtn`) so it highlights when Help is shown, and clears it when a content tab is selected.
- The help button is given `aria-pressed` reflecting the active state for accessibility.

A CSS edge-fade mask is added to the tab bar as a progressive enhancement: when the tab bar can scroll (content wider than container), a subtle gradient fade appears on the trailing edge. This is a minor enhancement for robustness if a fourth tab ever overflows on very narrow screens, but the primary fix is reducing the tab count to four.

### Changes in `srv/templates/index.html`

1. Remove the `<button class="tab" data-tab="help" ...>Help</button>` from the `<nav class="tabs">`.
2. Add a help button to the header, e.g. `<button id="helpBtn" class="help-btn" onclick="switchTab('help')" aria-label="Help" aria-pressed="false">❓</button>`, placed so it sits at the right of the header.
3. Bump `?v=21` to `?v=22` on the stylesheet, manifest, script, and sw.js references, and update the `<footer>` version text to `v22`.

### Changes in `srv/static/style.css`

1. Add `.help-btn` styling: a circular/rounded icon button matching the header scale, with an `.active` state (accent color / background) used when Help is open.
2. Add a `mask-image` (or `::after` gradient fade) to `.tabs` for a trailing-edge scroll hint. Use `-webkit-mask-image` for Safari.
3. Bump version if a version comment exists (none currently; no change needed beyond the query string).

### Changes in `srv/static/script.js`

1. In `switchTab(tab)`, after toggling `.tab` active states, toggle the help button active state: `helpBtn.classList.toggle('active', tab === 'help')` and set `aria-pressed` accordingly. (The existing `.tab` querySelectorAll loop already ignores the removed help tab, so no other change is needed there.)
2. No change to `currentTab` logic or `refreshCurrent`.

### Changes in `srv/static/sw.js`

1. Bump `CACHE_NAME` from `remember-everything-v21` to `remember-everything-v22` and update the precache `?v=` query strings to `v22`.

### Tests

No new automated tests are required for a layout/CSS change, but we verify manually via the mobile-emulated browser that:

- All four content tabs are visible without horizontal scroll on a 390px viewport.
- The help button is visible in the header and opens the Help view.
- The help button shows an active state while Help is open and clears when a content tab is tapped.
- The existing JS unit tests (`node ./srv/static/script.test.js`) and Go tests (`go test ./...`) still pass.

## Out of scope

- The inactive `web-app/` Quasar prototype is not changed.
- No tab content, data model, or review logic changes.
- No backend/API changes.
