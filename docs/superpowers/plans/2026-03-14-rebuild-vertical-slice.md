# Rebuild: Vertical Slice Bootstrap — Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rebuild Spacer PWA as a single vertical slice — a user can create a deck, add a card, and review it, with SM-2 state persisting to IndexedDB.

**Architecture:** Single-page Vue 3 app with three views (Home, Deck, Review) connected to Dexie.js for IndexedDB persistence. SM-2 is a pure function. No stores — views query Dexie directly to keep the slice thin. One integration test exercises the full data path.

**Tech Stack:** Vue 3 + TypeScript, Vite, Tailwind CSS v4, Dexie.js, Vitest + happy-dom + fake-indexeddb

---

## File Structure

```
spacer/
├── CLAUDE.md                    # Project conventions and build commands
├── index.html                   # SPA entry point
├── package.json                 # Dependencies and npm scripts
├── vite.config.ts               # Vite + Vue + Tailwind plugins
├── vitest.config.ts             # Vitest with happy-dom + fake-indexeddb setup
├── tsconfig.json                # Project references
├── tsconfig.app.json            # App TypeScript config (strict: true)
├── tsconfig.node.json           # Build tools TypeScript config
├── env.d.ts                     # Vue SFC type declarations
├── dev/
│   ├── dev                      # Start Vite dev server
│   ├── test                     # Run Vitest
│   ├── build                    # Type-check + production build
│   └── check                    # Type-check + tests (CI command)
└── src/
    ├── main.ts                  # App entry — creates app, router, mounts
    ├── main.css                 # Tailwind import
    ├── App.vue                  # Shell — header + router-view
    ├── db.ts                    # Dexie DB instance + Deck/Card interfaces
    ├── sm2.ts                   # SM-2 algorithm (pure function)
    ├── test-setup.ts            # fake-indexeddb/auto import for tests
    ├── views/
    │   ├── HomeView.vue         # List decks, create new deck
    │   ├── DeckView.vue         # Show cards in deck, add card, start review
    │   └── ReviewView.vue       # Flip cards, rate 1-5, SM-2 update
    └── __tests__/
        └── flow.test.ts         # Integration test: full user journey
```

**Key decisions:**
- No Pinia stores — views query Dexie directly. Stores are an abstraction layer with no value until we need cross-component reactivity.
- No component files — all UI is in the three views. Extract components later when there's a real reason.
- Single test file — one integration test that proves the vertical slice works. No unit tests on isolated pieces (those gave false confidence last time).

## Chunk 1: Scaffold + DB + SM-2

### Task 1: Project scaffold

**Files:**
- Create: `package.json`, `index.html`, `vite.config.ts`, `vitest.config.ts`, `tsconfig.json`, `tsconfig.app.json`, `tsconfig.node.json`, `env.d.ts`, `src/main.css`, `src/test-setup.ts`
- Create: `dev/dev`, `dev/test`, `dev/build`, `dev/check`

- [ ] **Step 1: Create `package.json`**

```json
{
  "name": "spacer",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc --noEmit && vite build",
    "test": "vitest run",
    "check": "vue-tsc --noEmit && vitest run"
  }
}
```

- [ ] **Step 2: Install dependencies**

```bash
npm install vue vue-router dexie
npm install -D vite @vitejs/plugin-vue typescript vue-tsc vitest @vue/test-utils happy-dom fake-indexeddb @tailwindcss/vite tailwindcss
```

- [ ] **Step 3: Create build configs**

`vite.config.ts`:
```ts
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [vue(), tailwindcss()],
});
```

`vitest.config.ts`:
```ts
import { defineConfig } from "vitest/config";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: "happy-dom",
    setupFiles: ["./src/test-setup.ts"],
  },
});
```

`tsconfig.json`:
```json
{
  "files": [],
  "references": [
    { "path": "./tsconfig.app.json" },
    { "path": "./tsconfig.node.json" }
  ]
}
```

