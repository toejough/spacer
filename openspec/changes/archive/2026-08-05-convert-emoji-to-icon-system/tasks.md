## 1. Fix the existing SVG/emoji inconsistency (no new icons needed)

- [x] 1.1 `renderNoteCard` (script.js): replace the raw `🏳️` abandon button with the same `abandonIcon`/reopen-icon pattern `renderTodoCard` already uses
- [x] 1.2 `renderNoteCard` (script.js): replace the raw `✏️` edit button with the existing `editIcon` SVG
- [x] 1.3 Stack-rename button (script.js line ~234): replace the raw `✏️` with the existing `editIcon` SVG
- [x] 1.4 Verify note cards and stack-rename now visually match todo cards (same icon, same size, same hover behavior) — no new CSS needed since `.btn-icon svg{width:20px;height:20px}` already covers it

## 2. Source or author the new icons

- [x] 2.1 Source Lucide path data (or hand-author in the matching stroke style if unavailable) for: `help-circle`, `clipboard-list`, `check-square`, a notes icon, `search`, `party-popper`, `scissors`, `lightbulb`, `save`, `download`, `upload`, `chevron-down`, `chevron-right`, `grip-vertical`, `brain`
- [x] 2.2 Confirm every sourced/authored icon uses `viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"`, matching the existing icons exactly

## 3. Convert index.html

- [x] 3.1 Header brand mark: replace `🧠` with the `brain` icon
- [x] 3.2 Help button: replace `❓` with `help-circle`
- [x] 3.3 Four tab buttons: replace `📋 ✅ 📝 🔍` with `clipboard-list` / `check-square` / notes icon / `search`
- [x] 3.4 Four matching help-content headings: reuse the exact same icon constants as the tabs (per design.md decision 2), not independently-added copies
- [x] 3.5 Empty-review state: replace `🎉` with `party-popper`
- [x] 3.6 Help-content headings: replace `✂️` (Cloze Deletions) with `scissors`, `💡` (Tips) with `lightbulb`, `💾` (Backup & Restore) with `save`
- [x] 3.7 Export/Import buttons: replace `⬇️`/`⬆️` with `download`/`upload`

## 4. Convert script.js structural icons

- [x] 4.1 Expand/collapse control: replace `▾`/`▸` with `chevron-down`/`chevron-right`
- [x] 4.2 Drag handle: replace `⠿` with `grip-vertical`

## 5. Favicon and app icon

- [x] 5.1 Build the `brain` SVG as a `data:image/svg+xml` URI using the app's theme color
- [x] 5.2 Add `<link rel="icon" type="image/svg+xml" href="...">` to `index.html`'s `<head>`, using that same data URI
- [x] 5.3 Add `<link rel="apple-touch-icon" href="...">` using the same data URI (best-effort, per design.md decision 4)
- [x] 5.4 Replace `manifest.json`'s `icons[0].src` (the emoji-as-SVG-text hack) with the same data URI, so the favicon, apple-touch-icon, and PWA icon are provably the same asset

## 6. Verify

- [x] 6.1 `grep -rnoP '[\x{1F000}-\x{1FFFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}]'` over `srv/templates/index.html` and `srv/static/script.js` (excluding `welcome.html`, out of scope per proposal.md) returns nothing
- [x] 6.2 Load the app in a browser: header, tabs, help button, empty-review state, and every help-content heading render SVG icons, not emoji or missing-glyph boxes
- [x] 6.3 Browser tab shows the new favicon; confirm no console errors
- [x] 6.4 Note cards and todo cards render visually matching abandon/edit icons
- [x] 6.5 `make test` and `make test-js` still pass
- [x] 6.6 Confirm `/openspec/`, `/repo/`, `/static/script.test.js` still 404 and nothing else broke in `srv/server.go`'s serving (no code there should have changed)
