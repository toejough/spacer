# Tasks: fix-mobile-tab-overflow

## 1. Move Help out of the tab bar

- [x] 1.1 In `srv/templates/index.html`, remove the Help `<button class="tab" data-tab="help">` from `<nav class="tabs">`.
- [x] 1.2 Add a persistent help button to the `<header>`: `<button id="helpBtn" class="help-btn" onclick="switchTab('help')" aria-label="Help" aria-pressed="false">?</button>`.
- [x] 1.3 Verify the four remaining tabs (Review, Todos, Notes, Search) fit on a 390px viewport with no horizontal overflow.

## 2. Sync help button active state

- [x] 2.1 In `srv/static/script.js` `switchTab()`, toggle `.active` on `#helpBtn` (`tab === 'help'`) and update `aria-pressed`.
- [x] 2.2 Confirm the `.tab` active loop still clears all content tabs when Help is active.

## 3. Style the help button and tab-bar scroll hint

- [x] 3.1 In `srv/static/style.css`, add `.help-btn` base + `.active` styles consistent with the header.
- [x] 3.2 Add a trailing-edge fade `mask-image` to `.tabs` as a progressive scroll-overflow hint.

## 4. Bump version

- [x] 4.1 Update `srv/templates/index.html` `?v=21` -> `?v=22` on style.css, manifest, script.js, sw.js, and the footer text.
- [x] 4.2 Update `srv/static/sw.js` `CACHE_NAME` to `remember-everything-v22` and precache `?v=22`.

## 5. Verify

- [x] 5.1 Run `node ./srv/static/script.test.js` — all JS tests pass.
- [x] 5.2 Run `go test ./...` — all Go tests pass.
- [x] 5.3 Manually verify on mobile (390px) via browser: all content tabs visible, help button visible, help view opens, active state toggles correctly.
- [x] 5.4 Rebuild and restart the systemd service; confirm via the proxy.

## 6. OpenSpec

- [x] 6.1 Write `proposal.md`, `design.md`, `tasks.md`, and `specs/navigation/spec.md` delta.
- [x] 6.2 Run `openspec validate --all`.
- [x] 6.3 Archive the change with `openspec archive fix-mobile-tab-overflow`.
- [x] 6.4 Commit OpenSpec artifacts alongside code.
