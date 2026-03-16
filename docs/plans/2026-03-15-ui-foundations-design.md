# Cycle 13 — UI Foundations (#033 + #036)

Establish design tokens and view state patterns. Views were rewritten for data access in cycle 11; now standardize their styling and loading/error/empty representation. All subsequent UX issues (#032, #034, #035) inherit these patterns.

## 1. Design Tokens (#033)

**Mechanism:** Tailwind CSS v4 `@theme` directive in `main.css`. Defines CSS custom properties that Tailwind utilities consume automatically (e.g., `bg-primary` maps to `--color-primary`).

**Palette — Warm Stone (solarized-inspired, cyan/green primary, dark magenta secondary):**

| Token | Value | Usage |
|-------|-------|-------|
| `--color-bg` | `#e8e4dc` | Page background |
| `--color-surface` | `#ddd8ce` | Card/list item backgrounds |
| `--color-surface-raised` | `#ece8e0` | Elevated surfaces (inputs, form areas) |
| `--color-primary` | `#0891b2` | Cyan — primary buttons, links |
| `--color-primary-hover` | `#0e7490` | Cyan darkened for hover |
| `--color-secondary` | `#65a30d` | Green — success accents, deck indicators |
| `--color-accent` | `#86195e` | Dark magenta — review CTA, emphasis |
| `--color-accent-hover` | `#701a50` | Magenta darkened for hover |
| `--color-text` | `#4a4a4a` | Primary text |
| `--color-text-muted` | `#6b6560` | Secondary text |
| `--color-text-faint` | `#9b958e` | Hint/placeholder text |
| `--color-border` | `#d0cabe` | Input borders, dividers |

**Typography:** Use Tailwind's default type scale (`text-sm`, `text-lg`, `font-semibold`, etc.). No custom typography tokens — the defaults are fine for this app's scale.

**Spacing/shape tokens:**

| Token | Value | Usage |
|-------|-------|-------|
| `--radius-sm` | `6px` | Buttons, inputs, small elements |
| `--radius-md` | `8px` | Cards, containers |
| `--shadow-card` | `0 1px 3px rgba(0,0,0,0.08)` | Card/list item elevation |

**Accessibility note:** `--color-text-faint` (#9b958e) on `--color-bg` (#e8e4dc) is ~2.5:1 contrast, below WCAG AA. This is acceptable for hint text that supplements visible labels/content, but must never be the sole means of conveying information. A full accessibility pass is out of scope for this cycle.

**Where it lives:** All in `main.css` via `@theme`. No component library, no extracted components. Views use Tailwind utility classes that reference the tokens. Changing a color in one place changes it everywhere.

## 2. View State Pattern (#036)

**Type:** A discriminated union `ViewState<T>` in a new `src/view-state.ts`:

```ts
type ViewState<T> =
  | { status: 'loading' }
  | { status: 'loaded'; data: T }
  | { status: 'empty' }
  | { status: 'error'; message: string }
  | { status: 'not-found' }
```

### Composable: `useLiveViewState`

Subscribes to Dexie's `liveQuery` observable directly (does **not** wrap `useLiveQuery`). This is because `useLiveQuery` has no error callback — `useLiveViewState` needs error handling built in.

Signature: `useLiveViewState<T>(querier: () => Promise<T>, isEmpty: (data: T) => boolean): Ref<ViewState<T>>`

- Starts as `{ status: 'loading' }`
- On first emission: transitions to `loaded` or `empty` based on the `isEmpty` predicate
- On error: transitions to `{ status: 'error', message }`
- Cleans up subscription via `onScopeDispose`

The `isEmpty` predicate is caller-supplied so the composable doesn't need to know whether `T` is an array, an object, or undefined.

### Per-view integration

**HomeView** — Single query (all decks). One `useLiveViewState` call with `isEmpty: (decks) => decks.length === 0`.

**DeckView** — Three queries with different semantics:
- `deck` query → determines `not-found` (undefined) vs loaded. Uses a dedicated `useLiveViewState` call with a sentinel: if the querier resolves to `undefined`, map to `not-found` instead of `empty`. Implementation: DeckView's primary state is driven by the deck query. A helper `useLiveViewStateOrNotFound<T>` variant (or a flag/option on `useLiveViewState`) handles the undefined → not-found mapping.
- `cards` query → separate `useLiveViewState` for the card list, determines the empty-cards sub-state.
- `dueCards` query → remains a plain `useLiveQuery` (just a count for the review button, no view state needed).

The template switches on `deckState.status` first (loading/error/not-found), then on `cardsState.status` for the content area (loaded/empty).

**ReviewView** — Does not use `useLiveQuery`; it loads cards once imperatively via `onMounted`. Integration path: wrap the `load()` call in a `ref<ViewState<Card[]>>` managed manually:
- Initialize as `{ status: 'loading' }`
- After `getDueCards` resolves: set to `loaded` or `empty`
- On catch: set to `{ status: 'error', message }`
- The existing `currentIndex`, `flipped`, `done` refs remain alongside the view state (they manage review session state, not data-loading state)

### Skeleton loading

Each view gets a skeleton template — grey shimmer bars matching the layout shape. Shown during `loading` state. Implemented as inline template blocks (no extracted components, per architecture convention). Uses a CSS `@keyframes pulse` animation defined once in `main.css`.

### Empty states

Centered layout with a Lucide icon, title, and hint text.

| State | Icon | Title | Hint | `data-testid` |
|-------|------|-------|------|---------------|
| No decks | `Library` | No decks yet | Create your first deck to start studying | `empty-state` |
| No cards | `Layers` | No cards yet | Add cards using the form above | `no-cards` |
| Not found | `Search` | Deck not found | This deck may have been deleted | `not-found` |
| Error | `AlertTriangle` | Something went wrong | Could not load your data | `error-state` |
| Loading | — | — | — | `loading-skeleton` |

### Not-found handling

DeckView and ReviewView check for invalid deck ID → `not-found` state with a "Back to decks" link.

## 3. Icon Integration

**Package:** `lucide-vue-next`. Tree-shakeable — only imported icons are bundled. Each view imports the specific icons it needs directly.

## 4. Test Plan

**Unit tests** (`tests/unit/`):
- `ViewState` type: verify discriminated union narrows correctly via type-level tests
- `useLiveViewState`: test loading → loaded, loading → empty, loading → error transitions using a mock observable. Test that `isEmpty` predicate is respected. Test subscription cleanup.

**Behavioral tests** (`tests/behavior/`):
- Each view renders skeleton during loading state
- Each view renders correct empty state (icon + message) when data is empty
- DeckView renders not-found state for invalid deck ID
- Error state renders with message and retry/back link

**Integration tests** (`tests/integration/`):
- Full data path: create deck → view deck → see loaded state (not empty, not skeleton)
- Navigate to nonexistent deck ID → not-found state displayed

## 5. File Changes

| File | Change |
|------|--------|
| `main.css` | Add `@theme` block with all tokens, pulse keyframe animation |
| `src/view-state.ts` | New — `ViewState<T>` type, `useLiveViewState` composable |
| `src/use-live-query.ts` | No change — `useLiveViewState` subscribes to `liveQuery` directly |
| `src/views/HomeView.vue` | Use token classes, view state, skeleton, empty state with Lucide icon |
| `src/views/DeckView.vue` | Use token classes, dual view state (deck + cards), skeleton, empty + not-found |
| `src/views/ReviewView.vue` | Use token classes, manual view state ref, skeleton, done/empty states |
| `package.json` | Add `lucide-vue-next` dependency |

## 6. What This Unlocks

All subsequent UX issues (#032 action feedback, #034 navigation, #035 rating labels) inherit the token system and view state pattern — no more ad-hoc styling or state handling per view.
