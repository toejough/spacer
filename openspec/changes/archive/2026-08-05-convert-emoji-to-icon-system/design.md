## Context

See proposal.md for motivation. One fact settles most of the "how": the existing icons aren't an original design — `completeIcon`, `abandonIcon`, `reopenIcon`, and `editIcon` in `script.js` are byte-for-byte [Feather Icons](https://feathericons.com) path data (`check-circle`, `x-circle`, `rotate-ccw`, `edit-2`), embedded as inline SVG strings, MIT-licensed. That's the established convention this change extends, not a fresh choice.

Feather itself has no "brain" icon, which matters for the header/favicon mark. [Lucide](https://lucide.dev) is Feather's direct, actively-maintained successor — same MIT license, same 24×24/`stroke-width:2`/round-cap visual language, kept Feather's existing paths unchanged for icons it inherited, and added many more including `brain`. It's the natural place to source everything this change needs without breaking visual consistency with what's already there.

## Goals / Non-Goals

**Goals:**
- One icon language app-wide: Feather/Lucide-style inline SVG, `stroke="currentColor"` so icons inherit the theme like the existing set already does.
- Stay inline, not separate icon files — matches the existing embedding pattern and avoids touching `sw.js`'s `PRECACHE` list (new asset files would need adding there for offline support; inline SVG in already-cached HTML/JS doesn't).
- A real favicon and PWA icon, replacing the emoji-drawn-as-SVG-text hack in `manifest.json`.

**Non-Goals:**
- `srv/templates/welcome.html`'s emoji — dead code, tracked separately as GitHub issue #27.
- Redesigning the icon *system* (sizing, hover states, `.btn-icon` CSS) — those already work; this only swaps glyphs.

## Decisions

### 1. Source icons from Lucide (Feather's successor), inline SVG, hand-matched to the existing style

Pull path data for each needed icon from Lucide during implementation (`help-circle`, `clipboard-list`, `check-square`, `notebook-pen` or `file-text`, `search`, `party-popper`, `scissors`, `lightbulb`, `save`, `download`, `upload`, `chevron-down`, `chevron-right`, `grip-vertical`, `brain`) and embed them the same way the existing icons are embedded — inline `<svg>` strings in `index.html`/`script.js`, not `<img>` tags or a separate sprite file. If a specific glyph can't be sourced (no network access, icon renamed), hand-author a simple stroke-based path in the same `viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"` style rather than falling back to something visually inconsistent (an emoji, a filled icon, a different stroke width).

**Alternative considered**: an icon font (e.g., self-hosted Feather webfont). Rejected — inline SVG is what's already there, needs no extra `@font-face`/file, and avoids a flash-of-unstyled-icon on load.

### 2. Tab icons and their matching help-heading icons are the same icon, reused

`📋 Review` (tab) and `📋 Review` (help heading) currently use the identical emoji for the same concept. Reuse one SVG constant per concept (e.g., one `reviewIcon` used in both the tab button and the help heading), the same way `editIcon` is already reused across multiple buttons — not two independently-authored SVGs that could drift.

### 3. Favicon and manifest icon are the same asset, delivered as a data: URI

Replace `manifest.json`'s `data:image/svg+xml,<svg...><text>🧠</text></svg>` with a real inline SVG brain icon (Lucide's `brain`, adapted to the app's theme color) — still a `data:` URI, still no separate file, so `sw.js`'s `PRECACHE` doesn't need to change. Add `<link rel="icon" type="image/svg+xml" href="data:image/svg+xml,...">` to `index.html`'s `<head>` using the same SVG, so the browser tab and the PWA icon are provably the same asset, not two hand-maintained copies.

**Alternative considered**: a separate `/static/favicon.svg` file. Rejected for this change — would need adding to `PRECACHE` and a cache-busting `?v=NN` like the other static assets, more moving parts than a one-icon inline hack needs. Worth revisiting if the icon set grows large enough that inlining becomes unwieldy.

### 4. `apple-touch-icon` is added on a best-effort basis, not guaranteed

Some older Safari/iOS versions handle `data:` URI `apple-touch-icon` links inconsistently. Add it anyway (`<link rel="apple-touch-icon" href="data:image/svg+xml,...">`) since it costs nothing and is strictly better than today's total absence of one — but the PWA install icon (via `manifest.json`, which works reliably) remains the primary, tested path, not the `apple-touch-icon` link.

## Risks / Trade-offs

- **Hand-sourcing SVG path data during implementation without live verification against Lucide's actual files** → the fallback (hand-author in the same stroke style) keeps visual consistency even if a specific glyph can't be exactly matched; a slightly-off custom icon in the right style is a smaller risk than breaking consistency.
- **`apple-touch-icon` via `data:` URI may not render on some older iOS Safari versions** → acceptable; not a regression (nothing renders there today), and the manifest-based PWA icon is the reliable path.
- **Reworking `renderNoteCard` to match `renderTodoCard`'s icon usage touches active, tested UI code** → keep the edit strictly to icon substitution (swap the emoji strings for the existing `editIcon`/`abandonIcon` SVG constants), not a broader refactor of the two parallel render functions — that consolidation is a separate, larger concern not in scope here.
