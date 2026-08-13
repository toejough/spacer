## Context

`switchTab(tab)` in `srv/static/script.js` already special-cases one non-tab-row view: it looks up `#helpBtn` by id and toggles its `active` class / `aria-pressed` whenever `tab === 'help'`, separately from the generic `.tab[data-tab=...]` active-state loop that drives Review/Todos/Notes/Search today. Search is about to become a second view of the same kind — a header icon button that toggles a `.tab-content` section without being a row of `.tab` elements. See proposal.md - Why / What Changes.

## Goals / Non-Goals

**Goals:**
- Reuse the Help button's existing open/active mechanism for Search rather than inventing a second one.
- Keep the change confined to markup relocation + the small generalization needed to support two header-action buttons instead of one.

**Non-Goals:**
- No inline-expanding search bar (the Gmail/Material "search view" pattern) — out of scope per proposal, this is the "swap to a dedicated view" approach matching Help.
- No change to search behavior itself (`search` capability: matching, stack results) — only where its trigger lives.
- No change to Help's content, Review/Todos/Notes behavior, or other capabilities.

## Decisions

**Generalize header-action active-state handling instead of duplicating it.**
`switchTab` currently has one hardcoded block for `helpBtn`. Copy-pasting a near-identical block for `searchBtn` would double the maintenance surface for what is, behaviorally, the same pattern (per the new `navigation` delta, Help and Search now have symmetric active-state requirements). Instead, `switchTab` should toggle active/aria-pressed for both header-action buttons from a small local id list (e.g. `['help', 'search']` mapped to `['helpBtn', 'searchBtn']`), so a future third header action doesn't require a third copy-pasted block.
- Alternative considered: duplicate the `helpBtn` block verbatim for `searchBtn`. Rejected — it's the exact thing the new spec requirements say should behave identically, and duplication is how the two silently drift later (e.g. one gets a fix the other doesn't).

**Reuse the `.help-btn` CSS class for the Search button rather than a parallel class.**
The two buttons are visually identical (circular icon button in `.header-actions`). Applying the existing `.help-btn` class to both (keeping `#searchBtn` only as an id for JS hooks) avoids defining a second, easily-diverging copy of the same rules. If a future change gives the two buttons different visual treatment, that's the point at which to split the class — not before.

**Remove `#tabBadgeSearch` entirely rather than hide it.**
Per proposal, the badge never conveyed a real ambient count (it only reflected the currently-open Search view's own result count). Delete the badge element and its `updateTabCounts` entry rather than leaving dead markup/logic that displays "0" forever.

**Bump the static-asset cache-buster (`?v=42` → `?v=43`) on `style.css`, `script.js`, and `sw.js`/`manifest.json` references in `index.html`.**
This is the existing convention for shipping frontend changes (see current `?v=42` on all four references) — not a new decision, just carried forward so the deployed app actually picks up the new markup/JS/CSS instead of serving cached copies.

## Risks / Trade-offs

- [Two icon buttons plus the H1 title could crowd a very narrow header] → `header` already uses `flex-wrap`/`gap` and was verified at 375px width with the existing single Help button; two circular 38px buttons plus badge-free tabs was confirmed to fit during exploration. If a third header action is ever added, re-verify at the same width.
- [Existing JS tests or fixtures may assert Search is `.tab[data-tab="search"]` or check `#tabBadgeSearch`] → covered in tasks.md; `srv/static/script.test.js` needs a pass to update any such assertions.

## Open Questions

(none — the interaction pattern, markup location, and badge removal were all resolved during exploration; see proposal.md.)
