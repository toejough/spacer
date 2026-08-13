## 1. Markup (`srv/templates/index.html`)

- [x] 1.1 Remove the `<button class="tab" data-tab="search">...</button>` entry (with its `#tabBadgeSearch` badge span) from `nav.tabs`.
- [x] 1.2 Add a Search button to `.header-actions`, before or after `#helpBtn`, reusing the `.help-btn` class and the same attribute pattern (`id="searchBtn"`, `aria-label="Search"`, `aria-pressed="false"`, `title="Search"`, `onclick="switchTab('search')"`), using the existing search SVG icon (the one currently inside the removed tab button).
- [x] 1.3 Leave `<section id="tab-search">` and its contents (search input, submit button, `#searchResults`) unchanged — only the trigger moves.
- [x] 1.4 Bump the cache-busting query param (`?v=42` → `?v=43`) on the `style.css`, `manifest.json`, `script.js`, and `sw.js` references in `index.html` (and, since it's the same versioning scheme, `sw.js`'s own `CACHE_NAME` and `PRECACHE` entries and the footer's `v42` text).

## 2. Behavior (`srv/static/script.js`)

- [x] 2.1 In `switchTab`, replace the single hardcoded `helpBtn` active-state block with a loop over the header-action ids (`help` → `helpBtn`, `search` → `searchBtn`), toggling `.active` and `aria-pressed` for whichever one matches `tab`, per design.md's generalization decision.
- [x] 2.2 In `updateTabCounts`, remove the `searchCount` computation and the `{ id: 'tabBadgeSearch', count: searchCount }` entry from the `badges` array (badge element no longer exists in markup after task 1.1).

## 3. Styling (`srv/static/style.css`)

- [x] 3.1 Verify `.header-actions` (already `flex-wrap` + `gap`) accommodates two `.help-btn`-styled circular buttons without adjustment; add spacing only if visually cramped. (confirmed visually in task 5.1)
- [x] 3.2 Confirm no dangling CSS references `#tabBadgeSearch` or `.tab[data-tab="search"]` after the markup change — grep confirms zero matches in `style.css`.

## 4. Tests

- [x] 4.1 Confirm `srv/static/script.test.js` has no assertions tied to Search being a `.tab`/tab-row entry or to `#tabBadgeSearch` — grep confirms zero matches for `switchTab`/`tabBadge`/`helpBtn`/`searchBtn` in the test file.
- [x] 4.2 Run `make test-js` and `make test` to confirm no regressions — both pass (38/38 JS tests, Go build has no test files).

## 5. Manual verification

- [x] 5.1 Using a headless-browser pass at a 375px-wide viewport with one item each in Review, Todos, and Notes: confirmed `tabs.scrollWidth === tabs.clientWidth === 355` (no overflow) with tab texts `["Review 2", "Todos 1", "Notes 1"]`.
- [x] 5.2 Confirmed the Search button is visible in the header next to Help, and clicking it opens `#tab-search` (`display: block`) with the search input/results.
- [x] 5.3 Confirmed active-state mutual exclusivity: opening Help while Search is active clears `searchBtn`'s active class/aria-pressed and sets `helpBtn`'s; opening Search then selecting a content tab clears `searchBtn`'s active state and activates the tab. Computed style confirms `.help-btn.active` renders the accent-blue fill on `#searchBtn` exactly as it does on `#helpBtn`.
- [x] 5.4 Confirmed existing search behavior is unaffected: searching "milk" after relocating the trigger still returns 1 matching result.