`tsconfig.app.json`:
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "jsx": "preserve",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "esModuleInterop": true,
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "skipLibCheck": true,
    "noEmit": true
  },
  "include": ["src/**/*.ts", "src/**/*.vue", "env.d.ts"]
}
```

`tsconfig.node.json`:
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "isolatedModules": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "noEmit": true
  },
  "include": ["vite.config.ts", "vitest.config.ts"]
}
```

`env.d.ts`:
```ts
/// <reference types="vite/client" />

declare module "*.vue" {
  import type { DefineComponent } from "vue";
  const component: DefineComponent<{}, {}, any>;
  export default component;
}
```

`index.html`:
```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Spacer</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
```

`src/main.css`:
```css
@import "tailwindcss";
```

`src/test-setup.ts`:
```ts
import "fake-indexeddb/auto";
```

- [ ] **Step 4: Create dev scripts**

`dev/dev`:
```bash
#!/usr/bin/env bash
set -euo pipefail
exec npx vite
```

`dev/test`:
```bash
#!/usr/bin/env bash
set -euo pipefail
exec npx vitest run "$@"
```

`dev/build`:
```bash
#!/usr/bin/env bash
set -euo pipefail
npx vue-tsc --noEmit && npx vite build
```

`dev/check`:
```bash
#!/usr/bin/env bash
set -euo pipefail
npx vue-tsc --noEmit && npx vitest run
```

```bash
chmod +x dev/dev dev/test dev/build dev/check
```

- [ ] **Step 5: Verify scaffold**

