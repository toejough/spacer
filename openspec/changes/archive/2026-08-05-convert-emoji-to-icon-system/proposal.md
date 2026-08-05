## Why

The app already has a real icon system — feather-style outline SVGs (24×24 viewBox, `stroke="currentColor"`, 20px rendered) used for todo-card actions (complete, abandon, reopen, edit, delete). But it's only half-applied: note cards render their abandon/edit buttons with raw emoji (🏳️, ✏️) instead of the same SVGs the todo cards already use, the stack-rename button does the same (✏️), and everything else in the app — the header brand, all four tab labels, the help button, the empty-review state, and every heading in the help content — is still emoji. There's no dedicated favicon or app icon either; the PWA manifest's only icon is a hack, an inline SVG that draws the 🧠 emoji as `<text>`.

Emoji render inconsistently across platforms/fonts, don't take the app's theme color, and — per the note-card/stack-rename gap — the app is already inconsistent with itself about which icon language it uses. Finishing the conversion makes the app visually coherent and gives it an actual icon identity, including a real favicon, instead of leaning on whatever emoji font the OS happens to supply.

## What Changes

- **Fix the existing inconsistency**: `renderNoteCard` (script.js) and the stack-rename button pick up the same `editIcon`/`abandonIcon` SVGs `renderTodoCard` already defines, instead of raw emoji.
- **Design and add new SVG icons**, in the same feather-style outline language as the existing set, for every remaining decorative/branding emoji spot in `srv/templates/index.html`: header brand mark (replacing 🧠), help button (❓), the four tab icons (📋 ✅ 📝 🔍) and their matching help-content headings, empty-review state (🎉), and the help-content icons (✂️ 💡 💾 ⬇️ ⬆️).
- **Convert the two Unicode-symbol icons** used as icons rather than typography — the expand/collapse triangles (▾/▸) and the drag-handle braille pattern (⠿) in `script.js` — to matching SVGs, for full consistency with "the current icon system," per the request.
- **Add a real favicon and app icon**: design a dedicated icon (based on the new header brand mark) as an SVG, wire it up with `<link rel="icon">` (and `apple-touch-icon`) in `index.html`, and replace `manifest.json`'s emoji-as-text SVG hack with the real icon.
- **Not in scope**: `srv/templates/welcome.html`'s emoji (🔑💻🤖🏠) — that template is unreferenced by any route or link (confirmed via grep), so it's dead code, not part of the live app. Filed as GitHub issue #27 to investigate and likely delete separately, rather than touched here.

## Capabilities

### New Capabilities
- `iconography`: the app SHALL present a single, consistent SVG icon system for all in-app iconography and the favicon/app icon — no emoji glyphs used as icons anywhere in the served app.

## Impact

- `srv/templates/index.html` (all emoji replaced with inline SVG, new `<link rel="icon">`/`apple-touch-icon`)
- `srv/static/script.js` (`renderNoteCard`, stack-rename button, expand/collapse triangles, drag-handle glyph)
- `srv/static/manifest.json` (icon replaced with the real favicon/app icon, no longer an emoji-as-text hack)
- `srv/static/style.css` (only if new icon classes need sizing beyond the existing `.btn-icon svg` rule)
- No change to `srv/server.go`, routes, or app behavior — purely visual/iconography.
