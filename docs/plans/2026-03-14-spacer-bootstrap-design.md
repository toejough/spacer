# Spacer PWA Bootstrap — Design Spec

**Date:** 2026-03-14
**Status:** Approved
**Scope:** Full scaffold — Vite + Vue 3 + Tailwind + Vitest + Dexie + PWA manifest + smoke test

## Tech Stack

| Layer | Choice |
|-------|--------|
| Framework | Vue 3 + Composition API + TypeScript |
| Build | Vite |
| Styling | Tailwind CSS |
| State | Pinia |
| Storage | Dexie.js (IndexedDB) |
| Testing | Vitest + @vue/test-utils |
| PWA | vite-plugin-pwa |
| Algorithm | SM-2 |

## Architecture

```
src/
├── App.vue                 # Shell + router-view
├── main.ts                 # Entry point + PWA registration
├── router/index.ts         # Routes (home, deck, review)
├── stores/
│   ├── decks.ts            # Deck CRUD
│   └── reviews.ts          # Review session state
├── db/
│   ├── index.ts            # Dexie instance + schema
│   └── models.ts           # Card, Deck, ReviewLog types
├── lib/
│   └── sm2.ts              # SM-2 algorithm (pure function)
├── components/
│   ├── CardEditor.vue      # Create/edit card form
│   ├── ReviewCard.vue      # Flip card during review
│   └── DeckList.vue        # Deck overview
├── views/
│   ├── HomeView.vue        # Dashboard: decks + due count
│   ├── DeckView.vue        # Cards in a deck
│   └── ReviewView.vue      # Active review session
└── assets/
    └── main.css            # Tailwind imports
```

## Data Model

```typescript
interface Card {
  id: string
  deckId: string
  front: string
  back: string
  easeFactor: number    // starts at 2.5
  interval: number      // days
  repetitions: number   // consecutive correct
  nextReview: Date
  createdAt: Date
}

interface Deck {
  id: string
  name: string
  description?: string
  createdAt: Date
}
```

## Screens

1. **Home** — List of decks, each showing name + cards due today. "Review" button when cards are due.
2. **Deck** — Cards in a deck, add/edit/delete.
3. **Review** — One card at a time. Tap to flip. Rate 1-5. Progress bar. Summary at end.

## SM-2 Algorithm

Pure function: `sm2(quality: number, card: SM2State) => SM2State`

- quality 0-2: reset repetitions, interval = 1 day
- quality 3+: increment repetitions, multiply interval by easeFactor
- easeFactor adjusted: `EF + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02))`, min 1.3

## Bootstrap Deliverables

1. Vite + Vue 3 + TypeScript project created
2. Tailwind CSS configured
3. Vitest configured with one passing smoke test
4. Dexie.js wired with Card + Deck schema
5. Vue Router with 3 route stubs
6. Pinia installed
7. PWA manifest + vite-plugin-pwa configured
8. `dev/` scripts: `dev/test` (run vitest), `dev/dev` (run dev server)
9. App runs, tests pass, PWA manifest serves

## Out of Scope

- Cloud sync, rich text, import/export, tags, search, settings
- These become future increment issues