Run: `npx vite build 2>&1 | tail -3`
Expected: Build succeeds (may warn about missing entry, that's fine — no src/main.ts yet)

- [ ] **Step 6: Commit**

```
git add package.json package-lock.json index.html vite.config.ts vitest.config.ts tsconfig.json tsconfig.app.json tsconfig.node.json env.d.ts src/main.css src/test-setup.ts dev/
git commit -m "feat: scaffold Vite + Vue 3 + TypeScript + Tailwind project"
```

### Task 2: Database and SM-2

**Files:**
- Create: `src/db.ts`, `src/sm2.ts`

- [ ] **Step 1: Create `src/db.ts`**

```ts
import Dexie, { type EntityTable } from "dexie";

export interface Deck {
  id: number;
  name: string;
  createdAt: Date;
}

export interface Card {
  id: number;
  deckId: number;
  front: string;
  back: string;
  easeFactor: number;
  interval: number;
  repetitions: number;
  nextReview: Date;
}

export class SpacerDB extends Dexie {
  decks!: EntityTable<Deck, "id">;
  cards!: EntityTable<Card, "id">;

  constructor(name = "spacer") {
    super(name);
    this.version(1).stores({
      decks: "++id, name",
      cards: "++id, deckId, nextReview",
    });
  }
}

export const db = new SpacerDB();
```

Note: `SpacerDB` class is exported so tests can create isolated instances with unique names. The default `db` export is for production use.

- [ ] **Step 2: Create `src/sm2.ts`**

```ts
export interface SM2State {
  easeFactor: number;
  interval: number;
  repetitions: number;
}

export interface SM2Result extends SM2State {
  nextReview: Date;
}

export function sm2(state: SM2State, quality: number, now = new Date()): SM2Result {
  const q = Math.max(0, Math.min(5, Math.round(quality)));

  let { easeFactor, interval, repetitions } = state;

  if (q < 3) {
    repetitions = 0;
    interval = 0;
  } else {
    if (repetitions === 0) {
      interval = 1;
    } else if (repetitions === 1) {
      interval = 6;
    } else {
      interval = Math.round(interval * easeFactor);
    }
    repetitions += 1;
  }

  easeFactor = Math.max(
    1.3,
    easeFactor + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
  );

  const nextReview = new Date(now);
  nextReview.setDate(nextReview.getDate() + interval);

  return { easeFactor, interval, repetitions, nextReview };
}
```

- [ ] **Step 3: Commit**

```
git add src/db.ts src/sm2.ts
git commit -m "feat: add Dexie DB schema and SM-2 algorithm"
```

## Chunk 2: App + Views (the vertical slice)

### Task 3: App shell and router

**Files:**
- Create: `src/main.ts`, `src/App.vue`

- [ ] **Step 1: Create `src/main.ts`**

```ts
import { createApp } from "vue";
import { createRouter, createWebHistory } from "vue-router";
import App from "./App.vue";
import "./main.css";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", component: () => import("./views/HomeView.vue") },
    { path: "/deck/:id", component: () => import("./views/DeckView.vue") },
    { path: "/review/:deckId", component: () => import("./views/ReviewView.vue") },
  ],
});

createApp(App).use(router).mount("#app");
```

- [ ] **Step 2: Create `src/App.vue`**

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-indigo-600 text-white p-4">
      <router-link to="/" class="text-xl font-bold">Spacer</router-link>
    </header>
    <main class="max-w-lg mx-auto p-4">
      <router-view />
    </main>
  </div>
</template>
```

- [ ] **Step 3: Commit**

```
git add src/main.ts src/App.vue
git commit -m "feat: add app shell with router"
```

### Task 4: HomeView — list and create decks

**Files:**
- Create: `src/views/HomeView.vue`

- [ ] **Step 1: Create `src/views/HomeView.vue`**

```vue
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { db, type Deck } from "../db";

const decks = ref<Deck[]>([]);
const newDeckName = ref("");

async function loadDecks() {
  decks.value = await db.decks.toArray();
}

async function createDeck() {
  const name = newDeckName.value.trim();
  if (!name) return;
  await db.decks.add({ name, createdAt: new Date() } as Deck);
  newDeckName.value = "";
  await loadDecks();
}

onMounted(loadDecks);
</script>

<template>
  <div>
    <h2 class="text-lg font-semibold mb-4">Your Decks</h2>

    <form @submit.prevent="createDeck" class="flex gap-2 mb-4">
      <input
        v-model="newDeckName"
        placeholder="New deck name"
        class="flex-1 border rounded px-3 py-2"
        data-testid="deck-name-input"
      />
      <button
        type="submit"
        class="bg-indigo-600 text-white px-4 py-2 rounded"
        data-testid="create-deck-btn"
      >
        Create
      </button>
    </form>

    <p v-if="decks.length === 0" class="text-gray-500" data-testid="empty-state">
      No decks yet. Create one above.
    </p>

    <ul class="space-y-2">
      <li v-for="deck in decks" :key="deck.id">
        <router-link
          :to="`/deck/${deck.id}`"
          class="block p-3 bg-white rounded shadow hover:shadow-md"
          :data-testid="`deck-${deck.id}`"
        >
          {{ deck.name }}
        </router-link>
      </li>
    </ul>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```
git add src/views/HomeView.vue
git commit -m "feat: add HomeView — list and create decks"
```

### Task 5: DeckView — show cards, add card, start review

**Files:**
- Create: `src/views/DeckView.vue`

- [ ] **Step 1: Create `src/views/DeckView.vue`**

```vue
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { db, type Card, type Deck } from "../db";

const route = useRoute();
const deckId = Number(route.params.id);

const deck = ref<Deck | undefined>();
const cards = ref<Card[]>([]);
const front = ref("");
const back = ref("");

async function load() {
  deck.value = await db.decks.get(deckId);
  cards.value = await db.cards.where("deckId").equals(deckId).toArray();
}

async function addCard() {
  const f = front.value.trim();
  const b = back.value.trim();
  if (!f || !b) return;
  await db.cards.add({
    deckId,
    front: f,
    back: b,
    easeFactor: 2.5,
    interval: 0,
    repetitions: 0,
    nextReview: new Date(),
  } as Card);
  front.value = "";
  back.value = "";
  await load();
}

function dueCount() {
  const now = new Date();
  return cards.value.filter((c) => c.nextReview <= now).length;
}

onMounted(load);
</script>

<template>
  <div v-if="deck">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold">{{ deck.name }}</h2>
      <router-link
        :to="`/review/${deckId}`"
        v-if="dueCount() > 0"
        class="bg-indigo-600 text-white px-4 py-2 rounded text-sm"
        data-testid="start-review-btn"
      >
        Review ({{ dueCount() }} due)
      </router-link>
    </div>

    <form @submit.prevent="addCard" class="space-y-2 mb-4">
      <input
        v-model="front"
        placeholder="Front (question)"
        class="w-full border rounded px-3 py-2"
        data-testid="card-front-input"
      />
      <input
        v-model="back"
        placeholder="Back (answer)"
        class="w-full border rounded px-3 py-2"
        data-testid="card-back-input"
      />
      <button
        type="submit"
        class="bg-indigo-600 text-white px-4 py-2 rounded"
        data-testid="add-card-btn"
      >
        Add Card
      </button>
    </form>

    <p v-if="cards.length === 0" class="text-gray-500" data-testid="no-cards">
      No cards yet. Add one above.
    </p>

    <ul class="space-y-2">
      <li
        v-for="card in cards"
        :key="card.id"
        class="p-3 bg-white rounded shadow"
        :data-testid="`card-${card.id}`"
      >
        <div class="font-medium">{{ card.front }}</div>
        <div class="text-gray-500 text-sm">{{ card.back }}</div>
      </li>
    </ul>

    <router-link to="/" class="inline-block mt-4 text-indigo-600 text-sm">
      &larr; Back to decks
    </router-link>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```
git add src/views/DeckView.vue
git commit -m "feat: add DeckView — show cards, add card, start review"
```

### Task 6: ReviewView — flip, rate, SM-2 update

**Files:**
- Create: `src/views/ReviewView.vue`

- [ ] **Step 1: Create `src/views/ReviewView.vue`**

```vue
<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { db, type Card } from "../db";
import { sm2 } from "../sm2";

const route = useRoute();
const router = useRouter();
const deckId = Number(route.params.deckId);

const dueCards = ref<Card[]>([]);
const currentIndex = ref(0);
const flipped = ref(false);
const done = ref(false);

const current = computed(() => dueCards.value[currentIndex.value]);

async function load() {
  const now = new Date();
  const all = await db.cards.where("deckId").equals(deckId).toArray();
  dueCards.value = all.filter((c) => c.nextReview <= now);
  if (dueCards.value.length === 0) {
    done.value = true;
  }
}

function flip() {
  flipped.value = true;
}

async function rate(quality: number) {
  const card = current.value;
  if (!card) return;

  const result = sm2(
    {
      easeFactor: card.easeFactor,
      interval: card.interval,
      repetitions: card.repetitions,
    },
    quality
  );

  await db.cards.update(card.id, {
    easeFactor: result.easeFactor,
    interval: result.interval,
    repetitions: result.repetitions,
    nextReview: result.nextReview,
  });

  flipped.value = false;
  if (currentIndex.value < dueCards.value.length - 1) {
    currentIndex.value++;
  } else {
    done.value = true;
  }
}

onMounted(load);
</script>

<template>
  <div>
    <div v-if="done" data-testid="review-done">
      <h2 class="text-lg font-semibold mb-2">Review Complete</h2>
      <p class="text-gray-600 mb-4">
        Reviewed {{ dueCards.length }} card{{ dueCards.length === 1 ? "" : "s" }}.
      </p>
      <button
        @click="router.push(`/deck/${deckId}`)"
        class="bg-indigo-600 text-white px-4 py-2 rounded"
      >
        Back to Deck
      </button>
    </div>

    <div v-else-if="current">
      <p class="text-sm text-gray-500 mb-2" data-testid="review-progress">
        Card {{ currentIndex + 1 }} of {{ dueCards.length }}
      </p>

      <div class="bg-white rounded shadow p-6 min-h-[200px] flex items-center justify-center">
        <div class="text-center">
          <p class="text-xl" data-testid="card-front">{{ current.front }}</p>
          <p
            v-if="flipped"
            class="mt-4 text-gray-600 border-t pt-4"
            data-testid="card-back"
          >
            {{ current.back }}
          </p>
        </div>
      </div>

      <div v-if="!flipped" class="mt-4 text-center">
        <button
          @click="flip"
          class="bg-indigo-600 text-white px-6 py-2 rounded"
          data-testid="flip-btn"
        >
          Show Answer
        </button>
      </div>

      <div v-else class="mt-4 flex justify-center gap-2">
        <button
          v-for="q in [1, 2, 3, 4, 5]"
          :key="q"
          @click="rate(q)"
          class="px-4 py-2 rounded border hover:bg-gray-100"
          :data-testid="`rate-${q}`"
        >
          {{ q }}
        </button>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Commit**

```
git add src/views/ReviewView.vue
git commit -m "feat: add ReviewView — flip card, rate, SM-2 update"
```

## Chunk 3: Integration test + CLAUDE.md + docs

### Task 7: Integration test

**Files:**
- Create: `src/__tests__/flow.test.ts`

The test must exercise the full data path without monkey-patching module exports (which fails with ES module getters). Instead, create a test-scoped `SpacerDB` instance and interact with it directly — this tests the same code paths the views use.

- [ ] **Step 1: Create `src/__tests__/flow.test.ts`**

```ts
import { describe, it, expect, beforeEach } from "vitest";
import { SpacerDB } from "../db";
import { sm2 } from "../sm2";

let testDbCounter = 0;

function createTestDB() {
  return new SpacerDB(`test-spacer-${++testDbCounter}-${Date.now()}`);
}

describe("Full flow: create deck → add card → review → SM-2 update", () => {
  let db: SpacerDB;

  beforeEach(() => {
    db = createTestDB();
  });

  it("completes the entire user journey", async () => {
    // Create a deck
    const deckId = await db.decks.add({
      name: "Test Deck",
      createdAt: new Date(),
    } as any);

    // Add a card with default SM-2 state
    const now = new Date();
    const cardId = await db.cards.add({
      deckId,
      front: "What is 2+2?",
      back: "4",
      easeFactor: 2.5,
      interval: 0,
      repetitions: 0,
      nextReview: now,
    } as any);

    // Verify card is queryable by deckId
    const cards = await db.cards.where("deckId").equals(deckId).toArray();
    expect(cards).toHaveLength(1);
    expect(cards[0].front).toBe("What is 2+2?");

    // Filter due cards (nextReview <= now)
    const dueCards = cards.filter((c) => c.nextReview <= now);
    expect(dueCards).toHaveLength(1);

    // Review: apply SM-2 with quality=4 (good)
    const card = dueCards[0];
    const result = sm2(
      {
        easeFactor: card.easeFactor,
        interval: card.interval,
        repetitions: card.repetitions,
      },
      4
    );

    // Persist SM-2 result
    await db.cards.update(card.id, {
      easeFactor: result.easeFactor,
      interval: result.interval,
      repetitions: result.repetitions,
      nextReview: result.nextReview,
    });

    // Verify persisted state
    const updated = await db.cards.get(card.id);
    expect(updated!.repetitions).toBe(1);
    expect(updated!.interval).toBe(1);
    expect(updated!.easeFactor).toBeGreaterThanOrEqual(2.5);
    expect(updated!.nextReview.getTime()).toBeGreaterThan(now.getTime());

    // Card should no longer be due (nextReview is tomorrow)
    const stillDue = (await db.cards.where("deckId").equals(deckId).toArray())
      .filter((c) => c.nextReview <= now);
    expect(stillDue).toHaveLength(0);
  });

  it("handles failed review (quality < 3) by resetting repetitions", async () => {
    const deckId = await db.decks.add({
      name: "Fail Deck",
      createdAt: new Date(),
    } as any);

    await db.cards.add({
      deckId,
      front: "Hard question",
      back: "Hard answer",
      easeFactor: 2.5,
      interval: 6,
      repetitions: 2,
      nextReview: new Date(),
    } as any);

    const card = (await db.cards.where("deckId").equals(deckId).toArray())[0];
    const result = sm2(
      {
        easeFactor: card.easeFactor,
        interval: card.interval,
        repetitions: card.repetitions,
      },
      1 // fail
    );

    await db.cards.update(card.id, {
      easeFactor: result.easeFactor,
      interval: result.interval,
      repetitions: result.repetitions,
      nextReview: result.nextReview,
    });

    const updated = await db.cards.get(card.id);
    expect(updated!.repetitions).toBe(0);
    expect(updated!.interval).toBe(0);
  });
});
```

- [ ] **Step 2: Run tests**

Run: `dev/test`
Expected: 2 tests pass

- [ ] **Step 3: Commit**

```
git add src/__tests__/flow.test.ts
git commit -m "test: add integration test for full create-review flow"
```

### Task 8: CLAUDE.md and issue/status docs

**Files:**
- Create: `CLAUDE.md`
- Create: `docs/issues/012-rebuild-vertical-slice.md`
- Modify: `docs/status.md`

- [ ] **Step 1: Create `CLAUDE.md`**

```markdown
# Spacer PWA

Spaced repetition flashcard app. Vue 3 + TypeScript + Vite + Dexie (IndexedDB) + Tailwind CSS v4.

## Build Commands

```bash
dev/dev          # Start Vite dev server
dev/test         # Run Vitest
dev/build        # Type-check + production build
dev/check        # Type-check + tests
```

## Conventions

- **TDD always** — write failing tests first, then minimum code to pass
- **Vertical slices** — every increment delivers a working user-facing feature
- **No empty stubs** — only create files with real, working code
- **Integration tests** — at least one test per feature that exercises the full data path
- **Commits** — use `/commit`, conventional commits format
- **Issues** — `docs/issues/{number}-{slug}.md`
- **Status** — `docs/status.md` updated every cycle

## Tech Notes

- DB: Dexie with EntityTable for type-safe IndexedDB
- SM-2: pure function in `src/sm2.ts`
- Test env: happy-dom + fake-indexeddb
- Each test gets its own DB instance (no shared mutable state)
```

- [ ] **Step 2: Create issue file `docs/issues/012-rebuild-vertical-slice.md`**

```markdown
# Rebuild: vertical slice bootstrap

**Status:** in-progress
**Priority:** p0
**Labels:** feature
**Created:** 2026-03-14
**Closed:**

## Description
Rebuild the Spacer PWA with a vertical slice: create deck, add card, review with SM-2.

## Acceptance Criteria
- [ ] User can create a deck, add a card, and review it
- [ ] SM-2 state persists to IndexedDB
- [ ] Integration test verifies full flow
- [ ] targ (dev/ scripts) as build interface
- [ ] Project CLAUDE.md exists

## Notes
Replaces original bootstrap (#1). See docs/prompts/2026-03-14-reset-and-rebuild.md.
```

- [ ] **Step 3: Update `docs/status.md`**

```markdown
# Spacer — Project Status

**Last updated:** 2026-03-14
**Current increment:** 2
**Streak:** 0

## Done
- #1 Bootstrap project with full stack (2026-03-14) — reverted in reset

## Closed (wont-fix)
- #2-#11 — superseded by rebuild

## In Progress
- #12 Rebuild: vertical slice bootstrap

## Up Next

## Blocked
```

- [ ] **Step 4: Run full check**

Run: `dev/check`
Expected: Type-check passes, 2 tests pass

- [ ] **Step 5: Commit**

```
git add CLAUDE.md docs/issues/012-rebuild-vertical-slice.md docs/status.md
git commit -m "docs: add CLAUDE.md, create rebuild issue #12, update status"
```

### Task 9: Close issue and retro

- [ ] **Step 1: Update issue #12 status to done**

Set `**Status:** done` and `**Closed:** 2026-03-14` in `docs/issues/012-rebuild-vertical-slice.md`.
Check all AC boxes.

- [ ] **Step 2: Update `docs/status.md`**

Move #12 to Done, increment streak to 1.

- [ ] **Step 3: Append retro to `docs/retros.md`**

```markdown
### Increment #2: Rebuild vertical slice bootstrap — pass
**What worked:** Light process (plan → execute → review) kept the rebuild focused and fast.
**What to improve:** [Fill in based on actual experience.]
**Action items:**
- [Fill in based on actual experience.]
```

- [ ] **Step 4: Commit**

```
git add docs/issues/012-rebuild-vertical-slice.md docs/status.md docs/retros.md
git commit -m "docs: close #12 rebuild, update status and retro"
```
